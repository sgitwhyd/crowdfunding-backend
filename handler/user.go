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