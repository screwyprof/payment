package command_handler

import (
	"context"
	"fmt"
	"github.com/screwyprof/payment/pkg/domain/account"
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
	store := &accountStorageStub{}
	store.ReturnedError = expected

	h := NewOpenAccount(store, nil)

	// act
	err := h.Handle(context.Background(), command.OpenAccount{})

	// assert
	assert.EqualError(t, err, "cannot open account: an error occurred")
}

func TestOpenAccountHandle_ValidCommandGiven_AccountOpened(t *testing.T) {
	t.Parallel()

	// arrange
	expectedNumber := account.Number("777")
	expectedBalance := *money.New(10000, "USD")

	expectedAccount := &account.Account{
		Number:  expectedNumber,
		Balance: expectedBalance,
	}

	expectedEvent := event.AccountOpened{
		Number:  string(expectedNumber),
		Balance: expectedBalance,
	}

	store := &accountStorageStub{}
	notifier := &notifierStub{}
	h := NewOpenAccount(store, notifier)

	// act
	err := h.Handle(context.Background(), command.OpenAccount{Number: expectedNumber, Balance: expectedBalance})
	require.NoError(t, err)

	// assert
	e := notifier.Event.(event.AccountOpened)

	assert.Equal(t, expectedAccount, store.AddedAccount)
	assert.Equal(t, expectedEvent.Balance, e.Balance)
}
