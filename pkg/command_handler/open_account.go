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
	store    account.AccountStorage
	notifier observer.Notifier
}

// NewOpenAccount Creates a new instance of OpenAccount.
func NewOpenAccount(eventStore account.AccountStorage, notifier observer.Notifier) cqrs.CommandHandler {
	return &OpenAccount{store: eventStore, notifier: notifier}
}

// Handle Opens an account.
func (h *OpenAccount) Handle(ctx context.Context, c cqrs.Command) error {
	cmd, ok := c.(command.OpenAccount)
	if !ok {
		return fmt.Errorf("invalid command %+#v given", c)
	}

	acc := account.CreateEmpty()
	accountOpened := acc.OpenAccount(cmd.Number, cmd.Balance)
	if err := h.store.Add(acc); err != nil {
		return fmt.Errorf("cannot open account: %v", err)
	}

	h.notifier.Notify(accountOpened)
	return nil
}
