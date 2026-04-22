package cloud

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	vision "cloud.google.com/go/vision/apiv1"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"google.golang.org/api/compute/v1"
)

var (
	REDIS_HOST     = os.Getenv("REDIS_GCP_HOST")
	REDIS_PORT     = os.Getenv("REDIS_GCP_PORT")
	REDIS_PASSWORD = os.Getenv("REDIS_GCP_PASSWORD")

	REDIS_KEY_PAYLOAD = os.Getenv("REDIS_KEY_PAYLOAD")

	REDIS_CLIENT = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", REDIS_HOST, REDIS_PORT),
		Password: REDIS_PASSWORD,
	})
)

func ProcessImage(ctx context.Context, e GCSEvent) error {
	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create storage client: %v", err)
	}

	visionClient, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return err
	}
	defer visionClient.Close()

	// 1. Use Cloud Vision API to check the image
	image := vision.NewImageFromURI("gs://" + e.Bucket + "/" + e.Name)
	props, err := visionClient.DetectSafeSearch(ctx, image, nil)
	if err != nil {
		return err
	}

	// 2. Check if the image is likely to contain adult content, violence, medical content, or racy content
	// level: VERY_UNLIKELY, UNLIKELY, POSSIBLE, LIKELY, VERY_LIKELY <=> 1, 2, 3, 4, 5
	if props.Adult >= 3 || props.Violence >= 3 || props.Medical >= 3 || props.Racy >= 3 { // 3 là 'POSSIBLE' (Có khả năng)
		log.Printf("Detecting image violates community standards... Adult: %s, Violence: %s, Medical: %s, Racy: %s \nDeleting file %s from bucket %s", props.Adult, props.Violence, props.Medical, props.Racy, e.Name, e.Bucket)

		err = invalidateCache(ctx, e.Name)
		if err != nil {
			return fmt.Errorf("failed to invalidate cache: %v", err)
		}

		notifyViaRedis(ctx, "rejected", e.Name)

		// delete image immediately if it violates community standards
		return storageClient.Bucket(e.Bucket).Object(e.Name).Delete(ctx)
	}

	// 3. If safe to use, copy file to bucket processed
	dstBucket := os.Getenv("GOOGLE_CLOUD_STORAGE_BUCKET_PROCESSED")
	if dstBucket == "" {
		return fmt.Errorf("GOOGLE_CLOUD_STORAGE_BUCKET_PROCESSED is empty")
	}

	// Copy file from raw bucket to processed bucket
	srcReader, err := storageClient.Bucket(e.Bucket).Object(e.Name).NewReader(ctx)
	if err != nil {
		return fmt.Errorf("failed to create reader for source file: %v", err)
	}
	defer srcReader.Close()

	obj := storageClient.Bucket(dstBucket).Object(e.Name)
	wc := obj.NewWriter(ctx)
	defer wc.Close()

	if _, err := io.Copy(wc, srcReader); err != nil {
		return fmt.Errorf("failed to copy file to GCS: %v", err)
	}

	if err := wc.Close(); err != nil {
		return fmt.Errorf("failed to close writer: %v", err)
	}

	err = invalidateCache(ctx, e.Name)
	if err != nil {
		return fmt.Errorf("failed to invalidate cache: %v", err)
	}
	notifyViaRedis(ctx, "processed", e.Name)

	if err := storageClient.Bucket(e.Bucket).Object(e.Name).Delete(ctx); err != nil {
		log.Printf("warning: failed to delete raw object %s from bucket %s: %v", e.Name, e.Bucket, err)
	}

	return nil
}

func notifyViaRedis(ctx context.Context, status, filePath string) {

	//file name is UserID_uuid
	userID := strings.TrimSpace(strings.Split(filePath, "_")[0])
	if userID == "" {
		log.Printf("Failed to extract user ID from file path: %s", filePath)
		return
	}

	payload := jwt.MapClaims{
		"file_path": filePath,
		"user_id":   userID,
		"status":    status,
		"iss":       "image-processing-service-cloud-function",
		"exp":       time.Now().Add(5 * time.Minute).Unix(),
		"iat":       time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	signedToken, err := token.SignedString([]byte(REDIS_KEY_PAYLOAD))
	if err != nil {
		log.Printf("Failed to sign JWT: %v", err)
		return
	}

	err = REDIS_CLIENT.Publish(ctx, "image-processing-results", signedToken).Err()
	if err != nil {
		log.Printf("Failed to publish message to Redis: %v", err)
		return
	}

	log.Printf("Published message to Redis channel 'image-processing-results' for file %s with status %s", filePath, status)
}

func invalidateCache(ctx context.Context, filePath string) error {
	service, err := compute.NewService(ctx)
	if err != nil {
		return fmt.Errorf("failed to create compute service: %v", err)
	}

	op, err := service.UrlMaps.InvalidateCache(os.Getenv("PROJECT_ID"), os.Getenv("PROJECT_BALANCER_NAME"), &compute.CacheInvalidationRule{
		Path: os.Getenv("PROJECT_PATH_BUCKET_RAW") + filePath,
	}).Do()

	if err != nil {
		return err
	}

	log.Printf("Invalidating cache for file %s from CDN... Operation: %s", filePath, op.Name)

	return nil
}
