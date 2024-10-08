package campaign

import "strings"

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

type UserFormat struct {
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}

type Images struct {
	ImageUrl  string `json:"image_url"`
	IsPrimary int   `json:"is_primary"`
}

type CampaignDetailFormat struct {
	ID               int        `json:"id"`
	UserID           int        `json:"user_id"`
	Name             string     `json:"name"`
	ShortDescription string     `json:"short_description"`
	Description      string     `json:"description"`
	Perks            []string   `json:"perks"`
	BackerCount      int        `json:"backer_count"`
	GoalAmount       int        `json:"goal_amount"`
	CurrentAmount    int        `json:"current_amount"`
	Slug             string     `json:"slug"`
	ImageURL         string		   `json:"image_url"`
	User             UserFormat `json:"user"`
	Images           []Images   `json:"images"`
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

func FormatDetailCampaign(campaign Campaign) CampaignDetailFormat {
	campaignDetail := CampaignDetailFormat{}
	campaignDetail.ID = campaign.ID
	campaignDetail.UserID = campaign.UserID
	campaignDetail.Name = campaign.Name
	campaignDetail.ShortDescription = campaign.ShortDescription
	campaignDetail.Description = campaign.Description
	campaignDetail.BackerCount = campaign.BackerCount
	campaignDetail.GoalAmount = campaign.GoalAmount
	campaignDetail.CurrentAmount = campaign.CurrentAmount
	campaignDetail.Slug = campaign.Slug
	campaignDetail.Perks = []string{}

	userFormat := UserFormat{
		Name:      campaign.User.Name,
		AvatarURL: campaign.User.AvatarFileName,
	}

	campaignDetail.User = userFormat

	if len(campaign.CampaignImages) > 0 {
		campaignDetail.ImageURL = campaign.CampaignImages[0].FileName
	}

	var perks []string
	for _, perk := range strings.Split(campaign.Perks, ","){
		perks = append(perks, strings.TrimSpace(perk))
	}

	campaignDetail.Perks = perks

	images := []Images{}

	for _, image := range campaign.CampaignImages {
		imageFormat := Images{}
		imageFormat.ImageUrl = image.FileName

		isPrimary := 0
		if image.IsPrimary == 1 {
			isPrimary = 1
		}

		imageFormat.IsPrimary = isPrimary
		images = append(images, imageFormat)
	}

	campaignDetail.Images = images

	return campaignDetail
}