package dto

import "github.com/user-manage/internal/models"

type UserDTO struct {
	UUID   string `json:"uuid"`
	Name   string `json:"full_name"`
	Email  string `json:"email_address"`
	Age    int    `json:"age"`
	Status string `json:"status"`
	Level  string `json:"level"`
}

type CreateUserInput struct {
	Name     string `json:"name" binding:"required"`
	Age      int    `json:"age" binding:"required,gt=18"`
	Email    string `json:"email" binding:"required,email,email_advanced"`
	Password string `json:"password" binding:"required,min=8,strong_password"`
	Status   int    `json:"status" binding:"required,oneof=1 2"`
	Level    int    `json:"level" binding:"required,oneof=1 2"`
}

type UpdateUserInput struct {
	Name     string `json:"name" binding:"omitempty"`
	Age      int    `json:"age" binding:"omitempty,gt=18"`
	Email    string `json:"email" binding:"omitempty,email,email_advanced"`
	Password string `json:"password" binding:"omitempty,min=8,strong_password"`
	Status   int    `json:"status" binding:"omitempty,oneof=1 2"`
	Level    int    `json:"level" binding:"omitempty,oneof=1 2"`
}

func (input *CreateUserInput) MapCreateInputToUserModel() models.User {
	return models.User{
		Name:     input.Name,
		Age:      input.Age,
		Email:    input.Email,
		Password: input.Password,
		Status:   input.Status,
		Level:    input.Level,
	}
}

func (input *UpdateUserInput) MapUpdateInputToUserModel() models.User {
	return models.User{
		Name:     input.Name,
		Age:      input.Age,
		Email:    input.Email,
		Password: input.Password,
		Status:   input.Status,
		Level:    input.Level,
	}
}

func MapUserToDTO(user models.User) *UserDTO {
	return &UserDTO{
		UUID:   user.UUID,
		Name:   user.Name,
		Email:  user.Email,
		Age:    user.Age,
		Status: mapStatusText(user.Status),
		Level:  mapLevelText(user.Level),
	}
}

func mapStatusText(status int) string {
	switch status {
	case 1:
		return "Show"
	case 2:
		return "Hide"
	default:
		return "None"
	}
}

func mapLevelText(level int) string {
	switch level {
	case 1:
		return "Admin"
	case 2:
		return "Member"
	default:
		return "None"
	}
}
