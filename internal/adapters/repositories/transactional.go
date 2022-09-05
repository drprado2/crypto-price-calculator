package repositories

import (
	"context"
	"crypto-price-calculator/internal/ctxutils"
	"github.com/jackc/pgx/v4/pgxpool"
)

type (
	Transactional struct {
		dbPool *pgxpool.Pool
	}
)

func (t *Transactional) BeginTx(ctx context.Context) (context.Context, error) {
	tx, err := t.dbPool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	return ctxutils.WithDbTx(ctx, tx), nil
}

func (t *Transactional) Commit(ctx context.Context) error {
	if tx := ctxutils.GetDbTx(ctx); tx != nil {
		return tx.Commit(ctx)
	}

	return nil
}

func (t *Transactional) Rollback(ctx context.Context) {
	if tx := ctxutils.GetDbTx(ctx); tx != nil {
		tx.Rollback(ctx)
	}
}
