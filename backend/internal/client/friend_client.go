package client

import (
	friend_proto "github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/gen/friend"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/utils"
)

type FriendClient struct {
	Client friend_proto.FriendServiceClient
}

func NewFriendClient(addr string, certFile string, keyFile string) (*FriendClient, error) {
	conn, err := NewGRPCConn(addr, utils.GetEnv("FRIEND_SERVER_NAME", ""), certFile, keyFile)
	if err != nil {
		return nil, err
	}

	return &FriendClient{
		Client: friend_proto.NewFriendServiceClient(conn),
	}, nil
}
