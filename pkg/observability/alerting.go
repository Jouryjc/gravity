package observability

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
)

// AlertLevel represents the severity level of an alert
type AlertLevel string

const (
	AlertLevelP0 AlertLevel = "P0" // Critical - Immediate action required
	AlertLevelP1 AlertLevel = "P1" // High - Action required within 1 hour
	AlertLevelP2 AlertLevel = "P2" // Medium - Action required within 4 hours
	AlertLevelP3 AlertLevel = "P3" // Low - Action required within 24 hours
)

// AlertRule represents a Prometheus alert rule
type AlertRule struct {
	Name        string            `yaml:"alert"`
	Expr        string            `yaml:"expr"`
	For         string            `yaml:"for"`
	Labels      map[string]string `yaml:"labels"`
	Annotations map[string]string `yaml:"annotations"`
}

// AlertGroup represents a group of related alert rules
type AlertGroup struct {
	Name  string      `yaml:"name"`
	Rules []AlertRule `yaml:"rules"`
}

// AlertConfig represents the complete alert configuration
type AlertConfig struct {
	Groups []AlertGroup `yaml:"groups"`
}

// AlertManager handles alert management and notification
type AlertManager struct {
	prometheusClient v1.API
	config          *AlertingConfig
}

// AlertingConfig holds alerting configuration
type AlertingConfig struct {
	PrometheusURL    string        `yaml:"prometheus_url"`
	EvaluationPeriod time.Duration `yaml:"evaluation_period"`
	WebhookURL       string        `yaml:"webhook_url"`
	SlackChannel     string        `yaml:"slack_channel"`
	EmailRecipients  []string      `yaml:"email_recipients"`
}

