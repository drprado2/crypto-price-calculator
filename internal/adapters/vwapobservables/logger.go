package vwapobservables

import (
	"context"
	"crypto-price-calculator/internal/core/vwap"
	"crypto-price-calculator/internal/observability/applog"
)

type (
	Logger struct{}
)

func NewLogger() *Logger {
	return &Logger{}
}

func (l *Logger) HandleNewVwap(ctx context.Context, data *vwap.VwapUpdatedEvent) {
	applog.Logger(ctx).
		WithField("ProductId", data.ProductId).
		WithField("Vwap", data.VolumeWeightedAveragePrice).
		Infof("Vwap updated for product %s, new value:%v", data.ProductId, data.VolumeWeightedAveragePrice)
}
