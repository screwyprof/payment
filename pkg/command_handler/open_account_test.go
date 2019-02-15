package command_handler

import (
	"context"
	"fmt"
	"testing"

	"github.com/rhymond/go-money"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/screwyprof/payment/pkg/command"
	"github.com/screwyprof/payment/pkg/event"
)

func TestOpenAccountHandle_InvalidCommandGiven_ErrorReturned(t *testing.T) {
	t.Parallel()

	// arrange
	h := NewOpenAccount(nil, nil)

	// act
	err := h.Handle(context.Background(), command.Unknown{})

	// assert
	assert.EqualError(t, err, "invalid command command.Unknown{} given")
}

func TestOpenAccountHandle_EventStoreErrorOccurred_ErrorReturned(t *testing.T) {
	t.Parallel()

	// arrange
	expected := fmt.Errorf("an error occurred")
	eventStore := &eventStoreStub{}
	eventStore.Error = expected

	h := NewOpenAccount(eventStore, nil)

	// act
	err := h.Handle(context.Background(), command.OpenAccount{})

	// assert
	assert.EqualError(t, err, "cannot open account: an error occurred")
}

func TestOpenAccountHandle_ValidCommandGiven_AccountOpened(t *testing.T) {
	t.Parallel()

	// arrange
	expectedEvent := event.AccountOpened{
		Balance: *money.New(10000, "USD"),
	}

	eventStore := &eventStoreStub{}
	notifier := &notifierStub{}
	h := NewOpenAccount(eventStore, notifier)

	// act
	err := h.Handle(context.Background(), command.OpenAccount{Balance: *money.New(10000, "USD")})
	require.NoError(t, err)

	// assert
	e := notifier.Event.(event.AccountOpened)
	ev := eventStore.Event.(event.AccountOpened)

	assert.Equal(t, expectedEvent.Balance, ev.Balance)
	assert.Equal(t, expectedEvent.Balance, e.Balance)
}
