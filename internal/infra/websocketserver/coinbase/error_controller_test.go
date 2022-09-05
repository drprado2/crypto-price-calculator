package coinbase

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestError_Handle(t *testing.T) {
	ctx := context.Background()
	controller := NewErrorController()

	t.Run("invalid json message", func(t *testing.T) {
		assert.NoError(t, controller.Handle(ctx, []byte(`{"type":"error","message":"Failed to subscribe","reason":"ticket is not a valid channel}`)))
	})

	t.Run("Acceptable error", func(t *testing.T) {
		assert.NoError(t, controller.Handle(ctx, []byte(`{"type":"error","message":"Failed to receive message","reason":"error to receive message"}`)))
	})

	t.Run("Subscription error", func(t *testing.T) {
		assert.ErrorIs(t, controller.Handle(ctx, []byte(`{"type":"error","message":"Failed to receive message","reason":"Failed to subscribe"}`)), SubscriptionErr)
	})

	t.Run("Channel error", func(t *testing.T) {
		assert.ErrorIs(t, controller.Handle(ctx, []byte(`{"type":"error","message":"Failed to receive message","reason":"Foo is not a valid channel"}`)), SubscriptionErr)
	})
}
