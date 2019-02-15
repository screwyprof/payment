package command_handler

import (
	"context"
	"fmt"
	"testing"

	"github.com/rhymond/go-money"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/screwyprof/payment/pkg/command"
	"github.com/screwyprof/payment/pkg/domain/account"
	"github.com/screwyprof/payment/pkg/event"
)

func TestReceiveMoneyHandle_InvalidCommandGiven_ErrorReturned(t *testing.T) {
	t.Parallel()

	// arrange
	h := NewReceiveMoney(nil, nil, nil)

	// act
	err := h.Handle(context.Background(), command.Unknown{})

	// assert
	assert.EqualError(t, err, "invalid command command.Unknown{} given")
}

func TestReceiveMoneyHandle_AccountProviderErrorOccurred_ErrorReturned(t *testing.T) {
	t.Parallel()

	// arrange
	expected := fmt.Errorf("an error occurred")
	accountProvider := accountProviderStub{
		ReturnedError: expected,
	}

	h := NewReceiveMoney(accountProvider, nil, nil)

	// act
	err := h.Handle(context.Background(), command.ReceiveMoney{From: account.Number("123"), To: account.Number("321")})

	// assert
	assert.EqualError(t, err, "cannot receive money from 123 to 321: an error occurred")
}

func TestReceiveMoneyHandle_ReceiveToTheSameAccountErrorOccurred_ErrorReturned(t *testing.T) {
	t.Parallel()

	// arrange
	accountProvider := accountProviderStub{
		ReturnedAccount: &account.Account{
			Number: account.Number("123"),
		},
	}

	h := NewReceiveMoney(accountProvider, nil, nil)

	// act
	err := h.Handle(context.Background(), command.ReceiveMoney{From: account.Number("123"), To: account.Number("123")})

	// assert
	assert.EqualError(t, err,
		"cannot receive money from 123 to 123: cannot transfer money to the same account: 123")
}

func TestReceiveMoneyHandle_EventStoreErrorOccurred_ErrorReturned(t *testing.T) {
	t.Parallel()

	// arrange
	accountProvider := accountProviderStub{
		ReturnedAccount: &account.Account{
			Number:  account.Number("321"),
			Balance: *money.New(10000, "USD"),
		},
	}

	expected := fmt.Errorf("an error occurred")
	eventStore := &eventStoreStub{}
	eventStore.Error = expected

	h := NewReceiveMoney(accountProvider, eventStore, nil)

	// act
	err := h.Handle(context.Background(), command.ReceiveMoney{
		From:   account.Number("123"),
		To:     account.Number("321"),
		Amount: *money.New(1000, "USD"),
	})

	// assert
	assert.EqualError(t, err, "cannot receive money from 123 to 321: an error occurred")
}

func TestReceiveMoneyHandle_ValidCommandGiven_MoneyReceived(t *testing.T) {
	t.Parallel()

	// arrange
	expectedEvent := event.MoneyReceived{
		From:    "123",
		To:      "321",
		Amount:  *money.New(1000, "USD"),
		Balance: *money.New(11000, "USD"),
	}

	accountProvider := accountProviderStub{
		ReturnedAccount: &account.Account{
			Number:  account.Number("321"),
			Balance: *money.New(10000, "USD"),
		},
	}

	eventStore := &eventStoreStub{}
	notifier := &notifierStub{}
	h := NewReceiveMoney(accountProvider, eventStore, notifier)

	// act
	err := h.Handle(context.Background(), command.ReceiveMoney{
		From:   account.Number("123"),
		To:     account.Number("321"),
		Amount: *money.New(1000, "USD"),
	})
	require.NoError(t, err)

	// assert
	e := notifier.Event.(event.MoneyReceived)
	ev := eventStore.Event.(event.MoneyReceived)

	assert.Equal(t, expectedEvent.Balance, ev.Balance)
	assert.Equal(t, expectedEvent.Balance, e.Balance)
}
