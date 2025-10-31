package report

import (
	"fmt"
	"github.com/Rhymond/go-money"
)

type Ledger struct {
	Action string
	Amount money.Money
}

func (l Ledger) ToString() string {
	return fmt.Sprintf("%s, %s", l.Action, l.Amount.Display())
}
