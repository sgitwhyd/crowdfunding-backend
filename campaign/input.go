package campaign

type GetCampaignInput struct {
	ID int `uri:"id" binding:"required"`
}