package transaction

type GetCampaignTransactionInput struct {
	CampaignID int `uri:"campaign_id" binding:"required"`
}

type CreateTransactionInput struct {
	CampaignID int `json:"campaign_id" binding:"required"`
	Amount int `json:"amount" binding:"required"`
}

type TransactionNotificationInput struct {
	TransactionStatus string `json:"transaction_status" binding:"required"`
	OrderID           string `json:"order_id"`
	PaymentType       string `json:"payment_type"`
	FraudStatus       string `json:"fraud_status"`
}