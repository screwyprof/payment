package command

import (
	"github.com/rhymond/go-money"
	"github.com/screwyprof/payment/pkg/domain/account"
)

type OpenAccount struct {
	Number  account.Number
	Balance money.Money
}

func (r OpenAccount) CommandID() string {
	return "OpenAccount"
}
