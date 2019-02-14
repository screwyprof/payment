package account

import (
	"fmt"
	"strconv"
	"time"

	"github.com/rhymond/go-money"

	"github.com/screwyprof/payment/pkg/event"
)

// Number A unique account ID.
type Number string

// generateAccNumber Generates a random Account number.
func generateAccNumber() Number {
	return Number(strconv.Itoa(int(time.Now().UnixNano())))
}

var numberGenerator = generateAccNumber

// Account An Account representation.
type Account struct {
	Number  Number
	Balance money.Money
}

// CreateEmpty Creates a new account instance.
func CreateEmpty() *Account {
	return &Account{}
}

// OpenAccount Opens a new account with optional starting balance.
func (a *Account) OpenAccount(balance money.Money) event.AccountOpened {
	a.Number = numberGenerator()
	a.Balance = balance

	return event.AccountOpened{
		Number:  string(a.Number),
		Balance: a.Balance,
	}
}

// Deposit Credits the account.
func (a *Account) Deposit(amount money.Money) (event.MoneyDeposited, error) {
	balance, err := a.Balance.Add(&amount)
	if err != nil {
		return event.MoneyDeposited{}, fmt.Errorf("cannot deposit account %s: %v", a.Number, err)
	}

	a.Balance = *balance

	return event.MoneyDeposited{
		Number:  string(a.Number),
		Amount:  amount,
		Balance: *balance,
	}, nil
}

// SendTransferTo Creates a transaction to transfer money from an account to another account.
func (a *Account) SendTransferTo(to Number, amount money.Money) (event.MoneyTransferred, error) {
	if err := a.ensureAccountsAreDifferent(a.Number, to); err != nil {
		return event.MoneyTransferred{}, err
	}

	newBalance, err := a.Balance.Subtract(&amount)
	if err != nil {
		return event.MoneyTransferred{}, fmt.Errorf("cannot send transfer from %s to %s: %v", a.Number, to, err)
	}

	if newBalance.IsNegative() || newBalance.IsZero() {
		return event.MoneyTransferred{},
			fmt.Errorf("cannot send transfer from %s to %s: balance %s is not high enough",
				a.Number, to, newBalance.Display())
	}

	a.Balance = *newBalance

	return event.MoneyTransferred{
		From:    string(a.Number),
		To:      string(to),
		Amount:  amount,
		Balance: *newBalance,
	}, nil
}

// ReceiveMoneyFrom Receives the money from the given account.
func (a *Account) ReceiveMoneyFrom(from Number, amount money.Money) (event.MoneyReceived, error) {
	if err := a.ensureAccountsAreDifferent(from, a.Number); err != nil {
		return event.MoneyReceived{}, err
	}

	newBalance, err := a.Balance.Add(&amount)
	if err != nil {
		return event.MoneyReceived{}, fmt.Errorf("cannot receive money from %s to %s: %v", from, a.Number, err)
	}

	a.Balance = *newBalance

	return event.MoneyReceived{
		From:    string(from),
		To:      string(a.Number),
		Amount:  amount,
		Balance: *newBalance,
	}, nil
}

func (a *Account) ensureAccountsAreDifferent(from Number, to Number) error {
	if from == to {
		return fmt.Errorf("cannot transfer money to the same account: %v", from)
	}
	return nil
}

// ToString Renders the Account as a string.
func (a *Account) ToString() string {
	return fmt.Sprintf("#%s: %s", a.Number, a.Balance.Display())
}
