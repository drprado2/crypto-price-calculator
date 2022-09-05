package vwap

import "context"

type (
	CalculatorMock struct {
		MockUpdateVwapPrice func(ctx context.Context, tradeEvent *TradeEvent)
		MockSetup           func(ctx context.Context) error
	}
)

func (c *CalculatorMock) UpdateVwapPrice(ctx context.Context, tradeEvent *TradeEvent) {
	if c.MockUpdateVwapPrice != nil {
		c.MockUpdateVwapPrice(ctx, tradeEvent)
	}
}

func (c *CalculatorMock) Setup(ctx context.Context) error {
	if c.MockSetup != nil {
		return c.MockSetup(ctx)
	}

	return nil
}
