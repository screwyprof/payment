package command_handler

import (
	"context"
	"fmt"
	"github.com/screwyprof/payment/internal/pkg/observer"
	"github.com/screwyprof/payment/pkg/event"
	"testing"

	"github.com/rhymond/go-money"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/screwyprof/payment/internal/pkg/cqrs"

	"github.com/screwyprof/payment/pkg/command"
)

type eventStoreStub struct {
	Event cqrs.Event
	Error error
}

func (m *eventStoreStub) StoreEvent(event cqrs.Event) error {
	m.Event = event
	if m.Error != nil {
		return m.Error
	}
	return nil
}

type notifierStub struct {
	observer.Notifier
	Event observer.Event
}

func (n *notifierStub) Notify(event observer.Event) {
	n.Event = event
}

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
	//eventStore.On("StoreEvent", mock.Anything).Return(expected)

	h := NewOpenAccount(eventStore, nil)

	// act
	err := h.Handle(context.Background(), command.OpenAccount{})

	// assert
	//eventStore.AssertExpectations(t)
	assert.EqualError(t, err, "cannot open account: an error occurred")
}

func TestOpenAccountHandle_ValidCommandGiven_AccountOpened(t *testing.T) {
	t.Parallel()

	// arrange
	expectedEvent := event.AccountOpened{
		Balance: *money.New(10000, "USD"),
	}

	eventStore := &eventStoreStub{}
	//eventStore.On("StoreEvent", mock.Anything).Return(nil)

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
