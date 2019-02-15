package command_handler

import (
	"context"
	"fmt"

	"github.com/screwyprof/payment/internal/pkg/cqrs"
	"github.com/screwyprof/payment/internal/pkg/observer"

	"github.com/screwyprof/payment/pkg/command"
	"github.com/screwyprof/payment/pkg/domain/account"
)

// ReceiveMoney Receives Money to the given another Account.
type ReceiveMoney struct {
	accountProvider account.GetAccountByNumber
	eventStore      account.EventStore
	notifier        observer.Notifier
}

// NewReceiveMoney Creates a new instance of ReceiveMoney.
func NewReceiveMoney(
	accountProvider account.GetAccountByNumber,
	eventStore account.EventStore,
	notifier observer.Notifier) cqrs.CommandHandler {
	return &ReceiveMoney{
		accountProvider: accountProvider,
		eventStore:      eventStore,
		notifier:        notifier,
	}
}

// Handle Receives Money from an Account to another Account.
func (h *ReceiveMoney) Handle(ctx context.Context, c cqrs.Command) error {
	cmd, ok := c.(command.ReceiveMoney)
	if !ok {
		return fmt.Errorf("invalid command %+#v given", c)
	}

	acc, err := h.accountProvider.ByNumber(cmd.To)
	if err != nil {
		return fmt.Errorf("cannot receive money from %s to %s: %v", cmd.From, cmd.To, err)
	}

	moneyReceived, err := acc.ReceiveMoneyFrom(cmd.From, cmd.Amount)
	if err != nil {
		return fmt.Errorf("cannot receive money from %s to %s: %v", cmd.From, cmd.To, err)
	}

	if err := h.eventStore.StoreEvent(moneyReceived); err != nil {
		return fmt.Errorf("cannot receive money from %s to %s: %v", cmd.From, cmd.To, err)
	}

	h.notifier.Notify(moneyReceived)
	return nil
}
