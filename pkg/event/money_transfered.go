package event

import (
	"github.com/rhymond/go-money"
)

type MoneyTransferred struct {
	From    string
	To      string
	Amount  money.Money
	Balance money.Money
}

func (e MoneyTransferred) EventID() string {
	return "MoneyTransferred"
}
