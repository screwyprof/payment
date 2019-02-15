package command_handler

import (
	"context"
	"fmt"

	"github.com/screwyprof/payment/internal/pkg/cqrs"
	"github.com/screwyprof/payment/internal/pkg/observer"

	"github.com/screwyprof/payment/pkg/command"
	"github.com/screwyprof/payment/pkg/domain/account"
)

// OpenAccount Opens an account.
type OpenAccount struct {
	eventStore account.EventStore
	notifier   observer.Notifier
}

// NewOpenAccount Creates a new instance of OpenAccount.
func NewOpenAccount(eventStore account.EventStore, notifier observer.Notifier) cqrs.CommandHandler {
	return &OpenAccount{eventStore: eventStore, notifier: notifier}
}

// Handle Opens an account.
func (h *OpenAccount) Handle(ctx context.Context, c cqrs.Command) error {
	cmd, ok := c.(command.OpenAccount)
	if !ok {
		return fmt.Errorf("invalid command %+#v given", c)
	}

	accountOpened := account.CreateEmpty().OpenAccount(cmd.Balance)
	if err := h.eventStore.StoreEvent(accountOpened); err != nil {
		return fmt.Errorf("cannot open account: %v", err)
	}

	h.notifier.Notify(accountOpened)
	return nil
}
