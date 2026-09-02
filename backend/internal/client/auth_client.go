package client

import (
	auth_proto "github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/gen/auth"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/utils"
)

type AuthClient struct {
	Client auth_proto.AuthServiceClient
}

func NewAuthClient(addr string, certFile string, keyFile string) (*AuthClient, error) {
	conn, err := NewGRPCConn(addr, utils.GetEnv("AUTH_SERVER_NAME", ""), certFile, keyFile)
	if err != nil {
		return nil, err
	}

	return &AuthClient{
		Client: auth_proto.NewAuthServiceClient(conn),
	}, nil
}
