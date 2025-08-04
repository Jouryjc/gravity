package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/moiot/gravity/pkg/observability"
)

// 这个示例展示了如何在Gravity项目中集成新的可观测性功能
// This example demonstrates how to integrate the new observability features in the Gravity project

func main() {
	// 1. 初始化可观测性系统
	// Initialize observability system
	config := observability.GetDefaultConfig()
	config.ServiceName = "gravity-example"
	config.ServiceVersion = "1.0.0"
	config.Environment = "development"
	config.TracingConfig.JaegerEndpoint = "http://localhost:14268/api/traces"
	config.MetricsConfig.PrometheusPort = 9090
	config.AlertingConfig.PrometheusURL = "http://localhost:9090"

	err := observability.Initialize(config)
	if err != nil {
		log.Fatalf("Failed to initialize observability: %v", err)
	}
	defer observability.Shutdown()

	// 2. 设置管道名称
	// Set pipeline name
	observability.SetPipelineName("mysql-to-kafka-pipeline")

	// 3. 演示分布式追踪
	// Demonstrate distributed tracing
	demonstrateTracing()

	// 4. 演示指标记录
	// Demonstrate metrics recording
	demonstrateMetrics()

	// 5. 演示结构化日志
	// Demonstrate structured logging
	demonstrateLogging()

	// 6. 演示告警检查
	// Demonstrate alert checking
	demonstrateAlerting()

	// 7. 获取健康状态
	// Get health status
	health := observability.GetHealthStatus()
	fmt.Printf("System Health: %+v\n", health)
}

// 演示分布式追踪功能
// Demonstrate distributed tracing functionality
func demonstrateTracing() {
	fmt.Println("=== Distributed Tracing Demo ===")

	// 创建根span
	// Create root span
	ctx := context.Background()
	span := observability.StartSpanFromContext(ctx, "mysql_binlog_read")
	defer observability.FinishSpan(span)

	// 添加span属性
	// Add span attributes
	observability.AddSpanAttributes(span, map[string]interface{}{
		"database.type":     "mysql",
		"database.endpoint": "localhost:3306",
		"schema":            "test_db",
		"table":             "users",
		"operation":         "INSERT",
	})

	// 模拟数据处理
	// Simulate data processing
	processData(observability.InjectTraceContext(ctx, span))

	fmt.Println("Tracing demo completed")
}

// 模拟数据处理过程
// Simulate data processing
func processData(ctx context.Context) {
	// 从上下文中提取trace信息并创建子span
	// Extract trace info from context and create child span
	span := observability.StartSpanFromContext(ctx, "data_transformation")
	defer observability.FinishSpan(span)

	// 模拟处理时间
	// Simulate processing time
	time.Sleep(100 * time.Millisecond)

	// 记录处理结果
	// Record processing result
	observability.AddSpanAttributes(span, map[string]interface{}{
		"rows_processed": 150,
		"processing_time": "100ms",
	})

	// 发送到Kafka
	// Send to Kafka
	sendToKafka(ctx)
}

// 模拟发送到Kafka
// Simulate sending to Kafka
func sendToKafka(ctx context.Context) {
	span := observability.StartSpanFromContext(ctx, "kafka_produce")
	defer observability.FinishSpan(span)

	// 模拟Kafka发送
	// Simulate Kafka sending
	time.Sleep(50 * time.Millisecond)

	observability.AddSpanAttributes(span, map[string]interface{}{
		"kafka.topic":     "gravity.test_db.users",
		"kafka.partition": 0,
		"kafka.offset":    12345,
		"message_size":    1024,
	})
}

// 演示指标记录功能
// Demonstrate metrics recording functionality
func demonstrateMetrics() {
	fmt.Println("=== Metrics Recording Demo ===")

	// 记录错误指标
	// Record error metrics
	observability.RecordMetric("error", map[string]string{
		"component":   "input",
		"error_type":  "connection_timeout",
		"schema":      "test_db",
		"table":       "users",
	}, 1)

	// 记录重试指标
	// Record retry metrics
	observability.RecordMetric("retry", map[string]string{
		"component":    "output",
		"retry_reason": "kafka_unavailable",
		"schema":       "test_db",
		"table":        "orders",
	}, 1)

	// 记录数据一致性指标
	// Record data consistency metrics
	observability.RecordMetric("data_consistency", map[string]string{
		"schema": "test_db",
		"table":  "products",
	}, 0.99)

	// 记录吞吐量指标
	// Record throughput metrics
	observability.RecordMetric("throughput", map[string]string{
		"schema":    "test_db",
		"table":     "users",
		"operation": "INSERT",
	}, 1500)

	// 记录延迟指标
	// Record latency metrics
	start := time.Now()
	time.Sleep(200 * time.Millisecond) // 模拟操作
	duration := time.Since(start)

	observability.RecordMetric("component_latency", map[string]string{
		"component": "filter",
		"operation": "data_validation",
	}, duration.Seconds())

	fmt.Println("Metrics recording demo completed")
}

