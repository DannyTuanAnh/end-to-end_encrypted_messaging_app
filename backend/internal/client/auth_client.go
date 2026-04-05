package client

import (
	"crypto/tls"
	"crypto/x509"
	"os"

	auth_proto "github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/gen/auth"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/interceptor"
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

	caCert, err := os.ReadFile(utils.GetEnv("PATH_CERT_CA", ""))
	if err != nil {
		return nil, err
	}

	caPool := x509.NewCertPool()
	caPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caPool,
	}

	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
		grpc.WithUnaryInterceptor(interceptor.AuthClientInterceptor(utils.GetEnv("PATH_KEY_AUTH_SERVICE", ""))),
	)

	if err != nil {
		return nil, err
	}

	return &AuthClient{
		Client: auth_proto.NewAuthServiceClient(conn),
	}, nil
}
