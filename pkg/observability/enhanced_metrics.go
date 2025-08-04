package observability

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/moiot/gravity/pkg/env"
)

const (
	// Metric label keys
	LabelPipeline    = "pipeline"
	LabelComponent   = "component"
	LabelSchema      = "schema"
	LabelTable       = "table"
	LabelOperation   = "operation"
	LabelErrorType   = "error_type"
	LabelReason      = "reason"
	LabelEndpoint    = "endpoint"
	LabelTopic       = "topic"
	LabelPartition   = "partition"
	LabelStatus      = "status"
)

// Error monitoring metrics
var (
	// ErrorRateCounter tracks error occurrences by component and type
	ErrorRateCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "gravity",
		Name:      "error_total",
		Help:      "Total number of errors by component and type",
	}, []string{LabelPipeline, LabelComponent, LabelErrorType})

	// RetryCounter tracks retry attempts
	RetryCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "gravity",
		Name:      "retry_total",
		Help:      "Total number of retries",
	}, []string{LabelPipeline, LabelComponent, LabelReason})

	// ConnectionFailureCounter tracks connection failures
	ConnectionFailureCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "gravity",
		Name:      "connection_failure_total",
		Help:      "Total number of connection failures",
	}, []string{LabelPipeline, LabelComponent, LabelEndpoint})
)

// Data quality metrics
var (
	// DataConsistencyGauge measures data consistency between source and target
	DataConsistencyGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "gravity",
		Name:      "data_consistency_score",
		Help:      "Data consistency score between source and target (0-1)",
	}, []string{LabelPipeline, LabelSchema, LabelTable})

	// DataLagHistogram measures data lag between source and target
	DataLagHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "gravity",
		Name:      "data_lag_seconds",
		Help:      "Data lag between source and target in seconds",
		Buckets:   prometheus.ExponentialBuckets(0.1, 2, 15), // 0.1s to ~54min
	}, []string{LabelPipeline, LabelSchema, LabelTable})

	// DataSkipCounter tracks skipped records
	DataSkipCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "gravity",
		Name:      "data_skip_total",
		Help:      "Total number of skipped records",
	}, []string{LabelPipeline, LabelSchema, LabelTable, LabelReason})

	// DataCorruptionCounter tracks data corruption events
	DataCorruptionCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "gravity",
		Name:      "data_corruption_total",
		Help:      "Total number of data corruption events",
	}, []string{LabelPipeline, LabelSchema, LabelTable, LabelErrorType})
)

// Business metrics
var (
	// ThroughputGauge measures data throughput in rows per second
	ThroughputGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "gravity",
		Name:      "throughput_rows_per_second",
		Help:      "Data throughput in rows per second",
	}, []string{LabelPipeline, LabelSchema, LabelTable, LabelOperation})

	// ConnectionHealthGauge tracks connection health status
	ConnectionHealthGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "gravity",
		Name:      "connection_health",
		Help:      "Connection health status (1=healthy, 0=unhealthy)",
	}, []string{LabelPipeline, LabelComponent, LabelEndpoint})

	// ReplicationLagGauge measures replication lag
	ReplicationLagGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "gravity",
		Name:      "replication_lag_seconds",
		Help:      "Replication lag in seconds",
	}, []string{LabelPipeline, LabelComponent})

	// TableSizeGauge tracks table sizes
	TableSizeGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "gravity",
		Name:      "table_size_bytes",
		Help:      "Table size in bytes",
	}, []string{LabelPipeline, LabelSchema, LabelTable})
)

// Performance metrics
var (
	// ComponentLatencyHistogram measures component processing latency
	ComponentLatencyHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "gravity",
		Name:      "component_latency_seconds",
		Help:      "Component processing latency in seconds",
		Buckets:   prometheus.ExponentialBuckets(0.0001, 2, 20), // 0.1ms to ~104s
	}, []string{LabelPipeline, LabelComponent, LabelOperation})

	// MemoryUsageGauge tracks memory usage by component
	MemoryUsageGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "gravity",
		Name:      "memory_usage_bytes",
		Help:      "Memory usage in bytes by component",
	}, []string{LabelPipeline, LabelComponent})

	// CPUUsageGauge tracks CPU usage by component
	CPUUsageGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "gravity",
		Name:      "cpu_usage_percent",
		Help:      "CPU usage percentage by component",
	}, []string{LabelPipeline, LabelComponent})
)

// Kafka-specific metrics
var (
	// KafkaProduceLatencyHistogram measures Kafka produce latency
	KafkaProduceLatencyHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "gravity",
		Name:      "kafka_produce_latency_seconds",
		Help:      "Kafka produce latency in seconds",
		Buckets:   prometheus.ExponentialBuckets(0.001, 2, 15), // 1ms to ~32s
	}, []string{LabelPipeline, LabelTopic, LabelPartition})

	// KafkaMessageSizeHistogram measures Kafka message sizes
	KafkaMessageSizeHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "gravity",
		Name:      "kafka_message_size_bytes",
		Help:      "Kafka message size in bytes",
		Buckets:   prometheus.ExponentialBuckets(100, 2, 15), // 100B to ~3.2MB
	}, []string{LabelPipeline, LabelTopic})

	// KafkaPartitionOffsetGauge tracks Kafka partition offsets
	KafkaPartitionOffsetGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "gravity",
		Name:      "kafka_partition_offset",
		Help:      "Current Kafka partition offset",
	}, []string{LabelPipeline, LabelTopic, LabelPartition})
)

