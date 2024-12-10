package campaign

import (
	"gorm.io/gorm"
)

type Repository interface {
	Save(campaign Campaign) (Campaign, error)
	FindByUserID(userID int) ([]Campaign, error)
	FindAll() ([]Campaign, error)
	FindByID(campaignID int) (Campaign, error)
	CreateCampaignImage(campaignImage CampaignImage) (CampaignImage, error)
	MarkAllImagesAsNonPrimary(campaignID int) (bool, error)
	Update(campaign Campaign) (Campaign, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) Save(campaign Campaign) (Campaign, error) {
	err := r.db.Create(&campaign).Error
	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

func (r *repository) FindByUserID(userID int) ([]Campaign, error) {
	var campaigns []Campaign

	err := r.db.Where("user_id = ?", userID).Preload("CampaignImages", "campaign_images.is_primary = 1").Find(&campaigns).Error
	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}

func (r *repository) FindAll() ([]Campaign, error) {
	var campaigns []Campaign

		err := r.db.Preload("CampaignImages" , "campaign_images.is_primary = 1").Preload("User").Find(&campaigns).Error
		if err != nil {
			return campaigns, err
		}
	
return campaigns, nil
	
}

func (r *repository) CreateCampaignImage(campaignImage CampaignImage) (CampaignImage, error) {
	err := r.db.Create(&campaignImage).Error
	if err != nil {
		return campaignImage, err
	}

	return campaignImage, nil
}

func (r *repository) FindByID(campaignID int) (Campaign, error) {
	var campaign Campaign

	err := r.db.Preload("User").Preload("CampaignImages").
   Preload("Transactions", func(db *gorm.DB) *gorm.DB {
       return db.Where("campaign_id = ?", campaignID).Preload("User") // or use the actual user ID variable here
   }).Where("id = ?", campaignID).First(&campaign).Error

	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

func (r *repository) MarkAllImagesAsNonPrimary(campaignID int) (bool, error) {
	err := r.db.Model(&CampaignImage{}).Where("campaign_id = ?", campaignID).Update("is_primary", 0).Error
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *repository) Update(campaign Campaign) (Campaign, error) {
	err := r.db.Save(&campaign).Error
	if err != nil {
		return campaign, err
	}

	return campaign, nil
}