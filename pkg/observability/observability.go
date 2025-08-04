package observability

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// ObservabilityConfig holds the configuration for observability components
type ObservabilityConfig struct {
	Tracing  *TracingConfig  `yaml:"tracing"`
	Alerting *AlertingConfig `yaml:"alerting"`
	Enabled  bool            `yaml:"enabled"`
	Pipeline string          `yaml:"pipeline"`
}

// ObservabilityManager manages all observability components
type ObservabilityManager struct {
	config       *ObservabilityConfig
	alertManager *AlertManager
	tracer       trace.Tracer
	logger       *StructuredLogger
	mu           sync.RWMutex
	initialized  bool
}

// Global observability manager instance
var (
	globalManager *ObservabilityManager
	once          sync.Once
	pipelineName  string = "default"
)

// Initialize initializes the global observability system
func Initialize(config *ObservabilityConfig) error {
	var initErr error
	once.Do(func() {
		globalManager = &ObservabilityManager{
			config: config,
		}
		initErr = globalManager.init()
	})
	return initErr
}

// GetManager returns the global observability manager
func GetManager() *ObservabilityManager {
	return globalManager
}

// SetPipelineName sets the global pipeline name
func SetPipelineName(name string) {
	pipelineName = name
}

// GetPipelineName returns the current pipeline name
func GetPipelineName() string {
	return pipelineName
}

// init initializes all observability components
func (om *ObservabilityManager) init() error {
	om.mu.Lock()
	defer om.mu.Unlock()

	if om.initialized {
		return nil
	}

	if !om.config.Enabled {
		return nil
	}

	// Set pipeline name
	if om.config.Pipeline != "" {
		SetPipelineName(om.config.Pipeline)
	}

	// Initialize tracing
	if om.config.Tracing != nil {
		if err := InitTracing(om.config.Tracing); err != nil {
			return fmt.Errorf("failed to initialize tracing: %w", err)
		}
		om.tracer = otel.Tracer("gravity")
	}

	// Initialize enhanced metrics
	if err := InitEnhancedMetrics(); err != nil {
		return fmt.Errorf("failed to initialize enhanced metrics: %w", err)
	}

	// Initialize structured logging
	InitStructuredLogging(GetPipelineName())
	om.logger = NewStructuredLogger(GetPipelineName(), "observability")

	// Initialize alert manager
	if om.config.Alerting != nil {
		alertManager, err := NewAlertManager(om.config.Alerting)
		if err != nil {
			return fmt.Errorf("failed to initialize alert manager: %w", err)
		}
		om.alertManager = alertManager
	}

	om.initialized = true
	om.logger.Info(context.Background(), "Observability system initialized successfully")

	return nil
}

// IsInitialized returns whether the observability system is initialized
func (om *ObservabilityManager) IsInitialized() bool {
	om.mu.RLock()
	defer om.mu.RUnlock()
	return om.initialized
}

// GetTracer returns the OpenTelemetry tracer
func (om *ObservabilityManager) GetTracer() trace.Tracer {
	om.mu.RLock()
	defer om.mu.RUnlock()
	return om.tracer
}

// GetAlertManager returns the alert manager
func (om *ObservabilityManager) GetAlertManager() *AlertManager {
	om.mu.RLock()
	defer om.mu.RUnlock()
	return om.alertManager
}

// GetLogger returns the structured logger
func (om *ObservabilityManager) GetLogger() *StructuredLogger {
	om.mu.RLock()
	defer om.mu.RUnlock()
	return om.logger
}

// StartSpanFromContext starts a new span from context
func StartSpanFromContext(ctx context.Context, operationName string) (context.Context, trace.Span) {
	if globalManager == nil || !globalManager.IsInitialized() {
		return ctx, trace.SpanFromContext(ctx)
	}
	return StartSpan(ctx, operationName)
}

