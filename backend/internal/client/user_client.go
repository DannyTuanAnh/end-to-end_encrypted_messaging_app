package client

import (
	user_proto "github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/gen/user"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/utils"
)

type UserClient struct {
	Client user_proto.UserServiceClient
}

func NewUserClient(addr string, certFile string, keyFile string) (*UserClient, error) {
	conn, err := NewGRPCConn(addr, utils.GetEnv("USER_SERVER_NAME", ""), certFile, keyFile)
	if err != nil {
		return nil, err
	}

	return &UserClient{
		Client: user_proto.NewUserServiceClient(conn),
	}, nil
}
