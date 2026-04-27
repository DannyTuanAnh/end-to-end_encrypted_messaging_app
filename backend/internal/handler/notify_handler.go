package handler

import (
	"errors"
	"log"
	"strconv"

	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/client"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/sse"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/utils"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/validation"
	"github.com/gin-gonic/gin"
)

type NotifyHandler struct {
	user_client *client.UserClient
}

func NewNotifyHandler(user_client *client.UserClient) *NotifyHandler {
	return &NotifyHandler{
		user_client: user_client,
	}
}

func (n *NotifyHandler) HandleSSE(ctx *gin.Context) {
	ctx.Writer.Header().Set("Content-Type", "text/event-stream")
	ctx.Writer.Header().Set("Cache-Control", "no-cache")
	ctx.Writer.Header().Set("Connection", "keep-alive")

	userId, exist := ctx.Get("user_id")
	log.Println("User ID from context:", userId, "Exist:", exist)
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

	userIDStr := strconv.FormatInt(userID, 10)

	messageChan := sse.MainBroker.AddClient(userIDStr)
	defer sse.MainBroker.RemoveClient(userIDStr)

	for {
		select {
		case message := <-messageChan:
			// Send the message to the client
			log.Println("Sending message to client:", message)
			_, err := ctx.Writer.Write([]byte("data: " + message + "\n\n"))
			if err != nil {
				return
			}
			ctx.Writer.Flush()
		case <-ctx.Request.Context().Done():
			return
		}
	}
}
