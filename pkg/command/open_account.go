package command

import (
	"github.com/rhymond/go-money"
)

type OpenAccount struct {
	Balance money.Money
}

func (r OpenAccount) CommandID() string {
	return "OpenAccount"
}
