package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/config"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/db/sqlc"
	user_proto "github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/gen/user"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/repository"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/service"
	"google.golang.org/grpc"
)

type UserServer struct {
	user_proto.UnimplementedUserServiceServer
	ctx    context.Context
	cfg    *config.Config
	server *grpc.Server
}

func NewUserServer(ctx context.Context, db sqlc.Querier) *UserServer {
	cfg := config.NewConfigUserService()

	user_repo := repository.NewUserRepository(db)
	user_service := service.NewUserService(user_repo)

	s := grpc.NewServer()

	user_proto.RegisterUserServiceServer(s, user_service)

	return &UserServer{
		ctx:    ctx,
		cfg:    cfg,
		server: s,
	}
}

func (as *UserServer) Run() (string, error) {
	listener, err := net.Listen("tcp", as.cfg.Service.UserServiceAddr)
	if err != nil {
		return "", fmt.Errorf("Failed to listen: %v", err)
	}

	errChan := make(chan error, 1)

	go func() {
		log.Printf("User server is listening on %s", listener.Addr())
		if err := as.server.Serve(listener); err != nil {
			errChan <- fmt.Errorf("Failed to serve: %v", err)
		}
	}()

	select {
	case err := <-errChan:
		return "", fmt.Errorf("User server error: %v", err)
	case <-as.ctx.Done():
		log.Println("User server is shutting down...")
		done := make(chan struct{})

		go func() {
			as.server.GracefulStop()
			close(done)
		}()

		select {
		case <-done:
			return "User server stopped gracefully", nil
		case <-time.After(5 * time.Second):
			log.Println("User server shutdown timed out, forcing stop")
			as.server.Stop()
		}
		return "User server stopped gracefully", nil
	}
}
