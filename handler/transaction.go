package handler

import (
	"be-bwastartup/campaign"
	"be-bwastartup/helper"
	"be-bwastartup/payment"
	"be-bwastartup/transaction"
	"be-bwastartup/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type transactionHandler struct {
	service transaction.Service
	campaignService campaign.Service
	paymentService payment.Service
}

func NewTransactionHandler(service transaction.Service, campaignService campaign.Service, paymentService payment.Service) *transactionHandler {
	return &transactionHandler{service, campaignService, paymentService}
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

func (h *transactionHandler) CreateTransaction(c *gin.Context){
	var input transaction.CreateTransactionInput

	err := c.ShouldBind(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Failed to create transaction", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)



	newTransaction, err := h.service.CreateTransaction(input, currentUser)
	if err != nil {
		data := gin.H{"errors": err.Error()}
		response := helper.APIResponse("Failed to create transaction", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Transaction has been created", http.StatusOK, "success", transaction.FormatTransaction(newTransaction))

	c.JSON(http.StatusOK, response)


}
// input ammount dari user
// maping dari input user ke struct input service
// panggil service buat transaksi, panggil sistem midtrans
// panggil repo create new transaksi