// MySQL/MongoDB specific metrics
var (
	// BinlogPositionGauge tracks MySQL binlog position
	BinlogPositionGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "gravity",
		Name:      "mysql_binlog_position",
		Help:      "Current MySQL binlog position",
	}, []string{LabelPipeline, "binlog_file"})

	// OplogTimestampGauge tracks MongoDB oplog timestamp
	OplogTimestampGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "gravity",
		Name:      "mongo_oplog_timestamp",
		Help:      "Current MongoDB oplog timestamp",
	}, []string{LabelPipeline})

	// DatabaseConnectionsGauge tracks database connections
	DatabaseConnectionsGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "gravity",
		Name:      "database_connections",
		Help:      "Number of active database connections",
	}, []string{LabelPipeline, LabelComponent, LabelStatus})
)

// Helper functions for common metric operations

// RecordError records an error with appropriate labels
func RecordError(pipeline, component, errorType string) {
	ErrorRateCounter.WithLabelValues(pipeline, component, errorType).Inc()
}

// RecordRetry records a retry attempt
func RecordRetry(pipeline, component, reason string) {
	RetryCounter.WithLabelValues(pipeline, component, reason).Inc()
}

// RecordConnectionFailure records a connection failure
func RecordConnectionFailure(pipeline, component, endpoint string) {
	ConnectionFailureCounter.WithLabelValues(pipeline, component, endpoint).Inc()
}

// UpdateConnectionHealth updates connection health status
func UpdateConnectionHealth(pipeline, component, endpoint string, healthy bool) {
	value := 0.0
	if healthy {
		value = 1.0
	}
	ConnectionHealthGauge.WithLabelValues(pipeline, component, endpoint).Set(value)
}

// RecordComponentLatency records component processing latency
func RecordComponentLatency(pipeline, component, operation string, duration time.Duration) {
	ComponentLatencyHistogram.WithLabelValues(pipeline, component, operation).Observe(duration.Seconds())
}

// UpdateThroughput updates throughput metrics
func UpdateThroughput(pipeline, schema, table, operation string, rowsPerSecond float64) {
	ThroughputGauge.WithLabelValues(pipeline, schema, table, operation).Set(rowsPerSecond)
}

// RecordDataLag records data lag between source and target
func RecordDataLag(pipeline, schema, table string, lag time.Duration) {
	DataLagHistogram.WithLabelValues(pipeline, schema, table).Observe(lag.Seconds())
}

// UpdateDataConsistency updates data consistency score
func UpdateDataConsistency(pipeline, schema, table string, score float64) {
	DataConsistencyGauge.WithLabelValues(pipeline, schema, table).Set(score)
}

// RecordDataSkip records a skipped record
func RecordDataSkip(pipeline, schema, table, reason string) {
	DataSkipCounter.WithLabelValues(pipeline, schema, table, reason).Inc()
}

// RecordKafkaProduceLatency records Kafka produce latency
func RecordKafkaProduceLatency(pipeline, topic, partition string, duration time.Duration) {
	KafkaProduceLatencyHistogram.WithLabelValues(pipeline, topic, partition).Observe(duration.Seconds())
}

// RecordKafkaMessageSize records Kafka message size
func RecordKafkaMessageSize(pipeline, topic string, size int) {
	KafkaMessageSizeHistogram.WithLabelValues(pipeline, topic).Observe(float64(size))
}

// UpdateBinlogPosition updates MySQL binlog position
func UpdateBinlogPosition(pipeline, binlogFile string, position float64) {
	BinlogPositionGauge.WithLabelValues(pipeline, binlogFile).Set(position)
}

// UpdateOplogTimestamp updates MongoDB oplog timestamp
func UpdateOplogTimestamp(pipeline string, timestamp float64) {
	OplogTimestampGauge.WithLabelValues(pipeline).Set(timestamp)
}

// GetPipelineName returns the current pipeline name from environment
func GetPipelineName() string {
	if env.PipelineName != "" {
		return env.PipelineName
	}
	return "unknown"
}

// init registers all metrics with Prometheus
func init() {
	prometheus.MustRegister(
		// Error metrics
		ErrorRateCounter,
		RetryCounter,
		ConnectionFailureCounter,
		
		// Data quality metrics
		DataConsistencyGauge,
		DataLagHistogram,
		DataSkipCounter,
		DataCorruptionCounter,
		
		// Business metrics
		ThroughputGauge,
		ConnectionHealthGauge,
		ReplicationLagGauge,
		TableSizeGauge,
		
		// Performance metrics
		ComponentLatencyHistogram,
		MemoryUsageGauge,
		CPUUsageGauge,
		
		// Kafka metrics
		KafkaProduceLatencyHistogram,
		KafkaMessageSizeHistogram,
		KafkaPartitionOffsetGauge,
		
		// Database metrics
		BinlogPositionGauge,
		OplogTimestampGauge,
		DatabaseConnectionsGauge,
	)
}