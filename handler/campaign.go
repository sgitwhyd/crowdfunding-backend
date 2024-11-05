package handler

import (
	"be-bwastartup/campaign"
	"be-bwastartup/helper"
	"be-bwastartup/user"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type campaignHandler struct {
	campaignService campaign.Service
}

func NewCampaignHandler(campaignService campaign.Service) *campaignHandler {
	return &campaignHandler{campaignService}
}

// @Tags Campaign
// @Summary Get All Campaign data
// @Description Get All Campaign data
// @Produce application/json
// @Success 200 {object} helper.response{data=[]campaign.CampaignFormatter}
// @Param user_id query string false "find by user_id" 
// @Router /campaigns [get]
func (h *campaignHandler) GetCampaigns(c *gin.Context) {
	
	userID, _ := strconv.Atoi(c.Query("user_id"))

	campaigns, err := h.campaignService.GetCampaigns(userID)
	if err != nil {
		errorsResponse := gin.H{"errors": err.Error()}
		response := helper.APIResponse("Get campaigns failed", http.StatusBadRequest, "error", errorsResponse)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("List of campaigns", http.StatusOK, "success", campaign.FormatCampaigns(campaigns))
	c.JSON(http.StatusOK, response)
}

// @Tags Campaign
// @Summary Create Campaign data
// @Description Create Campaign data
// @Produce application/json
// @Security BearerAuth
// @Param request body campaign.CreateCampaignInput true "Body Required"
// @Success 200 {object} helper.response{data=campaign.CampaignFormatter}
// @Router /campaigns [post]
func (h *campaignHandler) CreateCampaign(c *gin.Context) {
	var input = campaign.CreateCampaignInput{}

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := gin.H{"errors": err.Error()}

		response := helper.APIResponse("Create campaign failed", http.StatusUnprocessableEntity, "error", errors)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)



	newCampaign, err := h.campaignService.CreateCampaign(input, currentUser)
	if err != nil {
		errorsResponse := gin.H{"errors": err.Error()}
		response := helper.APIResponse("Create campaign failed", http.StatusBadRequest, "error", errorsResponse)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Campaign has been created", http.StatusOK, "success", campaign.FormatCampaign(newCampaign))
	c.JSON(http.StatusOK, response)
}

// @Tags Campaign
// @Summary Save Campaign Image data
// @Description Create Campaign Image data
// @Produce application/json
// @Security BearerAuth
// @Param request body campaign.CreateCampaignImageInput true "Body Required"
// @Success 200 {object} helper.response{data=helper.UploadImageResponse}
// @Router /campaigns/images [post]
func (h *campaignHandler) SaveCampaignImage(c *gin.Context){
	var input campaign.CreateCampaignImageInput

	err := c.ShouldBind(&input)
	if err != nil {
		errors := gin.H{"errors": err.Error()}

		response := helper.APIResponse("Save campaign image failed", http.StatusUnprocessableEntity, "error", errors)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		data := gin.H{"is_uploaded": false, "error": err.Error()}
		errorResponse := helper.APIResponse("Failed to upload campaign image", http.StatusUnprocessableEntity, "error", data)
		c.JSON(http.StatusUnprocessableEntity, errorResponse)
		return
	}
	
	path := fmt.Sprintf("images/campaign-images/%d-%s", input.CampaignID, file.Filename)

	err = c.SaveUploadedFile(file, path)
	if err != nil {
		data := gin.H{"is_uploaded": false, "error": err.Error()}
		errorResponse := helper.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)

	_, err = h.campaignService.UploadCampaignImage(input, path,currentUser)
	if err != nil {
		data := gin.H{"is_uploaded": false, "error": err.Error()}
		errorResponse := helper.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	data := gin.H{"is_uploaded": true}
	response := helper.APIResponse("Campaign image has been uploaded", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)

}

// @Tags Campaign
// @Summary Get Campaign Image data
// @Description Detail Campaign
// @Produce application/json
// @Security BearerAuth
// @Param id path  string true "Campaign ID"
// @Success 200 {object} helper.response{data=campaign.CampaignFormatter}
// @Router /campaigns/{id} [get]
func (h *campaignHandler) GetCampaign(c *gin.Context) {
	var input campaign.GetCampaignDetailInput

	err := c.ShouldBindUri(&input)
	if err != nil {
		errors := gin.H{"errors": err.Error()}

		response := helper.APIResponse("Get campaign failed", http.StatusBadRequest, "error", errors)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	campaignDetail, err := h.campaignService.GetCampaign(input)
	if err != nil {
		errorsResponse := gin.H{"errors": err.Error()}
		response := helper.APIResponse("Get campaign failed", http.StatusBadRequest, "error", errorsResponse)
		c.JSON(http.StatusBadRequest, response)
		return
	}
 
	response := helper.APIResponse("Campaign detail", http.StatusOK, "success",  campaign.FormatDetailCampaign(campaignDetail))
	c.JSON(http.StatusOK, response)
}

// @Tags Campaign
// @Summary Update Campaign data
// @Description Update Campaign
// @Produce application/json
// @Security BearerAuth
// @Param id path  string true "Campaign ID"
// @Param request body campaign.CreateCampaignInput true "Update Campaign Body"
// @Success 200 {object} helper.response{data=campaign.CampaignFormatter}
// @Router /campaigns/{id} [put]
func (h *campaignHandler) UpdateCampaign(c *gin.Context){
	var inputID campaign.GetCampaignDetailInput
	var inputData campaign.CreateCampaignInput

	err := c.ShouldBindUri(&inputID)
	if err != nil {
		errors := gin.H{"errors": err.Error()}

		response := helper.APIResponse("Update campaign failed", http.StatusBadRequest, "error", errors)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	err = c.ShouldBind(&inputData)
	if err != nil {
		errors := gin.H{"errors": err.Error()}

		response := helper.APIResponse("Update campaign failed", http.StatusBadRequest, "error", errors)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)

	updatedCampaign, err := h.campaignService.UpdateCampaign(inputID, inputData, currentUser)
	if err != nil {
		errorsResponse := gin.H{"errors": err.Error()}
		response := helper.APIResponse("Update campaign failed", http.StatusBadRequest, "error", errorsResponse)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Campaign has been updated", http.StatusOK, "success", campaign.FormatCampaign(updatedCampaign))
	c.JSON(http.StatusOK, response)
}