package query_handler

import (
	"context"

	"github.com/screwyprof/payment/internal/pkg/bus"
	"github.com/screwyprof/payment/internal/pkg/cqrs"
)

// ToDomain makes infra.QueryHandler compatible with cqrs.QueryHandler.
func ToDomain(h bus.QueryHandler) cqrs.QueryHandler {
	return domainAdaptor{h: h}
}

type domainAdaptor struct {
	h bus.QueryHandler
}

func (a domainAdaptor) Handle(ctx context.Context, q cqrs.Query, report interface{}) error {
	return a.h.Handle(ctx, q.(bus.Query), report)
}

// FromDomain makes usecase.QueryHandler compatible with infra.QueryHandler.
func FromDomain(h cqrs.QueryHandler) bus.QueryHandler {
	return infraAdapter{h: h}
}

type infraAdapter struct {
	h cqrs.QueryHandler
}

func (a infraAdapter) Handle(ctx context.Context, q bus.Query, report interface{}) error {
	return a.h.Handle(ctx, q.(cqrs.Query), report)
}
