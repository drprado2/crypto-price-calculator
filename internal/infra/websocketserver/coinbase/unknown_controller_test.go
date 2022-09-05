package coinbase

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnknownController_Handle(t *testing.T) {
	ctx := context.Background()
	controller := NewUnknownController()

	t.Run("invalid json message", func(t *testing.T) {
		assert.NoError(t, controller.Handle(ctx, []byte(`{"type":"error","message":"Failed to subscribe","reason":"ticket is not a valid channel}`)))
	})

	t.Run("valid json message", func(t *testing.T) {
		assert.NoError(t, controller.Handle(ctx, []byte(`{"type":"error","message":"Failed to receive message","reason":"error to receive message"}`)))
	})
}