// RecordMetric records a metric with the given name and value
func RecordMetric(metricName string, value float64, labels map[string]string) {
	if globalManager == nil || !globalManager.IsInitialized() {
		return
	}

	// This is a simplified metric recording function
	// In practice, you would route to the appropriate metric based on the name
	switch metricName {
	case "error_rate":
		if pipeline, ok := labels["pipeline"]; ok {
			if component, ok := labels["component"]; ok {
				if errorType, ok := labels["error_type"]; ok {
					RecordError(pipeline, component, errorType)
				}
			}
		}
	case "latency":
		if pipeline, ok := labels["pipeline"]; ok {
			if component, ok := labels["component"]; ok {
				if operation, ok := labels["operation"]; ok {
					RecordComponentLatency(pipeline, component, operation, time.Duration(value)*time.Millisecond)
				}
			}
		}
	case "throughput":
		if pipeline, ok := labels["pipeline"]; ok {
			if schema, ok := labels["schema"]; ok {
				if table, ok := labels["table"]; ok {
					if operation, ok := labels["operation"]; ok {
						UpdateThroughput(pipeline, schema, table, operation, value)
					}
				}
			}
		}
	}
}

// LogWithTrace logs a message with trace context
func LogWithTrace(ctx context.Context, level LogLevel, message string, logCtx *LogContext) {
	if globalManager == nil || !globalManager.IsInitialized() {
		return
	}

	logger := globalManager.GetLogger()
	if logger != nil {
		logger.LogWithSpan(ctx, level, message, logCtx)
	}
}

// CheckAlerts evaluates all alert conditions and returns active alerts
func CheckAlerts(ctx context.Context) ([]ActiveAlert, error) {
	if globalManager == nil || !globalManager.IsInitialized() {
		return nil, fmt.Errorf("observability system not initialized")
	}

	alertManager := globalManager.GetAlertManager()
	if alertManager == nil {
		return nil, fmt.Errorf("alert manager not configured")
	}

	return alertManager.EvaluateAlerts(ctx)
}

// GetHealthStatus returns the health status of the observability system
func GetHealthStatus() map[string]interface{} {
	status := map[string]interface{}{
		"initialized": false,
		"tracing":     false,
		"metrics":     false,
		"logging":     false,
		"alerting":    false,
	}

	if globalManager == nil {
		return status
	}

	status["initialized"] = globalManager.IsInitialized()

	if globalManager.IsInitialized() {
		status["tracing"] = globalManager.GetTracer() != nil
		status["metrics"] = true // Enhanced metrics are always available once initialized
		status["logging"] = globalManager.GetLogger() != nil
		status["alerting"] = globalManager.GetAlertManager() != nil
	}

	return status
}

// GetMetricsSnapshot returns a snapshot of current metrics
func GetMetricsSnapshot() (map[string]interface{}, error) {
	if globalManager == nil || !globalManager.IsInitialized() {
		return nil, fmt.Errorf("observability system not initialized")
	}

	// Gather metrics from Prometheus registry
	metricFamilies, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		return nil, fmt.Errorf("failed to gather metrics: %w", err)
	}

	snapshot := make(map[string]interface{})
	for _, mf := range metricFamilies {
		if mf.GetName() != "" {
			metrics := make([]map[string]interface{}, 0)
			for _, metric := range mf.GetMetric() {
				m := map[string]interface{}{
					"labels": make(map[string]string),
				}

				// Extract labels
				for _, label := range metric.GetLabel() {
					m["labels"].(map[string]string)[label.GetName()] = label.GetValue()
				}

				// Extract value based on metric type
				switch mf.GetType() {
				case prometheus.CounterValue:
					if metric.GetCounter() != nil {
						m["value"] = metric.GetCounter().GetValue()
					}
				case prometheus.GaugeValue:
					if metric.GetGauge() != nil {
						m["value"] = metric.GetGauge().GetValue()
					}
				case prometheus.HistogramValue:
					if metric.GetHistogram() != nil {
						m["count"] = metric.GetHistogram().GetSampleCount()
						m["sum"] = metric.GetHistogram().GetSampleSum()
					}
				}

				metrics = append(metrics, m)
			}
			snapshot[mf.GetName()] = metrics
		}
	}

	return snapshot, nil
}

