package campaign

type CampaignFormat struct {
	ID               int    `json:"id"`
	UserID           int    `json:"user_id"`
	Title            string `json:"title"`
	ShortDescription string `json:"short_description"`
	GoalAmount       int    `json:"goal_amount"`
	CurrentAmount    int    `json:"current_amount"`
	Slug             string `json:"slug"`
	ImageUrl         string `json:"image_url"`
}

func FormatCampaigns(campaigns []Campaign) []CampaignFormat {

	campaignsResponse := []CampaignFormat{}

	for _, campaign := range campaigns {
		imageUrl := ""

		if len(campaign.CampaignImages) > 0 {
			imageUrl = campaign.CampaignImages[0].FileName
		}

		campaignResponse := CampaignFormat{
			ID:               campaign.ID,
			Title:            campaign.Name,
			ShortDescription: campaign.ShortDescription,
			GoalAmount:       campaign.GoalAmount,
			CurrentAmount:    campaign.CurrentAmount,
			Slug:             campaign.Slug,
			ImageUrl:         imageUrl,
			UserID:           campaign.UserID,
		}
		campaignsResponse = append(campaignsResponse, campaignResponse)

	}
	return campaignsResponse
}

func FormatCampaign(campaign Campaign) CampaignFormat {
	imageUrl := ""

	if len(campaign.CampaignImages) > 0 {
		imageUrl = campaign.CampaignImages[0].FileName
	}

	campaignResponse := CampaignFormat{
		ID:               campaign.ID,
		Title:            campaign.Name,
		ShortDescription: campaign.ShortDescription,
		GoalAmount:       campaign.GoalAmount,
		CurrentAmount:    campaign.CurrentAmount,
		Slug:             campaign.Slug,
		ImageUrl:         imageUrl,
		UserID:           campaign.UserID,
	}
	return campaignResponse
}