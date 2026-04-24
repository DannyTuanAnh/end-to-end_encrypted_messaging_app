package server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"time"

	"cloud.google.com/go/storage"
	vision "cloud.google.com/go/vision/apiv1"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/client"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/config"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/db/sqlc"
	user_proto "github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/gen/user"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/interceptor"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/repository"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/service"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/utils"
	"github.com/redis/go-redis/v9"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var userPolicies = map[string][]string{
	"/proto.UserService/DisableUserByUserID": {
		"api-gateway",
	},
	"/proto.UserService/CreateProfile": {
		"auth-service",
	},
	"/proto.UserService/GetProfileByUserID": {
		"api-gateway",
	},
	"/proto.UserService/GetProfileByUserUUID": {
		"api-gateway",
	},
	"/proto.UserService/UpdateProfile": {
		"api-gateway",
	},
	"/proto.UserService/VerifyIDTokenOTP": {
		"api-gateway",
	},
	"/proto.UserService/ReportUserImage": {
		"api-gateway",
	},
}

type UserServer struct {
	user_proto.UnimplementedUserServiceServer
	ctx    context.Context
	cfg    *config.Config
	server *grpc.Server
}

func NewUserServer(ctx context.Context, db sqlc.Querier, rdb *redis.Client) (*UserServer, error) {
	userCertFile := utils.GetEnv("PATH_CERT_USER_SERVICE", "")
	userKeyFile := utils.GetEnv("PATH_KEY_USER_SERVICE", "")

	userCertPEM := []byte(userCertFile)
	userKeyPEM := []byte(userKeyFile)

	cert, err := tls.X509KeyPair(userCertPEM, userKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("Failed to load user service TLS credentials: %v", err)
	}

	caCert := []byte(utils.GetEnv("PATH_CERT_CA", ""))

	caPool := x509.NewCertPool()
	caPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},

		// ClientAuth: tls.RequireAndVerifyClientCert,

		// ClientCAs: caPool,

		ClientAuth: tls.VerifyClientCertIfGiven,
		ClientCAs:  caPool,
	}

	userCfg := config.NewConfigUserService()
	authCfg := config.NewConfigAuthService()

	cfg := &config.Config{}

	cfg.Service.AuthServiceAddr = authCfg.Service.AuthServiceAddr
	cfg.Service.UserServiceAddr = userCfg.Service.UserServiceAddr

	cfg.Service.UserServiceListenAddr = userCfg.Service.UserServiceListenAddr

	auth_client, err := client.NewAuthClient(cfg.Service.AuthServiceAddr, userCertFile, userKeyFile)
	if err != nil {
		return nil, fmt.Errorf("Failed to create auth client: %v", err)
	}

	vision_client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed to create Vision API client: %v", err)
	}

	compute_service, err := compute.NewService(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create compute service: %v", err)
	}

	user_repo := repository.NewUserRepository(db)
	user_service := service.NewUserService(user_repo, rdb, auth_client, connectAuthFirebase(ctx), vision_client, compute_service, connectGCS(ctx))

	s := grpc.NewServer(
		grpc.Creds(credentials.NewTLS(tlsConfig)),
		grpc.ChainUnaryInterceptor(
			interceptor.RBACInterceptor(userPolicies),
			interceptor.AuthServerInterceptor(userCertFile),
		),
	)

	user_proto.RegisterUserServiceServer(s, user_service)

	return &UserServer{
		ctx:    ctx,
		cfg:    userCfg,
		server: s,
	}, nil
}

func (as *UserServer) Run() (string, error) {
	listener, err := net.Listen("tcp", as.cfg.Service.UserServiceListenAddr)
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

func connectGCS(ctx context.Context) *storage.Client {
	serviceAccountKey := utils.GetEnv("GOOGLE_CLOUD_STORAGE_CREDENTIALS", "")

	var storage_client *storage.Client
	var err error
	if serviceAccountKey != "" {
		storage_client, err = storage.NewClient(ctx, option.WithAuthCredentialsFile(option.ServiceAccount, serviceAccountKey))
	} else {
		storage_client, err = storage.NewClient(ctx)
	}

	if err != nil {
		panic("Failed to initialize Google Cloud Storage client: " + err.Error())
	}

	return storage_client
}

func connectAuthFirebase(ctx context.Context) *auth.Client {
	// Initialize Firebase app for verify otp
	serviceAccountKey := utils.GetEnv("GOOGLE_APPLICATION_FIREBASE_CREDENTIALS", "")
	opt := option.WithAuthCredentialsJSON(option.ServiceAccount, []byte(serviceAccountKey))

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		panic("Failed to initialize Firebase app: " + err.Error())
	}

	auth_client_firebase, err := app.Auth(ctx)
	if err != nil {
		panic("Failed to initialize Firebase Auth client: " + err.Error())
	}

	return auth_client_firebase
}
