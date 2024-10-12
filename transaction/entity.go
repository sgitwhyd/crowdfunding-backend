package transaction

import (
	"be-bwastartup/campaign"
	"be-bwastartup/user"
	"time"
)

type Transaction struct {
	ID         int  
	CampaignID int  
	UserId     int  
	Amount     int  
	Status     string 
	Code       string 
	PaymentURL string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	User  	 user.User
	Campaign   campaign.Campaign `gorm:"foreignKey:CampaignID"`
}

