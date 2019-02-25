package command

import (
	"github.com/google/uuid"
	"github.com/rhymond/go-money"
)

type DepositMoney struct {
	AggID uuid.UUID

	Amount money.Money
}

func (c DepositMoney) AggregateID() uuid.UUID {
	return c.AggID
}

func (c DepositMoney) AggregateType() string {
	return "account.Account"
}
