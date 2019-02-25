package cqrs

import (
	"reflect"

	"github.com/google/uuid"

	"github.com/screwyprof/payment/pkg/domain"
)

type AggregateFactory func(uuid.UUID) domain.Aggregate

var aggregatefactories = make(map[string]AggregateFactory)

func RegisterAggregate(factory AggregateFactory) {
	agg := factory(uuid.New())
	aggregatefactories[reflect.TypeOf(agg).Elem().String()] = factory
}

func CreateAggregate(aggregateType string, id uuid.UUID) domain.Aggregate {
	factory, ok := aggregatefactories[aggregateType]
	if !ok {
		panic(aggregateType + " is not registered")
	}
	return factory(id)
}
