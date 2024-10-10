package handler

import (
	"be-bwastartup/campaign"
	"be-bwastartup/helper"
	"be-bwastartup/transaction"
	"be-bwastartup/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type transactionHandler struct {
	service transaction.Service
	campaignService campaign.Service
}

func NewTransactionHandler(service transaction.Service, campaignService campaign.Service) *transactionHandler {
	return &transactionHandler{service, campaignService}
}

func (h *transactionHandler) GetCampaignTransactions(c *gin.Context){
	var input transaction.GetTransactionByIDInput

	err := c.ShouldBindUri(&input)
	if err != nil {
		errorsResponse := helper.FormatValidationError(err)
		data := gin.H{"errors": errorsResponse}
		errorResponse := helper.APIResponse("Failed to get campaign's transactions", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	campaignInput := campaign.GetCampaignInput{ID: input.ID}
	campaign, err := h.campaignService.GetCampaign(campaignInput)
	if err != nil {
		errorResponse := helper.APIResponse("Failed to get campaign's transactions", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	if campaign.UserID != currentUser.ID {
		errorResponse := helper.APIResponse("You are not the owner of this campaign", http.StatusForbidden, "error", nil)
		c.JSON(http.StatusForbidden, errorResponse)
		return
	}

	transactions, err := h.service.GetTransactionByCampaignID(input)
	if err != nil {
		errorResponse := helper.APIResponse("Failed to get campaign's transactions", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	response := helper.APIResponse("Campaign's transactions", http.StatusOK, "success", transaction.FormatTransactions(transactions))
	c.JSON(http.StatusOK, response)
}

func (h *transactionHandler) GetUserTransactions(c *gin.Context){
	currentUser := c.MustGet("currentUser").(user.User)

	transactions, err := h.service.GetTransactionsUser(currentUser.ID)
	if err != nil {
		errorResponse := helper.APIResponse("Failed to get user's transactions", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	response := helper.APIResponse("List of backed campaign", http.StatusOK, "success", transaction.FormatBackeds(transactions))
	c.JSON(http.StatusOK, response)
}