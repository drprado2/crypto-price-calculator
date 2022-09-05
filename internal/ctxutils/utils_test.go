package ctxutils

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCid(t *testing.T) {
	ctx := context.Background()
	cid := "cid"
	ctx = WithCid(ctx, cid)
	assert.Equal(t, cid, GetCid(ctx))
}

func TestSpanId(t *testing.T) {
	ctx := context.Background()
	spanId := "spanId"
	ctx = WithSpanId(ctx, spanId)
	assert.Equal(t, spanId, GetSpanId(ctx))
}

func TestTraceId(t *testing.T) {
	ctx := context.Background()
	traceId := "traceId"
	ctx = WithTraceId(ctx, traceId)
	assert.Equal(t, traceId, GetTraceId(ctx))
}
