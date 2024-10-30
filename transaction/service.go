package transaction

import (
	"be-bwastartup/campaign"
	"be-bwastartup/payment"
	"be-bwastartup/user"
	"strconv"
)

type Service interface {
	GetByCampaignID(input GetCampaignTransactionInput) ([]Transaction, error)
	GetByUserID(userID int) ([]Transaction, error)
	CreateTransaction(input CreateTransactionInput, user user.User)(Transaction, error)
	ProsesPayment(input TransactionNotificationInput) error
}

type service struct {
	repository         Repository
	campaignRepository campaign.Repository
	paymentService payment.Service
}

func NewService(repository Repository, campaignRepository campaign.Repository, paymentService payment.Service) *service {
	return &service{repository, campaignRepository,paymentService}
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

func (s *service) CreateTransaction(input CreateTransactionInput, user user.User)(Transaction, error){


	transactionPayload := Transaction{}
	transactionPayload.Amount = input.Amount
	transactionPayload.CampaignID = input.CampaignID
	transactionPayload.User = user
	transactionPayload.Status = "pending"
	transactionPayload.Code = "example-code"


	transaction, err := s.repository.Create(transactionPayload)
	if err != nil {	
		return transaction, err
	}

	campaign, err := s.campaignRepository.FindByID(transaction.CampaignID)
	if err != nil {
		return transaction, err
	}

	paymentPayload := payment.Transaction{
		ID:    strconv.Itoa(transaction.ID),
		Amount: transaction.Amount,
		Product: campaign.Name,
	}

	paymentURL, err := s.paymentService.GeneratePaymentURL(paymentPayload, user)
	if err != nil {
		return transaction, err
	}

	transaction.PaymentURL = paymentURL

	updatedTransaction, err := s.repository.Update(transaction)
	if err != nil {
		return updatedTransaction, err
	}

	return updatedTransaction, nil
}

func (s *service) ProsesPayment(input TransactionNotificationInput) error {
	transaction_id, err := strconv.Atoi(input.OrderID)
	if err != nil {
		return err
	}

	transaction, err := s.repository.FindByID(transaction_id)
	if err != nil {
		return err
	}

	if input.PaymentType == "credit_card" && input.TransactionStatus == "capture" && input.FraudStatus == "accept" {
		transaction.Status = "paid"
	} else if input.TransactionStatus == "settlement" {
		transaction.Status = "paid"
	} else if input.TransactionStatus == "deny" || input.TransactionStatus == "expire" || input.TransactionStatus == "cancel" {
		transaction.Status = "cancelled"
	}

	updatedTransaction, err := s.repository.Update(transaction)
	if err != nil {
		return err
	}

	campaign, err := s.campaignRepository.FindByID(updatedTransaction.CampaignID)
	if err != nil {
		return err
	}

	if updatedTransaction.Status == "paid" {
		campaign.BackerCount = campaign.BackerCount + 1
		campaign.CurrentAmount = campaign.CurrentAmount + updatedTransaction.Amount

		_, err := s.campaignRepository.Update(campaign)
		if err != nil {
			return err
		}
	}

	return nil
}
