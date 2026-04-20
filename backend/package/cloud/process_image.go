package cloud

import (
	"context"
	"fmt"
	"io"
	"log"

	"cloud.google.com/go/storage"
	vision "cloud.google.com/go/vision/apiv1"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/utils"
	"google.golang.org/api/compute/v1"
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

		err = InvalidateCache(ctx, e.Name)
		if err != nil {
			return fmt.Errorf("failed to invalidate cache: %v", err)
		}

		// delete image immediately if it violates community standards
		return storageClient.Bucket(e.Bucket).Object(e.Name).Delete(ctx)
	}

	// 3. If safe to use, copy file to bucket processed
	dstBucket := utils.GetEnv("GOOGLE_CLOUD_STORAGE_BUCKET_PROCESSED", "")
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

	err = InvalidateCache(ctx, e.Name)
	if err != nil {
		return fmt.Errorf("failed to invalidate cache: %v", err)
	}

	return nil
}

func InvalidateCache(ctx context.Context, filePath string) error {
	service, err := compute.NewService(ctx)
	if err != nil {
		return fmt.Errorf("failed to create compute service: %v", err)
	}

	op, err := service.UrlMaps.InvalidateCache(utils.GetEnv("PROJECT_ID", ""), utils.GetEnv("PROJECT_BALANCER_NAME", ""), &compute.CacheInvalidationRule{
		Path: utils.GetEnv("PROJECT_PATH_BUCKET_RAW", "") + filePath,
	}).Do()

	if err != nil {
		return err
	}

	log.Printf("Invalidating cache for file %s from CDN... Operation: %s", filePath, op.Name)

	return nil
}
