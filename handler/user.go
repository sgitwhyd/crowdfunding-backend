package handler

import (
	"be-bwastartup/helper"
	"be-bwastartup/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) *userHandler {
	return &userHandler{userService}
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

	userToken := "token" // ini sementara


	registerResponse := user.FormatUser(newUser, userToken)

	response := helper.APIResponse("Account has been registered", http.StatusCreated, "success", registerResponse)

	c.JSON(http.StatusOK, response)
}

func (h *userHandler) Login(c *gin.Context) {
	// user memasukkan input (email dan password)
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
	// input ditangkap handler
	// mapping dari input user ke input struct

	// input struct passing service
	logged, err := h.userService.Login(loginPayload)
	
	if(err != nil){
		errorMessage := gin.H{"errors": err.Error()}
		errorResponse := helper.APIResponse("Login failed", http.StatusNotFound, "error", errorMessage)
		c.JSON(http.StatusNotFound, errorResponse)
		return
	}
	// generate jwt token
	token := "token"
	// kembalikan response
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

