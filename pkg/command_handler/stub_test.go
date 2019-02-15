package command_handler

import (
	"github.com/screwyprof/payment/internal/pkg/cqrs"
	"github.com/screwyprof/payment/internal/pkg/observer"

	"github.com/screwyprof/payment/pkg/domain/account"
)

type accountProviderStub struct {
	ReturnedError   error
	ReturnedAccount *account.Account
}

func (m accountProviderStub) ByNumber(number account.Number) (*account.Account, error) {
	if m.ReturnedError != nil {
		return &account.Account{}, m.ReturnedError
	}
	return m.ReturnedAccount, nil
}


type eventStoreStub struct {
	Event cqrs.Event
	Error error
}

func (m *eventStoreStub) StoreEvent(event cqrs.Event) error {
	m.Event = event
	if m.Error != nil {
		return m.Error
	}
	return nil
}

type notifierStub struct {
	observer.Notifier
	Event observer.Event
}

func (n *notifierStub) Notify(event observer.Event) {
	n.Event = event
}