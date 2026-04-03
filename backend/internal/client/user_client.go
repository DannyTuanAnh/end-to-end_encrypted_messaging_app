package client

import (
	"crypto/tls"
	"crypto/x509"
	"os"

	user_proto "github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/gen/user"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type UserClient struct {
	Client user_proto.UserServiceClient
}

func NewUserClient(addr string, certFile string, keyFile string) (*UserClient, error) {
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

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))

	if err != nil {
		return nil, err
	}

	return &UserClient{
		Client: user_proto.NewUserServiceClient(conn),
	}, nil
}
