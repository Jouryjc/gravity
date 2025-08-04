# Gravity 全链路可观测性增强方案

## 1. 现状分析

### 1.1 当前监控能力
Gravity 项目已具备基础监控能力：
- 使用 Prometheus + Grafana 进行监控
- 提供 `/metrics` 端点暴露指标
- 已有基础指标：
  - 输入组件指标：`gravity_input_*`
  - 调度器指标：`gravity_scheduler_*`
  - 输出组件指标：`gravity_output_*`
  - 端到端延迟指标：`gravity_event_time_latency`、`gravity_process_time_latency`
  - 队列长度指标：`gravity_queue_length`

### 1.2 观测能力缺口
1. **链路追踪缺失**：无法追踪单个事件在整个数据流中的完整路径
2. **错误监控不足**：缺乏详细的错误分类和错误率监控
3. **业务指标缺失**：缺乏数据一致性、数据质量等业务层面的监控
4. **告警机制不完善**：缺乏智能告警和故障自愈能力
5. **可视化不够丰富**：缺乏实时数据流拓扑图和详细的性能分析

## 2. 增强方案设计

### 2.1 分布式链路追踪

#### 2.1.1 引入 OpenTelemetry
- 集成 OpenTelemetry Go SDK
- 为每个数据事件生成唯一的 TraceID 和 SpanID
- 在 MySQL binlog、MongoDB oplog 到 Kafka 的完整链路中传递追踪信息

#### 2.1.2 Span 设计
```
Trace: MySQL/MongoDB -> Gravity -> Kafka
├── Span: Input (MySQL Binlog/MongoDB Oplog)
├── Span: Filter Processing
├── Span: Scheduler Processing
└── Span: Output (Kafka)
```

### 2.2 增强指标体系

#### 2.2.1 错误监控指标
```go
// 错误率指标
var ErrorRateCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
    Namespace: "gravity",
    Name:      "error_total",
    Help:      "Total number of errors by component and type",
}, []string{"pipeline", "component", "error_type"})

// 重试指标
var RetryCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
    Namespace: "gravity",
    Name:      "retry_total",
    Help:      "Total number of retries",
}, []string{"pipeline", "component", "reason"})
```

#### 2.2.2 数据质量指标
```go
// 数据一致性检查
var DataConsistencyGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
    Namespace: "gravity",
    Name:      "data_consistency_score",
    Help:      "Data consistency score between source and target",
}, []string{"pipeline", "schema", "table"})

// 数据延迟分布
var DataLagHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
    Namespace: "gravity",
    Name:      "data_lag_seconds",
    Help:      "Data lag between source and target in seconds",
    Buckets:   prometheus.ExponentialBuckets(0.1, 2, 15),
}, []string{"pipeline", "schema", "table"})
```

#### 2.2.3 业务指标
```go
// 数据吞吐量
var ThroughputGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
    Namespace: "gravity",
    Name:      "throughput_rows_per_second",
    Help:      "Data throughput in rows per second",
}, []string{"pipeline", "schema", "table", "operation"})

// 连接健康状态
var ConnectionHealthGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
    Namespace: "gravity",
    Name:      "connection_health",
    Help:      "Connection health status (1=healthy, 0=unhealthy)",
}, []string{"pipeline", "component", "endpoint"})
```

### 2.3 结构化日志增强

#### 2.3.1 统一日志格式
```go
type LogEntry struct {
    Timestamp   time.Time `json:"timestamp"`
    Level       string    `json:"level"`
    Pipeline    string    `json:"pipeline"`
    Component   string    `json:"component"`
    TraceID     string    `json:"trace_id"`
    SpanID      string    `json:"span_id"`
    Message     string    `json:"message"`
    Schema      string    `json:"schema,omitempty"`
    Table       string    `json:"table,omitempty"`
    Operation   string    `json:"operation,omitempty"`
    Duration    int64     `json:"duration_ms,omitempty"`
    Error       string    `json:"error,omitempty"`
    Metadata    map[string]interface{} `json:"metadata,omitempty"`
}
```

### 2.4 智能告警系统

#### 2.4.1 多级告警策略
- **P0 告警**：数据同步中断、连接断开
- **P1 告警**：延迟超过阈值、错误率过高
- **P2 告警**：性能下降、资源使用率高
- **P3 告警**：数据质量问题、配置变更

#### 2.4.2 告警规则示例
```yaml
groups:
- name: gravity.rules
  rules:
  - alert: GravityHighErrorRate
    expr: rate(gravity_error_total[5m]) > 0.1
    for: 2m
    labels:
      severity: critical
    annotations:
      summary: "High error rate in Gravity pipeline {{ $labels.pipeline }}"
      
  - alert: GravityHighLatency
    expr: histogram_quantile(0.95, gravity_event_time_latency_bucket) > 300
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "High latency in Gravity pipeline {{ $labels.pipeline }}"
```

### 2.5 可视化增强

#### 2.5.1 实时拓扑图
- 展示数据流向：MySQL/MongoDB -> Gravity -> Kafka
- 实时显示各组件状态和性能指标
- 支持点击查看详细信息

#### 2.5.2 新增 Grafana Dashboard
1. **链路追踪面板**：显示请求链路和性能瓶颈
2. **错误分析面板**：错误趋势、错误分类、错误热点
3. **数据质量面板**：一致性检查、延迟分析、吞吐量监控
4. **资源监控面板**：CPU、内存、网络、磁盘使用情况

## 3. 实施计划

### 3.1 第一阶段：基础增强（2周）
1. 集成 OpenTelemetry SDK
2. 添加链路追踪到核心组件
3. 增强错误监控指标
4. 优化结构化日志

### 3.2 第二阶段：业务监控（2周）
1. 实现数据质量监控
2. 添加业务指标收集
3. 开发健康检查机制
4. 创建新的 Grafana Dashboard

### 3.3 第三阶段：智能化（2周）
1. 实现智能告警系统
2. 开发故障自愈机制
3. 添加性能优化建议
4. 完善文档和培训材料

## 4. 技术实现要点

### 4.1 性能考虑
- 使用采样策略减少追踪开销
- 异步发送监控数据
- 合理设置指标保留策略

### 4.2 兼容性
- 保持向后兼容
- 支持渐进式部署
- 提供配置开关

### 4.3 扩展性
- 插件化架构支持自定义监控
- 支持多种后端存储
- 预留扩展接口

## 5. 预期收益

1. **故障定位时间减少 80%**：通过链路追踪快速定位问题
2. **系统可用性提升至 99.9%**：通过智能告警和自愈机制
3. **运维效率提升 50%**：通过自动化监控和告警
4. **数据质量可见性 100%**：实时监控数据一致性和延迟

这个方案将显著提升 Gravity 的可观测性，为 MySQL、MongoDB 到 Kafka 的数据同步链路提供全面的监控和故障排查能力。