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