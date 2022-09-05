package vwap

import (
	"context"
	"github.com/stretchr/testify/assert"
	"sync/atomic"
	"testing"
	"time"
)

func TestConsumer(t *testing.T) {
	ctx := context.Background()

	product := "prod"
	var calculateCount int32 = 0
	calculator := &CalculatorMock{
		MockUpdateVwapPrice: func(ctx context.Context, tradeEvent *TradeEvent) {
			if tradeEvent.Product == product {
				atomic.AddInt32(&calculateCount, 1)
			}
		},
	}
	tradesCh := make(chan *TradeEvent, 1)

	consumer := NewConsumer(calculator, tradesCh)
	consumer.concurrence = 1

	consumer.StartConsumer(ctx)

	tradesCh <- &TradeEvent{
		Product: product,
	}
	time.Sleep(time.Millisecond * 50)
	assert.Equal(t, int32(1), atomic.LoadInt32(&calculateCount))

	consumer.Close(ctx)

	time.Sleep(time.Millisecond * 50)

	tradesCh <- &TradeEvent{
		Product: product,
	}
	time.Sleep(time.Millisecond * 50)
	assert.Equal(t, int32(1), atomic.LoadInt32(&calculateCount))
}
