package coinbase

import (
	"context"
	"crypto-price-calculator/internal/core/repositories"
	"crypto-price-calculator/internal/core/usecases/registertradeorder"
	"crypto-price-calculator/internal/observability/applog"
	"crypto-price-calculator/internal/observability/apptracer"
	"encoding/json"
	"errors"
	"time"
)

const (
	MatchType = "match"
)

type (
	MatchEvent struct {
		Type         string    `json:"type"`
		TradeId      int       `json:"trade_id"`
		MakerOrderId string    `json:"maker_order_id"`
		TakerOrderId string    `json:"taker_order_id"`
		Side         string    `json:"side"`
		Size         string    `json:"size"`
		Price        string    `json:"price"`
		ProductId    string    `json:"product_id"`
		Sequence     int64     `json:"sequence"`
		Time         time.Time `json:"time"`
	}

	MatchController struct {
		useCase registertradeorder.HandlerInterface
	}
)

func NewMatchController(useCase registertradeorder.HandlerInterface) *MatchController {
	return &MatchController{
		useCase: useCase,
	}
}

func (r *MatchController) Handle(ctx context.Context, message []byte) error {
	ctx, span := apptracer.StartOperation(ctx, "MatchController:Handle", apptracer.SpanKindConsumer)
	defer span.Finish()

	logger := applog.Logger(ctx)

	input := new(MatchEvent)
	if err := json.Unmarshal(message, input); err != nil {
		logger.WithError(err).Error("error on unmarshal match event")
		span.SetError(err)
		return nil
	}

	_, err := r.useCase.Exec(ctx, input.ToUseCaseInput())
	if err != nil {
		if err, ok := err.(*registertradeorder.InvalidInputErr); ok {
			logger.WithError(err).Warn("invalid match event")
			return nil
		}
		if errors.Is(err, repositories.RegisterAlreadyExists) {
			logger.WithError(err).Warn("trade order already exists")
			return nil
		}
		logger.WithError(err).Error("internal error happened")
		return nil
	}

	logger.
		WithField("ProductId", input.ProductId).
		WithField("Type", MatchType).
		Info("match event successfully consumed")

	return nil
}

func (r *MatchEvent) ToUseCaseInput() *registertradeorder.Input {
	return &registertradeorder.Input{
		TradeId:      r.TradeId,
		MakerOrderId: r.MakerOrderId,
		TakerOrderId: r.TakerOrderId,
		Side:         r.Side,
		Size:         r.Size,
		Price:        r.Price,
		ProductId:    r.ProductId,
		Sequence:     r.Sequence,
		Time:         r.Time,
	}
}
