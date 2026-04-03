package client

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"

	auth_proto "github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/gen/auth"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type AuthClient struct {
	Client auth_proto.AuthServiceClient
}

func NewAuthClient(addr string, certFile string, keyFile string) (*AuthClient, error) {
	cert, err := tls.LoadX509KeyPair(
		certFile,
		keyFile,
	)
	if err != nil {
		return nil, err
	}

	log.Println("client cert:", certFile)

	caCert, err := os.ReadFile(utils.GetEnv("PATH_CERT_CA", ""))
	if err != nil {
		return nil, err
	}

	caPool := x509.NewCertPool()
	caPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caPool,
		ServerName:   "auth-service",
	}

	log.Printf("NewAuthClient connecting to %s with cert=%s key=%s", addr, certFile, keyFile)

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))

	if err != nil {
		return nil, err
	}

	return &AuthClient{
		Client: auth_proto.NewAuthServiceClient(conn),
	}, nil
}
