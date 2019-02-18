package reporting

import (
	"fmt"

	"github.com/screwyprof/payment/pkg/report"
)

type InMemoryAccountReporter struct {
	accounts map[string]*report.Account
}

func NewInMemoryAccountReporter() *InMemoryAccountReporter {
	return &InMemoryAccountReporter{
		accounts: make(map[string]*report.Account),
	}
}

func (r *InMemoryAccountReporter) ByNumber(number string) (*report.Account, error) {
	fmt.Printf("ReadSide: retreiving account %s\n", number)

	acc, ok := r.accounts[number]
	if !ok {
		return &report.Account{}, fmt.Errorf("account %s is not found", number)
	}

	fmt.Printf("ReadSide: account %s retrieved with balance %s\n", number, acc.Balance.Display())
	return acc, nil
}

func (r *InMemoryAccountReporter) All() (report.Accounts, error) {
	var accs []*report.Account
	for _, acc := range r.accounts {
		accs = append(accs, acc)
	}

	return accs, nil
}

func (r *InMemoryAccountReporter) Update(acc *report.Account) error {
	// add new account report
	current, ok := r.accounts[acc.Number]
	if !ok {
		r.accounts[acc.Number] = acc
		return nil
	}

	// update report
	current.Number = acc.Number
	current.Balance = acc.Balance
	current.Ledgers = append(current.Ledgers, acc.Ledgers...)

	r.accounts[acc.Number] = current

	fmt.Printf("ReadSide: updating account %s with balance %s\n", acc.Number, acc.Balance.Display())
	return nil
}
