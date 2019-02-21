package cqrs

import (
	"fmt"
	"reflect"

	"github.com/screwyprof/bank/pkg/domain"
)

type EventSourceHandler struct {
	eventStore domain.EventStore
}

func NewEventSourceHandler(eventStore domain.EventStore) *EventSourceHandler {
	return &EventSourceHandler{
		eventStore:eventStore,
	}
}

func (h *EventSourceHandler) Handle(c domain.Command) error {
	eventStream, err := h.eventStore.LoadEventStream(c.AggregateID())
	if err != nil {
		return err
	}

	agg := CreateAggregate(c.AggregateType(), c.AggregateID())
	if err := applyEvents(agg, eventStream); err != nil {
		return err
	}

	// handle command
	events, err := handle(agg, c)
	if err != nil {
		return err
	}

	if len(events) < 1 {
		return nil
	}

	// store events
	err = h.eventStore.Store(agg.AggregateID(), uint64(len(eventStream)), events)
	if err != nil {
		return err
	}

	// publish events

	return nil
}

func applyEvents(agg domain.Aggregate, eventStream []domain.DomainEvent) error {
	for _, event := range eventStream {
		if err := applyEvent(agg, event); err != nil {
			return err
		}
	}

	return nil
}

func applyEvent(agg domain.Aggregate, event domain.DomainEvent) error {
	eventType := reflect.TypeOf(event)
	aggType := reflect.TypeOf(agg)

	// aggregate's event handler for the given event
	method, ok := aggType.MethodByName("On" + eventType.Name())
	if !ok {
		return fmt.Errorf("event handler for %s is not found", eventType.Name())
	}

	// call handler
	aggValue := reflect.ValueOf(agg)
	eventValue := reflect.ValueOf(event)
	method.Func.Call([]reflect.Value{aggValue, eventValue})

	return nil
}

func handle(agg domain.Aggregate, cmd domain.Command) ([]domain.DomainEvent, error) {
	cmdType := reflect.TypeOf(cmd)
	aggType := reflect.TypeOf(agg)

	// aggregate's command handler for the given cmd
	method, ok := aggType.MethodByName(cmdType.Name())
	if !ok {
		return nil, fmt.Errorf("command handler for %s is not found", cmdType.Name())
	}

	// call handler
	aggValue := reflect.ValueOf(agg)
	cmdValue := reflect.ValueOf(cmd)
	result := method.Func.Call([]reflect.Value{aggValue, cmdValue})

	// checking for errors
	resErr := result[1].Interface()
	if resErr != nil {
		return nil, resErr.(error)
	}

	// extracting events
	eventsIntf := result[0].Interface()
	events := eventsIntf.([]domain.DomainEvent)

	return events, nil
}

