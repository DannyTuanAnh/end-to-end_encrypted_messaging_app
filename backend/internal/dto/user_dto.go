package dto

type GetProfileByUserUUID struct {
	UUID string `uri:"uuid" binding:"required,uuid"`
}

type UpdateUserRequest struct {
	Name     *string `form:"name" binding:"omitempty,min=1,max=100,not_blank"`
	Birthday *string `form:"birthday" binding:"omitempty,min=10,datetime"`
	Phone    *string `form:"phone" binding:"omitempty,is_phone_mobile"`
}

type VerifyIDTokenOTP struct {
	IDToken string `form:"id_token" binding:"required,not_blank"`
}

type ReportUserImageRequest struct {
	UUID string `json:"uuid" binding:"required,uuid"`
	Name string `json:"name" binding:"required,min=1,not_blank,uuid"`
}
