package vwap

import (
	"context"
	"crypto-price-calculator/internal/observability/applog"
)

type (
	Consumer struct {
		tradesCh    chan *TradeEvent
		concurrence int
		closeCh     chan interface{}
		closed      bool
		calculator  CalculatorInterface
	}
)

func NewConsumer(calculator CalculatorInterface, tradesCh chan *TradeEvent) *Consumer {
	return &Consumer{
		tradesCh:    tradesCh,
		concurrence: 5,
		closeCh:     make(chan interface{}, 1),
		calculator:  calculator,
	}
}

func (c *Consumer) StartConsumer(ctx context.Context) {
	for i := 0; i < c.concurrence; i++ {
		go c.Consume(ctx)
	}
}

func (c *Consumer) Consume(ctx context.Context) {
	logger := applog.Logger(ctx)
	logger.Info("Starting VWAP consumer")

	for {
		select {
		case why := <-ctx.Done():
			logger.Infof("Closing calculator consumer, due to ctx done %v", why)
			return
		case <-c.closeCh:
			logger.Info("Closing calculator consumer, due to close channel")
			return
		case tradeEvent := <-c.tradesCh:
			logger.Info("Receive trade event on calculator")
			c.calculator.UpdateVwapPrice(ctx, tradeEvent)
		}
	}
}

func (c *Consumer) Close(ctx context.Context) {
	if c.closed {
		return
	}
	applog.Logger(ctx).Info("Closing VWAP consumer")

	c.closed = true
	c.closeCh <- true
}
