package handler

import (
	"be-bwastartup/auth"
	"be-bwastartup/campaign"
	"be-bwastartup/helper"
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