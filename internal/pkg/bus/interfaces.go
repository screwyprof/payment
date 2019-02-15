package bus

import "context"

// Command defines command parameters.
type Command interface {
	CommandID() string
}

// CommandHandler executes a use case.
type CommandHandler interface {
	// Handle handles the given command and returns.
	Handle(ctx context.Context, c Command) error
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

type Event interface {
	EventID() string
}
