package transaction

import "be-bwastartup/user"

type GetTransactionByIDInput struct {
	ID int `uri:"id" binding:"required"`
}

type CreateTransactionInput struct {
	Amount int `json:"amount" binding:"required"`
	User   user.User `json:"user" binding:"required"`
	CampaignID int `json:"campaign_id" binding:"required"`
}

type UpdateTransactionInput struct {
	ID 	 int `json:"id" binding:"required"`
	Status string `json:"status" binding:"required"`
	
}