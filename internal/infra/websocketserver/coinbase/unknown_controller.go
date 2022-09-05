package coinbase

import (
	"context"
	"crypto-price-calculator/internal/observability/applog"
	"crypto-price-calculator/internal/observability/apptracer"
)

const (
	UnknownType = "unknown"
)

type (
	UnknownController struct {
	}
)

func NewUnknownController() *UnknownController {
	return new(UnknownController)
}

func (r *UnknownController) Handle(ctx context.Context, message []byte) error {
	ctx, span := apptracer.StartOperation(ctx, "UnknownController:Handle", apptracer.SpanKindConsumer)
	defer span.Finish()

	logger := applog.Logger(ctx)

	logger.
		WithField("Type", UnknownType).
		Infof("Unknown message received: %v", string(message))

	return nil
}
