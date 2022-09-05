package repositories

import (
	"context"
	"crypto-price-calculator/internal/core/entities"
)

type (
	TradeOrderRepository interface {
		Create(ctx context.Context, trade *entities.TradeOrder) error
		RetrieveLastNTradesByProduct(ctx context.Context, productId string, limit int) ([]*entities.TradeOrder, error)
		Transactional
	}
)