// Shutdown gracefully shuts down the observability system
func Shutdown(ctx context.Context) error {
	if globalManager == nil {
		return nil
	}

	globalManager.mu.Lock()
	defer globalManager.mu.Unlock()

	if !globalManager.initialized {
		return nil
	}

	// Log shutdown
	if globalManager.logger != nil {
		globalManager.logger.Info(ctx, "Shutting down observability system")
	}

	// Shutdown tracing
	if err := ShutdownTracing(ctx); err != nil {
		return fmt.Errorf("failed to shutdown tracing: %w", err)
	}

	globalManager.initialized = false
	return nil
}

// Helper functions for common observability patterns

// TraceOperation wraps an operation with tracing and logging
func TraceOperation(ctx context.Context, operationName string, fn func(ctx context.Context) error) error {
	if globalManager == nil || !globalManager.IsInitialized() {
		return fn(ctx)
	}

	// Start span
	ctx, span := StartSpanFromContext(ctx, operationName)
	defer FinishSpan(span)

	// Log operation start
	logger := globalManager.GetLogger()
	if logger != nil {
		logger.Info(ctx, fmt.Sprintf("Starting operation: %s", operationName),
			WithContext().WithOperation(operationName))
	}

	// Execute operation
	start := time.Now()
	err := fn(ctx)
	duration := time.Since(start)

	// Record metrics and logs
	if err != nil {
		RecordError(GetPipelineName(), "operation", "execution_error")
		RecordErrorOnSpan(span, err)
		if logger != nil {
			logger.Error(ctx, fmt.Sprintf("Operation failed: %s", operationName),
				WithContext().WithOperation(operationName).WithDuration(duration).WithError(err))
		}
	} else {
		if logger != nil {
			logger.Info(ctx, fmt.Sprintf("Operation completed: %s", operationName),
				WithContext().WithOperation(operationName).WithDuration(duration))
		}
	}

	// Record latency
	RecordComponentLatency(GetPipelineName(), "operation", operationName, duration)

	return err
}

// TraceDataProcessing wraps data processing operations with comprehensive observability
func TraceDataProcessing(ctx context.Context, schema, table, operation string, rowCount int, fn func(ctx context.Context) error) error {
	if globalManager == nil || !globalManager.IsInitialized() {
		return fn(ctx)
	}

	operationName := fmt.Sprintf("%s.%s.%s", schema, table, operation)

	// Start span with attributes
	ctx, span := StartSpanFromContext(ctx, operationName)
	AddSpanAttributes(span, map[string]interface{}{
		"db.schema":    schema,
		"db.table":     table,
		"db.operation": operation,
		"row.count":    rowCount,
	})
	defer FinishSpan(span)

	// Execute operation with logging
	logger := globalManager.GetLogger()
	if logger != nil {
		return logger.LogOperation(ctx, operationName, func() error {
			start := time.Now()
			err := fn(ctx)
			duration := time.Since(start)

			// Log data processing details
			LogDataProcessing(ctx, logger, operation, schema, table, rowCount, duration, err)

			return err
		})
	}

	return fn(ctx)
}

// Default configuration for quick setup
func GetDefaultConfig() *ObservabilityConfig {
	return &ObservabilityConfig{
		Enabled:  true,
		Pipeline: "gravity-pipeline",
		Tracing: &TracingConfig{
			Enabled:     true,
			ServiceName: "gravity",
			JaegerEndpoint: "http://localhost:14268/api/traces",
			SampleRate:  1.0,
		},
		Alerting: &AlertingConfig{
			PrometheusURL:    "http://localhost:9090",
			EvaluationPeriod: 30 * time.Second,
		},
	}
}