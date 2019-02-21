package command

import (
	"github.com/google/uuid"
)

type DepositMoney struct {
	AggID uuid.UUID
	AggType string

	Amount  int64
}

func (c DepositMoney) AggregateID() uuid.UUID {
	return c.AggID
}

func (c DepositMoney) AggregateType() string{
	return "account.Account"
}
