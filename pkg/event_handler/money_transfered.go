package event_handler

import (
	"fmt"

	"github.com/screwyprof/payment/internal/pkg/observer"
	"github.com/screwyprof/payment/pkg/event"
	"github.com/screwyprof/payment/pkg/report"
)

type MoneyTransferred struct {
	accountReporter report.AccountUpdater
}

func NewMoneyTransfered(accountReporter report.AccountUpdater) *MoneyTransferred {
	return &MoneyTransferred{
		accountReporter: accountReporter,
	}
}

func (h *MoneyTransferred) Handle(e observer.Event) {
	evn, ok := e.(event.MoneyTransferred)
	if !ok {
		return
	}
	//fmt.Printf("MoneyTransferedEventHandler: %s=%s, %s => %s\n",
	//	evn.From, evn.Balance.Display(), evn.Amount.Display(), evn.To)

	rep := &report.Account{
		ID:      evn.AggID,
		Number:  evn.From,
		Balance: evn.Balance,
		Ledgers: []report.Ledger{
			{
				Action: fmt.Sprintf("Transfer to %s", evn.To),
				Amount: evn.Amount,
			},
		},
	}

	h.accountReporter.Update(rep)
}
