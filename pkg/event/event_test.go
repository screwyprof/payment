package event

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccountOpened_EventID(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "AccountOpened", AccountOpened{}.EventID())
}

func TestMoneyDeposited_EventID(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "MoneyDeposited", MoneyDeposited{}.EventID())
}

func TestMoneyTransferred_EventID(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "MoneyTransferred", MoneyTransferred{}.EventID())
}

func TestMoneyReceived_EventID(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "MoneyReceived", MoneyReceived{}.EventID())
}
