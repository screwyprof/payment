package report

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/rhymond/go-money"
)

// Account An Account representation.
type Account struct {
	ID      uuid.UUID
	Number  string
	Balance money.Money
	Ledgers []Ledger
}

type Accounts []*Account

// ToString Renders the Account as a string.
func (a *Account) ToString() string {
	return fmt.Sprintf("#%s: %s", a.Number, a.Balance.Display())
}
