package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/storage"
	vision "cloud.google.com/go/vision/apiv1"
	auth_firebase "firebase.google.com/go/auth"

	"buf.build/go/protovalidate"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/client"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/db/sqlc"
	auth_proto "github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/gen/auth"
	user_proto "github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/gen/user"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/interceptor"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/models"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/repository"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/utils"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/validation"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"google.golang.org/api/compute/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type userService struct {
	user_proto.UnimplementedUserServiceServer
	user_repo    repository.UserRepository
	redis_memory *redis.Client
	validator    protovalidate.Validator

	auth_client *client.AuthClient

	auth_client_firebase *auth_firebase.Client
	vision_client        *vision.ImageAnnotatorClient

	compute_service *compute.Service
	gcs_client      *storage.Client
}

func NewUserService(user_repo repository.UserRepository, rdb *redis.Client, auth_client *client.AuthClient, auth_client_firebase *auth_firebase.Client, vision_client *vision.ImageAnnotatorClient, compute_service *compute.Service, gcs_client *storage.Client) *userService {
	v, err := protovalidate.New()
	if err != nil {
		panic(fmt.Sprintf("Failed to create validator: %v", err))
	}
	return &userService{
		user_repo:            user_repo,
		validator:            v,
		redis_memory:         rdb,
		auth_client:          auth_client,
		auth_client_firebase: auth_client_firebase,
		vision_client:        vision_client,
		compute_service:      compute_service,
		gcs_client:           gcs_client,
	}
}

func (s *userService) GetProfileByUserID(ctx context.Context, req *user_proto.GetProfileByUserIDRequest) (*user_proto.GetProfileByUserIDResponse, error) {
	caller := utils.GetCaller(ctx)

	if caller != ctx.Value(interceptor.CtxCallerKey).(string) {
		return nil, status.Errorf(codes.PermissionDenied, "Unauthorized: Caller in context does not match expected caller")
	}

	if req.UserId != ctx.Value(interceptor.CtxUserIDKey).(int64) {
		return nil, status.Errorf(codes.PermissionDenied, "Unauthorized: User ID in context does not match User ID in request")
	}

	if err := s.validator.Validate(req); err != nil {
		return nil, validation.BuildValidationError(err)
	}

	profileData, err := s.redis_memory.Get(ctx, fmt.Sprintf("user_profile:%d", req.UserId)).Bytes()
	if err == nil && len(profileData) > 0 {
		var cachedProfile models.ProfileRedis
		if err := json.Unmarshal(profileData, &cachedProfile); err == nil {
			log.Println("Profile data found in Redis in user-service for user ID: ", req.UserId)
			return &user_proto.GetProfileByUserIDResponse{
				Uuid:      cachedProfile.UserUUID.String(),
				Name:      cachedProfile.Name,
				Email:     cachedProfile.Email,
				Phone:     cachedProfile.Phone,
				AvatarUrl: cachedProfile.AvatarUrl,
				Birthday:  cachedProfile.Birthday,
			}, nil
		}
	}

	data, err := s.user_repo.GetProfileByUserID(ctx, req.UserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, validation.BuildBusinessError("USER_NOT_FOUND", "User not found with the given ID. Please check the user ID and try again.")
		}

		return nil, status.Errorf(codes.Internal, "Failed to get profile: %v", err)
	}

	var avatar_url, birthday *string

	if data.AvatarUrl.String == "" {
		avatar_url = nil
	} else {
		v := fmt.Sprintf("%s?v=%d", data.AvatarUrl.String, data.AvatarVersion)
		avatar_url = &v
	}

	if data.Birthday.Valid {
		birthdayStr := data.Birthday.Time.Format("2006-01-02")
		birthday = &birthdayStr
	} else {
		birthday = nil
	}

	profileRedis := models.ProfileRedis{
		UserUUID:  data.Uuid,
		Name:      &data.Name,
		Email:     data.Email.String,
		Phone:     &data.Phone.String,
		AvatarUrl: avatar_url,
		Birthday:  birthday,
	}

	profileBytes, err := json.Marshal(profileRedis)
	if err != nil {
		log.Printf("Failed to marshal profile data (in user service layer): %v", err)
	} else {
		if err := s.redis_memory.Set(ctx, fmt.Sprintf("user_profile:%d", req.UserId), profileBytes, 24*7*time.Hour).Err(); err != nil {
			log.Printf("Failed to set profile data in Redis (in user service layer): %v", err)
		} else {
			log.Println("Profile data cached in Redis for user ID: ", req.UserId)
		}
	}

	return &user_proto.GetProfileByUserIDResponse{
		Uuid:      data.Uuid.String(),
		Name:      &data.Name,
		Email:     data.Email.String,
		Phone:     &data.Phone.String,
		AvatarUrl: avatar_url,
		Birthday:  birthday,
	}, nil
}

