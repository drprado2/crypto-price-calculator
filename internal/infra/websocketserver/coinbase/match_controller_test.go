package coinbase

import (
	"context"
	"crypto-price-calculator/internal/core/entities"
	"crypto-price-calculator/internal/core/usecases/registertradeorder"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMatchController_Handle(t *testing.T) {
	ctx := context.Background()
	useCase := &registertradeorder.HandlerMock{}
	controller := NewMatchController(useCase)

	t.Run("invalid json message", func(t *testing.T) {
		assert.NoError(t, controller.Handle(ctx, []byte(`{"type":"match","trade_id":407115390,"maker_order_id":"df315bb9-c940-4f55-be69-ea165f987e35","taker_order_id":"6d5be16b-2936-49ad-883f-78860c94891c","side":"buy","size":"0.580621","price":"18924.09","product_id":"BTC-USD","sequence":44805966539,"time":"2022-09-07T14:38:51.473515Z}`)))
	})

	t.Run("Use case input error", func(t *testing.T) {
		useCase.MockExec = func(ctx context.Context, input *registertradeorder.Input) (*entities.TradeOrder, error) {
			return nil, &registertradeorder.InvalidInputErr{
				Details: map[string]string{
					"price":   "invalid price",
					"product": "invalid product",
				},
			}
		}

		assert.NoError(t, controller.Handle(ctx, []byte(`{"type":"match","trade_id":407115390,"maker_order_id":"df315bb9-c940-4f55-be69-ea165f987e35","taker_order_id":"6d5be16b-2936-49ad-883f-78860c94891c","side":"buy","size":"0.580621","price":"18924.09","product_id":"BTC-USD","sequence":44805966539,"time":"2022-09-07T14:38:51.473515Z"}`)))
	})

	t.Run("Use case internal error", func(t *testing.T) {
		useCase.MockExec = func(ctx context.Context, input *registertradeorder.Input) (*entities.TradeOrder, error) {
			return nil, errors.New("internal error")
		}

		assert.NoError(t, controller.Handle(ctx, []byte(`{"type":"match","trade_id":407115390,"maker_order_id":"df315bb9-c940-4f55-be69-ea165f987e35","taker_order_id":"6d5be16b-2936-49ad-883f-78860c94891c","side":"buy","size":"0.580621","price":"18924.09","product_id":"BTC-USD","sequence":44805966539,"time":"2022-09-07T14:38:51.473515Z"}`)))
	})

	t.Run("Success case", func(t *testing.T) {
		useCase.MockExec = func(ctx context.Context, input *registertradeorder.Input) (*entities.TradeOrder, error) {
			return &entities.TradeOrder{
				ProductId: "BTC-USD",
			}, nil
		}

		assert.NoError(t, controller.Handle(ctx, []byte(`{"type":"match","trade_id":407115390,"maker_order_id":"df315bb9-c940-4f55-be69-ea165f987e35","taker_order_id":"6d5be16b-2936-49ad-883f-78860c94891c","side":"buy","size":"0.580621","price":"18924.09","product_id":"BTC-USD","sequence":44805966539,"time":"2022-09-07T14:38:51.473515Z"}`)))
	})
}
