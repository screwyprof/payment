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

func TestTransferMoneyHandle_InvalidCommandGiven_ErrorReturned(t *testing.T) {
	t.Parallel()

	// arrange
	h := NewTransferMoney(nil, nil, nil)

	// act
	err := h.Handle(context.Background(), command.Unknown{})

	// assert
	assert.EqualError(t, err, "invalid command command.Unknown{} given")
}

func TestTransferMoneyHandle_AccountProviderErrorOccurred_ErrorReturned(t *testing.T) {
	t.Parallel()

	// arrange
	expected := fmt.Errorf("an error occurred")
	accountProvider := accountProviderStub{
		ReturnedError: expected,
	}

	h := NewTransferMoney(accountProvider, nil, nil)

	// act
	err := h.Handle(context.Background(), command.TransferMoney{From: account.Number("123"), To: account.Number("321")})

	// assert
	assert.EqualError(t, err, "cannot transfer money from 123 to 321: an error occurred")
}

func TestTransferMoneyHandle_TransferToTheSameAccountErrorOccurred_ErrorReturned(t *testing.T) {
	t.Parallel()

	// arrange
	accountProvider := accountProviderStub{
		ReturnedAccount: &account.Account{
			Number: account.Number("123"),
		},
	}

	h := NewTransferMoney(accountProvider, nil, nil)

	// act
	err := h.Handle(context.Background(), command.TransferMoney{From: account.Number("123"), To: account.Number("123")})

	// assert
	assert.EqualError(t, err,
		"cannot transfer money from 123 to 123: cannot transfer money to the same account: 123")
}

func TestTransferMoneyHandle_EventStoreErrorOccurred_ErrorReturned(t *testing.T) {
	t.Parallel()

	// arrange
	accountProvider := accountProviderStub{
		ReturnedAccount: &account.Account{
			Number:  account.Number("123"),
			Balance: *money.New(10000, "USD"),
		},
	}

	expected := fmt.Errorf("an error occurred")
	store := &accountStorageStub{ReturnedError: expected}

	h := NewTransferMoney(accountProvider, store, nil)

	// act
	err := h.Handle(context.Background(), command.TransferMoney{
		From:   account.Number("123"),
		To:     account.Number("321"),
		Amount: *money.New(1000, "USD"),
	})

	// assert
	assert.EqualError(t, err, "cannot transfer money from 123 to 321: an error occurred")
}

func TestTransferMoneyHandle_ValidCommandGiven_MoneyTransferred(t *testing.T) {
	t.Parallel()

	// arrange
	expectedAccount := &account.Account{
		Number:  account.Number("123"),
		Balance: *money.New(9000, "USD"),
	}

	expectedEvent := event.MoneyTransferred{
		From:    "123",
		To:      "321",
		Amount:  *money.New(1000, "USD"),
		Balance: *money.New(9000, "USD"),
	}

	accountProvider := accountProviderStub{
		ReturnedAccount: &account.Account{
			Number:  account.Number("123"),
			Balance: *money.New(10000, "USD"),
		},
	}

	store := &accountStorageStub{}
	notifier := &notifierStub{}
	h := NewTransferMoney(accountProvider, store, notifier)

	// act
	err := h.Handle(context.Background(), command.TransferMoney{
		From:   account.Number("123"),
		To:     account.Number("321"),
		Amount: *money.New(1000, "USD"),
	})
	require.NoError(t, err)

	// assert
	e := notifier.Event.(event.MoneyTransferred)

	assert.Equal(t, expectedAccount, store.AddedAccount)
	assert.Equal(t, expectedEvent.Balance, e.Balance)
}
