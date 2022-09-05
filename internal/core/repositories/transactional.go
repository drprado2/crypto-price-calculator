package repositories

import (
	"context"
)

type (
	Transactional interface {
		BeginTx(ctx context.Context) (context.Context, error)
		Commit(ctx context.Context) error
		Rollback(ctx context.Context)
	}
)
