package service

import (
	"context"
	"log"

	proto "github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/grpc/auth"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/repository"
)

type authService struct {
	proto.UnimplementedAuthServiceServer
	auth_repo repository.AuthRepository
}

func NewAuthService(auth_repo repository.AuthRepository) *authService {
	return &authService{auth_repo: auth_repo}
}

func (s *authService) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {

	log.Println("LoginGoogle called in AuthService with authCode:", req.AuthorCode)

	return &proto.LoginResponse{}, nil
}
