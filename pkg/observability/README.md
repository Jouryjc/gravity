# Gravity 可观测性增强包

本包为 Gravity 项目提供了全面的可观测性能力，包括分布式追踪、增强指标、结构化日志和智能告警。

## 功能特性

### 🔍 分布式追踪
- 基于 OpenTelemetry 的分布式追踪
- 支持 Jaeger 导出
- 自动 TraceID 和 SpanID 传播
- 跨组件的端到端追踪

### 📊 增强指标
- 丰富的 Prometheus 指标
- 错误率、重试率、连接健康度监控
- 数据质量和一致性指标
- 性能和资源使用指标
- Kafka 和数据库特定指标

### 📝 结构化日志
- 统一的日志格式
- 自动集成 TraceID 和 SpanID
- 组件级别的日志分类
- 支持多种日志级别

### 🚨 智能告警
- 基于 Prometheus 的告警规则
- 多级别告警策略 (P0-P3)
- 自动告警检查和通知
- 可配置的告警阈值

### 📈 可视化
- 预定义的 Grafana 仪表板
- 实时拓扑图
- 多维度数据展示

## 快速开始

### 1. 初始化可观测性系统

```go
package main

import (
    "github.com/moiot/gravity/pkg/observability"
)

func main() {
    // 使用默认配置
    config := observability.GetDefaultConfig()
    config.ServiceName = "gravity-pipeline"
    config.TracingConfig.JaegerEndpoint = "http://localhost:14268/api/traces"
    
    err := observability.Initialize(config)
    if err != nil {
        panic(err)
    }
    defer observability.Shutdown()
    
    // 设置管道名称
    observability.SetPipelineName("mysql-to-kafka")
}
```

### 2. 使用分布式追踪

```go
// 创建 span
ctx := context.Background()
span := observability.StartSpanFromContext(ctx, "mysql_binlog_read")
defer observability.FinishSpan(span)

// 添加属性
observability.AddSpanAttributes(span, map[string]interface{}{
    "database.type": "mysql",
    "schema": "test_db",
    "table": "users",
})

// 传播上下文
childCtx := observability.InjectTraceContext(ctx, span)
processData(childCtx)
```

### 3. 记录指标

```go
// 记录错误
observability.RecordMetric("error", map[string]string{
    "component": "input",
    "error_type": "connection_timeout",
}, 1)

// 记录吞吐量
observability.RecordMetric("throughput", map[string]string{
    "schema": "test_db",
    "table": "users",
}, 1500)

// 记录延迟
start := time.Now()
// ... 执行操作
duration := time.Since(start)
observability.RecordMetric("component_latency", map[string]string{
    "component": "filter",
}, duration.Seconds())
```

### 4. 结构化日志

```go
// 带 trace 的日志
observability.LogWithTrace(ctx, "info", "Processing binlog event", map[string]interface{}{
    "component": "input",
    "schema": "test_db",
    "table": "users",
    "operation": "INSERT",
})

// 错误日志
observability.LogWithTrace(ctx, "error", "Database connection failed", map[string]interface{}{
    "component": "input",
    "error": "connection timeout",
    "endpoint": "localhost:3306",
})
```

### 5. 便捷的操作追踪

```go
// 追踪完整的数据处理操作
result := observability.TraceDataProcessing(
    "mysql_to_kafka_sync",
    "test_db",
    "users",
    "INSERT",
    func(ctx context.Context) error {
        // 你的数据处理逻辑
        return processUserData(ctx)
    },
)

if result.Error != nil {
    log.Printf("Processing failed: %v", result.Error)
} else {
    log.Printf("Processing completed in %v", result.Duration)
}
```

## 配置说明

### TracingConfig
```go
type TracingConfig struct {
    Enabled         bool   `yaml:"enabled"`
    JaegerEndpoint  string `yaml:"jaeger_endpoint"`
    SamplingRate    float64 `yaml:"sampling_rate"`
}
```

### MetricsConfig
```go
type MetricsConfig struct {
    Enabled        bool   `yaml:"enabled"`
    PrometheusPort int    `yaml:"prometheus_port"`
    MetricsPath    string `yaml:"metrics_path"`
}
```

