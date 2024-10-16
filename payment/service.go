package payment

import (
	"be-bwastartup/user"
	"os"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type Service interface{
	GeneratePaymentURL(transaction Transaction, user user.User)(string, error)
}

type service struct {}

func NewService() *service {
	return &service{}
}

var sm snap.Client

func (s *service) GeneratePaymentURL(transaction Transaction,  user user.User)(string, error){
	sm.New(os.Getenv("MIDTRANS_SERVER_KEY"), midtrans.Sandbox)

	snapReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID: transaction.ID,
			GrossAmt: int64(transaction.Amount),
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: user.Name,
			Email: user.Email,
		},
		Items: &[]midtrans.ItemDetails{
			{
				ID: transaction.ID,
				Name: transaction.Product,
				Qty: 1,
				Price: int64(transaction.Amount),
			},
		},
	}

	resp, err := sm.CreateTransaction(snapReq)
	if err != nil {
		return "", err
	}

	return resp.RedirectURL, nil
}

