package dto

type GetProfileByUserUUID struct {
	UUID string `uri:"uuid" binding:"required,uuid"`
}
