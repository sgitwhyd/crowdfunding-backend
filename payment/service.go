package payment

import (
	"be-bwastartup/user"
	"context"
	"os"
	"strconv"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type service struct{
}

type Service interface {
	GetPaymentURL(transaction Transaction, user user.User) (string, error)
}

func NewService() *service {
	return &service{
		
	}
}


var sn snap.Client


func (s *service) GetPaymentURL(transaction Transaction, user user.User) (string, error) {
		SERVER_KEY := os.Getenv("MIDTRANS_SERVER_KEY")

	sn.New(SERVER_KEY, midtrans.Sandbox)
	sn.Options.SetContext(context.Background())

		snapReq := &snap.Request{
			TransactionDetails: midtrans.TransactionDetails{
				OrderID: "order-" + strconv.Itoa(transaction.ID),
				GrossAmt: int64(transaction.Amount),
			},
			CustomerDetail: &midtrans.CustomerDetails{
				FName: user.Name,
				Email: user.Email,		
			},
			Items: &[]midtrans.ItemDetails{
			{
					ID: strconv.Itoa(transaction.ID),
					Name: transaction.Campaign.Name,
					Price: int64(transaction.Amount),
					Qty: 1,		
			},
			},
		}

		snapTokenRes, err := sn.CreateTransaction(snapReq)
		if err != nil {
			return "", err
		}


		return snapTokenRes.RedirectURL, nil

}