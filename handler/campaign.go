package handler

import (
	"be-bwastartup/auth"
	"be-bwastartup/campaign"
	"be-bwastartup/helper"
	"be-bwastartup/user"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type campaignHandler struct {
	campaignService campaign.Service
	authService 	 auth.Service
}

func NewCampaignHandler(campaignService campaign.Service, authService auth.Service) *campaignHandler {
	return &campaignHandler{campaignService, authService}
}

func (h *campaignHandler) FindCampaigns(c *gin.Context){
	userID, _ := strconv.Atoi(c.Query("user_id"))
	
	campaigns, err := h.campaignService.GetCampaigns(userID)
	if err != nil {
			errorResponse := helper.APIResponse("Failed to get campaigns", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	response := helper.APIResponse("List of campaigns", http.StatusOK, "success", campaign.FormatCampaigns(campaigns))
	c.JSON(http.StatusOK, response)
}


func (h *campaignHandler) FindCampaign(c *gin.Context){
	var input campaign.GetCampaignInput

	err := c.ShouldBindUri(&input)
	if err != nil {
		error := gin.H{"errors": err.Error()}
		errorResponse := helper.APIResponse("Failed to get campaign", http.StatusBadRequest, "error", error)
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	data, err := h.campaignService.GetCampaign(input)
	if err != nil {
		error := gin.H{"errors": err.Error()}
		errorResponse := helper.APIResponse("Failed to get campaign", http.StatusNotFound, "error", error)
		c.JSON(http.StatusNotFound, errorResponse)
		return
	}

	response := helper.APIResponse("Campaign detail", http.StatusOK, "success", campaign.FormatDetailCampaign(data))
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) FindCampaignBySlug(c *gin.Context){
	var input campaign.GetCampaignBySlugInput

	err := c.ShouldBindUri(&input)
	if err != nil {
		error := gin.H{"errors": err.Error()}
		errorResponse := helper.APIResponse("Failed to get campaign", http.StatusBadRequest, "error", error)
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	data, err := h.campaignService.GetCampaignBySlug(input)
	if err != nil {
		error := gin.H{"errors": err.Error()}
		errorResponse := helper.APIResponse("Failed to get campaign", http.StatusNotFound, "error", error)
		c.JSON(http.StatusNotFound, errorResponse)
		return
	}

	response := helper.APIResponse("Campaign detail", http.StatusOK, "success", campaign.FormatDetailCampaign(data))
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) CreateCampaign(c *gin.Context){
	var createCampaignPayload campaign.CreateCampaignInput

	err := c.ShouldBindJSON(&createCampaignPayload)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Failed to create campaign", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	createCampaignPayload.User = currentUser

	newCampaign, err := h.campaignService.CreateCampaign(createCampaignPayload)
	if err != nil {
		errors := gin.H{"errors": err.Error()}
		response := helper.APIResponse("Failed to create campaign", http.StatusBadRequest, "error", errors)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Campaign has been created", http.StatusOK, "success", campaign.FormatCampaign(newCampaign))
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) UpdateCampaign(c *gin.Context){
	var updateCampaignPayload campaign.CreateCampaignInput
	err := c.ShouldBindJSON(&updateCampaignPayload)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Failed to update campaign", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	updateCampaignPayload.User = currentUser

	var inputSlug campaign.GetCampaignBySlugInput

	err = c.ShouldBindUri(&inputSlug)
	if err != nil {
		errors := gin.H{"errors": err.Error()}
		response := helper.APIResponse("Failed to update campaign", http.StatusBadRequest, "error", errors)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	updatedCampaign, err := h.campaignService.UpdateCampaign(inputSlug, updateCampaignPayload)
	if err != nil {
		errors := gin.H{"errors": err.Error()}
		response := helper.APIResponse("Failed to update campaign", http.StatusBadRequest, "error", errors)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Campaign has been updated", http.StatusOK, "success", campaign.FormatCampaign(updatedCampaign))
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) UploadImage(c *gin.Context){
	var input campaign.UploadCampaignImageInput

	err := c.ShouldBind(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{
			"is_uploaded": false,
			"errors": errors,
		}
		response := helper.APIResponse("Failed to upload campaign image", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser
	userID := currentUser.ID

	file, err := c.FormFile("file")
	if err != nil {
		errorMessage := gin.H{
			"is_uploaded": false,
			"errors": errors.New("file is required"),
		}
		response := helper.APIResponse("Failed to upload campaign image", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	path := fmt.Sprintf("images/%d-%s", userID, file.Filename)

	err = c.SaveUploadedFile(file, path)
	if err != nil {
		errorMessage := gin.H{
			"is_uploaded": false,
			"errors": errors.New("failed to save uploaded file"),
		}
		response := helper.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	_,err = h.campaignService.UploadCampaignImage(input, path)
	if err != nil {
		if err.Error() == "Campaign not found" {
			errorMessage := gin.H{
				"errors": "Campaign not found",
			}
			response := helper.APIResponse("Failed to upload campaign image", http.StatusNotFound, "error", errorMessage)
			c.JSON(http.StatusNotFound, response)
			return
		}

		errorMessage := gin.H{
			"is_uploaded": false,
			"errors": err.Error(),
		}
		response := helper.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{"is_uploaded": true}
	response := helper.APIResponse("Campaign image has been uploaded", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)

	
}