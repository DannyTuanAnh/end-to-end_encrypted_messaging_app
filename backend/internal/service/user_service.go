package service

import (
	"context"
	"fmt"
	"time"

	"buf.build/go/protovalidate"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/db/sqlc"
	user_proto "github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/gen/user"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/repository"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/utils"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/validation"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type userService struct {
	user_proto.UnimplementedUserServiceServer
	user_repo repository.UserRepository
	validator protovalidate.Validator
}

func NewUserService(user_repo repository.UserRepository) *userService {
	v, err := protovalidate.New()
	if err != nil {
		panic(fmt.Sprintf("Failed to create validator: %v", err))
	}
	return &userService{
		user_repo: user_repo,
		validator: v,
	}
}

func (s *userService) CreateProfile(ctx context.Context, req *user_proto.CreateProfileRequest) (*user_proto.CreateProfileResponse, error) {
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
