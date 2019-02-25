package query_handler

import (
	"context"

	"github.com/screwyprof/payment/internal/pkg/bus"

	"github.com/screwyprof/payment/pkg/domain"
)

// ToDomain makes infra.QueryHandler compatible with domain.QueryHandler.
func ToDomain(h bus.QueryHandler) domain.QueryHandler {
	return domainAdaptor{h: h}
}

type domainAdaptor struct {
	h bus.QueryHandler
}

func (a domainAdaptor) Handle(ctx context.Context, q domain.Query, report interface{}) error {
	return a.h.Handle(ctx, q.(bus.Query), report)
}

// FromDomain makes usecase.QueryHandler compatible with infra.QueryHandler.
func FromDomain(h domain.QueryHandler) bus.QueryHandler {
	return infraAdapter{h: h}
}

type infraAdapter struct {
	h domain.QueryHandler
}

func (a infraAdapter) Handle(ctx context.Context, q bus.Query, report interface{}) error {
	return a.h.Handle(ctx, q.(domain.Query), report)
}
