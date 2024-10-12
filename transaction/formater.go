package transaction

import (
	"time"
)

type TransactionFormat struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Amount     int    `json:"amount"`
	CreatedAt  time.Time `json:"created_at"`
}

type TransactionUserFormat struct {
	Name string `json:"name"`
}

type BackedFormat struct {
	ID int `json:"id"`
	Amount int `json:"amount"`
	Status string `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	Campaign CampaignFormat `json:"campaign"`
}

type CampaignFormat struct {
	Name string `json:"name"`
	ImageURL string `json:"image_url"`
}

func FormatTransactions(transactions []Transaction) []TransactionFormat {
	transactionFormat := []TransactionFormat{}

	for _, transaction := range transactions {
		transactionFormat = append(transactionFormat, TransactionFormat{
			ID:         int(transaction.ID),
			Name: 		 transaction.User.Name,
			Amount:     transaction.Amount,
			CreatedAt: transaction.CreatedAt,
		})

	}

	return transactionFormat
}

func FormatBackeds(transactions []Transaction) []BackedFormat {
	backedFormat := []BackedFormat{}

	if len(transactions) == 0 {
		return backedFormat
	}

	ImageUrl := ""
	

	for _, transaction := range transactions {

		if len(transaction.Campaign.CampaignImages) > 0 {
			ImageUrl = transaction.Campaign.CampaignImages[0].FileName
		}

		backedFormat = append(backedFormat, BackedFormat{
			ID: int(transaction.ID),
			Amount: transaction.Amount,
			Status: transaction.Status,
			CreatedAt: transaction.CreatedAt,
			Campaign: CampaignFormat{
				Name: transaction.Campaign.Name,
				ImageURL: ImageUrl,
			},
		})
	}

	return backedFormat
}

type TransactionFormatter struct {
	ID int `json:"id"`
	CampaignID int `json:"campaign_id"`
	UserID int `json:"user_id"`
	Amount int `json:"amount"`
	Status string `json:"status"`
	Code string `json:"code"`
	PaymentURL string `json:"payment_url"`
}

func FormatTransaction(transaction Transaction) TransactionFormatter {
	transactionFormat := TransactionFormatter{}

	transactionFormat.ID = int(transaction.ID)
	transactionFormat.CampaignID = transaction.CampaignID
	transactionFormat.UserID = transaction.User.ID
	transactionFormat.Amount = transaction.Amount
	transactionFormat.Status = transaction.Status
	transactionFormat.Code = transaction.Code
	transactionFormat.PaymentURL = transaction.PaymentURL

	return transactionFormat
}