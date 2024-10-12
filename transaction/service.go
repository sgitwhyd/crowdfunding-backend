package transaction

import (
	"be-bwastartup/campaign"
	"be-bwastartup/payment"
	"be-bwastartup/user"
)

type service struct {
	repository Repositoy
	paymentService payment.Service
	campaignService campaign.Service
}

type Service interface {
	GetTransactionByCampaignID(input GetTransactionByIDInput) ([]Transaction, error)
	GetTransactionsUser(userID int) ([]Transaction, error)
	CreateTransaction(input CreateTransactionInput, currentUser user.User) (Transaction, error)
}

func NewService(repository Repositoy, paymentService payment.Service, campaignService campaign.Service) *service {
	return &service{repository, paymentService, campaignService}
}

func (s *service) GetTransactionByCampaignID(input GetTransactionByIDInput) ([]Transaction, error) {

	transactions, err := s.repository.GetByCampaignID(input.ID)
	if err != nil {
		return transactions, err
	}

	if len(transactions) == 0 {
		return []Transaction{}, nil
	}

	return transactions, nil
}

func (s *service) GetTransactionsUser(userID int) ([]Transaction, error) {

	transactions, err := s.repository.GetByUserID(userID)
	if err != nil {
		return transactions, err
	}

	if len(transactions) == 0 {
		return []Transaction{}, nil
	}

	return transactions, nil
}

func (s *service) CreateTransaction(input CreateTransactionInput, currentUser user.User) (Transaction, error) {
	transaction := Transaction{}

	transaction.Amount = input.Amount
	transaction.CampaignID = input.CampaignID
	transaction.Status = "pending"
	transaction.Code = "testingCode"
	transaction.User = currentUser

	
	// check if campaign id not found
	campaign, err := s.campaignService.GetCampaign(campaign.GetCampaignInput{ID: input.CampaignID})
	if err != nil {
		return transaction, err
	}

	newTransaction, err := s.repository.Create(transaction)
	if err != nil {
		return newTransaction, err
	}

	newTransaction.Campaign = campaign

	transactionPayload := payment.Transaction{
		ID:     newTransaction.ID,
		Amount: newTransaction.Amount,
		Campaign: newTransaction.Campaign,
	}

	paymentURL, err := s.paymentService.GetPaymentURL(transactionPayload, currentUser)
	if err != nil {
		return newTransaction, err
	}

	newTransaction.PaymentURL = paymentURL

	newTransaction, err = s.repository.Update(newTransaction)
	if err != nil {
		return newTransaction, err
	}


	return newTransaction, nil
}