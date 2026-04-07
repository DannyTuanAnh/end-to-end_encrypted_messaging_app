package models

import "github.com/google/uuid"

type ProfileRedis struct {
	UserUUID  uuid.UUID `json:"user_uuid"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	AvatarUrl string    `json:"avatar_url"`
	Birthday  string    `json:"birthday"`
}
