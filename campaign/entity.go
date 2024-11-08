package campaign

import (
	"be-bwastartup/user"
	"time"
)

type Campaign struct {
	ID               int
	UserID           int
	Name             string
	ShortDescription string
	Description      string
	Perks            string
	BackerCount      int
	GoalAmount       int
	CurrentAmount    int
	Slug             string
	User             user.User
	CampaignImages  []CampaignImage
	Transactions 	[]Transaction `gorm:"foreignKey:CampaignID"`
	CreatedAt     	 time.Time
	UpdatedAt     	 time.Time
}

type CampaignImage struct {
	ID  int
	CampaignID int
	FileName string
	IsPrimary int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Transaction struct {
	ID int
	UserID int
	CampaignID int 
	User user.User
	Amount int
	CreatedAt time.Time
}