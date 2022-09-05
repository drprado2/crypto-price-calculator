package ctxutils

import (
	"context"
	"github.com/jackc/pgx/v4"
)

const (
	cidKey           = "x-cid"
	customHeadersKey = "x-custom-headers"
	dbTxKey          = "x-db-tx"

	spanIdKey  = "x-span-id"
	traceIdKey = "x-trace-id"
)

func GetCid(ctx context.Context) string {
	return getString(ctx, cidKey)
}

func GetSpanId(ctx context.Context) string {
	return getString(ctx, spanIdKey)
}

func GetTraceId(ctx context.Context) string {
	return getString(ctx, traceIdKey)
}

func GetDbTx(ctx context.Context) pgx.Tx {
	if v := ctx.Value(dbTxKey); v != nil {
		return v.(pgx.Tx)
	}

	return nil
}

func WithCid(ctx context.Context, cid string) context.Context {
	return context.WithValue(ctx, cidKey, cid)
}

func WithSpanId(ctx context.Context, spanId string) context.Context {
	return context.WithValue(ctx, spanIdKey, spanId)
}

func WithDbTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, dbTxKey, tx)
}

func WithTraceId(ctx context.Context, traceId string) context.Context {
	return context.WithValue(ctx, traceIdKey, traceId)
}

func getString(ctx context.Context, param string) string {
	if v := ctx.Value(param); v != nil {
		return v.(string)
	}

	return ""
}
