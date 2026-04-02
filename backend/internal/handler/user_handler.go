package handler

import (
	"errors"
	"net/http"

	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/client"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/utils"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/validation"
	"github.com/gin-gonic/gin"

	user_proto "github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/gen/user"
)

type UserHandler struct {
	user_client *client.UserClient
}

func NewUserHandler(user_client *client.UserClient) *UserHandler {
	return &UserHandler{
		user_client: user_client,
	}
}

func (h *UserHandler) UpdateProfile(ctx *gin.Context) {}

func (h *UserHandler) DisableUser(ctx *gin.Context) {
	userId, exist := ctx.Get("user_id")
	if !exist {
		utils.ResponseErrorAbort(ctx, utils.NewError("User ID not found in context", utils.ErrCodeNotFound))
	}

	userID, ok := userId.(int64)
	if !ok {
		utils.ResponseErrorAbort(ctx, utils.NewError("User ID in context has invalid type", utils.ErrCodeInternal))
	}

	if userID <= 0 {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(errors.New("UserID must greater than 0")))
	}

	_, err := h.user_client.Client.DisableUserByUserID(ctx, &user_proto.DisableUserRequest{UserId: userID})
	if err != nil {
		utils.WriteGRPCErrorToGin(ctx, err)
		return
	}

	utils.ResponseStatusCode(ctx, http.StatusNoContent)

}
