package service

import (
	"context"
	"fmt"

	"buf.build/go/protovalidate"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/client"
	sqlc "github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/db/sqlc/friend"
	friend_proto "github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/gen/friend"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/repository"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/validation"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type friendService struct {
	friend_proto.UnimplementedFriendServiceServer
	friend_repo repository.FriendRepository
	user_client *client.UserClient
	validator   protovalidate.Validator
}

func NewFriendService(friend_repo repository.FriendRepository, user_client *client.UserClient) *friendService {
	v, err := protovalidate.New()
	if err != nil {
		panic(fmt.Sprintf("Failed to create validator: %v", err))
	}

	return &friendService{
		friend_repo: friend_repo,
		user_client: user_client,
		validator:   v,
	}
}

func (fs *friendService) GetRelationship(ctx context.Context, req *friend_proto.GetRelationshipRequest) (*friend_proto.GetRelationshipResponse, error) {
	if err := fs.validator.Validate(req); err != nil {
		return nil, validation.BuildValidationError(err)
	}

	info, err := fs.friend_repo.GetInfoRelationship(ctx, sqlc.GetInfoRelationshipParams{
		CurrentUserID: req.CurrentUserId,
		TargetUserID:  req.TargetUserId,
	})
	if err != nil {
		if err == repository.ErrSameUser {
			return nil, status.Errorf(codes.InvalidArgument, "Cannot perform action on the same user")
		}

		return nil, status.Errorf(codes.Internal, "Failed to get relationship info: %v", err)
	}

	return &friend_proto.GetRelationshipResponse{
		FriendRequestDirection: info.FriendRequestDirection,
		IsFriend:               info.IsFriend,
	}, nil
}
