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

func (r *InMemoryAccountReporter) Update(acc *report.Account) error {
	r.accounts[acc.Number] = acc

	fmt.Printf("ReadSide: updating account %s with balance %s\n", acc.Number, acc.Balance.Display())
	return nil
}
