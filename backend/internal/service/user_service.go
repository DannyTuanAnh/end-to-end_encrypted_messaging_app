package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"buf.build/go/protovalidate"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/client"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/db/sqlc"
	auth_proto "github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/gen/auth"
	user_proto "github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/gen/user"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/repository"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/utils"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/validation"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type userService struct {
	user_proto.UnimplementedUserServiceServer
	user_repo    repository.UserRepository
	redis_memory *redis.Client
	validator    protovalidate.Validator

	auth_client *client.AuthClient
}

func NewUserService(user_repo repository.UserRepository, rdb *redis.Client, auth_client *client.AuthClient) *userService {
	v, err := protovalidate.New()
	if err != nil {
		panic(fmt.Sprintf("Failed to create validator: %v", err))
	}
	return &userService{
		user_repo:    user_repo,
		validator:    v,
		redis_memory: rdb,
		auth_client:  auth_client,
	}
}

func (s *userService) CreateProfile(ctx context.Context, req *user_proto.CreateProfileRequest) (*user_proto.CreateProfileResponse, error) {
	caller := utils.GetCaller(ctx)

	if caller != "auth-service" {
		return nil, status.Errorf(codes.PermissionDenied, "Unauthorized: Only auth service can call CreateProfile")
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

func (s *userService) DisableUserByUserID(ctx context.Context, req *user_proto.DisableUserRequest) (*user_proto.DisableUserResponse, error) {
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

	_, err = s.auth_client.Client.LogoutAll(ctx, reqLogoutAll)
	if err != nil {
		return nil, validation.MapUserServiceError(err, "auth")
	}

	return &user_proto.DisableUserResponse{
		Success: true,
		Message: "User disabled successfully",
	}, nil
}
