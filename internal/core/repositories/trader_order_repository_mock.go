package repositories

import (
	"context"
	"crypto-price-calculator/internal/core/entities"
)

type (
	TraderOrderRepositoryMock struct {
		MockCreate                       func(ctx context.Context, trade *entities.TradeOrder) error
		MockRetrieveLastNTradesByProduct func(ctx context.Context, productId string, limit int) ([]*entities.TradeOrder, error)
		MockBeginTx                      func(ctx context.Context) (context.Context, error)
		MockCommit                       func(ctx context.Context) error
		MockRollback                     func(ctx context.Context)
	}
)

func (t *TraderOrderRepositoryMock) Create(ctx context.Context, trade *entities.TradeOrder) error {
	if t.MockCreate != nil {
		return t.MockCreate(ctx, trade)
	}

	return nil
}

func (t *TraderOrderRepositoryMock) RetrieveLastNTradesByProduct(ctx context.Context, productId string, limit int) ([]*entities.TradeOrder, error) {
	if t.MockRetrieveLastNTradesByProduct != nil {
		return t.MockRetrieveLastNTradesByProduct(ctx, productId, limit)
	}

	return make([]*entities.TradeOrder, 0), nil

}

func (t *TraderOrderRepositoryMock) BeginTx(ctx context.Context) (context.Context, error) {
	if t.MockBeginTx != nil {
		return t.MockBeginTx(ctx)
	}

	return context.Background(), nil
}

func (t *TraderOrderRepositoryMock) Commit(ctx context.Context) error {
	if t.MockCommit != nil {
		return t.MockCommit(ctx)
	}

	return nil
}

func (t *TraderOrderRepositoryMock) Rollback(ctx context.Context) {
	if t.MockCreate != nil {
		t.Rollback(ctx)
	}
}
