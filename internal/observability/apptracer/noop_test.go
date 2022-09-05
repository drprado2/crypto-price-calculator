package apptracer

import (
	"context"
	"crypto-price-calculator/internal/ctxutils"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStartNoop(t *testing.T) {
	ctx := context.Background()

	ctx, span := Start(ctx)

	assert.NotEmpty(t, span.TraceID())
	assert.NotEqual(t, span.TraceID(), ctxutils.GetTraceId(ctx))
	assert.NotEmpty(t, span.SpanID())
	assert.NotEqual(t, span.SpanID(), ctxutils.GetSpanId(ctx))

	assert.NotEmpty(t, span.Name())
	assert.NotNil(t, span.With("value-1", "1"))
	assert.NotNil(t, span.SetError(errors.New("error")))
	assert.NotNil(t, span.SetOk())

	assert.Empty(t, span.Attributes(), "should be empty because fill attributes was not set")
}
