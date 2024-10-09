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