### LoggingConfig
```go
type LoggingConfig struct {
    Level      string `yaml:"level"`
    Format     string `yaml:"format"`
    Output     string `yaml:"output"`
    MaxSize    int    `yaml:"max_size"`
    MaxBackups int    `yaml:"max_backups"`
    MaxAge     int    `yaml:"max_age"`
}
```

## 部署配置

### Prometheus 告警规则

将 `configs/prometheus/gravity-alerts.yml` 添加到你的 Prometheus 配置中：

```yaml
rule_files:
  - "gravity-alerts.yml"
```

### Grafana 仪表板

导入以下仪表板 JSON 文件到 Grafana：
- `configs/grafana/dashboards/gravity-tracing.json` - 分布式追踪
- `configs/grafana/dashboards/gravity-error-analysis.json` - 错误分析
- `configs/grafana/dashboards/gravity-data-quality.json` - 数据质量
- `configs/grafana/dashboards/gravity-resource-monitoring.json` - 资源监控

### Jaeger 部署

使用 Docker 快速启动 Jaeger：

```bash
docker run -d --name jaeger \
  -p 16686:16686 \
  -p 14268:14268 \
  jaegertracing/all-in-one:latest
```

## 指标说明

### 错误监控指标
- `gravity_error_total` - 错误总数
- `gravity_retry_total` - 重试总数
- `gravity_connection_failure_total` - 连接失败总数

### 数据质量指标
- `gravity_data_consistency_ratio` - 数据一致性比率
- `gravity_data_lag_bucket` - 数据延迟分布
- `gravity_data_skip_total` - 跳过的数据记录数
- `gravity_data_corruption_total` - 数据损坏事件数

### 业务指标
- `gravity_throughput_rows_per_second` - 吞吐量（行/秒）
- `gravity_connection_health` - 连接健康度
- `gravity_replication_lag_seconds` - 复制延迟
- `gravity_table_size_bytes` - 表大小

### 性能指标
- `gravity_component_latency_bucket` - 组件延迟分布
- `gravity_memory_usage_bytes` - 内存使用量
- `gravity_cpu_usage_percent` - CPU 使用率

## 告警规则

### P0 级别（严重）
- 管道停止运行
- 错误率 > 5%
- 数据库连接完全失败

### P1 级别（高）
- 错误率 > 1%
- 延迟 > 10 秒
- 复制延迟 > 60 秒

### P2 级别（中）
- 错误率 > 0.1%
- 延迟 > 5 秒
- 内存使用 > 2GB

### P3 级别（低）
- 重试率 > 1%
- CPU 使用 > 80%
- 连接数 > 80% 最大值

## 最佳实践

### 1. Span 命名规范
- 使用动词_名词格式：`read_binlog`, `transform_data`, `send_kafka`
- 保持简洁明了
- 使用小写和下划线

### 2. 属性添加
- 添加有意义的业务属性
- 包含足够的上下文信息
- 避免添加敏感信息

### 3. 错误处理
- 记录错误到 span
- 同时记录到日志
- 包含错误类型和详细信息

### 4. 性能考虑
- 合理设置采样率
- 避免在高频路径中创建过多 span
- 定期清理过期的指标

## 故障排查

### 常见问题

1. **Jaeger 连接失败**
   - 检查 Jaeger 服务是否运行
   - 验证端点配置是否正确
   - 检查网络连接

2. **指标未显示**
   - 确认 Prometheus 配置正确
   - 检查指标端点是否可访问
   - 验证指标名称和标签

3. **日志格式问题**
   - 检查日志配置
   - 验证输出路径权限
   - 确认日志级别设置

### 调试模式

启用调试模式获取更多信息：

```go
config := observability.GetDefaultConfig()
config.LoggingConfig.Level = "debug"
```

## 贡献指南

欢迎提交 Issue 和 Pull Request 来改进这个可观测性包。请确保：

1. 遵循现有的代码风格
2. 添加适当的测试
3. 更新相关文档
4. 提供清晰的提交信息

## 许可证

本项目采用与 Gravity 主项目相同的许可证。