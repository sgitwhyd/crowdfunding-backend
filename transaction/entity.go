package transaction

import (
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
	CreatedAt  time.Time
	UpdatedAt  time.Time
	User  	 user.User
}