// NewAlertManager creates a new alert manager
func NewAlertManager(config *AlertingConfig) (*AlertManager, error) {
	client, err := api.NewClient(api.Config{
		Address: config.PrometheusURL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Prometheus client: %w", err)
	}

	return &AlertManager{
		prometheusClient: v1.NewAPI(client),
		config:          config,
	}, nil
}

// GetGravityAlertRules returns predefined alert rules for Gravity
func GetGravityAlertRules() *AlertConfig {
	return &AlertConfig{
		Groups: []AlertGroup{
			{
				Name: "gravity-critical",
				Rules: []AlertRule{
					{
						Name: "GravityHighErrorRate",
						Expr: "rate(gravity_error_total[5m]) > 0.1",
						For:  "2m",
						Labels: map[string]string{
							"severity": string(AlertLevelP0),
							"service":  "gravity",
						},
						Annotations: map[string]string{
							"summary":     "High error rate detected in Gravity pipeline",
							"description": "Error rate is {{ $value | humanizePercentage }} for pipeline {{ $labels.pipeline }}",
							"runbook":     "https://docs.gravity.io/runbooks/high-error-rate",
						},
					},
					{
						Name: "GravityHighLatency",
						Expr: "histogram_quantile(0.95, rate(gravity_event_time_latency_bucket[5m])) > 300",
						For:  "5m",
						Labels: map[string]string{
							"severity": string(AlertLevelP0),
							"service":  "gravity",
						},
						Annotations: map[string]string{
							"summary":     "High latency detected in Gravity pipeline",
							"description": "95th percentile latency is {{ $value }}s for pipeline {{ $labels.pipeline }}",
							"runbook":     "https://docs.gravity.io/runbooks/high-latency",
						},
					},
					{
						Name: "GravityPipelineDown",
						Expr: "up{job=\"gravity\"} == 0",
						For:  "1m",
						Labels: map[string]string{
							"severity": string(AlertLevelP0),
							"service":  "gravity",
						},
						Annotations: map[string]string{
							"summary":     "Gravity pipeline is down",
							"description": "Gravity pipeline {{ $labels.instance }} has been down for more than 1 minute",
							"runbook":     "https://docs.gravity.io/runbooks/pipeline-down",
						},
					},
				},
			},
			{
				Name: "gravity-database",
				Rules: []AlertRule{
					{
						Name: "GravityDatabaseConnectionFailure",
						Expr: "rate(gravity_connection_failure_total[5m]) > 0.05",
						For:  "3m",
						Labels: map[string]string{
							"severity": string(AlertLevelP1),
							"service":  "gravity",
						},
						Annotations: map[string]string{
							"summary":     "High database connection failure rate",
							"description": "Database connection failure rate is {{ $value | humanizePercentage }} for {{ $labels.database_type }} {{ $labels.endpoint }}",
							"runbook":     "https://docs.gravity.io/runbooks/db-connection-failure",
						},
					},
					{
						Name: "GravityHighReplicationLag",
						Expr: "gravity_replication_lag_seconds > 60",
						For:  "5m",
						Labels: map[string]string{
							"severity": string(AlertLevelP1),
							"service":  "gravity",
						},
						Annotations: map[string]string{
							"summary":     "High replication lag detected",
							"description": "Replication lag is {{ $value }}s for {{ $labels.database_type }} {{ $labels.endpoint }}",
							"runbook":     "https://docs.gravity.io/runbooks/replication-lag",
						},
					},
					{
						Name: "GravityBinlogPositionStuck",
						Expr: "increase(gravity_binlog_position[10m]) == 0 and gravity_binlog_position > 0",
						For:  "10m",
						Labels: map[string]string{
							"severity": string(AlertLevelP1),
							"service":  "gravity",
						},
						Annotations: map[string]string{
							"summary":     "MySQL binlog position appears to be stuck",
							"description": "Binlog position has not advanced for 10 minutes on {{ $labels.endpoint }}",
							"runbook":     "https://docs.gravity.io/runbooks/binlog-stuck",
						},
					},
				},
			},
			{
				Name: "gravity-kafka",
				Rules: []AlertRule{
					{
						Name: "GravityKafkaProduceLatency",
						Expr: "histogram_quantile(0.95, rate(gravity_kafka_produce_latency_bucket[5m])) > 1000",
						For:  "5m",
						Labels: map[string]string{
							"severity": string(AlertLevelP2),
							"service":  "gravity",
						},
						Annotations: map[string]string{
							"summary":     "High Kafka produce latency",
							"description": "95th percentile Kafka produce latency is {{ $value }}ms for topic {{ $labels.topic }}",
							"runbook":     "https://docs.gravity.io/runbooks/kafka-latency",
						},
					},
					{
						Name: "GravityKafkaMessageSizeLarge",
						Expr: "histogram_quantile(0.95, rate(gravity_kafka_message_size_bucket[5m])) > 1048576",
						For:  "10m",
						Labels: map[string]string{
							"severity": string(AlertLevelP2),
							"service":  "gravity",
						},
						Annotations: map[string]string{
							"summary":     "Large Kafka message size detected",
							"description": "95th percentile message size is {{ $value | humanizeBytes }} for topic {{ $labels.topic }}",
							"runbook":     "https://docs.gravity.io/runbooks/large-messages",
						},
					},
				},
			},
			{
				Name: "gravity-data-quality",
				Rules: []AlertRule{
					{
						Name: "GravityDataInconsistency",
						Expr: "gravity_data_consistency_ratio < 0.99",
						For:  "5m",
						Labels: map[string]string{
							"severity": string(AlertLevelP1),
							"service":  "gravity",
						},
						Annotations: map[string]string{
							"summary":     "Data consistency issue detected",
							"description": "Data consistency ratio is {{ $value | humanizePercentage }} for {{ $labels.schema }}.{{ $labels.table }}",
							"runbook":     "https://docs.gravity.io/runbooks/data-consistency",
						},
					},
					{
						Name: "GravityHighDataSkipRate",
						Expr: "rate(gravity_data_skip_total[5m]) > 0.01",
						For:  "5m",
						Labels: map[string]string{
							"severity": string(AlertLevelP2),
							"service":  "gravity",
						},
						Annotations: map[string]string{
							"summary":     "High data skip rate detected",
							"description": "Data skip rate is {{ $value | humanizePercentage }} for {{ $labels.schema }}.{{ $labels.table }}",
							"runbook":     "https://docs.gravity.io/runbooks/data-skip",
						},
					},
					{
						Name: "GravityDataCorruption",
						Expr: "increase(gravity_data_corruption_total[1h]) > 0",
						For:  "0m",
						Labels: map[string]string{
							"severity": string(AlertLevelP0),
							"service":  "gravity",
						},
						Annotations: map[string]string{
							"summary":     "Data corruption detected",
							"description": "Data corruption detected in {{ $labels.schema }}.{{ $labels.table }}",
							"runbook":     "https://docs.gravity.io/runbooks/data-corruption",
						},
					},
				},
			},
			{
				Name: "gravity-performance",
				Rules: []AlertRule{
					{
						Name: "GravityHighMemoryUsage",
						Expr: "gravity_memory_usage_bytes / gravity_memory_limit_bytes > 0.8",
						For:  "10m",
						Labels: map[string]string{
							"severity": string(AlertLevelP2),
							"service":  "gravity",
						},
						Annotations: map[string]string{
							"summary":     "High memory usage detected",
							"description": "Memory usage is {{ $value | humanizePercentage }} for component {{ $labels.component }}",
							"runbook":     "https://docs.gravity.io/runbooks/high-memory",
						},
					},
					{
						Name: "GravityHighCPUUsage",
						Expr: "gravity_cpu_usage_percent > 80",
						For:  "15m",
						Labels: map[string]string{
							"severity": string(AlertLevelP3),
							"service":  "gravity",
						},
						Annotations: map[string]string{
							"summary":     "High CPU usage detected",
							"description": "CPU usage is {{ $value }}% for component {{ $labels.component }}",
							"runbook":     "https://docs.gravity.io/runbooks/high-cpu",
						},
					},
					{
						Name: "GravityLowThroughput",
						Expr: "gravity_throughput_rows_per_second < 100",
						For:  "10m",
						Labels: map[string]string{
							"severity": string(AlertLevelP3),
							"service":  "gravity",
						},
						Annotations: map[string]string{
							"summary":     "Low throughput detected",
							"description": "Throughput is {{ $value }} rows/sec for {{ $labels.schema }}.{{ $labels.table }}",
							"runbook":     "https://docs.gravity.io/runbooks/low-throughput",
						},
					},
				},
			},
		},
	}
}

// EvaluateAlerts evaluates alert conditions and returns active alerts
func (am *AlertManager) EvaluateAlerts(ctx context.Context) ([]ActiveAlert, error) {
	alertRules := GetGravityAlertRules()
	var activeAlerts []ActiveAlert

	for _, group := range alertRules.Groups {
		for _, rule := range group.Rules {
			result, _, err := am.prometheusClient.Query(ctx, rule.Expr, time.Now())
			if err != nil {
				return nil, fmt.Errorf("failed to evaluate alert %s: %w", rule.Name, err)
			}

			// Check if alert condition is met
			if am.isAlertConditionMet(result) {
				activeAlerts = append(activeAlerts, ActiveAlert{
					Name:        rule.Name,
					Severity:    AlertLevel(rule.Labels["severity"]),
					Expression:  rule.Expr,
					Labels:      rule.Labels,
					Annotations: rule.Annotations,
					TriggeredAt: time.Now(),
				})
			}
		}
	}

	return activeAlerts, nil
}

// ActiveAlert represents an active alert
type ActiveAlert struct {
	Name        string            `json:"name"`
	Severity    AlertLevel        `json:"severity"`
	Expression  string            `json:"expression"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
	TriggeredAt time.Time         `json:"triggered_at"`
}

// isAlertConditionMet checks if the query result indicates an alert condition
func (am *AlertManager) isAlertConditionMet(result interface{}) bool {
	// This is a simplified implementation
	// In practice, you would need to parse the Prometheus query result
	// and check if any values exceed the threshold
	return false // Placeholder implementation
}

// SendAlert sends an alert notification
func (am *AlertManager) SendAlert(ctx context.Context, alert ActiveAlert) error {
	// Implementation would depend on your notification system
	// This could include Slack, email, PagerDuty, etc.
	message := am.formatAlertMessage(alert)
	
	// Log the alert
	if OutputLogger != nil {
		OutputLogger.Error(ctx, "Alert triggered", WithContext().
			WithOperation("alert").
			WithMetadata("alert_name", alert.Name).
			WithMetadata("severity", string(alert.Severity)).
			WithMetadata("message", message))
	}

	return nil
}

// formatAlertMessage formats an alert for notification
func (am *AlertManager) formatAlertMessage(alert ActiveAlert) string {
	var builder strings.Builder
	
	builder.WriteString(fmt.Sprintf("🚨 **%s Alert: %s**\n", alert.Severity, alert.Name))
	builder.WriteString(fmt.Sprintf("**Summary:** %s\n", alert.Annotations["summary"]))
	builder.WriteString(fmt.Sprintf("**Description:** %s\n", alert.Annotations["description"]))
	builder.WriteString(fmt.Sprintf("**Triggered At:** %s\n", alert.TriggeredAt.Format(time.RFC3339)))
	
	if runbook, exists := alert.Annotations["runbook"]; exists {
		builder.WriteString(fmt.Sprintf("**Runbook:** %s\n", runbook))
	}
	
	builder.WriteString("**Labels:**\n")
	for k, v := range alert.Labels {
		builder.WriteString(fmt.Sprintf("  - %s: %s\n", k, v))
	}
	
	return builder.String()
}

// GetAlertSummary returns a summary of alert statistics
func (am *AlertManager) GetAlertSummary(ctx context.Context) (*AlertSummary, error) {
	activeAlerts, err := am.EvaluateAlerts(ctx)
	if err != nil {
		return nil, err
	}

	summary := &AlertSummary{
		Total: len(activeAlerts),
		ByLevel: make(map[AlertLevel]int),
	}

	for _, alert := range activeAlerts {
		summary.ByLevel[alert.Severity]++
	}

	return summary, nil
}

// AlertSummary provides statistics about active alerts
type AlertSummary struct {
	Total   int                   `json:"total"`
	ByLevel map[AlertLevel]int    `json:"by_level"`
}

// Helper functions for common alert patterns

// CheckErrorRateThreshold checks if error rate exceeds threshold
func CheckErrorRateThreshold(pipeline, component string, threshold float64) bool {
	// This would query Prometheus for the actual error rate
	// For now, return false as placeholder
	return false
}

// CheckLatencyThreshold checks if latency exceeds threshold
func CheckLatencyThreshold(pipeline, component string, thresholdMs float64) bool {
	// This would query Prometheus for the actual latency
	// For now, return false as placeholder
	return false
}

// CheckDataConsistency checks data consistency ratio
func CheckDataConsistency(pipeline, schema, table string, threshold float64) bool {
	// This would query Prometheus for the actual consistency ratio
	// For now, return false as placeholder
	return false
}

// Alert notification templates
const (
	SlackAlertTemplate = `{
		"text": "Gravity Alert: {{ .Name }}",
		"attachments": [
			{
				"color": "{{ if eq .Severity \"P0\" }}danger{{ else if eq .Severity \"P1\" }}warning{{ else }}good{{ end }}",
				"fields": [
					{
						"title": "Severity",
						"value": "{{ .Severity }}",
						"short": true
					},
					{
						"title": "Pipeline",
						"value": "{{ .Labels.pipeline }}",
						"short": true
					},
					{
						"title": "Description",
						"value": "{{ .Annotations.description }}",
						"short": false
					}
				]
			}
		]
	}`

	EmailAlertTemplate = `
	Subject: [{{ .Severity }}] Gravity Alert: {{ .Name }}
	
	Dear Operations Team,
	
	An alert has been triggered in the Gravity data pipeline:
	
	Alert Name: {{ .Name }}
	Severity: {{ .Severity }}
	Triggered At: {{ .TriggeredAt.Format "2006-01-02 15:04:05 UTC" }}
	
	Summary: {{ .Annotations.summary }}
	Description: {{ .Annotations.description }}
	
	{{ if .Annotations.runbook }}Runbook: {{ .Annotations.runbook }}{{ end }}
	
	Labels:
	{{ range $key, $value := .Labels }}- {{ $key }}: {{ $value }}
	{{ end }}
	
	Please investigate and take appropriate action.
	
	Best regards,
	Gravity Monitoring System
	`
)