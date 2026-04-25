package server

import (
	"context"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/client"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/config"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/db/sqlc"
	auth_proto "github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/gen/auth"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/interceptor"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/utils"

	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/repository"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/service"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var authPolicies = map[string][]string{
	"/proto.AuthService/LoginGoogle": {
		os.Getenv("API_GATEWAY_NAME"),
	},
	"/proto.AuthService/Logout": {
		os.Getenv("API_GATEWAY_NAME"),
	},
	"/proto.AuthService/LogoutAll": {
		os.Getenv("API_GATEWAY_NAME"),
		os.Getenv("USER_SERVICE_NAME"),
	},
}

type AuthServer struct {
	auth_proto.UnimplementedAuthServiceServer
	ctx    context.Context
	cfg    *config.Config
	server *grpc.Server
}

func NewAuthServer(ctx context.Context, db sqlc.Querier, rdb *redis.Client) (*AuthServer, error) {
	authCertFile := utils.GetEnv("PATH_CERT_AUTH_SERVICE", "")
	authKeyFile := utils.GetEnv("PATH_KEY_AUTH_SERVICE", "")

	authCertPEM := []byte(authCertFile)
	authKeyPEM := []byte(authKeyFile)

	cert, err := tls.X509KeyPair(authCertPEM, authKeyPEM)

	x509Cert, _ := x509.ParseCertificate(cert.Certificate[0])

	log.Printf("Auth service cert CN=%s DNS=%v Issuer=%s",
		x509Cert.Subject.CommonName,
		x509Cert.DNSNames,
		x509Cert.Issuer.CommonName,
	)

	log.Printf("Auth leaf Subject=%s", x509Cert.Subject)
	log.Printf("Auth leaf Issuer=%s", x509Cert.Issuer)
	log.Printf("Auth leaf Serial=%s", x509Cert.SerialNumber.String())
	log.Printf("Auth leaf DNSNames=%v", x509Cert.DNSNames)
	log.Printf("Auth leaf SHA256=%x", sha256.Sum256(x509Cert.Raw))

	if err != nil {
		return nil, fmt.Errorf("Failed to load auth service TLS credentials: %v", err)
	}

	caCert := []byte(utils.GetEnv("PATH_CERT_CA", ""))

	block, _ := pem.Decode(caCert)
	if block == nil {
		log.Fatal("failed to parse auth CA PEM")
	}

	caParsed, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Auth trusted CA Subject=%s", caParsed.Subject)
	log.Printf("Auth trusted CA Issuer=%s", caParsed.Issuer)
	log.Printf("Auth trusted CA Serial=%s", caParsed.SerialNumber.String())
	log.Printf("Auth trusted CA IsCA=%v", caParsed.IsCA)
	log.Printf("Auth trusted CA SHA256=%x", sha256.Sum256(caParsed.Raw))

	caPool := x509.NewCertPool()
	caPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},

		ClientAuth: tls.RequireAndVerifyClientCert,
		ClientCAs:  caPool,

		MinVersion: tls.VersionTLS12,
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

	s := grpc.NewServer(
		grpc.Creds(credentials.NewTLS(tlsConfig)),
		grpc.ChainUnaryInterceptor(
			interceptor.MTLSIdentityInterceptor(),
			interceptor.RBACInterceptor(authPolicies),
			interceptor.JWTAuthServerInterceptor(authCertFile),
		),
	)

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
