package query_handler

import "github.com/screwyprof/payment/pkg/report"

type accountProviderStub struct {
	ReturnedError    error
	ReturnedAccount  *report.Account
	ReturnedAccounts report.Accounts
}

func (m accountProviderStub) ByNumber(number string) (*report.Account, error) {
	if m.ReturnedError != nil {
		return &report.Account{}, m.ReturnedError
	}
	return m.ReturnedAccount, nil
}

func (m accountProviderStub) All() (report.Accounts, error) {
	if m.ReturnedError != nil {
		return nil, m.ReturnedError
	}
	return m.ReturnedAccounts, nil
}
