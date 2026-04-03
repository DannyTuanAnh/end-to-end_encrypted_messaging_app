package server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/client"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/config"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/db/sqlc"
	auth_proto "github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/gen/auth"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/utils"

	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/repository"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/service"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type AuthServer struct {
	auth_proto.UnimplementedAuthServiceServer
	ctx    context.Context
	cfg    *config.Config
	server *grpc.Server
}

func NewAuthServer(ctx context.Context, db sqlc.Querier, rdb *redis.Client) (*AuthServer, error) {
	authCertFile := utils.GetEnv("PATH_CERT_AUTH_SERVICE", "")
	authKeyFile := utils.GetEnv("PATH_KEY_AUTH_SERVICE", "")

	cert, err := tls.LoadX509KeyPair(
		authCertFile,
		authKeyFile,
	)

	if err != nil {
		return nil, fmt.Errorf("Failed to load auth service TLS credentials: %v", err)
	}

	caCert, err := os.ReadFile(utils.GetEnv("PATH_CERT_CA", ""))
	if err != nil {
		return nil, fmt.Errorf("Failed to read CA certificate: %v", err)
	}

	caPool := x509.NewCertPool()
	caPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},

		ClientAuth: tls.RequireAndVerifyClientCert,

		ClientCAs: caPool,
	}

	authCfg := config.NewConfigAuthService()
	userCfg := config.NewConfigUserService()

	cfg := &config.Config{}

	cfg.Service.AuthServiceAddr = authCfg.Service.AuthServiceAddr
	cfg.Service.UserServiceAddr = userCfg.Service.UserServiceAddr

	cfg.Service.AuthServiceListenAddr = authCfg.Service.AuthServiceListenAddr

	user_client, err := client.NewUserClient(cfg.Service.UserServiceAddr, authCertFile, authKeyFile)
	if err != nil {
		return nil, fmt.Errorf("Failed to create user client: %v", err)
	}

	auth_repo := repository.NewAuthRepository(db)
	auth_service := service.NewAuthService(auth_repo, user_client, rdb)

	s := grpc.NewServer(grpc.Creds(credentials.NewTLS(tlsConfig)))

	auth_proto.RegisterAuthServiceServer(s, auth_service)

	return &AuthServer{
		ctx:    ctx,
		cfg:    cfg,
		server: s,
	}, nil
}

func (as *AuthServer) Run() (string, error) {
	listener, err := net.Listen("tcp", as.cfg.Service.AuthServiceListenAddr)
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
