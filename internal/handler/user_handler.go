package handler

import (
	"net/http"

	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/dto"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/service"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/utils"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/validation"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	user_service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{
		user_service: service,
	}
}

type GetUserByUUID struct {
	UUID string `uri:"uuid" binding:"required,uuid"`
}

type GetUsersParams struct {
	Search string `form:"search" binding:"omitempty,min=3,max=10,search"`
	Page   int    `form:"page" binding:"omitempty,gte=1,lte=100"`
	Limit  int    `form:"limit" binding:"omitempty,gte=1,lte=100"`
}

// GetAllUsers is a handler function that retrieves all users from the user service and responds with a JSON array of user data.
// It calls the GetAllUsers method of the user service to get a map of users, converts the map to a slice, and then maps each user to a UserDTO before responding with the data.
// If there are any errors during the retrieval process, it responds with an appropriate error message.
func (uh *UserHandler) GetAllUsers(ctx *gin.Context) {
	var params GetUsersParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	if params.Page == 0 {
		params.Page = 1
	}

	if params.Limit == 0 {
		params.Limit = 10
	}

	usersMap, err := uh.user_service.GetAllUsers(params.Search, params.Page, params.Limit)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	userDTOs := utils.ConvertMapToSliceWithTransform(usersMap, dto.MapUserToDTO)
	utils.ResponseSuccess(ctx, http.StatusOK, userDTOs)
}

// GetUserByUUID is a handler function that retrieves a user by their UUID.
// It binds the UUID from the URI, validates it, and then calls the GetUserByUUID method of the user service to retrieve the user.
// If there are any errors during binding or retrieval, it responds with an appropriate error message.
// If the user is found successfully, it responds with a success message and the user data.
func (uh *UserHandler) GetUserByUUID(ctx *gin.Context) {
	var param GetUserByUUID
	if err := ctx.ShouldBindUri(&param); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	user, err := uh.user_service.GetUserByUUID(param.UUID)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	userDTO := dto.MapUserToDTO(user)
	utils.ResponseSuccess(ctx, http.StatusOK, &userDTO)
}

// CreateUser is a handler function that handles the creation of a new user.
// It binds the incoming JSON request to a User struct, validates the input,
// and then calls the CreateUser method of the user service to create the user.
// If there are any errors during binding or creation, it responds with an appropriate error message.
// If the user is created successfully, it responds with a success message and the created user data.
func (uh *UserHandler) CreateUser(ctx *gin.Context) {
	var input dto.CreateUserInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	user := input.MapCreateInputToUserModel()

	createdUser, err := uh.user_service.CreateUser(user)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	userDTO := dto.MapUserToDTO(createdUser)
	utils.ResponseSuccess(ctx, http.StatusCreated, &userDTO)
}

// UpdateUser is a handler function that handles the updating of an existing user.
// It binds the UUID from the URI and the incoming JSON request to a User struct, validates the input,
// and then calls the UpdateUser method of the user service to update the user.
// If there are any errors during binding or updating, it responds with an appropriate error message.
// If the user is updated successfully, it responds with a success message and the updated user data.
func (uh *UserHandler) UpdateUser(ctx *gin.Context) {
	var param GetUserByUUID
	if err := ctx.ShouldBindUri(&param); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	var input dto.UpdateUserInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	user := input.MapUpdateInputToUserModel()
	updatedUser, err := uh.user_service.UpdateUser(param.UUID, user)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	userDTO := dto.MapUserToDTO(updatedUser)
	utils.ResponseSuccess(ctx, http.StatusOK, &userDTO)
}
func (uh *UserHandler) DeleteUser(ctx *gin.Context) {
	var param GetUserByUUID
	if err := ctx.ShouldBindUri(&param); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	if err := uh.user_service.DeleteUser(param.UUID); err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseStatusCode(ctx, http.StatusNoContent)
}
