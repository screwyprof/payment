package command

import (
	"github.com/rhymond/go-money"

	"github.com/screwyprof/payment/pkg/domain/account"
)

type TransferMoney struct {
	From   account.Number
	To     account.Number
	Amount money.Money
}

func (r TransferMoney) CommandID() string {
	return "TransferMoney"
}
