package transaction

import (
	"be-bwastartup/campaign"
	"be-bwastartup/user"
	"time"
)

type Transaction struct {
	ID         int
	CampaignID int
	UserID     int
	Amount     int
	Status     string
	Code       string
	PaymentURL string
	User  		 user.User
	Campaign 	 campaign.Campaign
	CreatedAt  time.Time
	UpdatedAt  time.Time
}