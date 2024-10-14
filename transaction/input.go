package transaction

import "be-bwastartup/user"

type GetCampaignTransactionInput struct {
	CampaignID int `uri:"campaign_id" binding:"required"`
}

type CreateTransactionInput struct {
	CampaignID int `json:"campaign_id" binding:"required"`
	Amount int `json:"amount" binding:"required"`
	User       user.User
}