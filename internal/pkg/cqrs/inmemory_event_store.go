package cqrs

import (
	"fmt"
	"sync"

	"github.com/google/uuid"

	"github.com/screwyprof/payment/pkg/domain"
)

// InMemoryEventStore implements EventStore as an in memory structure.
type InMemoryEventStore struct {
	eventStreams   map[uuid.UUID][]domain.DomainEvent
	eventStreamsMu sync.RWMutex
}

func NewInMemoryEventStore() *InMemoryEventStore {
	return &InMemoryEventStore{
		eventStreams: make(map[uuid.UUID][]domain.DomainEvent),
	}
}

func (s *InMemoryEventStore) Store(aggregateID uuid.UUID, version uint64, events []domain.DomainEvent) error {
	if len(events) < 1 {
		return fmt.Errorf("no events given")
	}

	eventStream, err := s.LoadEventStream(aggregateID)
	if err != nil {
		return err
	}

	if uint64(len(eventStream)) != version {
		return fmt.Errorf("EventStream has already been modified (concurrently)")
	}

	s.eventStreamsMu.Lock()
	defer s.eventStreamsMu.Unlock()

	eventStream = append(eventStream, events...)
	s.eventStreams[aggregateID] = eventStream

	return nil
}

func (s *InMemoryEventStore) LoadEventStream(aggregateID uuid.UUID) ([]domain.DomainEvent, error) {
	s.eventStreamsMu.RLock()
	defer s.eventStreamsMu.RUnlock()

	eventStream, ok := s.eventStreams[aggregateID]
	if !ok {
		return nil, nil
	}

	return eventStream, nil
}
