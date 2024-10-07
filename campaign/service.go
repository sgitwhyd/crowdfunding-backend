package campaign

import (
	"errors"

	"github.com/gosimple/slug"
)

type Service interface {
	GetCampaigns(userID int) ([]Campaign, error)
	GetCampaign(input GetCampaignInput) (Campaign, error)
	GetCampaignBySlug(input GetCampaignBySlugInput)(Campaign, error)
	CreateCampaign(input CreateCampaignInput)(Campaign, error)
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

func (s *service) GetCampaignBySlug(input GetCampaignBySlugInput) (Campaign, error) {
	campaign, err := s.repository.FindBySlug(input.Slug)
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

	slug := slug.Make(input.Name)

	// check slug already use or not
	findedCampaign, err := s.repository.FindBySlug(slug)
	if err != nil {
		return findedCampaign, err
	}

	if findedCampaign.ID != 0 {
		return findedCampaign, errors.New("title already taken")
	}
	
	campaign.Slug = slug

	savedCampaign, err := s.repository.Save(campaign)
	if err != nil {
		return campaign, err
	}

	return savedCampaign, nil
}
