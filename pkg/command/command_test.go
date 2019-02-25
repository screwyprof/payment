package command

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpenAccount_AggregateType(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "account.Account", OpenAccount{}.AggregateType())
}

func TestTransferMoney_AggregateType(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "account.Account", TransferMoney{}.AggregateType())
}

func TestReceiveMoney_AggregateType(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "account.Account", ReceiveMoney{}.AggregateType())
}

func TestUnknown_AggregateType(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "unknown.Unknown", Unknown{}.AggregateType())
}
