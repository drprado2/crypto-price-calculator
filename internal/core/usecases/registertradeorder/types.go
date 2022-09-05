package registertradeorder

import (
	"context"
	"crypto-price-calculator/internal/core/entities"
	"fmt"
	"golang.org/x/exp/slices"
	"time"
)

type (
	Input struct {
		TradeId      int
		MakerOrderId string
		TakerOrderId string
		Side         string
		Size         string
		Price        string
		ProductId    string
		Sequence     int64
		Time         time.Time
	}

	HandlerInterface interface {
		Exec(ctx context.Context, input *Input) (*entities.TradeOrder, error)
	}
)

func (i *Input) Validate(expectedProductIds []string) error {
	err := make(map[string]string)

	if i.TradeId < 1 {
		err["trade_id"] = "invalid trade ID"
	}
	if i.MakerOrderId == "" {
		err["marker_order_id"] = "invalid marker order ID"
	}
	if i.TakerOrderId == "" {
		err["taker_order_id"] = "invalid taker order ID"
	}
	if i.Side != "sell" && i.Side != "buy" {
		err["side"] = fmt.Sprintf("invalid side, received: %s", i.Side)
	}
	if i.Size == "" {
		err["size"] = fmt.Sprintf("invalid size, received: %s", i.Size)
	}
	if i.Price == "" {
		err["price"] = fmt.Sprintf("invalid price, received: %s", i.Price)
	}
	if i.ProductId == "" || !slices.Contains(expectedProductIds, i.ProductId) {
		err["product_id"] = fmt.Sprintf("invalid product ID, received: %s", i.ProductId)
	}
	if i.Sequence < 1 {
		err["sequence"] = "invalid sequence"
	}

	if len(err) > 0 {
		return &InvalidInputErr{
			Details: err,
		}
	}

	return nil
}