func (s *userService) GetProfileByUserUUID(ctx context.Context, req *user_proto.GetProfileByUserUUIDRequest) (*user_proto.GetProfileByUserUUIDResponse, error) {
	caller := utils.GetCaller(ctx)

	if caller != ctx.Value(interceptor.CtxCallerKey).(string) {
		return nil, status.Errorf(codes.PermissionDenied, "Unauthorized: Caller in context does not match expected caller")
	}

	if req.UserId != ctx.Value(interceptor.CtxUserIDKey).(int64) {
		return nil, status.Errorf(codes.PermissionDenied, "Unauthorized: User ID in context does not match User ID in request")
	}

	if err := s.validator.Validate(req); err != nil {
		return nil, validation.BuildValidationError(err)
	}

	targetUUID, err := uuid.Parse(req.Uuid)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid UUID format: %v", err)
	}

	request := sqlc.GetProfileByUserUUIDParams{
		UserID: req.UserId,
		Uuid:   targetUUID,
	}

	data, err := s.user_repo.GetProfileByUserUUID(ctx, request)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, validation.BuildBusinessError("USER_NOT_FOUND", "User not found with the given UUID. Please check the UUID and try again.")
		}

		return nil, status.Errorf(codes.Internal, "Failed to get profile: %v", err)
	}

	return &user_proto.GetProfileByUserUUIDResponse{
		Name:      data.Name,
		Birthday:  data.Birthday.Time.Format("2006-01-02"),
		AvatarUrl: fmt.Sprintf("%s?v=%d", data.AvatarUrl.String, data.AvatarVersion),
	}, nil

}

