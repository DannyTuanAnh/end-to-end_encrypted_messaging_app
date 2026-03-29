package dto

type RequestLoginGoogle struct {
	AuthCode string `json:"auth_code" binding:"required"`
}
