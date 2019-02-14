package report

import (
	"testing"

	"github.com/rhymond/go-money"
	"github.com/stretchr/testify/assert"
)

func TestAccountToString_MethodCalled_ValidStringReturned(t *testing.T) {
	t.Parallel()

	// arrange
	acc := &Account{
		Number:  "123",
		Balance: *money.New(10000, "USD"),
	}

	expected := "#123: $100.00"

	// act
	obtained := acc.ToString()

	// assert
	assert.Equal(t, expected, obtained)
}
