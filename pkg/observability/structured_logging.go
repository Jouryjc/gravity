package observability

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

// LogLevel represents the log level
type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
	LogLevelFatal LogLevel = "fatal"
)

// LogEntry represents a structured log entry
type LogEntry struct {
	Timestamp time.Time            `json:"timestamp"`
	Level     string               `json:"level"`
	Pipeline  string               `json:"pipeline"`
	Component string               `json:"component"`
	TraceID   string               `json:"trace_id,omitempty"`
	SpanID    string               `json:"span_id,omitempty"`
	Message   string               `json:"message"`
	Schema    string               `json:"schema,omitempty"`
	Table     string               `json:"table,omitempty"`
	Operation string               `json:"operation,omitempty"`
	Duration  int64                `json:"duration_ms,omitempty"`
	Error     string               `json:"error,omitempty"`
	Caller    string               `json:"caller,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// StructuredLogger provides structured logging with tracing integration
type StructuredLogger struct {
	pipeline  string
	component string
	logger    *log.Logger
}

// NewStructuredLogger creates a new structured logger
func NewStructuredLogger(pipeline, component string) *StructuredLogger {
	logger := log.New()
	logger.SetFormatter(&log.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
		FieldMap: log.FieldMap{
			log.FieldKeyTime:  "timestamp",
			log.FieldKeyLevel: "level",
			log.FieldKeyMsg:   "message",
		},
	})

	return &StructuredLogger{
		pipeline:  pipeline,
		component: component,
		logger:    logger,
	}
}

// LogContext holds contextual information for logging
type LogContext struct {
	Schema    string
	Table     string
	Operation string
	Duration  time.Duration
	Error     error
	Metadata  map[string]interface{}
}

// WithContext creates a new LogContext
func WithContext() *LogContext {
	return &LogContext{
		Metadata: make(map[string]interface{}),
	}
}

// WithSchema sets the schema name
func (lc *LogContext) WithSchema(schema string) *LogContext {
	lc.Schema = schema
	return lc
}

// WithTable sets the table name
func (lc *LogContext) WithTable(table string) *LogContext {
	lc.Table = table
	return lc
}

// WithOperation sets the operation name
func (lc *LogContext) WithOperation(operation string) *LogContext {
	lc.Operation = operation
	return lc
}

// WithDuration sets the operation duration
func (lc *LogContext) WithDuration(duration time.Duration) *LogContext {
	lc.Duration = duration
	return lc
}

// WithError sets the error
func (lc *LogContext) WithError(err error) *LogContext {
	lc.Error = err
	return lc
}

// WithMetadata adds metadata
func (lc *LogContext) WithMetadata(key string, value interface{}) *LogContext {
	if lc.Metadata == nil {
		lc.Metadata = make(map[string]interface{})
	}
	lc.Metadata[key] = value
	return lc
}

// createLogEntry creates a structured log entry
func (sl *StructuredLogger) createLogEntry(ctx context.Context, level LogLevel, message string, logCtx *LogContext) *LogEntry {
	entry := &LogEntry{
		Timestamp: time.Now(),
		Level:     string(level),
		Pipeline:  sl.pipeline,
		Component: sl.component,
		Message:   message,
		Caller:    getCaller(),
	}

	// Add trace information if available
	if ctx != nil {
		entry.TraceID = GetTraceID(ctx)
		entry.SpanID = GetSpanID(ctx)
	}

	// Add context information if provided
	if logCtx != nil {
		entry.Schema = logCtx.Schema
		entry.Table = logCtx.Table
		entry.Operation = logCtx.Operation
		entry.Metadata = logCtx.Metadata

		if logCtx.Duration > 0 {
			entry.Duration = logCtx.Duration.Nanoseconds() / 1e6 // Convert to milliseconds
		}

		if logCtx.Error != nil {
			entry.Error = logCtx.Error.Error()
		}
	}

	return entry
}

// getCaller returns the caller information
func getCaller() string {
	_, file, line, ok := runtime.Caller(3) // Skip 3 frames to get the actual caller
	if !ok {
		return "unknown"
	}

	// Extract just the filename from the full path
	parts := strings.Split(file, "/")
	filename := parts[len(parts)-1]

	return fmt.Sprintf("%s:%d", filename, line)
}

// logWithLevel logs a message at the specified level
func (sl *StructuredLogger) logWithLevel(ctx context.Context, level LogLevel, message string, logCtx *LogContext) {
	entry := sl.createLogEntry(ctx, level, message, logCtx)

	// Convert to logrus fields
	fields := log.Fields{
		"pipeline":  entry.Pipeline,
		"component": entry.Component,
		"caller":    entry.Caller,
	}

	if entry.TraceID != "" {
		fields["trace_id"] = entry.TraceID
	}
	if entry.SpanID != "" {
		fields["span_id"] = entry.SpanID
	}
	if entry.Schema != "" {
		fields["schema"] = entry.Schema
	}
	if entry.Table != "" {
		fields["table"] = entry.Table
	}
	if entry.Operation != "" {
		fields["operation"] = entry.Operation
	}
	if entry.Duration > 0 {
		fields["duration_ms"] = entry.Duration
	}
	if entry.Error != "" {
		fields["error"] = entry.Error
	}
	if entry.Metadata != nil {
		for k, v := range entry.Metadata {
			fields[k] = v
		}
	}

	// Log with appropriate level
	switch level {
	case LogLevelDebug:
		sl.logger.WithFields(fields).Debug(message)
	case LogLevelInfo:
		sl.logger.WithFields(fields).Info(message)
	case LogLevelWarn:
		sl.logger.WithFields(fields).Warn(message)
	case LogLevelError:
		sl.logger.WithFields(fields).Error(message)
	case LogLevelFatal:
		sl.logger.WithFields(fields).Fatal(message)
	}
}

// Debug logs a debug message
func (sl *StructuredLogger) Debug(ctx context.Context, message string, logCtx ...*LogContext) {
	var lc *LogContext
	if len(logCtx) > 0 {
		lc = logCtx[0]
	}
	sl.logWithLevel(ctx, LogLevelDebug, message, lc)
}

// Info logs an info message
func (sl *StructuredLogger) Info(ctx context.Context, message string, logCtx ...*LogContext) {
	var lc *LogContext
	if len(logCtx) > 0 {
		lc = logCtx[0]
	}
	sl.logWithLevel(ctx, LogLevelInfo, message, lc)
}

// Warn logs a warning message
func (sl *StructuredLogger) Warn(ctx context.Context, message string, logCtx ...*LogContext) {
	var lc *LogContext
	if len(logCtx) > 0 {
		lc = logCtx[0]
	}
	sl.logWithLevel(ctx, LogLevelWarn, message, lc)
}

// Error logs an error message
func (sl *StructuredLogger) Error(ctx context.Context, message string, logCtx ...*LogContext) {
	var lc *LogContext
	if len(logCtx) > 0 {
		lc = logCtx[0]
	}
	sl.logWithLevel(ctx, LogLevelError, message, lc)
}

// Fatal logs a fatal message and exits
func (sl *StructuredLogger) Fatal(ctx context.Context, message string, logCtx ...*LogContext) {
	var lc *LogContext
	if len(logCtx) > 0 {
		lc = logCtx[0]
	}
	sl.logWithLevel(ctx, LogLevelFatal, message, lc)
}

// LogOperation logs the start and end of an operation with timing
func (sl *StructuredLogger) LogOperation(ctx context.Context, operation string, fn func() error) error {
	start := time.Now()
	sl.Info(ctx, fmt.Sprintf("Starting %s", operation), WithContext().WithOperation(operation))

	err := fn()
	duration := time.Since(start)

	logCtx := WithContext().WithOperation(operation).WithDuration(duration)
	if err != nil {
		logCtx = logCtx.WithError(err)
		sl.Error(ctx, fmt.Sprintf("Failed %s", operation), logCtx)
	} else {
		sl.Info(ctx, fmt.Sprintf("Completed %s", operation), logCtx)
	}

	return err
}

// LogWithSpan logs a message and adds it as a span event
func (sl *StructuredLogger) LogWithSpan(ctx context.Context, level LogLevel, message string, logCtx *LogContext) {
	sl.logWithLevel(ctx, level, message, logCtx)

	// Also add as span event if tracing is enabled
	span := trace.SpanFromContext(ctx)
	if span != nil {
		attrs := []trace.EventOption{}
		if logCtx != nil {
			if logCtx.Schema != "" {
				attrs = append(attrs, trace.WithAttributes(AttrSchema.String(logCtx.Schema)))
			}
			if logCtx.Table != "" {
				attrs = append(attrs, trace.WithAttributes(AttrTable.String(logCtx.Table)))
			}
			if logCtx.Operation != "" {
				attrs = append(attrs, trace.WithAttributes(AttrOperation.String(logCtx.Operation)))
			}
			if logCtx.Error != nil {
				attrs = append(attrs, trace.WithAttributes(AttrErrorType.String(logCtx.Error.Error())))
			}
		}
		span.AddEvent(message, attrs...)
	}
}

// Global logger instances for common components
var (
	InputLogger     *StructuredLogger
	FilterLogger    *StructuredLogger
	SchedulerLogger *StructuredLogger
	OutputLogger    *StructuredLogger
	EmitterLogger   *StructuredLogger
)

// InitStructuredLogging initializes global structured loggers
func InitStructuredLogging(pipelineName string) {
	InputLogger = NewStructuredLogger(pipelineName, "input")
	FilterLogger = NewStructuredLogger(pipelineName, "filter")
	SchedulerLogger = NewStructuredLogger(pipelineName, "scheduler")
	OutputLogger = NewStructuredLogger(pipelineName, "output")
	EmitterLogger = NewStructuredLogger(pipelineName, "emitter")
}

// Helper functions for common logging patterns

// LogDatabaseOperation logs database operations with connection info
func LogDatabaseOperation(ctx context.Context, logger *StructuredLogger, operation, endpoint string, duration time.Duration, err error) {
	logCtx := WithContext().
		WithOperation(operation).
		WithDuration(duration).
		WithMetadata("endpoint", endpoint)

	if err != nil {
		logCtx = logCtx.WithError(err)
		logger.Error(ctx, fmt.Sprintf("Database operation failed: %s", operation), logCtx)
		RecordError(GetPipelineName(), logger.component, "database_error")
	} else {
		logger.Info(ctx, fmt.Sprintf("Database operation completed: %s", operation), logCtx)
	}

	RecordComponentLatency(GetPipelineName(), logger.component, operation, duration)
}

// LogKafkaOperation logs Kafka operations
func LogKafkaOperation(ctx context.Context, logger *StructuredLogger, operation, topic string, partition int, duration time.Duration, err error) {
	logCtx := WithContext().
		WithOperation(operation).
		WithDuration(duration).
		WithMetadata("topic", topic).
		WithMetadata("partition", partition)

	if err != nil {
		logCtx = logCtx.WithError(err)
		logger.Error(ctx, fmt.Sprintf("Kafka operation failed: %s", operation), logCtx)
		RecordError(GetPipelineName(), logger.component, "kafka_error")
	} else {
		logger.Info(ctx, fmt.Sprintf("Kafka operation completed: %s", operation), logCtx)
	}

	RecordComponentLatency(GetPipelineName(), logger.component, operation, duration)
	RecordKafkaProduceLatency(GetPipelineName(), topic, fmt.Sprintf("%d", partition), duration)
}

// LogDataProcessing logs data processing operations
func LogDataProcessing(ctx context.Context, logger *StructuredLogger, operation, schema, table string, rowCount int, duration time.Duration, err error) {
	logCtx := WithContext().
		WithSchema(schema).
		WithTable(table).
		WithOperation(operation).
		WithDuration(duration).
		WithMetadata("row_count", rowCount)

	if err != nil {
		logCtx = logCtx.WithError(err)
		logger.Error(ctx, fmt.Sprintf("Data processing failed: %s", operation), logCtx)
		RecordError(GetPipelineName(), logger.component, "processing_error")
	} else {
		logger.Info(ctx, fmt.Sprintf("Data processing completed: %s", operation), logCtx)
	}

	RecordComponentLatency(GetPipelineName(), logger.component, operation, duration)
	if duration.Seconds() > 0 {
		rowsPerSecond := float64(rowCount) / duration.Seconds()
		UpdateThroughput(GetPipelineName(), schema, table, operation, rowsPerSecond)
	}
}