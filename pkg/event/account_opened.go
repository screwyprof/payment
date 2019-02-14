package event

import (
	"github.com/rhymond/go-money"
)

type AccountOpened struct {
	Number  string
	Balance money.Money
}

func (e AccountOpened) EventID() string {
	return "AccountOpened"
}
