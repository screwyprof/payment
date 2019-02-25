package event

import (
	"github.com/google/uuid"
	"github.com/rhymond/go-money"
)

type MoneyTransferred struct {
	DomainEvent

	From    string
	To      string
	Amount  money.Money
	Balance money.Money
}

func NewMoneyTransferred(aggID uuid.UUID, from, to string, amount, balance money.Money) MoneyTransferred {
	return MoneyTransferred{
		DomainEvent: NewDomainEvent(aggID),

		From:    from,
		To:      to,
		Amount:  amount,
		Balance: balance,
	}
}
