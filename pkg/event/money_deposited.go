package event

import (
	"github.com/google/uuid"
	"github.com/rhymond/go-money"
)

type MoneyDeposited struct {
	DomainEvent

	Amount  money.Money
	Balance money.Money
}

func NewMoneyDeposited(aggID uuid.UUID, amount, balance money.Money) MoneyDeposited {
	return MoneyDeposited{
		DomainEvent: NewDomainEvent(aggID),
		Amount:      amount,
		Balance:     balance,
	}
}
