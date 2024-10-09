package campaign

import (
	"errors"
	"fmt"

	"github.com/gosimple/slug"
)

type Service interface {
	GetCampaigns(userID int) ([]Campaign, error)
	GetCampaign(input GetCampaignInput) (Campaign, error)
	CreateCampaign(input CreateCampaignInput)(Campaign, error)
	UpdateCampaign(slug GetCampaignInput, data CreateCampaignInput)(Campaign, error)
	UploadCampaignImage(input UploadCampaignImageInput, fileLocation string)(CampaignImage, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
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

func (s *service) GetCampaign(input GetCampaignInput) (Campaign, error) {
	campaign, err := s.repository.FindByID(input.ID)
	if err != nil {
		return campaign, err
	}

	if campaign.ID == 0 {
		return campaign, errors.New("Campaign not found")
	}

	return campaign, nil
}

func (s *service) GetCampaignByID(input GetCampaignInput) (Campaign, error) {
	campaign, err := s.repository.FindByID(input.ID)
	if err != nil {
		return campaign, err
	}

	if campaign.ID == 0 {
		return campaign, errors.New("Campaign not found")
	}

	return campaign, nil
}

func (s *service) CreateCampaign(input CreateCampaignInput)(Campaign, error){
	campaign := Campaign{}
	campaign.Name = input.Name
	campaign.ShortDescription = input.ShortDescription
	campaign.Description = input.Description
	campaign.GoalAmount = input.GoalAmount
	campaign.Perks = input.Perks
	campaign.UserID = input.User.ID


	// check slug already use or not
	slugCandidate := fmt.Sprintf("%s %d", input.Name, input.User.ID)
	campaign.Slug = slug.Make(slugCandidate)


	savedCampaign, err := s.repository.Save(campaign)
	if err != nil {
		return campaign, err
	}

	return savedCampaign, nil
}

func (s *service) UpdateCampaign(input GetCampaignInput, data CreateCampaignInput)(Campaign, error){
	campaign, err := s.repository.FindByID(input.ID)
	if err != nil {
		return campaign, err
	}

	campaign.Name = data.Name
	campaign.ShortDescription = data.ShortDescription
	campaign.Description = data.Description
	campaign.Perks = data.Perks
	campaign.GoalAmount = data.GoalAmount

	if campaign.User.ID != data.User.ID {
		return campaign, errors.New("not an owner of the campaign")
	}


	updatedCampaign, err := s.repository.Save(campaign)
	if err != nil {
		return updatedCampaign, err
	}

	return updatedCampaign, nil

}

func (s *service) UploadCampaignImage(input UploadCampaignImageInput, fileLocation string)(CampaignImage, error){
	campaign, err := s.repository.FindByID(input.CampaignID)
	if err != nil {
		return CampaignImage{}, err
	}

	if input.User.ID != campaign.UserID {
		return CampaignImage{}, errors.New("not an owner of the campaign")
	}


	if campaign.ID == 0 {
		return CampaignImage{}, errors.New("Campaign not found")
	}

	isPrimary := 0
	if input.IsPrimary {
		isPrimary = 1

		_, err := s.repository.MarkAllImageAsNonPrimary(input.CampaignID)
		if err != nil {
			return CampaignImage{}, err
		}
	}

	campaignImage := CampaignImage{}
	campaignImage.CampaignID = input.CampaignID
	campaignImage.FileName = fileLocation
	campaignImage.IsPrimary = isPrimary


	savedCampaignImage, err := s.repository.CreateImage(campaignImage)
	if err != nil {
		return savedCampaignImage, err
	}

	return savedCampaignImage, nil
}
