package command

import (
	"github.com/rhymond/go-money"

	"github.com/screwyprof/payment/pkg/domain/account"
)

type ReceiveMoney struct {
	From   account.Number
	To     account.Number
	Amount money.Money
}

func (r ReceiveMoney) CommandID() string {
	return "ReceiveMoney"
}
