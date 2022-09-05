package apptracer

import (
	"context"
	"crypto-price-calculator/internal/configs"
	"crypto-price-calculator/internal/ctxutils"
	"crypto-price-calculator/internal/observability"
	"crypto-price-calculator/internal/observability/applog"
	"go.opentelemetry.io/contrib"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"runtime"
	"strings"
)

const (
	InstrumentationName = "grpc/otelgrpc"

	SpanKindUnspecified = SpanKind(trace.SpanKindUnspecified)
	SpanKindInternal    = SpanKind(trace.SpanKindInternal)
	SpanKindServer      = SpanKind(trace.SpanKindServer)
	SpanKindClient      = SpanKind(trace.SpanKindClient)
	SpanKindProducer    = SpanKind(trace.SpanKindProducer)
	SpanKindConsumer    = SpanKind(trace.SpanKindConsumer)

	ContextSpanIDKey  = "x-span-id"
	ContextTraceIDKey = "x-trace-id"
)

var (
	tp                         Provider
	defaultAttributesExtractor = func(ctx context.Context) observability.Attributes {
		attrs := observability.NewAttributes()

		if cid := ctxutils.GetCid(ctx); cid != "" {
			attrs = attrs.With("Cid", cid)
		}

		return attrs
	}
)

type (
	SpanKind int

	Provider interface {
		Close(ctx context.Context) error
		WithContextAttributeExtractor(extractor observability.ContextFieldExtractor)
		Start(ctx context.Context, name string, kind SpanKind, attributes observability.Attributes) (context.Context, Span)
		FillAttributes(ctx context.Context, attributes observability.Attributes) observability.Attributes
		Collector() *sdktrace.TracerProvider
	}

	Span interface {
		TraceID() string
		SpanID() string
		Name() string
		Attributes() observability.Attributes
		With(name string, value interface{}) Span
		SetError(err error) Span
		SetOk() Span
		Finish()
	}

	provider struct {
		tracer         trace.Tracer
		collector      *sdktrace.TracerProvider
		fieldExtractor observability.ContextFieldExtractor
	}

	span struct {
		ctx      context.Context
		name     string
		kind     SpanKind
		attrs    observability.Attributes
		internal trace.Span
	}
)

func init() {
	tp = NewTracerProviderNoop()
}

func Setup(ctx context.Context, config *configs.Configuration) error {
	logger := applog.Logger(ctx)

	if !config.EnableTracing {
		logger.Info("tracer will not be enabled")
		return nil
	}

	jaegerExporter, err := jaeger.New(
		jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://localhost:14268/api/traces")))
	if err != nil {
		return err
	}

	collector := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(jaegerExporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(config.ServiceName),
			semconv.DeploymentEnvironmentKey.String(config.ServerEnvironment),
		)),
	)

	otel.SetTracerProvider(collector)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}))

	tracer := otel.GetTracerProvider().Tracer(
		InstrumentationName,
		trace.WithInstrumentationVersion(contrib.SemVersion()),
	)

	tp = &provider{
		collector: collector,
		tracer:    tracer,
	}

	tp.WithContextAttributeExtractor(defaultAttributesExtractor)

	return nil
}

func GetProvider() Provider {
	return tp
}

func (t *provider) Close(ctx context.Context) error {
	return t.collector.Shutdown(ctx)
}

func Close(ctx context.Context) error {
	return tp.Close(ctx)
}

func (t *provider) WithContextAttributeExtractor(extractor observability.ContextFieldExtractor) {
	t.fieldExtractor = extractor
}

func (t *provider) FillAttributes(ctx context.Context, attributes observability.Attributes) observability.Attributes {
	if t.fieldExtractor != nil {
		attributes = attributes.Add(t.fieldExtractor(ctx)...)
	}

	return attributes
}

func (t *provider) Collector() *sdktrace.TracerProvider {
	return t.collector
}

func (t *provider) Start(ctx context.Context, name string, kind SpanKind, attributes observability.Attributes) (context.Context, Span) {
	attributes = t.FillAttributes(ctx, attributes)

	ctx, internal := t.tracer.Start(
		ctx,
		name,
		trace.WithSpanKind(trace.SpanKind(kind)))

	ctx = context.WithValue(ctx, ContextTraceIDKey, internal.SpanContext().TraceID().String())
	ctx = context.WithValue(ctx, ContextSpanIDKey, internal.SpanContext().SpanID().String())

	span := &span{
		ctx:      ctx,
		name:     name,
		kind:     kind,
		attrs:    attributes,
		internal: internal,
	}

	return ctx, span
}

func Start(ctx context.Context, attributes ...observability.Attribute) (context.Context, Span) {
	funcName := ""
	if pc, _, _, ok := runtime.Caller(1); ok {
		funcName = runtime.FuncForPC(pc).Name()
		lastDot := strings.LastIndexByte(funcName, '.')
		if lastDot < 0 {
			lastDot = 0
		}
		funcName = funcName[lastDot+1:]
	}

	return tp.Start(ctx, funcName, SpanKindUnspecified, attributes)
}

func StartOperation(ctx context.Context, name string, kind SpanKind, attributes ...observability.Attribute) (context.Context, Span) {
	return tp.Start(ctx, name, kind, attributes)
}

func (s *span) TraceID() string {
	return s.internal.SpanContext().TraceID().String()
}

func (s *span) SpanID() string {
	return s.internal.SpanContext().SpanID().String()
}

func (s *span) Name() string {
	return s.name
}

func (s *span) Attributes() observability.Attributes {
	return s.attrs
}

func (s *span) With(name string, value interface{}) Span {
	s.attrs = s.attrs.With(name, value)
	return s
}

func (s *span) SetOk() Span {
	s.internal.SetStatus(codes.Ok, s.name)
	return s
}

func (s *span) SetError(err error) Span {
	s.internal.SetStatus(codes.Error, err.Error())
	s.internal.RecordError(err)
	return s
}

func (s *span) Finish() {
	s.internal.End()
}
