package command

import (
	"github.com/google/uuid"
	"github.com/rhymond/go-money"
)

type TransferMoney struct {
	AggID   uuid.UUID
	AggType string

	From   string
	To     string
	Amount money.Money
}

func (c TransferMoney) AggregateID() uuid.UUID {
	return c.AggID
}

func (c TransferMoney) AggregateType() string {
	return "account.Account"
}
