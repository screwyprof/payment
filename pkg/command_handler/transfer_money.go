package command_handler

import (
	"context"
	"fmt"

	"github.com/screwyprof/payment/internal/pkg/cqrs"
	"github.com/screwyprof/payment/internal/pkg/observer"

	"github.com/screwyprof/payment/pkg/command"
	"github.com/screwyprof/payment/pkg/domain/account"
)

// TransferMoney Transfers Money from an Account to another Account.
type TransferMoney struct {
	accountProvider account.GetAccountByNumber
	eventStore      account.EventStore
	notifier        observer.Notifier
}

// NewTransferMoney Creates a new instance of TransferMoney.
func NewTransferMoney(
	accountProvider account.GetAccountByNumber,
	eventStore account.EventStore,
	notifier observer.Notifier) cqrs.CommandHandler {
	return &TransferMoney{
		accountProvider: accountProvider,
		eventStore:      eventStore,
		notifier:        notifier,
	}
}

// Handle Transfers Money from an Account to another Account.
func (h *TransferMoney) Handle(ctx context.Context, c cqrs.Command) error {
	cmd, ok := c.(command.TransferMoney)
	if !ok {
		return fmt.Errorf("invalid command %+#v given", c)
	}

	fromAcc, err := h.accountProvider.ByNumber(cmd.From)
	if err != nil {
		return fmt.Errorf("cannot transfer money from %s to %s: %v", cmd.From, cmd.To, err)
	}

	moneyTransferred, err := fromAcc.SendTransferTo(cmd.To, cmd.Amount)
	if err != nil {
		return fmt.Errorf("cannot transfer money from %s to %s: %v", cmd.From, cmd.To, err)
	}

	if err := h.eventStore.StoreEvent(moneyTransferred); err != nil {
		return fmt.Errorf("cannot transfer money from %s to %s: %v", cmd.From, cmd.To, err)
	}

	h.notifier.Notify(moneyTransferred)
	return nil
}
