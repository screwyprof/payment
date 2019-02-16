package report

import (
	"fmt"
	"github.com/rhymond/go-money"
)

type Ledger struct {
	Action string
	Amount money.Money
}

func (l Ledger) ToString() string {
	return fmt.Sprintf("%s, %s", l.Action, l.Amount.Display())
}
