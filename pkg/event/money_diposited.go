package event

import (
	"github.com/rhymond/go-money"
)

type MoneyDeposited struct {
	Number  string
	Amount  money.Money
	Balance money.Money
}

func (e MoneyDeposited) EventID() string {
	return "MoneyDeposited"
}
