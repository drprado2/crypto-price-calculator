package coinbase

import (
	"context"
	"crypto-price-calculator/internal/configs"
	"crypto-price-calculator/internal/observability/applog"
	"crypto-price-calculator/internal/observability/apptracer"
	"encoding/json"
	"errors"
	"golang.org/x/exp/slices"
)

const (
	SubscriptionsType = "subscriptions"
)

type (
	SubscriptionResponse struct {
		Type     string    `json:"type"`
		Channels []Channel `json:"channels"`
	}

	SubscriptionsController struct {
		expectedProductIds []string
		matchesCh          string
	}
)

var (
	InvalidSubscriptionErr = errors.New("invalid subscription")
)

func NewSubscriptionsController(config *configs.Configuration) *SubscriptionsController {
	return &SubscriptionsController{
		expectedProductIds: config.GetProductIds(),
		matchesCh:          config.MatchesChannel,
	}
}

func (r *SubscriptionsController) Handle(ctx context.Context, message []byte) error {
	ctx, span := apptracer.StartOperation(ctx, "SubscriptionsController:Handle", apptracer.SpanKindConsumer)
	defer span.Finish()

	logger := applog.Logger(ctx)

	model := new(SubscriptionResponse)
	if err := json.Unmarshal(message, model); err != nil {
		logger.WithError(err).Error("error on unmarshal subscription response")
		span.SetError(err)
		return err
	}

	if len(model.Channels) == 0 {
		logger.
			WithError(InvalidSubscriptionErr).
			Error("Subscription message received empty channels")
		span.SetError(InvalidSubscriptionErr)

		return InvalidSubscriptionErr
	}

	for _, channel := range model.Channels {
		if channel.Name != r.matchesCh || !r.validProducts(channel.ProductIds) {
			logger.
				WithError(InvalidSubscriptionErr).
				WithField("ChannelReceived", channel.Name).
				WithField("ProductsReceived", channel.ProductIds).
				Error("Subscription message doesn't contains all expected channels and products")
			span.SetError(InvalidSubscriptionErr)

			return InvalidSubscriptionErr
		}
	}

	logger.
		WithField("Message", message).
		WithField("Type", SubscriptionsType).
		Info("subscription message successfully consumed")

	return nil
}

func (r *SubscriptionsController) validProducts(procutsReceived []string) bool {
	if len(r.expectedProductIds) < len(procutsReceived) {
		return false
	}

	for _, p := range procutsReceived {
		if !slices.Contains(r.expectedProductIds, p) {
			return false
		}
	}

	return true
}
