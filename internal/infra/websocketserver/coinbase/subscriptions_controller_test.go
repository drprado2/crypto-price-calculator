package coinbase

import (
	"context"
	"crypto-price-calculator/internal/configs"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSubscriptionsController_Handle(t *testing.T) {
	ctx := context.Background()
	config := &configs.Configuration{
		MatchesChannel: "matches",
		ProductIds:     "BTC-USD,ETH-USD,ETH-BTC",
	}
	controller := NewSubscriptionsController(config)

	t.Run("invalid json message", func(t *testing.T) {
		assert.Error(t, controller.Handle(ctx, []byte(`{"type":"error","message":"Failed to subscribe","reason":"ticket is not a valid channel}`)))
	})

	t.Run("Empty channels", func(t *testing.T) {
		assert.ErrorIs(t, controller.Handle(ctx, []byte(`{"type":"subscriptions","channels":[]}`)), InvalidSubscriptionErr)
	})

	t.Run("Wrong channel", func(t *testing.T) {
		assert.ErrorIs(t, controller.Handle(ctx, []byte(`{"type":"subscriptions","channels":[{"name":"invalid","product_ids":["BTC-USD","ETH-USD","ETH-BTC"]}]}`)), InvalidSubscriptionErr)
	})

	t.Run("Invalid products", func(t *testing.T) {
		assert.ErrorIs(t, controller.Handle(ctx, []byte(`{"type":"subscriptions","channels":[{"name":"matches","product_ids":["BTC-BRL","ETH-USD","ETH-BTC"]}]}`)), InvalidSubscriptionErr)
	})

	t.Run("Successfully message", func(t *testing.T) {
		assert.NoError(t, controller.Handle(ctx, []byte(`{"type":"subscriptions","channels":[{"name":"matches","product_ids":["BTC-USD","ETH-USD","ETH-BTC"]}]}`)))
	})
}
