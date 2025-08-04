package observability

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv/v1.4.0"
	trace2 "go.opentelemetry.io/otel/trace"
	log "github.com/sirupsen/logrus"
)

const (
	// Span names for different components
	SpanNameInput     = "gravity.input"
	SpanNameFilter    = "gravity.filter"
	SpanNameScheduler = "gravity.scheduler"
	SpanNameOutput    = "gravity.output"
	SpanNameEmitter   = "gravity.emitter"
)

var (
	tracer trace2.Tracer
)

// TracingConfig holds the configuration for distributed tracing
type TracingConfig struct {
	Enabled     bool   `mapstructure:"enabled" json:"enabled" toml:"enabled"`
	ServiceName string `mapstructure:"service-name" json:"service-name" toml:"service-name"`
	JaegerURL   string `mapstructure:"jaeger-url" json:"jaeger-url" toml:"jaeger-url"`
	SampleRate  float64 `mapstructure:"sample-rate" json:"sample-rate" toml:"sample-rate"`
}

// InitTracing initializes the OpenTelemetry tracing
func InitTracing(config TracingConfig) error {
	if !config.Enabled {
		log.Info("[tracing] Tracing is disabled")
		return nil
	}

	if config.ServiceName == "" {
		config.ServiceName = "gravity"
	}

	if config.SampleRate == 0 {
		config.SampleRate = 0.1 // Default 10% sampling
	}

	// Create Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(config.JaegerURL)))
	if err != nil {
		return fmt.Errorf("failed to create Jaeger exporter: %w", err)
	}

	// Create trace provider
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(config.ServiceName),
			semconv.ServiceVersionKey.String("2.0"),
		)),
		trace.WithSampler(trace.TraceIDRatioBased(config.SampleRate)),
	)

	// Set global trace provider
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// Get tracer instance
	tracer = otel.Tracer("gravity")

	log.Infof("[tracing] Initialized with service name: %s, sample rate: %.2f", config.ServiceName, config.SampleRate)
	return nil
}

// StartSpan starts a new span with the given name and context
func StartSpan(ctx context.Context, spanName string, attrs ...attribute.KeyValue) (context.Context, trace2.Span) {
	if tracer == nil {
		return ctx, trace2.SpanFromContext(ctx)
	}
	return tracer.Start(ctx, spanName, trace2.WithAttributes(attrs...))
}

// AddSpanAttributes adds attributes to the current span
func AddSpanAttributes(span trace2.Span, attrs ...attribute.KeyValue) {
	if span != nil {
		span.SetAttributes(attrs...)
	}
}

// AddSpanEvent adds an event to the current span
func AddSpanEvent(span trace2.Span, name string, attrs ...attribute.KeyValue) {
	if span != nil {
		span.AddEvent(name, trace2.WithAttributes(attrs...))
	}
}

// RecordError records an error in the current span
func RecordError(span trace2.Span, err error, attrs ...attribute.KeyValue) {
	if span != nil && err != nil {
		span.RecordError(err, trace2.WithAttributes(attrs...))
		span.SetStatus(trace2.StatusError, err.Error())
	}
}

// FinishSpan finishes the span with optional attributes
func FinishSpan(span trace2.Span, attrs ...attribute.KeyValue) {
	if span != nil {
		if len(attrs) > 0 {
			span.SetAttributes(attrs...)
		}
		span.End()
	}
}

// GetTraceID returns the trace ID from the current span context
func GetTraceID(ctx context.Context) string {
	span := trace2.SpanFromContext(ctx)
	if span != nil {
		return span.SpanContext().TraceID().String()
	}
	return ""
}

// GetSpanID returns the span ID from the current span context
func GetSpanID(ctx context.Context) string {
	span := trace2.SpanFromContext(ctx)
	if span != nil {
		return span.SpanContext().SpanID().String()
	}
	return ""
}

// InjectTraceContext injects trace context into a map (for message headers)
func InjectTraceContext(ctx context.Context, carrier map[string]string) {
	if carrier == nil {
		carrier = make(map[string]string)
	}
	otel.GetTextMapPropagator().Inject(ctx, propagation.MapCarrier(carrier))
}

// ExtractTraceContext extracts trace context from a map (from message headers)
func ExtractTraceContext(ctx context.Context, carrier map[string]string) context.Context {
	if carrier == nil {
		return ctx
	}
	return otel.GetTextMapPropagator().Extract(ctx, propagation.MapCarrier(carrier))
}

// WithTimeout creates a span with timeout monitoring
func WithTimeout(ctx context.Context, spanName string, timeout time.Duration, attrs ...attribute.KeyValue) (context.Context, trace2.Span, context.CancelFunc) {
	ctx, span := StartSpan(ctx, spanName, attrs...)
	ctx, cancel := context.WithTimeout(ctx, timeout)
	
	// Monitor for timeout
	go func() {
		<-ctx.Done()
		if ctx.Err() == context.DeadlineExceeded {
			AddSpanEvent(span, "timeout", attribute.String("timeout", timeout.String()))
			span.SetStatus(trace2.StatusError, "operation timed out")
		}
	}()
	
	return ctx, span, cancel
}

// Common attribute keys for Gravity components
var (
	AttrPipeline   = attribute.Key("gravity.pipeline")
	AttrComponent  = attribute.Key("gravity.component")
	AttrSchema     = attribute.Key("gravity.schema")
	AttrTable      = attribute.Key("gravity.table")
	AttrOperation  = attribute.Key("gravity.operation")
	AttrRowCount   = attribute.Key("gravity.row_count")
	AttrBinlogFile = attribute.Key("gravity.binlog_file")
	AttrBinlogPos  = attribute.Key("gravity.binlog_pos")
	AttrOplogTS    = attribute.Key("gravity.oplog_ts")
	AttrKafkaTopic = attribute.Key("gravity.kafka_topic")
	AttrKafkaPartition = attribute.Key("gravity.kafka_partition")
	AttrErrorType  = attribute.Key("gravity.error_type")
	AttrRetryCount = attribute.Key("gravity.retry_count")
)