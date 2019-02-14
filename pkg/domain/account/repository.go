package account

import "github.com/screwyprof/payment/internal/pkg/cqrs"

type GetAccountByNumber interface {
	ByNumber(number Number) (*Account, error)
}

type EventStore interface {
	StoreEvent(event cqrs.Event) error
}