func (s *userService) CreateProfile(ctx context.Context, req *user_proto.CreateProfileRequest) (*user_proto.CreateProfileResponse, error) {
	caller := utils.GetCaller(ctx)

	if caller != ctx.Value(interceptor.CtxCallerKey).(string) {
		return nil, status.Errorf(codes.PermissionDenied, "Unauthorized: Caller in context does not match expected caller")
	}

	if req.UserId != ctx.Value(interceptor.CtxUserIDKey).(int64) {
		return nil, status.Errorf(codes.PermissionDenied, "Unauthorized: User ID in context does not match User ID in request")
	}

	if err := s.validator.Validate(req); err != nil {
		return nil, validation.BuildValidationError(err)
	}

	if req.Birthday != "" {
		birthday, _ := time.Parse("2006-01-02", req.Birthday)

		today := time.Now().Truncate(24 * time.Hour)
		if birthday.After(today) {
			return nil, status.Errorf(codes.InvalidArgument, "Birthday cannot be in the future")
		}
	}

	_, err := s.user_repo.CreateProfile(ctx, sqlc.CreateProfileParams{
		UserID:    req.UserId,
		Name:      req.Name,
		Email:     utils.ConvertToPgTypeText(req.Email),
		AvatarUrl: utils.ConvertToPgTypeText(req.AvatarUrl),
		Birthday:  utils.ConvertToPgTypeDate(req.Birthday),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create profile: %v", err)
	}

	return &user_proto.CreateProfileResponse{
		Success: true,
		Message: "Profile created successfully",
	}, nil
}

func (s *userService) UpdateProfile(ctx context.Context, req *user_proto.UpdateProfileRequest) (*user_proto.UpdateProfileResponse, error) {
	caller := utils.GetCaller(ctx)

	if caller != ctx.Value(interceptor.CtxCallerKey).(string) {
		return nil, status.Errorf(codes.PermissionDenied, "Unauthorized: Caller in context does not match expected caller")
	}

	if req.UserId != ctx.Value(interceptor.CtxUserIDKey).(int64) {
		return nil, status.Errorf(codes.PermissionDenied, "Unauthorized: User ID in context does not match User ID in request")
	}

	log.Println("Req: ", req)

	if err := s.validator.Validate(req); err != nil {
		return nil, validation.BuildValidationError(err)
	}

	var phone string
	var isPhone bool
	if req.Phone != nil {
		phone, isPhone = utils.IsPhoneNumber(req.Phone)
		if !isPhone {
			return nil, status.Errorf(codes.InvalidArgument, "Invalid phone number format")
		}
	}

	profile, err := s.user_repo.UpdateProfile(ctx, sqlc.UpdateProfileByUserIdParams{
		UserID:   req.UserId,
		Name:     utils.ConvertToPgTypeTextPtr(req.Name),
		Birthday: utils.ConvertToPgTypeDatePtr(req.Birthday),
		Phone:    utils.ConvertToPgTypeText(phone),
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, validation.BuildBusinessError("USER_NOT_FOUND", "User not found with the given ID. Please check the user ID and try again.")
		}
		return nil, status.Errorf(codes.Internal, "Failed to update profile: %v", err)
	}

	var avatar_url, birthday *string

	if profile.AvatarUrl.String == "" {
		avatar_url = nil
	} else {
		v := fmt.Sprintf("%s?v=%d", profile.AvatarUrl.String, profile.AvatarVersion)
		avatar_url = &v
	}

	if profile.Birthday.Valid {
		birthdayStr := profile.Birthday.Time.Format("2006-01-02")
		birthday = &birthdayStr
	} else {
		birthday = nil
	}

	profileRedis := models.ProfileRedis{
		UserUUID:  profile.Uuid,
		Name:      &profile.Name,
		Email:     profile.Email.String,
		Phone:     &profile.Phone.String,
		AvatarUrl: avatar_url,
		Birthday:  birthday,
	}

	profileBytes, err := json.Marshal(profileRedis)
	if err != nil {
		log.Printf("Failed to marshal profile data (in user service layer): %v", err)
	} else {
		if err := s.redis_memory.Set(ctx, fmt.Sprintf("user_profile:%d", req.UserId), profileBytes, 24*7*time.Hour).Err(); err != nil {
			log.Printf("Failed to set profile data in Redis (in user service layer): %v", err)
		} else {
			log.Println("Profile data cached in Redis for user ID: ", req.UserId)
		}
	}

	return &user_proto.UpdateProfileResponse{
		Message: "Profile updated successfully",
	}, nil
}

func (s *userService) DisableUserByUserID(ctx context.Context, req *user_proto.DisableUserRequest) (*user_proto.DisableUserResponse, error) {
	caller := utils.GetCaller(ctx)

	if caller != ctx.Value(interceptor.CtxCallerKey).(string) {
		return nil, status.Errorf(codes.PermissionDenied, "Unauthorized: Caller in context does not match expected caller")
	}

	if req.UserId != ctx.Value(interceptor.CtxUserIDKey).(int64) {
		return nil, status.Errorf(codes.PermissionDenied, "Unauthorized: User ID in context does not match User ID in request")
	}

	if err := s.validator.Validate(req); err != nil {
		return nil, validation.BuildValidationError(err)
	}

	err := s.user_repo.DisableUserByUserID(ctx, req.UserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, validation.BuildBusinessError("USER_NOT_FOUND", "User not found with the given ID. Please check the user ID and try again.")
		}
		return nil, status.Errorf(codes.Internal, "Failed to disable user: %v", err)
	}

	reqLogoutAll := &auth_proto.LogoutAllRequest{
		UserId: req.UserId,
	}

	ctx = context.WithValue(ctx, interceptor.CtxUserIDKey, req.UserId)
	ctx = context.WithValue(ctx, interceptor.CtxCallerKey, utils.GetEnv("USER_SERVICE_NAME", ""))
	ctx = context.WithValue(ctx, interceptor.CtxAudKey, utils.GetEnv("AUTH_SERVICE_NAME", ""))

	_, err = s.auth_client.Client.LogoutAll(ctx, reqLogoutAll)
	if err != nil {
		return nil, validation.MapServiceError(err, "auth")
	}

	return &user_proto.DisableUserResponse{
		Success: true,
		Message: "User disabled successfully",
	}, nil
}

