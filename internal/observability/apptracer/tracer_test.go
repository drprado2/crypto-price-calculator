package apptracer

import (
	"context"
	"crypto-price-calculator/internal/configs"
	"crypto-price-calculator/internal/ctxutils"
	"crypto-price-calculator/internal/observability"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProvider(t *testing.T) {
	ctx := ctxutils.WithCid(context.Background(), "x-cid-1234")

	err := Setup(ctx, configs.Get())
	assert.NoError(t, err)

	tp.WithContextAttributeExtractor(func(ctx context.Context) observability.Attributes {
		attributes := observability.NewAttributes()

		if cid := ctxutils.GetCid(ctx); cid != "" {
			attributes = attributes.With("Cid", cid)
		}

		return attributes
	})

	attrs := tp.FillAttributes(ctx, observability.NewAttributes().With("Latency", "2ms"))
	assert.Len(t, attrs, 2)
	assert.Contains(t, attrs, observability.Attribute{Name: "Cid", Value: "x-cid-1234"})
	assert.Contains(t, attrs, observability.Attribute{Name: "Latency", Value: "2ms"})

	ctx, span := tp.Start(ctx, "teste-provider", SpanKindUnspecified, attrs)
	assert.NotNil(t, span)

	assert.NotEmpty(t, ctxutils.GetSpanId(ctx))
	assert.NotEmpty(t, ctxutils.GetTraceId(ctx))

	assert.NoError(t, tp.Close(ctx))

	span.Finish()
}

func TestStart(t *testing.T) {
	ctx := ctxutils.WithCid(context.Background(), "x-cid-1234")

	err := Setup(ctx, configs.Get())
	assert.NoError(t, err)

	tp.WithContextAttributeExtractor(func(ctx context.Context) observability.Attributes {
		attributes := observability.NewAttributes()

		if cid := ctxutils.GetCid(ctx); cid != "" {
			attributes = attributes.With("Cid", cid)
		}

		return attributes
	})

	ctx, span := Start(ctx)
	assert.NotEmpty(t, ctxutils.GetSpanId(ctx))
	assert.NotEmpty(t, ctxutils.GetTraceId(ctx))

	assert.Equal(t, span.Name(), "TestStart")

	span.Finish()
}
