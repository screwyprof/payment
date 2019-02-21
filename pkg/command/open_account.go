package command

import (
	"github.com/google/uuid"
	"github.com/rhymond/go-money"
)

type OpenAccount struct {
	AggID   uuid.UUID
	AggType string

	Balance money.Money
	Number  string
}

func (c OpenAccount) AggregateID() uuid.UUID {
	return c.AggID
}

func (c OpenAccount) AggregateType() string {
	return "account.Account"
}
