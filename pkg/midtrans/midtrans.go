package midtrans

import (
	"errors"
	"fmt"
	"os"

	"github.com/CRobinDev/BCCGembira_Nusastra/internal/dto"
	"github.com/google/uuid"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type IMidtrans interface {
	NewTransactionToken(req dto.PaymentRequest) (*snap.Response, error)
}

type Midtrans struct {
	Client snap.Client
}

func NewMidtrans() IMidtrans {
	client := snap.Client{}
	client.New(os.Getenv("MIDTRANS_SERVER_KEY"), midtrans.Sandbox)

	return &Midtrans{
		Client: client,
	}
}

func (m *Midtrans) NewTransactionToken(req dto.PaymentRequest) (*snap.Response, error) {
	request := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  req.OrderID,
			GrossAmt: req.Amount,
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
		Items: &[]midtrans.ItemDetails{
			{
				ID:    uuid.NewString(),
				Name:  req.Type + fmt.Sprintf(" for %s", req.CustomerName),
				Price: req.Amount,
				Qty:   1,
			},
		},
		EnabledPayments: snap.AllSnapPaymentType,
		CustomerDetail: &midtrans.CustomerDetails{
			FName: req.CustomerName,
			Email: req.CustomerEmail,
		},
		Expiry: &snap.ExpiryDetails{
			Duration: 15,
			Unit:     "minute",
		},
	}

	snapResp, err := m.Client.CreateTransaction(request)
	var midtransErr *midtrans.Error
	if errors.As(err, &midtransErr) && midtransErr == nil {
		return snapResp, nil
	}
	return snapResp, err
}
