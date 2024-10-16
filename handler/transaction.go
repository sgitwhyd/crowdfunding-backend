package handler

import (
	"be-bwastartup/helper"
	"be-bwastartup/transaction"
	"be-bwastartup/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type transactionHandler struct {
	service transaction.Service
}

func NewTransactionHandler(service transaction.Service) *transactionHandler {
	return &transactionHandler{service}
}

func (h *transactionHandler) GetTransactionsByCampaignID(c *gin.Context){
	var input transaction.GetCampaignTransactionInput

	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helper.APIResponse("Failed to get campaign's transactions", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	transactions, err := h.service.GetByCampaignID(input)
	if err != nil {
		errorsResponse := gin.H{"errors": err.Error()}
		response := helper.APIResponse("Failed to get campaign's transactions", http.StatusBadRequest, "error", errorsResponse)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Campaign's transactions", http.StatusOK, "success", transaction.FormatCampaignTransactions(transactions))
	c.JSON(http.StatusOK, response)
}

func (h *transactionHandler) GetTransactionByUserID(c *gin.Context){
	currentUser := c.MustGet("currentUser").(user.User)
	userID := currentUser.ID

	transactions, err := h.service.GetByUserID(userID)
	if err != nil {
		errorResponse := gin.H{
			"error": err.Error(),
		}

		response := helper.APIResponse("Failed Get User Transaction", http.StatusBadRequest, "error", errorResponse)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Successfuly Get User Transaction", http.StatusOK, "success", transaction.FormatUserTransactions(transactions))
	c.JSON(http.StatusOK, response)
}

func (h *transactionHandler) CreateTransaction(c *gin.Context){
	var input transaction.CreateTransactionInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errorResponse := gin.H{
			"error": err.Error(),
			}

		response := helper.APIResponse("Failed Create Transaction", http.StatusBadRequest, "error", errorResponse)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser

	createdTransaction, err := h.service.CreateTransaction(input)
	if err != nil {
		errorResponse := gin.H{
			"error": "Campaign Not Found",
			}

		response := helper.APIResponse("Failed Create Transaction", http.StatusBadRequest, "error", errorResponse)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Transaction Created", http.StatusOK, "success", transaction.FormatTransaction(createdTransaction))
	c.JSON(http.StatusOK, response)
}

func (h *transactionHandler) GetNotification(c *gin.Context){
	var input transaction.TransactionNotificationInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		response := helper.APIResponse("Failed to get notification", http.StatusUnprocessableEntity, "error", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	err = h.service.ProsesPayment(input)
	if err != nil {
		response := helper.APIResponse("Failed to get notification", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Successfuly send notification", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}