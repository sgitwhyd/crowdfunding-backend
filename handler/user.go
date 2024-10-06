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
	
	
	token , err := h.authService.GenerateToken(newUser.ID)
	if err != nil {
		errorResponse := helper.APIResponse("Register account failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}



	registerResponse := user.FormatUser(newUser, token)

	response := helper.APIResponse("Account has been registered", http.StatusCreated, "success", registerResponse)

	c.JSON(http.StatusOK, response)
}

func (h *userHandler) Login(c *gin.Context) {
	var loginPayload user.LoginUserInput
	err := c.ShouldBindJSON(&loginPayload)
	if err != nil {
		// error handling
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		errorResponse := helper.APIResponse("Login failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, errorResponse)
		return
	}
	logged, err := h.userService.Login(loginPayload)
	
	if(err != nil){
		errorMessage := gin.H{"errors": err.Error()}
		errorResponse := helper.APIResponse("Login failed", http.StatusNotFound, "error", errorMessage)
		c.JSON(http.StatusNotFound, errorResponse)
		return
	}

	token, err := h.authService.GenerateToken(logged.ID)
	if(err != nil){
		errorMessage := gin.H{"errors": err.Error()}
		errorResponse := helper.APIResponse("Login failed", http.StatusNotFound, "error", errorMessage)
		c.JSON(http.StatusNotFound, errorResponse)
		return
	}

	loginResponse := helper.APIResponse("Successfuly loggedin", http.StatusOK, "success", user.FormatUser(logged, token))
	c.JSON(http.StatusOK, loginResponse)
}

func (h *userHandler) GetUsers(c *gin.Context){
	users, err := h.userService.GetUsers()
	if err != nil {
		errorResponse := helper.APIResponse("Failed to get users", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	response := helper.APIResponse("List of users", http.StatusOK, "success", user.FormatUsers(users))
	c.JSON(http.StatusOK, response)
}

func (h *userHandler) CheckEmailAvailability(c *gin.Context) {
	var emailInput user.CheckEmailInput
	err := c.ShouldBindJSON(&emailInput)
	if err != nil {
		// error handling
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		errorResponse := helper.APIResponse("Email checking failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, errorResponse)
		return
	}

	isEmailAvailable, err := h.userService.IsEmailAvailable(emailInput)
	if err != nil {
		response := helper.APIResponse("Email checking failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	errorMessage := "Email has been registered"
		if isEmailAvailable {
			errorMessage = "Email is available"
		}

		data := gin.H{"is_available": isEmailAvailable}
		response := helper.APIResponse(errorMessage, http.StatusOK, "success", data)
		c.JSON(http.StatusOK, response)
}

func (s *userHandler) UploadAvatar(c *gin.Context) {
	// JWT (sementara hardcode, seakan-akan user yang login ID = 1)
	userId := 5

	file, err := c.FormFile("avatar")
	if err != nil {
		data := gin.H{
			"is_uploaded": false,
			"error": "Required file avatar",
		}
		response := helper.APIResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
	}


	path := fmt.Sprintf("images/%d-%s", userId, file.Filename)

	err = c.SaveUploadedFile(file, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
	}

	
	user, err := s.userService.SaveAvatar(userId, path)
	if  err != nil {
		errorResponse := helper.APIResponse("Failed to upload avatar image", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	data := gin.H{"avatar_url": user.AvatarFileName}
	response := helper.APIResponse("Avatar successfuly uploaded", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)


}