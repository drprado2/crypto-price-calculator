package apptracer

import (
	"context"
	"crypto-price-calculator/internal/observability"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type (
	tracerProviderNoop struct{}
	spanNoop           struct{}
)

func (t tracerProviderNoop) Collector() *sdktrace.TracerProvider {
	panic("implement me")
}

func (t tracerProviderNoop) Start(ctx context.Context, name string, kind SpanKind, attributes observability.Attributes) (context.Context, Span) {
	span := &spanNoop{}
	return ctx, span
}

func (t tracerProviderNoop) FillAttributes(ctx context.Context, attributes observability.Attributes) observability.Attributes {
	return attributes
}

func NewTracerProviderNoop() Provider {
	return &tracerProviderNoop{}
}

func (t tracerProviderNoop) Close(context.Context) error {
	return nil
}

func (t tracerProviderNoop) WithContextAttributeExtractor(extractor observability.ContextFieldExtractor) {
}

func (s spanNoop) TraceID() string {
	return "1"
}

func (s spanNoop) SpanID() string {
	return "1"
}

func (s spanNoop) Name() string {
	return "Noop"
}

func (s spanNoop) With(string, interface{}) Span {
	return s
}

func (s spanNoop) Attributes() observability.Attributes {
	return observability.NewAttributes()
}

func (s spanNoop) SetError(error) Span {
	return s
}

func (s spanNoop) SetOk() Span {
	return s
}

func (s spanNoop) Finish() {}
