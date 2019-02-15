package repository

import (
	"fmt"
	"github.com/rhymond/go-money"

	"github.com/screwyprof/payment/internal/pkg/cqrs"

	"github.com/screwyprof/payment/pkg/domain/account"
	"github.com/screwyprof/payment/pkg/event"
)

type InMemoryAccountProvider struct {
	accounts map[account.Number]*account.Account
}

func NewInMemoryAccountReporter() *InMemoryAccountProvider {
	return &InMemoryAccountProvider{
		accounts: make(map[account.Number]*account.Account),
	}
}

func (r InMemoryAccountProvider) ByNumber(number account.Number) (*account.Account, error) {
	fmt.Printf("WriteSide: retrieving account %s\n", number)

	acc, ok := r.accounts[number]
	if !ok {
		return &account.Account{}, fmt.Errorf("account %s is not found", number)
	}

	fmt.Printf("WriteSide: account %s retrieved with balance %s\n", number, acc.Balance.Display())
	return acc, nil
}

func (r *InMemoryAccountProvider) StoreEvent(e cqrs.Event) error {
	switch evn := e.(type) {
	case event.AccountOpened:
		r.handleAccountOpened(evn)
	case event.MoneyTransferred:
		r.handleMoneyTransfered(evn)
	case event.MoneyReceived:
		r.handleMoneyReceived(evn)
	default:
		fmt.Printf("WriteSide: UNKNOWN event: %+#v\n", e)
	}

	return nil
}

func (r *InMemoryAccountProvider) handleAccountOpened(e event.AccountOpened) {
	fmt.Printf("WriteSide: handling AccountOpen event: %s = %s\n", e.Number, e.Balance.Display())
	r.UpdateBalance(e.Number, e.Balance)
}

func (r *InMemoryAccountProvider) handleMoneyTransfered(e event.MoneyTransferred) {
	fmt.Printf("WriteSide: handling MoneyTransfered event: %s => %s => %s, %s = %s\n",
		e.From, e.Amount.Display(), e.To, e.From, e.Balance.Display())

	r.UpdateBalance(e.From, e.Balance)
}

func (r *InMemoryAccountProvider) handleMoneyReceived(e event.MoneyReceived) {
	fmt.Printf("WriteSide: handling MoneyReceived event: %s <= %s <= %s, %s = %s\n",
		e.To, e.Amount.Display(), e.From, e.To, e.Balance.Display())

	r.UpdateBalance(e.From, e.Balance)
}

func (r *InMemoryAccountProvider) UpdateBalance(number string, balance money.Money) {
	acc := account.CreateEmpty()
	acc.Number = account.Number(number)
	acc.Balance = balance

	r.accounts[acc.Number] = acc
}