func (s *userService) VerifyIDTokenOTP(ctx context.Context, req *user_proto.VerifyIDTokenOTPRequest) (*user_proto.VerifyIDTokenOTPResponse, error) {
	caller := utils.GetCaller(ctx)

	if caller != ctx.Value(interceptor.CtxCallerKey).(string) {
		return nil, status.Errorf(codes.PermissionDenied, "Unauthorized: Caller in context does not match expected caller")
	}

	if err := s.validator.Validate(req); err != nil {
		return nil, validation.BuildValidationError(err)
	}

	token, err := s.auth_client_firebase.VerifyIDTokenAndCheckRevoked(ctx, req.GetIdToken())
	if err != nil {
		if auth_firebase.IsIDTokenRevoked(err) {
			return &user_proto.VerifyIDTokenOTPResponse{
				Valid:   false,
				Message: "ID token has been revoked. Please try again.",
			}, nil
		}
		return nil, status.Errorf(codes.Unauthenticated, "Failed to verify ID token: %v", err)
	}

	phoneClaim, ok := token.Claims["phone_number"].(string)
	if !ok || phoneClaim == "" {
		return &user_proto.VerifyIDTokenOTPResponse{
			Valid:   false,
			Message: "ID token does not contain a valid phone number claim.",
		}, nil
	}

	err = s.auth_client_firebase.RevokeRefreshTokens(ctx, token.UID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to revoke refresh tokens: %v", err)
	}

	return &user_proto.VerifyIDTokenOTPResponse{
		Valid:         true,
		VerifiedPhone: phoneClaim,
	}, nil
}

func (s *userService) ReportUserImage(ctx context.Context, req *user_proto.ReportUserImageRequest) (*user_proto.ReportUserImageResponse, error) {
	caller := utils.GetCaller(ctx)

	if caller != ctx.Value(interceptor.CtxCallerKey).(string) {
		return nil, status.Errorf(codes.PermissionDenied, "Unauthorized: Caller in context does not match expected caller")
	}

	// 1. Use Cloud Vision API to check the image
	image := vision.NewImageFromURI("gs://" + req.Bucket + "/" + req.Name)
	props, err := s.vision_client.DetectSafeSearch(ctx, image, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to analyze image with Vision API: %v", err)
	}

	// 2. Check if the image is likely to contain adult content, violence, medical content, or racy content
	// level: VERY_UNLIKELY, UNLIKELY, POSSIBLE, LIKELY, VERY_LIKELY <=> 1, 2, 3, 4, 5
	if props.Adult >= 2 || props.Violence >= 2 || props.Medical >= 2 || props.Racy >= 2 { // 2 is 'UNLIKELY' (have a slight possibility)
		log.Printf("Detecting image violates community standards... Adult: %s, Violence: %s, Medical: %s, Racy: %s \nDeleting file %s from bucket %s", props.Adult, props.Violence, props.Medical, props.Racy, req.Name, req.Bucket)

		// Deleting cache for image from CDN immediately if it violates community standards
		err := s.InvalidateCacheOfProcessedBucket(ctx, req.Name)
		if err != nil {
			log.Fatalf("Failed to invalidate cache for inappropriate image: %v", err)
		}

		// Delete image immediately if it violates community standards
		err = s.gcs_client.Bucket(req.Bucket).Object(req.Name).Delete(ctx)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Failed to delete inappropriate image from GCS: %v", err)
		}

		return &user_proto.ReportUserImageResponse{
			Message: "Image reported and removed due to inappropriate content",
		}, nil
	}

	return &user_proto.ReportUserImageResponse{
		Message: "Image reported successfully and does not violate community standards",
	}, nil
}

func (s *userService) InvalidateCacheOfProcessedBucket(ctx context.Context, filePath string) error {
	op, err := s.compute_service.UrlMaps.InvalidateCache(utils.GetEnv("PROJECT_ID", ""), utils.GetEnv("PROJECT_BALANCER_NAME", ""), &compute.CacheInvalidationRule{
		Path: utils.GetEnv("PROJECT_PATH_BUCKET_PROCESSED", "") + filePath,
	}).Do()

	if err != nil {
		return err
	}

	log.Printf("Invalidating cache for file %s from CDN... Operation: %s", filePath, op.Name)

	return nil
}
