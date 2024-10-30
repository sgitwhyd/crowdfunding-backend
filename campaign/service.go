package campaign

import (
	"be-bwastartup/user"
	"errors"
	"fmt"

	"github.com/gosimple/slug"
)

type Service interface {
	CreateCampaign(input CreateCampaignInput, user user.User) (Campaign, error)
	GetCampaigns(userID int) ([]Campaign, error)
	UploadCampaignImage(input CreateCampaignImageInput, fileLocation string, user user.User) (CampaignImage, error)
	GetCampaign(input GetCampaignDetailInput) (Campaign, error)
	UpdateCampaign(inputID GetCampaignDetailInput, input CreateCampaignInput, user user.User) (Campaign, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) CreateCampaign(input CreateCampaignInput, user user.User) (Campaign, error) {
	campaign := Campaign{}
	campaign.Name = input.Name
	campaign.ShortDescription = input.ShortDescription
	campaign.Description = input.Description
	campaign.Perks = input.Perks
	campaign.GoalAmount = input.GoalAmount
	campaign.User = user

	slugFormat := fmt.Sprintf("%s %d", input.Name, user.ID)
	campaign.Slug = slug.Make(slugFormat)

	newCampaign, err := s.repository.Save(campaign)
	if err != nil {
		return newCampaign, err
	}

	return newCampaign, nil
}

func (s *service) GetCampaigns(userID int) ([]Campaign, error) {
	if userID != 0 {
		campaigns, err := s.repository.FindByUserID(userID)
		if err != nil {
			return campaigns, err
		}

		return campaigns, nil
	}

	campaigns, err := s.repository.FindAll()
	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}

func (s *service) UploadCampaignImage(input CreateCampaignImageInput, fileLocation string,  user user.User) (CampaignImage, error) {

	campaign, err := s.repository.FindByID(input.CampaignID)
	if err != nil {
		return CampaignImage{}, errors.New("campaign not found")
	}

	if campaign.User.ID != user.ID {
		return CampaignImage{}, errors.New("not an owner of the campaign")
	}

	isPrimary := 0
	if input.IsPrimary {
		isPrimary = 1
		_, err := s.repository.MarkAllImagesAsNonPrimary(input.CampaignID)
		if err != nil {
			return CampaignImage{}, err
		}
	}

	campaignImage := CampaignImage{}

	campaignImage.CampaignID = input.CampaignID
	campaignImage.FileName = fileLocation
	campaignImage.IsPrimary = isPrimary

	newCampaignImage, err := s.repository.CreateCampaignImage(campaignImage)
	if err != nil {
		return newCampaignImage, err
	}

	return newCampaignImage, nil
}

func (s *service) GetCampaign(input GetCampaignDetailInput) (Campaign, error) {
	campaign, err := s.repository.FindByID(input.ID)
	if err != nil {
		return campaign, err
	}


	return campaign, nil
}

func (s *service) UpdateCampaign(inputID GetCampaignDetailInput, input CreateCampaignInput, user user.User) (Campaign, error){
	campaign, err := s.repository.FindByID(inputID.ID)
	if err != nil {
		return campaign, err
	}

	if campaign.UserID != user.ID {
		return campaign, errors.New("not an owner of the campaign")
	}

	campaign.Name = input.Name
	campaign.ShortDescription = input.ShortDescription
	campaign.Description = input.Description
	campaign.Perks = input.Perks
	campaign.GoalAmount = input.GoalAmount
	campaign.User = user

	updatedCampaign, err := s.repository.Update(campaign)
	if err != nil {
		return updatedCampaign, err
	}

	return updatedCampaign, nil
}