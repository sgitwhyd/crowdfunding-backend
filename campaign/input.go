package campaign

import "be-bwastartup/user"

type GetCampaignInput struct {
	ID int `uri:"id" binding:"required"`
}

type GetCampaignBySlugInput struct {
	Slug string `uri:"slug" binding:"required"`
}

type CreateCampaignInput struct {
	Name             string `json:"name" binding:"required"`
	ShortDescription string `json:"short_description" binding:"required"`
	Description      string `json:"description" binding:"required"`
	GoalAmount       int    `json:"goal_amount" binding:"required"`
	Perks            string `json:"perks" binding:"required"`
	User             user.User
}