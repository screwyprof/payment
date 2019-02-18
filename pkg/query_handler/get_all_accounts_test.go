package query_handler

import (
	"context"
	"fmt"
	"testing"

	"github.com/rhymond/go-money"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/screwyprof/payment/pkg/query"
	"github.com/screwyprof/payment/pkg/report"
)

func TestGetAllAccountsHandle_InvalidQueryGiven_ErrorReturned(t *testing.T) {
	t.Parallel()

	// arrange
	h := NewGetAllAccounts(nil)

	// act
	err := h.Handle(context.Background(), query.Unknown{}, nil)

	// assert
	assert.EqualError(t, err, "invalid query query.Unknown{} given")
}

func TestGetAllAccountsHandle_InvalidResultGiven_ErrorReturned(t *testing.T) {
	t.Parallel()

	// arrange
	h := NewGetAllAccounts(nil)

	// act
	err := h.Handle(context.Background(), query.GetAllAccounts{}, nil)

	// assert
	assert.EqualError(t, err, "invalid report <nil> given")
}

func TestGetAllAccountsHandle_AccountProviderErrorOccurred_ErrorReturned(t *testing.T) {
	t.Parallel()

	// arrange
	expected := fmt.Errorf("an error occurred")
	accountProvider := accountProviderStub{
		ReturnedError: expected,
	}

	h := NewGetAllAccounts(accountProvider)

	// act
	err := h.Handle(context.Background(), query.GetAllAccounts{}, &report.Accounts{})

	// assert
	assert.EqualError(t, err, "cannot retrieve accounts: an error occurred")
}

func TestGetAllAccountsHandle_ValidNumberGiven_AccountReturned(t *testing.T) {
	t.Parallel()

	// arrange
	expected := report.Accounts{
		{
			Number:  "123",
			Balance: *money.New(10000, "USD"),
		},
	}

	accountProvider := accountProviderStub{
		ReturnedAccounts: expected,
	}

	h := NewGetAllAccounts(accountProvider)

	// act
	obtained := report.Accounts{}
	err := h.Handle(context.Background(), query.GetAllAccounts{}, &obtained)
	require.NoError(t, err)

	// assert
	assert.Equal(t, expected, obtained)
}
