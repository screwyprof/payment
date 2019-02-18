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

func TestGetAccountShortInfoHandle_InvalidQueryGiven_ErrorReturned(t *testing.T) {
	t.Parallel()

	// arrange
	h := NewGetAccountShortInfo(nil)

	// act
	err := h.Handle(context.Background(), query.Unknown{}, nil)

	// assert
	assert.EqualError(t, err, "invalid query query.Unknown{} given")
}

func TestGetAccountShortInfoHandle_InvalidResultGiven_ErrorReturned(t *testing.T) {
	t.Parallel()

	// arrange
	h := NewGetAccountShortInfo(nil)

	// act
	err := h.Handle(context.Background(), query.GetAccountShortInfo{Number: "123"}, nil)

	// assert
	assert.EqualError(t, err, "invalid report <nil> given")
}

func TestGetAccountShortInfoHandle_AccountProviderErrorOccurred_ErrorReturned(t *testing.T) {
	t.Parallel()

	// arrange
	expected := fmt.Errorf("an error occurred")
	accountProvider := accountProviderStub{
		ReturnedError: expected,
	}

	h := NewGetAccountShortInfo(accountProvider)

	// act
	err := h.Handle(context.Background(), query.GetAccountShortInfo{Number: "123"}, &report.Account{})

	// assert
	assert.EqualError(t, err, "cannot retrieve account: an error occurred")
}

func TestGetAccountShortInfoHandle_ValidNumberGiven_AccountReturned(t *testing.T) {
	t.Parallel()

	// arrange
	expected := &report.Account{
		Number:  "123",
		Balance: *money.New(10000, "USD"),
	}

	accountProvider := accountProviderStub{
		ReturnedAccount: expected,
	}

	h := NewGetAccountShortInfo(accountProvider)

	// act
	obtained := &report.Account{}
	err := h.Handle(context.Background(), query.GetAccountShortInfo{Number: "123"}, obtained)
	require.NoError(t, err)

	// assert
	assert.Equal(t, expected, obtained)
}
