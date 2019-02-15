package event_handler

import (
	"github.com/screwyprof/payment/internal/pkg/observer"

	"github.com/screwyprof/payment/pkg/event"
	"github.com/screwyprof/payment/pkg/report"
)

type AccountOpened struct {
	accountReporter report.AccountUpdater
}

func NewAccountOpened(accountReporter report.AccountUpdater) *AccountOpened {
	return &AccountOpened{
		accountReporter: accountReporter,
	}
}

func (h *AccountOpened) Handle(e observer.Event) {
	evn, ok := e.(event.AccountOpened)
	if !ok {
		return
	}

	//fmt.Printf("AccountOpenedEventHandler: %s = %s\n", evn.Number, evn.Balance.Display())
	h.accountReporter.Update(&report.Account{Number: string(evn.Number), Balance: evn.Balance})
}
