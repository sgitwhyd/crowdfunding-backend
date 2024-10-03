package handler

import (
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
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	user, nil := h.userService.RegisterUser(input)
	if(err != nil){
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	c.JSON(http.StatusOK, user)
	// map input dari user ke struct RegisterUserInput
	// struct di atas kita passing sebagai parameter service
	// di file user.go
}