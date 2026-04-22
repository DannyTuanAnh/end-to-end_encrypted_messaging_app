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
	notify_proto "github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/gen/notify"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/interceptor"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/repository"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/service"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/utils"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var notifyPolicies = map[string][]string{}

type NotifyServer struct {
	notify_proto.UnimplementedNotifyServiceServer
	ctx    context.Context
	cfg    *config.Config
	server *grpc.Server
}

func NewNotifyServer(ctx context.Context, db sqlc.Querier, rdb *redis.Client) (*NotifyServer, error) {
	notifyCertFile := utils.GetEnv("PATH_CERT_NOTIFY_SERVICE", "")
	notifyKeyFile := utils.GetEnv("PATH_KEY_NOTIFY_SERVICE", "")

	cert, err := tls.LoadX509KeyPair(
		notifyCertFile,
		notifyKeyFile,
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to load notify service TLS credentials: %v", err)
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

	notifyCfg := config.NewConfigNotifyService()
	userCfg := config.NewConfigUserService()

	cfg := &config.Config{}

	cfg.Service.UserServiceAddr = userCfg.Service.UserServiceAddr
	cfg.Service.NotifyServiceAddr = notifyCfg.Service.NotifyServiceAddr

	cfg.Service.NotifyServiceListenAddr = notifyCfg.Service.NotifyServiceListenAddr

	user_client, err := client.NewUserClient(cfg.Service.UserServiceAddr, notifyCertFile, notifyKeyFile)
	if err != nil {
		return nil, fmt.Errorf("Failed to create auth client: %v", err)
	}

	notify_repo := repository.NewNotifyRepository(db)
	notify_service := service.NewNotifyService(notify_repo, rdb, user_client)

	s := grpc.NewServer(
		grpc.Creds(credentials.NewTLS(tlsConfig)),
		grpc.ChainUnaryInterceptor(
			interceptor.RBACInterceptor(notifyPolicies),
			interceptor.AuthServerInterceptor(notifyCertFile),
		),
	)

	notify_proto.RegisterNotifyServiceServer(s, notify_service)

	return &NotifyServer{
		ctx:    ctx,
		cfg:    notifyCfg,
		server: s,
	}, nil
}

func (as *NotifyServer) Run() (string, error) {
	listener, err := net.Listen("tcp", as.cfg.Service.NotifyServiceListenAddr)
	if err != nil {
		return "", fmt.Errorf("Failed to listen: %v", err)
	}

	errChan := make(chan error, 1)

	go func() {
		log.Printf("Notify server is listening on %s", listener.Addr())
		if err := as.server.Serve(listener); err != nil {
			errChan <- fmt.Errorf("Failed to serve: %v", err)
		}
	}()

	select {
	case err := <-errChan:
		return "", fmt.Errorf("Notify server error: %v", err)
	case <-as.ctx.Done():
		log.Println("Notify server is shutting down...")
		done := make(chan struct{})

		go func() {
			as.server.GracefulStop()
			close(done)
		}()

		select {
		case <-done:
			return "Notify server stopped gracefully", nil
		case <-time.After(5 * time.Second):
			log.Println("Notify server shutdown timed out, forcing stop")
			as.server.Stop()
		}
		return "Notify server stopped gracefully", nil
	}
}
