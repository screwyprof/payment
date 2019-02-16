package command_handler

import (
	"context"

	"github.com/screwyprof/payment/internal/pkg/bus"
	"github.com/screwyprof/payment/internal/pkg/cqrs"
)

// ToDomain makes infra.CommandHandler compatible with cqrs.CommandHandler.
func ToDomain(h bus.CommandHandler) cqrs.CommandHandler {
	return domainAdaptor{h: h}
}

type domainAdaptor struct {
	h bus.CommandHandler
}

func (a domainAdaptor) Handle(ctx context.Context, c cqrs.Command) error {
	return a.h.Handle(ctx, c.(bus.Command))
}

// FromDomain makes usecase.CommandHandler compatible with infra.CommandHandler.
func FromDomain(h cqrs.CommandHandler) bus.CommandHandler {
	return infraAdapter{h: h}
}

type infraAdapter struct {
	h cqrs.CommandHandler
}

func (a infraAdapter) Handle(ctx context.Context, request bus.Command) error {
	return a.h.Handle(ctx, request.(cqrs.Command))
}
