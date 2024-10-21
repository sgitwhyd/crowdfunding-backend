package handler

import (
	"be-bwastartup/auth"
	"be-bwastartup/helper"
	"be-bwastartup/user"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
	authService auth.Service
}

func NewUserHandler(userService user.Service, authService auth.Service) *userHandler {
	return &userHandler{userService, authService}
}

func (h *userHandler) RegisterUser(c *gin.Context) {
	// tangkap input dari user
	var input user.RegisterUserInput


	err := c.ShouldBindJSON(&input)
	if err != nil {
			// validation in here
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		errorReponse := helper.APIResponse("Register account failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusBadRequest, errorReponse)
		return
	}

	newUser, nil := h.userService.RegisterUser(input)
	if err != nil {
		errorResponse := helper.APIResponse("Register account failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	userToken, err := h.authService.GenerateToken(newUser.ID)
	if err != nil {
		errorsResponse := gin.H{"errors": err.Error()}
		errorResponse := helper.APIResponse("Register account failed", http.StatusBadRequest, "error", errorsResponse)
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}


	registerResponse := user.FormatUser(newUser, userToken)

	response := helper.APIResponse("Account has been registered", http.StatusCreated, "success", registerResponse)

	c.JSON(http.StatusOK, response)

}

func (h *userHandler) Login(c *gin.Context){
	var input user.LoginUserInput
	
	err := c.ShouldBind(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		errorReponse := helper.APIResponse("Login failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusBadRequest, errorReponse)
		return
	}

	loggedinUser, err := h.userService.Login(input)
	if err != nil {
		errorsResponse := gin.H{"errors": err.Error()}
		errorResponse := helper.APIResponse("Login failed", http.StatusBadRequest, "error", errorsResponse)
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	t, err := h.authService.GenerateToken(loggedinUser.ID)
	if err != nil {
		errorsResponse := gin.H{"errors": err.Error()}
		errorResponse := helper.APIResponse("Login failed", http.StatusBadRequest, "error", errorsResponse)
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	response := helper.APIResponse("Successfuly loggedin", http.StatusOK, "success", user.FormatUser(loggedinUser, t))
	c.JSON(http.StatusOK, response)
}

func (h *userHandler) CheckEmailAvailability(c *gin.Context){
	var input user.CheckEmailInput

	err := c.ShouldBind(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		errorReponse := helper.APIResponse("Email checking failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusBadRequest, errorReponse)
		return
	}

	isEmailAvailable, err := h.userService.IsEmailAvailable(input)
	if err != nil {
		errorsResponse := gin.H{
			"is_available": isEmailAvailable,
		}
		errorResponse := helper.APIResponse("Email checking failed", http.StatusBadRequest, "error", errorsResponse)
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	data := gin.H{
		"is_available": isEmailAvailable,
	}

	metaMessage := "Email has been registered"
	if isEmailAvailable {
		metaMessage = "Email is available"
	}

	response := helper.APIResponse(metaMessage, http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}

func (h *userHandler) UploadAvatar(c *gin.Context){
	file, err := c.FormFile("avatar")
	if err != nil {
		data := gin.H{"is_uploaded": false}
		errorResponse := helper.APIResponse("Failed to upload avatar image", http.StatusUnprocessableEntity, "error", data)
		c.JSON(http.StatusUnprocessableEntity, errorResponse)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)

	path := fmt.Sprintf("images/%d-%s", currentUser.ID, file.Filename)
	err = c.SaveUploadedFile(file, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		errorResponse := helper.APIResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	_, err = h.userService.UploadAvatar(currentUser.ID, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		errorResponse := helper.APIResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	data := gin.H{"is_uploaded": true}
	response := helper.APIResponse("Avatar successfuly uploaded", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}

func (h *userHandler) UpdateUser(c *gin.Context){
	var input user.FormUpdateUserInput

	err := c.ShouldBind(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		errorReponse := helper.APIResponse("Update failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusBadRequest, errorReponse)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)

	userID	:= currentUser.ID
	updatedUser, err := h.userService.UpdateUser(userID, input)
	if err != nil {
		errorsResponse := gin.H{"errors": err.Error()}
		errorResponse := helper.APIResponse("Update failed", http.StatusBadRequest, "error", errorsResponse)
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	response := helper.APIResponse("Successfuly updated user data", http.StatusOK, "success", updatedUser)
	c.JSON(http.StatusOK, response)
}

func (h *userHandler) GetCurrentUser(c *gin.Context){
	currentUser := c.MustGet("currentUser").(user.User)

	user := gin.H{
		"name": currentUser.Name,
		"email": currentUser.Email,
		"avatar_url": currentUser.AvatarFileName,
	}

	response := helper.APIResponse("Successfuly Get Current User", http.StatusOK, "success", user)
	c.JSON(http.StatusOK, response)
}