// 演示结构化日志功能
// Demonstrate structured logging functionality
func demonstrateLogging() {
	fmt.Println("=== Structured Logging Demo ===")

	// 创建带有trace信息的上下文
	// Create context with trace information
	ctx := context.Background()
	span := observability.StartSpanFromContext(ctx, "logging_demo")
	defer observability.FinishSpan(span)

	// 使用带trace的日志记录
	// Use logging with trace
	observability.LogWithTrace(ctx, "info", "Processing binlog event", map[string]interface{}{
		"component":      "input",
		"schema":         "test_db",
		"table":          "users",
		"operation":      "INSERT",
		"binlog_file":    "mysql-bin.000001",
		"binlog_position": 12345,
		"rows_affected":  1,
	})

	// 记录错误日志
	// Record error log
	observability.LogWithTrace(ctx, "error", "Failed to connect to database", map[string]interface{}{
		"component":      "input",
		"error":          "connection timeout",
		"endpoint":       "localhost:3306",
		"retry_count":    3,
		"max_retries":    5,
	})

	// 记录性能日志
	// Record performance log
	observability.LogWithTrace(ctx, "info", "Batch processing completed", map[string]interface{}{
		"component":        "processor",
		"batch_size":       1000,
		"processing_time":  "2.5s",
		"throughput":       400,
		"memory_usage":     "256MB",
	})

	fmt.Println("Structured logging demo completed")
}

// 演示告警检查功能
// Demonstrate alerting functionality
func demonstrateAlerting() {
	fmt.Println("=== Alerting Demo ===")

	// 检查告警
	// Check alerts
	alerts := observability.CheckAlerts()
	if len(alerts) > 0 {
		fmt.Printf("Active alerts found: %d\n", len(alerts))
		for _, alert := range alerts {
			fmt.Printf("- %s: %s\n", alert.Name, alert.Description)
		}
	} else {
		fmt.Println("No active alerts")
	}

	// 获取指标快照
	// Get metrics snapshot
	metrics := observability.GetMetricsSnapshot()
	fmt.Printf("Current metrics snapshot: %d metrics collected\n", len(metrics))

	fmt.Println("Alerting demo completed")
}

// 演示数据处理操作的完整追踪
// Demonstrate complete tracing of data processing operations
func demonstrateDataProcessingTrace() {
	fmt.Println("=== Complete Data Processing Trace Demo ===")

	// 使用便捷的追踪操作函数
	// Use convenient trace operation function
	result := observability.TraceDataProcessing(
		"mysql_to_kafka_sync",
		"test_db",
		"users",
		"INSERT",
		func(ctx context.Context) error {
			// 模拟完整的数据处理流程
			// Simulate complete data processing flow
			
			// 1. 读取binlog
			// Read binlog
			span1 := observability.StartSpanFromContext(ctx, "read_binlog")
			time.Sleep(50 * time.Millisecond)
			observability.AddSpanAttributes(span1, map[string]interface{}{
				"binlog_file":     "mysql-bin.000001",
				"binlog_position": 12345,
			})
			observability.FinishSpan(span1)

			// 2. 数据转换
			// Data transformation
			span2 := observability.StartSpanFromContext(ctx, "transform_data")
			time.Sleep(30 * time.Millisecond)
			observability.AddSpanAttributes(span2, map[string]interface{}{
				"transformation_rules": 5,
				"fields_mapped":       12,
			})
			observability.FinishSpan(span2)

			// 3. 发送到Kafka
			// Send to Kafka
			span3 := observability.StartSpanFromContext(ctx, "send_to_kafka")
			time.Sleep(40 * time.Millisecond)
			observability.AddSpanAttributes(span3, map[string]interface{}{
				"kafka_topic":     "gravity.test_db.users",
				"kafka_partition": 0,
				"kafka_offset":    67890,
			})
			observability.FinishSpan(span3)

			return nil
		},
	)

	if result.Error != nil {
		fmt.Printf("Data processing failed: %v\n", result.Error)
	} else {
		fmt.Printf("Data processing completed successfully in %v\n", result.Duration)
	}

	fmt.Println("Complete data processing trace demo completed")
}