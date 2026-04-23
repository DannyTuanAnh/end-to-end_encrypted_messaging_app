package client

import (
	user_proto "github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/gen/user"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/utils"
)

type UserClient struct {
	Client user_proto.UserServiceClient
}

func NewUserClient(addr string, certFile string, keyFile string) (*UserClient, error) {
	keyClient := utils.GetEnv("PATH_KEY_USER_SERVICE", "")
	conn, err := NewGRPCConn(addr, certFile, keyFile, keyClient)
	if err != nil {
		return nil, err
	}

	return &UserClient{
		Client: user_proto.NewUserServiceClient(conn),
	}, nil
}
