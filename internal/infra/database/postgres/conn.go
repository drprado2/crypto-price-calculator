package postgres

import (
	"context"
	"crypto-price-calculator/internal/configs"
	"crypto-price-calculator/internal/observability/applog"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	dbPool *pgxpool.Pool
)

func Setup(ctx context.Context, config *configs.Configuration) (*pgxpool.Pool, error) {
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%v/%s?connect_timeout=10&application_name=crypto-price-calculator&search_path=%s",
		config.DbUser,
		config.DbPass,
		config.DbHost,
		config.DbPort,
		config.DbName,
		config.DbSchema)
	logger := applog.Logger(ctx)
	logger.Info("Connecting postgres DB")
	dbconn, err := pgxpool.Connect(ctx, connectionString)
	if err != nil {
		logger.WithError(err).Error("error creating DB connection")
		return nil, err
	}

	dbPool = dbconn

	logger.Info("DB connected successfully")
	return dbconn, nil
}

func Close() {
	dbPool.Close()
}

func BeginTx(ctx context.Context) (pgx.Tx, error) {
	return dbPool.Begin(ctx)
}

func Commit(ctx context.Context, tx pgx.Tx) error {
	return tx.Commit(ctx)
}

func Rollback(ctx context.Context, tx pgx.Tx) error {
	return tx.Rollback(ctx)
}
