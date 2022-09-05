package vwapobservables

import (
	"context"
	"crypto-price-calculator/internal/core/vwap"
	"crypto-price-calculator/internal/observability/applog"
)

type (
	PublishSns struct{}
)

func NewPublishSns() *PublishSns {
	return &PublishSns{}
}

func (l *PublishSns) HandleNewVwap(ctx context.Context, data *vwap.VwapUpdatedEvent) {
	// TODO we can publish for a broker here, like SNS, Rabbit, kafka
	applog.Logger(ctx).
		WithField("ProductId", data.ProductId).
		WithField("Vwap", data.VolumeWeightedAveragePrice).
		Infof("Send SNS for product %s, new value:%v", data.ProductId, data.VolumeWeightedAveragePrice)

}
