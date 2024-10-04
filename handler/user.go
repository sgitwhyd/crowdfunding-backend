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
	// map input dari user ke struct RegisterUserInput
	// struct di atas kita passing sebagai parameter service
	// di file user.go
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