package event

import "github.com/google/uuid"

type DomainEvent struct {
	ID    uuid.UUID
	AggID uuid.UUID
}

func NewDomainEvent(aggID uuid.UUID) DomainEvent {
	return DomainEvent{
		ID:    uuid.New(),
		AggID: aggID,
	}
}

func (e DomainEvent) EventID() uuid.UUID {
	return e.ID
}

func (e DomainEvent) AggregateID() uuid.UUID {
	return e.AggID
}
