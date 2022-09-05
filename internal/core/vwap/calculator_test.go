package vwap

import (
	"context"
	"crypto-price-calculator/internal/configs"
	"crypto-price-calculator/internal/core/entities"
	"crypto-price-calculator/internal/core/repositories"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCalculator_UpdateVwapPrice(t *testing.T) {
	ctx := context.Background()

	btcProduct := "BTC-USD"
	ethProduct := "ETH-USD"
	config := &configs.Configuration{
		VwapWindowSize: 5,
		ProductIds:     fmt.Sprintf("%s,%s", btcProduct, ethProduct),
	}
	tradeOrderRepository := &repositories.TraderOrderRepositoryMock{
		MockRetrieveLastNTradesByProduct: func(ctx context.Context, productId string, limit int) ([]*entities.TradeOrder, error) {
			if productId == btcProduct && limit == 5 {
				return []*entities.TradeOrder{
					{Size: 0.0024722, Price: 19221.22},
					{Size: 0.00181211, Price: 19221.09},
					{Size: 0.00086249, Price: 19221.03},
				}, nil
			}
			if productId == ethProduct && limit == 5 {
				return []*entities.TradeOrder{
					{Size: 0.31543783, Price: 1624.56},
					{Size: 0.001, Price: 1624.57},
					{Size: 0.30776996, Price: 1624.59},
					{Size: 0.00001052, Price: 1624.59},
					{Size: 0.01578169, Price: 1624.6},
				}, nil
			}

			return nil, errors.New("error filters")
		},
	}
	obs1EventsReceived := make([]*VwapUpdatedEvent, 0)
	observable1 := &VwapObservableMock{
		MockHandleNewVwap: func(ctx context.Context, event *VwapUpdatedEvent) {
			obs1EventsReceived = append(obs1EventsReceived, event)
		},
	}
	obs2EventsReceived := make([]*VwapUpdatedEvent, 0)
	observable2 := &VwapObservableMock{
		MockHandleNewVwap: func(ctx context.Context, event *VwapUpdatedEvent) {
			obs2EventsReceived = append(obs2EventsReceived, event)
		},
	}

	calculator := NewCalculator(config, tradeOrderRepository, observable1, observable2)

	calculator.Setup(ctx)

	calculator.UpdateVwapPrice(ctx, &TradeEvent{
		Price:   1624.59,
		Size:    0.30775944,
		Product: ethProduct,
	})
	calculator.UpdateVwapPrice(ctx, &TradeEvent{
		Price:   19221.01,
		Size:    0.00024102,
		Product: btcProduct,
	})
	calculator.UpdateVwapPrice(ctx, &TradeEvent{
		Price:   1624.61,
		Size:    0.00000056,
		Product: ethProduct,
	})
	calculator.UpdateVwapPrice(ctx, &TradeEvent{
		Price:   19220.88,
		Size:    0.02678371,
		Product: btcProduct,
	})
	calculator.UpdateVwapPrice(ctx, &TradeEvent{
		Price:   19220.87,
		Size:    0.0001,
		Product: btcProduct,
	})

	assert.Len(t, obs1EventsReceived, 5)
	assert.Len(t, obs2EventsReceived, 5)
	assert.Equal(t, obs1EventsReceived[0].ProductId, ethProduct)
	assert.Equal(t, float64(1624.58), obs1EventsReceived[0].VolumeWeightedAveragePrice)
	assert.Equal(t, obs1EventsReceived[1].ProductId, btcProduct)
	assert.Equal(t, float64(19221.14), obs1EventsReceived[1].VolumeWeightedAveragePrice)
	assert.Equal(t, obs1EventsReceived[2].ProductId, ethProduct)
	assert.Equal(t, float64(1624.58), obs1EventsReceived[2].VolumeWeightedAveragePrice)
	assert.Equal(t, obs1EventsReceived[3].ProductId, btcProduct)
	assert.Equal(t, float64(19220.92), obs1EventsReceived[3].VolumeWeightedAveragePrice)
	assert.Equal(t, obs1EventsReceived[4].ProductId, btcProduct)
	assert.Equal(t, float64(19220.92), obs1EventsReceived[4].VolumeWeightedAveragePrice)
}
