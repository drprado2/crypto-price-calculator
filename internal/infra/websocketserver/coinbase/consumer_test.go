package coinbase

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResponseHandler_getTypeFromMessage(t *testing.T) {
	handler := &Consumer{}

	assert.Equal(t, "error", handler.getTypeFromMessage([]byte(`{"type":"error","message":"Failed to subscribe","reason":"ticket is not a valid channel"}`)))
	assert.Equal(t, "other-type", handler.getTypeFromMessage([]byte(`{"message":"Failed to subscribe","reason":"ticket is not a valid channel","type":"other-type"}`)))
	assert.Equal(t, "subscriptions", handler.getTypeFromMessage([]byte(`{"channels":[{"name":"matches","product_ids":["BTC-USD","ETH-USD","ETH-BTC"]}]},"type":"subscriptions","field": "value"`)))
	assert.Equal(t, "", handler.getTypeFromMessage([]byte(`{"channels":[{"name":"matches","product_ids":["BTC-USD","ETH-USD","ETH-BTC"]}]},"type":subscriptions","field": "value"`)))
	assert.Equal(t, "", handler.getTypeFromMessage([]byte(`{"channels":[{"name":"matches","product_ids":["BTC-USD","ETH-USD","ETH-BTC"]}]},"type":"subscriptions,"field": "value"`)))
}

func TestConsumer_Consume(t *testing.T) {
	ctx := context.Background()

	con1message := `{"type":"con1"}`
	con1CallCount := 0
	controller1 := &ControllerMock{
		MockHandle: func(ctx context.Context, message []byte) error {
			con1CallCount++
			if string(message) != con1message {
				return errors.New("invalid con1 message")
			}
			return nil
		},
	}
	con2message := `{"type":"con2"}`
	con2CallCount := 0
	controller2 := &ControllerMock{
		MockHandle: func(ctx context.Context, message []byte) error {
			con2CallCount++
			if string(message) != con2message {
				return errors.New("invalid con2 message")
			}
			return nil
		},
	}
	conUnkmessage := `{"type":"con3"}`
	conUnkCallCount := 0
	unkController := &ControllerMock{
		MockHandle: func(ctx context.Context, message []byte) error {
			conUnkCallCount++
			if string(message) != conUnkmessage {
				return errors.New("invalid conUnk message")
			}
			return nil
		},
	}

	consumer := NewConsumer(unkController, map[string]Controller{
		"con1": controller1,
		"con2": controller2,
	})

	t.Run("Controller 1 message", func(t *testing.T) {
		assert.NoError(t, consumer.Consume(ctx, []byte(con1message)))
		assert.Equal(t, 1, con1CallCount)
	})

	t.Run("Controller 2 message", func(t *testing.T) {
		assert.NoError(t, consumer.Consume(ctx, []byte(con2message)))
		assert.Equal(t, 1, con2CallCount)
	})

	t.Run("Controller unk message", func(t *testing.T) {
		assert.NoError(t, consumer.Consume(ctx, []byte(conUnkmessage)))
		assert.Equal(t, 1, conUnkCallCount)
	})

	t.Run("Propagate error", func(t *testing.T) {
		controller1.MockHandle = func(ctx context.Context, message []byte) error {
			return errors.New("error")
		}

		assert.Errorf(t, consumer.Consume(ctx, []byte(con1message)), "error")
	})
}
