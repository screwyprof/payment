package command_handler

import (
	"github.com/screwyprof/payment/internal/pkg/observer"

	"github.com/screwyprof/payment/pkg/domain/account"
)

type accountProviderStub struct {
	ReturnedError   error
	ReturnedAccount *account.Account
}

func (s accountProviderStub) ByNumber(number account.Number) (*account.Account, error) {
	if s.ReturnedError != nil {
		return &account.Account{}, s.ReturnedError
	}
	return s.ReturnedAccount, nil
}

type accountStorageStub struct {
	ReturnedError error
	AddedAccount  *account.Account
}

func (s *accountStorageStub) Add(acc *account.Account) error {
	if s.ReturnedError != nil {
		return s.ReturnedError
	}

	s.AddedAccount = acc
	return nil
}

type notifierStub struct {
	observer.Notifier
	Event observer.Event
}

func (n *notifierStub) Notify(event observer.Event) {
	n.Event = event
}
