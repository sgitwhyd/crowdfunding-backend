package handler

import (
	"be-bwastartup/auth"
	"be-bwastartup/helper"
	"be-bwastartup/user"
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

// @Tags Auth
// @Summary Register Example
// @Description Register API
// @Produce application/json
// @Param request body user.RegisterUserInput true "Body Required"
// @Success 200 {object} helper.response{data=user.RegisterUserResponse}
// @Router /users [post]
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

// @Tags Auth
// @Summary Login Example
// @Description Login API
// @Produce application/json
// @Param request body user.LoginUserInput true "Body Required"
// @Success 200 {object} helper.response{data=user.RegisterUserResponse}
// @Router /sessions [post]
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

// @Tags Auth
// @Summary Check Email Avaiability Example
// @Description Check Email Avaiability API
// @Produce application/json
// @Param request body user.CheckEmailInput true "Body Required"
// @Success 200 {object} helper.response{data=user.CheckEmailAvailabilityResponse}
// @Router /email_checker [post]
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
		errorsResponse := user.CheckEmailAvailabilityResponse{
			IsAvailable: isEmailAvailable,
		}
		errorResponse := helper.APIResponse("Email checking failed", http.StatusBadRequest, "error", errorsResponse)
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	data :=  user.CheckEmailAvailabilityResponse{
			IsAvailable: isEmailAvailable,
		}

	metaMessage := "Email has been registered"
	if isEmailAvailable {
		metaMessage = "Email is available"
	}

	response := helper.APIResponse(metaMessage, http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}

// @Tags Auth
// @Summary Upload avatar Example
// @Description Upload avatar API
// @Produce application/json
// @Security BearerAuth
// @Accept multipart/form-data
// @Param file formData file true "File to upload"
// @Success 200 {object} helper.response{data=user.UploadAvatarResponse}
// @Router /avatars [post]
func (h *userHandler) UploadAvatar(c *gin.Context){
	file, err := c.FormFile("file")
	if err != nil {
		data := gin.H{"is_uploaded": false}
		errorResponse := helper.APIResponse("Failed to upload avatar image", http.StatusUnprocessableEntity, "error", data)
		c.JSON(http.StatusUnprocessableEntity, errorResponse)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)


	_, err = h.userService.UploadAvatar(currentUser.ID, file)
	if err != nil {
		data := user.UploadAvatarResponse{
			IsUploaded: false,
		}
		errorResponse := helper.APIResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	data := user.UploadAvatarResponse{
			IsUploaded: true,
		}
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

// @Tags Auth
// @Summary Current User Example
// @Description Current User API
// @Produce application/json
// @Success 200 {object} helper.response{data=handler.GetCurrentUser.UserResponse}
// @Router /users/current [get]
// @Security BearerAuth
func (h *userHandler) GetCurrentUser(c *gin.Context){

	type UserResponse struct {
		ID int `json:"id"`
		Name string `json:"name"`
		Email string `json:"email"`
		AvatarURL string `json:"avatar_url"`

	}

	currentUser := c.MustGet("currentUser").(user.User)

	user := UserResponse{
		ID: currentUser.ID,
		Name: currentUser.Name,
		Email: currentUser.Email,
		AvatarURL: currentUser.AvatarFileName,
	}

	response := helper.APIResponse("Successfuly Get Current User", http.StatusOK, "success", user)
	c.JSON(http.StatusOK, response)
}