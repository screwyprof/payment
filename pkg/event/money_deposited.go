package event

import "github.com/google/uuid"

type MoneyDeposited struct {
	DomainEvent

	Amount  int64  // public get, private set
	Balance int64  // public get, private set
}

func NewMoneyDeposited(aggID uuid.UUID, amount, balance int64) MoneyDeposited {
	return MoneyDeposited{
		DomainEvent:NewDomainEvent(aggID),
		Amount:amount,
		Balance:balance,
	}
}