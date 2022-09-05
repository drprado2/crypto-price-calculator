package coinbase

import (
	"context"
	"crypto-price-calculator/internal/observability/applog"
	"crypto-price-calculator/internal/observability/apptracer"
	"encoding/json"
	"errors"
	"strings"
)

const (
	ErrorType          = "error"
	SubscriptionReason = "Failed to subscribe"
	InvalidChReason    = "is not a valid channel"
)

type (
	ErrorResponse struct {
		Type    string `json:"type"`
		Message string `json:"message"`
		Reason  string `json:"reason"`
	}

	Error struct {
		expectedProductIds []string
	}
)

var (
	SubscriptionErr = errors.New("invalid subscription message")
)

func NewErrorController() *Error {
	return new(Error)
}

func (r *Error) Handle(ctx context.Context, message []byte) error {
	ctx, span := apptracer.StartOperation(ctx, "Error:Handle", apptracer.SpanKindConsumer)
	defer span.Finish()

	logger := applog.Logger(ctx)

	model := new(ErrorResponse)
	if err := json.Unmarshal(message, model); err != nil {
		logger.WithError(err).Error("error on unmarshal error response")
		span.SetError(err)
		return nil
	}

	logger.
		WithField("Message", model.Message).
		WithField("Reason", model.Reason).
		WithField("Type", ErrorType).
		Infof("Error message received: %s", model.Reason)

	if model.Reason == SubscriptionReason || strings.Contains(model.Reason, InvalidChReason) {
		return SubscriptionErr
	}

	return nil
}
