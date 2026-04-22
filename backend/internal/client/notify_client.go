package client

import (
	notify_proto "github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/gen/notify"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/utils"
)

type NotifyClient struct {
	Client notify_proto.NotifyServiceClient
}

func NewNotifyClient(addr, certFile, keyFile string) (*NotifyClient, error) {
	keyClient := utils.GetEnv("PATH_KEY_NOTIFY_SERVICE", "")
	conn, err := NewGRPCConn(addr, certFile, keyFile, keyClient)
	if err != nil {
		return nil, err
	}

	return &NotifyClient{
		Client: notify_proto.NewNotifyServiceClient(conn),
	}, nil
}
