package transaction

import (
	"be-bwastartup/campaign"
	"be-bwastartup/payment"
	"be-bwastartup/user"
	"strconv"
)

type service struct {
	repository Repositoy
	paymentService payment.Service
	campaignRepository campaign.Repository
	
}

type Service interface {
	GetTransactionByCampaignID(input GetTransactionByIDInput) ([]Transaction, error)
	GetTransactionsUser(userID int) ([]Transaction, error)
	CreateTransaction(input CreateTransactionInput, currentUser user.User) (Transaction, error)
	ProsesPayment(input TransactionNotificationStatusInput) error
}

func NewService(repository Repositoy, paymentService payment.Service,campaignRepository campaign.Repository) *service {
	return &service{repository, paymentService, campaignRepository}
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
	campaign, err := s.campaignRepository.FindByID(input.CampaignID)
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

func (s *service) ProsesPayment(input TransactionNotificationStatusInput)  error {
	transacation_id, _ := strconv.Atoi(input.OrderID)

	transaction, err := s.repository.FindByID(transacation_id)
	if err != nil {
		return err
	}

	if input.PaymentType == "credit_card" && input.TransactionStatus == "capture" && input.FraudStatus == "accept" {
		transaction.Status = "paid"
	} else if  input.TransactionStatus == "settlement" {
		transaction.Status = "paid"
	} else if input.TransactionStatus == "cancel" || input.TransactionStatus == "deny" || input.FraudStatus == "expire" {
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
		campaign.BackerCount += 1
		campaign.CurrentAmount = campaign.CurrentAmount + updatedTransaction.Amount


		_, err := s.campaignRepository.Update(campaign)
		if err != nil {
			return err
		}
	}

	return nil


}