package bus

import (
	"context"
	"fmt"
)

// CommandBus a command handler bus.
type CommandBus struct {
	handlers map[string]CommandHandler
}

// NewCommandBus creates a new instance of CommandBus.
func NewCommandBus() *CommandBus {
	return &CommandBus{
		handlers: make(map[string]CommandHandler),
	}
}

// Handle dispatches the given command to its appropriate CommandHandler and handles it.
func (b *CommandBus) Handle(ctx context.Context, c Command) error {
	h, err := b.resolve(c.CommandID())
	if err != nil {
		return err
	}
	return h.Handle(ctx, c)
}

// Register registers an  CommandHandler for the given CommandID.
func (b *CommandBus) Register(CommandID string, i CommandHandler) *CommandBus {
	b.handlers[CommandID] = i
	return b
}

func (b *CommandBus) resolve(commandID string) (CommandHandler, error) {
	h, ok := b.handlers[commandID]
	if !ok {
		return nil, fmt.Errorf("handler for %q is not found", commandID)
	}
	return h, nil
}
