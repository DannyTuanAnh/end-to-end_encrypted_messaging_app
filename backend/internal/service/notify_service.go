package service

import (
	"fmt"

	"buf.build/go/protovalidate"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/client"
	notify_proto "github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/gen/notify"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/repository"
	"github.com/redis/go-redis/v9"
)

type notifyService struct {
	notify_proto.UnimplementedNotifyServiceServer
	validator protovalidate.Validator

	notify_repo repository.NotifyRepository

	redis_memory *redis.Client

	user_client *client.UserClient
}

func NewNotifyService(notify_repo repository.NotifyRepository, rdb *redis.Client, user_client *client.UserClient) *notifyService {
	v, err := protovalidate.New()
	if err != nil {
		panic(fmt.Sprintf("Failed to create validator: %v", err))
	}
	return &notifyService{
		notify_repo:  notify_repo,
		redis_memory: rdb,
		validator:    v,
	}
}
