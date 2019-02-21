package event

import (
	"github.com/google/uuid"
	"github.com/rhymond/go-money"
)

type AccountOpened struct {
	DomainEvent

	Number  string
	Balance money.Money
}

func NewAccountOpened(aggID uuid.UUID, number string, balance money.Money) AccountOpened {
	return AccountOpened{
		DomainEvent: NewDomainEvent(aggID),
		Number:      number,
		Balance:     balance,
	}
}
