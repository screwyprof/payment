package account

import (
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/rhymond/go-money"

	"github.com/screwyprof/payment/pkg/command"
	"github.com/screwyprof/payment/pkg/domain"
	"github.com/screwyprof/payment/pkg/event"
)

// Number A unique account ID.
type Number string

// GenerateNumber Generates a random Account number.
func GenerateNumber() Number {
	return Number("ACC" + strconv.Itoa(int(time.Now().UnixNano())))
}

// Account An Account representation.
type Account struct {
	AggID uuid.UUID

	Number  Number
	Balance money.Money
}

// Create Creates a new account instance.
func Create(aggID uuid.UUID) *Account {
	return &Account{
		AggID: aggID,
	}
}

// AggregateID Returns aggregate ID.
func (a *Account) AggregateID() uuid.UUID {
	return a.AggID
}

// OpenAccount Opens a new account with optional starting balance.
func (a *Account) OpenAccount(c command.OpenAccount) ([]domain.DomainEvent, error) {
	return []domain.DomainEvent{event.NewAccountOpened(a.AggID, string(c.Number), c.Balance)}, nil
}

// DepositMoney Credits the account.
func (a *Account) DepositMoney(c command.DepositMoney) ([]domain.DomainEvent, error) {
	balance, err := a.Balance.Add(&c.Amount)
	if err != nil {
		return nil, fmt.Errorf("cannot deposit account %s: %v", a.Number, err)
	}

	return []domain.DomainEvent{event.NewMoneyDeposited(a.AggID, c.Amount, *balance)}, nil
}

// TransferMoney Creates a transaction to transfer money from an account to another account.
func (a *Account) TransferMoney(c command.TransferMoney) ([]domain.DomainEvent, error) {
	if err := a.ensureAccountsAreDifferent(a.Number, Number(c.To)); err != nil {
		return nil, err
	}

	newBalance, err := a.Balance.Subtract(&c.Amount)
	if err != nil {
		return nil, fmt.Errorf("cannot send transfer from %s to %s: %v", a.Number, c.To, err)
	}

	if newBalance.IsNegative() || newBalance.IsZero() {
		return nil, fmt.Errorf("cannot send transfer from %s to %s: balance %s is not high enough",
			a.Number, c.To, newBalance.Display())
	}

	return []domain.DomainEvent{
		event.NewMoneyTransferred(a.AggregateID(), string(a.Number), c.To, c.Amount, *newBalance),
	}, nil
}

// ReceiveMoney Receives the money from the given account.
func (a *Account) ReceiveMoney(c command.ReceiveMoney) ([]domain.DomainEvent, error) {
	if err := a.ensureAccountsAreDifferent(Number(c.From), a.Number); err != nil {
		return nil, err
	}

	newBalance, err := a.Balance.Add(&c.Amount)
	if err != nil {
		return nil, fmt.Errorf("cannot receive money from %s to %s: %v", c.From, a.Number, err)
	}

	return []domain.DomainEvent{event.NewMoneyReceived(a.AggregateID(), c.From, c.To, c.Amount, *newBalance)}, nil
}

func (a *Account) OnAccountOpened(e event.AccountOpened) {
	a.Number = Number(e.Number)
	a.Balance = e.Balance
}

func (a *Account) OnMoneyDeposited(e event.MoneyDeposited) {
	a.Balance = e.Balance
	//  _ledgers.Add(new CreditMutation(cashDepositedEvent.Amount, new AccountNumber(string.Empty)));
}

func (a *Account) OnMoneyTransfered(e event.MoneyTransferred) {
	//_ledgers.Add(new CreditTransfer(moneyTransferSendEvent.Amount, new AccountNumber(moneyTransferSendEvent.TargetAccount)));
	a.Balance = e.Balance
}

func (a *Account) OnMoneyReceived(e event.MoneyReceived) {
	//_ledgers.Add(new DebitTransfer(moneyTransferReceivedEvent.Amount, new AccountNumber(moneyTransferReceivedEvent.TargetAccount)));
	a.Balance = e.Balance
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
