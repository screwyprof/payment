package event

import (
	"github.com/google/uuid"
	"github.com/rhymond/go-money"
)

type MoneyReceived struct {
	DomainEvent

	From    string
	To      string
	Amount  money.Money
	Balance money.Money
}

func NewMoneyReceived(aggID uuid.UUID, from, to string, amount, balance money.Money) MoneyReceived {
	return MoneyReceived{
		DomainEvent: NewDomainEvent(aggID),

		From:    from,
		To:      to,
		Amount:  amount,
		Balance: balance,
	}
}
