package event

import (
	"github.com/rhymond/go-money"
)

type MoneyReceived struct {
	From    string
	To      string
	Amount  money.Money
	Balance money.Money
}

func (e MoneyReceived) EventID() string {
	return "MoneyReceived"
}
