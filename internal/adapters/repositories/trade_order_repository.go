package repositories

import (
	"context"
	"crypto-price-calculator/internal/core/entities"
	"crypto-price-calculator/internal/core/repositories"
	"crypto-price-calculator/internal/observability/applog"
	"crypto-price-calculator/internal/observability/apptracer"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
)

type (
	TradeOrderRepository struct {
		Transactional
	}
)

func NewTradeOrderRepository(dbPool *pgxpool.Pool) *TradeOrderRepository {
	return &TradeOrderRepository{
		Transactional{dbPool: dbPool},
	}
}

func (t *TradeOrderRepository) Create(ctx context.Context, trade *entities.TradeOrder) error {
	ctx, span := apptracer.StartOperation(ctx, "TradeOrderRepository:Create", apptracer.SpanKindInternal)
	defer span.Finish()

	logger := applog.Logger(ctx)

	tag, err := t.dbPool.Exec(ctx, insertTradeOrderQuery,
		trade.TradeOrderId,
		trade.TradeId,
		trade.MakerOrderId,
		trade.TakerOrderId,
		trade.Side,
		trade.Size,
		trade.Price,
		trade.ProductId,
		trade.Sequence,
		trade.Time,
	)
	if err != nil {
		if err, ok := err.(*pgconn.PgError); ok {
			if err.Code == "23505" {
				logger.
					WithError(err).
					WithField("ProductId", trade.ProductId).
					WithField("TradeId", trade.TradeId).
					Warn("trade order already exists on storage")
				return repositories.RegisterAlreadyExists
			}
		}

		logger.
			WithError(err).
			WithField("ProductId", trade.ProductId).
			WithField("TradeId", trade.TradeId).
			Error("unexpected error creating trade order")
		return err
	}

	if tag.RowsAffected() != 1 {
		return repositories.NoChangesHappenedError
	}

	return nil
}

func (t *TradeOrderRepository) RetrieveLastNTradesByProduct(ctx context.Context, productId string, limit int) ([]*entities.TradeOrder, error) {
	ctx, span := apptracer.StartOperation(ctx, "TradeOrderRepository:Create", apptracer.SpanKindInternal)
	defer span.Finish()

	logger := applog.Logger(ctx)

	rows, err := t.dbPool.Query(ctx, getLastNTradesByProductQuery, productId, limit)
	if err != nil {
		logger.WithError(err).Error("error executing query")
		return nil, err
	}
	defer rows.Close()

	result := make([]*entities.TradeOrder, 0)
	for rows.Next() {
		trade := new(entities.TradeOrder)
		if err := rows.Scan(&trade.Size, &trade.Price); err != nil {
			return nil, err
		}

		result = append(result, trade)
	}

	return result, nil
}
