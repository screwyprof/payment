package bus

import (
	"context"
	"fmt"
)

// QueryHandlerBus a query handler bus.
type QueryHandlerBus struct {
	queryHandlers map[string]QueryHandler
}

// NewQueryHandlerBus creates a new instance of QueryHandlerBus.
func NewQueryHandlerBus() *QueryHandlerBus {
	return &QueryHandlerBus{
		queryHandlers: make(map[string]QueryHandler),
	}
}

// Handle dispatches the given query to its appropriate QueryHandler and handles it.
func (b *QueryHandlerBus) Handle(ctx context.Context, q Query, report interface{}) error {
	h, err := b.resolve(q.QueryID())
	if err != nil {
		return err
	}
	return h.Handle(ctx, q, report)
}

// Register registers an  QueryHandler for the given QueryID.
func (b *QueryHandlerBus) Register(QueryID string, h QueryHandler) *QueryHandlerBus {
	b.queryHandlers[QueryID] = h
	return b
}

func (b *QueryHandlerBus) resolve(queryID string) (QueryHandler, error) {
	h, ok := b.queryHandlers[queryID]
	if !ok {
		return nil, fmt.Errorf("handler for %q is not found", queryID)
	}
	return h, nil
}
