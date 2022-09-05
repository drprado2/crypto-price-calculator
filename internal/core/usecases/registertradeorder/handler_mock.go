package registertradeorder

import (
	"context"
	"crypto-price-calculator/internal/core/entities"
)

type (
	HandlerMock struct {
		MockExec func(ctx context.Context, input *Input) (*entities.TradeOrder, error)
	}
)

func (h *HandlerMock) Exec(ctx context.Context, input *Input) (*entities.TradeOrder, error) {
	if h.MockExec != nil {
		return h.MockExec(ctx, input)
	}

	return new(entities.TradeOrder), nil
}
