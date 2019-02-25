package domain

import (
	"context"

	"github.com/google/uuid"
)

type Command interface {
	AggregateID() uuid.UUID
	AggregateType() string
}

type CommandHandler interface {
	Handle(ctx context.Context, c Command) error
}

type DomainEvent interface {
	EventID() uuid.UUID
	AggregateID() uuid.UUID
}

type Aggregate interface {
	AggregateID() uuid.UUID
}

type EventStore interface {
	LoadEventStream(aggregateID uuid.UUID) ([]DomainEvent, error)
	Store(aggregateID uuid.UUID, version uint64, events []DomainEvent) error
}

// Query defines query parameters.
type Query interface {
	QueryID() string
}

// QueryHandler handles a query.
type QueryHandler interface {
	// Handle handles the given query and returns the report.
	Handle(ctx context.Context, q Query, report interface{}) error
}
