package service

import (
	"context"
	"fmt"

	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/db/sqlc"
	user_proto "github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/grpc/user"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/repository"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/utils"
)

type userService struct {
	user_proto.UnimplementedUserServiceServer
	user_repo repository.UserRepository
}

func NewUserService(user_repo repository.UserRepository) *userService {
	return &userService{user_repo: user_repo}
}

func (s *userService) CreateProfile(ctx context.Context, req *user_proto.CreateProfileRequest) (*user_proto.CreateProfileResponse, error) {
	if req.UserId == 0 {
		return nil, fmt.Errorf("User ID is required")
	}

	_, err := s.user_repo.CreateProfile(ctx, sqlc.CreateProfileParams{
		UserID:    req.UserId,
		Name:      req.Name,
		Email:     utils.ConvertToPgTypeText(req.Email),
		AvatarUrl: utils.ConvertToPgTypeText(req.AvatarUrl),
		Birthday:  utils.ConvertToPgTypeDate(req.Birthday),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create profile: %w", err)
	}

	return &user_proto.CreateProfileResponse{
		Success: true,
		Message: "Profile created successfully",
	}, nil
}
