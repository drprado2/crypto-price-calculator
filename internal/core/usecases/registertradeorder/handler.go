package registertradeorder

import (
	"context"
	"crypto-price-calculator/internal/core/entities"
	"crypto-price-calculator/internal/core/repositories"
	"crypto-price-calculator/internal/core/vwap"
	"crypto-price-calculator/internal/observability/apptracer"
	"fmt"
	"strconv"
)

type (
	Handler struct {
		expectedProductIds    []string
		traderOrderRepository repositories.TradeOrderRepository
		tradeChangeCh         chan *vwap.TradeEvent
	}
)

func NewHandler(expectedProductIds []string, traderOrderRepository repositories.TradeOrderRepository, tradeChangeCh chan *vwap.TradeEvent) HandlerInterface {
	return &Handler{
		expectedProductIds:    expectedProductIds,
		traderOrderRepository: traderOrderRepository,
		tradeChangeCh:         tradeChangeCh,
	}
}

func (h *Handler) Exec(ctx context.Context, input *Input) (*entities.TradeOrder, error) {
	ctx, span := apptracer.StartOperation(ctx, "registertradeorder:Exec", apptracer.SpanKindInternal)
	defer span.Finish()

	if err := input.Validate(h.expectedProductIds); err != nil {
		return nil, err
	}

	size, err := strconv.ParseFloat(input.Size, 64)
	if err != nil {
		return nil, &InvalidInputErr{
			Details: map[string]string{"size": fmt.Sprintf("invalid size, received: %s", input.Size)},
		}
	}
	price, err := strconv.ParseFloat(input.Price, 64)
	if err != nil {
		return nil, &InvalidInputErr{
			Details: map[string]string{"price": fmt.Sprintf("invalid price, received: %s", input.Price)},
		}
	}

	entity := entities.NewTradeOrder()
	entity.TradeId = input.TradeId
	entity.MakerOrderId = input.MakerOrderId
	entity.TakerOrderId = input.TakerOrderId
	entity.Side = input.Side
	entity.Size = size
	entity.Price = price
	entity.ProductId = input.ProductId
	entity.Sequence = input.Sequence
	entity.Time = input.Time

	if err := h.traderOrderRepository.Create(ctx, entity); err != nil {
		return nil, err
	}

	h.tradeChangeCh <- &vwap.TradeEvent{
		Price:   entity.Price,
		Size:    entity.Size,
		Product: entity.ProductId,
	}

	return entity, nil
}
