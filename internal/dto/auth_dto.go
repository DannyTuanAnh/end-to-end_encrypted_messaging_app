package dto

import auth_proto "github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/grpc/auth"

type RequestLoginGoogle struct {
	AuthCode string `json:"auth_code" binding:"required"`
}

type ResponseLoginGoogle struct {
	SessionId string `json:"session_id"`
	UserId    int64  `json:"user_id"`
	DeviceId  string `json:"device_id"`
}

func MapLoginGoogleResponseToDTO(response *auth_proto.LoginResponse) *ResponseLoginGoogle {
	return &ResponseLoginGoogle{
		SessionId: response.Session,
		UserId:    response.UserId,
		DeviceId:  response.DeviceId,
	}
}
