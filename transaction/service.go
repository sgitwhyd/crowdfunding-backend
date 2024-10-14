package transaction

import (
	"be-bwastartup/campaign"
)

type Service interface {
	GetByCampaignID(input GetCampaignTransactionInput) ([]Transaction, error)
	GetByUserID(userID int) ([]Transaction, error)
	CreateTransaction(input CreateTransactionInput)(Transaction, error)
}

type service struct {
	repository         Repository
	campaignRepository campaign.Repository
}

func NewService(repository Repository, campaignRepository campaign.Repository) *service {
	return &service{repository, campaignRepository}
}

func (s *service) GetByCampaignID(input GetCampaignTransactionInput) ([]Transaction, error) {

	campaign, err := s.campaignRepository.FindByID(input.CampaignID)
	if err != nil {
		return []Transaction{}, err
	}

	transactions, err := s.repository.FindByCampaignID(campaign.ID)
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (s *service) GetByUserID(userID int) ([]Transaction, error){
	transactions, err := s.repository.FindByUserID(userID)
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (s *service) CreateTransaction(input CreateTransactionInput)(Transaction, error){


	transactionPayload := Transaction{}
	transactionPayload.Amount = input.Amount
	transactionPayload.CampaignID = input.CampaignID
	transactionPayload.User = input.User
	transactionPayload.Status = "pending"
	transactionPayload.Code = "example-code"


	transaction, err := s.repository.Create(transactionPayload)
	if err != nil {	
		return transaction, err
	}

	return transaction, nil
}