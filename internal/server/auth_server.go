package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/config"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/db/sqlc"
	proto "github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/grpc/auth"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/repository"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/service"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

type AuthServer struct {
	proto.UnimplementedAuthServiceServer
	ctx    context.Context
	cfg    *config.Config
	server *grpc.Server
}

func NewAuthServer(ctx context.Context, cfg *config.Config, db sqlc.Querier, rdb *redis.Client) *AuthServer {
	auth_repo := repository.NewAuthRepository(db)
	auth_service := service.NewAuthService(auth_repo)

	s := grpc.NewServer()

	proto.RegisterAuthServiceServer(s, auth_service)

	return &AuthServer{
		ctx:    ctx,
		cfg:    cfg,
		server: s,
	}
}

func (as *AuthServer) Run() (string, error) {
	listener, err := net.Listen("tcp", as.cfg.Service.AuthServiceAddr)
	if err != nil {
		return "", fmt.Errorf("Failed to listen: %v", err)
	}

	errChan := make(chan error, 1)

	go func() {
		log.Printf("Auth server is listening on %s", listener.Addr())
		if err := as.server.Serve(listener); err != nil {
			errChan <- fmt.Errorf("Failed to serve: %v", err)
		}
	}()

	select {
	case err := <-errChan:
		return "", fmt.Errorf("Auth server error: %v", err)
	case <-as.ctx.Done():
		log.Println("Auth server is shutting down...")
		done := make(chan struct{})

		go func() {
			as.server.GracefulStop()
			close(done)
		}()

		select {
		case <-done:
			return "Auth server stopped gracefully", nil
		case <-time.After(5 * time.Second):
			log.Println("Auth server shutdown timed out, forcing stop")
			as.server.Stop()
		}
		return "Auth server stopped gracefully", nil
	}
}
