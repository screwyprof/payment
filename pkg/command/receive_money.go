package command

import (
	"github.com/google/uuid"
	"github.com/rhymond/go-money"
)

type ReceiveMoney struct {
	AggID   uuid.UUID
	AggType string

	From   string
	To     string
	Amount money.Money
}

func (c ReceiveMoney) AggregateID() uuid.UUID {
	return c.AggID
}

func (c ReceiveMoney) AggregateType() string {
	return "account.Account"
}
