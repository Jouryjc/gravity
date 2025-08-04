package observability

import (
	"encoding/json"
	"fmt"
	"time"
)

// DashboardConfig represents a Grafana dashboard configuration
type DashboardConfig struct {
	Dashboard Dashboard `json:"dashboard"`
	FolderID  int       `json:"folderId,omitempty"`
	Overwrite bool      `json:"overwrite"`
}

// Dashboard represents a Grafana dashboard
type Dashboard struct {
	ID              int         `json:"id,omitempty"`
	UID             string      `json:"uid,omitempty"`
	Title           string      `json:"title"`
	Description     string      `json:"description,omitempty"`
	Tags            []string    `json:"tags,omitempty"`
	Timezone        string      `json:"timezone"`
	Editable        bool        `json:"editable"`
	GraphTooltip    int         `json:"graphTooltip"`
	Time            TimeRange   `json:"time"`
	Timepicker      Timepicker  `json:"timepicker"`
	Refresh         string      `json:"refresh"`
	SchemaVersion   int         `json:"schemaVersion"`
	Version         int         `json:"version"`
	Panels          []Panel     `json:"panels"`
	Templating      Templating  `json:"templating"`
	Annotations     Annotations `json:"annotations"`
}

// TimeRange represents the time range for a dashboard
type TimeRange struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// Timepicker represents the timepicker configuration
type Timepicker struct {
	RefreshIntervals []string `json:"refresh_intervals"`
	TimeOptions      []string `json:"time_options"`
}

// Panel represents a dashboard panel
type Panel struct {
	ID          int                    `json:"id"`
	Title       string                 `json:"title"`
	Type        string                 `json:"type"`
	GridPos     GridPos                `json:"gridPos"`
	Targets     []Target               `json:"targets,omitempty"`
	FieldConfig FieldConfig            `json:"fieldConfig,omitempty"`
	Options     map[string]interface{} `json:"options,omitempty"`
	Alert       *Alert                 `json:"alert,omitempty"`
	Datasource  *Datasource            `json:"datasource,omitempty"`
}

// GridPos represents panel position and size
type GridPos struct {
	H int `json:"h"`
	W int `json:"w"`
	X int `json:"x"`
	Y int `json:"y"`
}

// Target represents a query target
type Target struct {
	Expr           string `json:"expr,omitempty"`
	Format         string `json:"format,omitempty"`
	IntervalFactor int    `json:"intervalFactor,omitempty"`
	LegendFormat   string `json:"legendFormat,omitempty"`
	RefID          string `json:"refId"`
	Step           int    `json:"step,omitempty"`
}

// FieldConfig represents field configuration
type FieldConfig struct {
	Defaults  FieldDefaults           `json:"defaults"`
	Overrides []map[string]interface{} `json:"overrides,omitempty"`
}

// FieldDefaults represents default field settings
type FieldDefaults struct {
	Color       map[string]interface{} `json:"color,omitempty"`
	Custom      map[string]interface{} `json:"custom,omitempty"`
	Mappings    []interface{}          `json:"mappings,omitempty"`
	Thresholds  Thresholds             `json:"thresholds,omitempty"`
	Unit        string                 `json:"unit,omitempty"`
	Min         *float64               `json:"min,omitempty"`
	Max         *float64               `json:"max,omitempty"`
	Decimals    *int                   `json:"decimals,omitempty"`
	DisplayName string                 `json:"displayName,omitempty"`
}

// Thresholds represents threshold configuration
type Thresholds struct {
	Mode  string      `json:"mode"`
	Steps []Threshold `json:"steps"`
}

// Threshold represents a single threshold
type Threshold struct {
	Color string   `json:"color"`
	Value *float64 `json:"value"`
}

// Alert represents panel alert configuration
type Alert struct {
	Conditions      []AlertCondition `json:"conditions"`
	ExecutionErrorState string       `json:"executionErrorState"`
	For             string           `json:"for"`
	Frequency       string           `json:"frequency"`
	Handler         int              `json:"handler"`
	Name            string           `json:"name"`
	NoDataState     string           `json:"noDataState"`
	Notifications   []interface{}    `json:"notifications"`
}

// AlertCondition represents an alert condition
type AlertCondition struct {
	Evaluator   Evaluator   `json:"evaluator"`
	Operator    Operator    `json:"operator"`
	Query       QueryRef    `json:"query"`
	Reducer     Reducer     `json:"reducer"`
	Type        string      `json:"type"`
}

// Evaluator represents alert evaluator
type Evaluator struct {
	Params []float64 `json:"params"`
	Type   string    `json:"type"`
}

// Operator represents alert operator
type Operator struct {
	Type string `json:"type"`
}

// QueryRef represents query reference
type QueryRef struct {
	Params []string `json:"params"`
}

// Reducer represents alert reducer
type Reducer struct {
	Params []interface{} `json:"params"`
	Type   string        `json:"type"`
}

// Datasource represents a data source
type Datasource struct {
	Type string `json:"type"`
	UID  string `json:"uid"`
}

// Templating represents dashboard templating
type Templating struct {
	List []Template `json:"list"`
}

// Template represents a dashboard template variable
type Template struct {
	AllValue    string                 `json:"allValue,omitempty"`
	Current     map[string]interface{} `json:"current"`
	Datasource  *Datasource            `json:"datasource,omitempty"`
	Definition  string                 `json:"definition,omitempty"`
	Hide        int                    `json:"hide"`
	IncludeAll  bool                   `json:"includeAll"`
	Label       string                 `json:"label,omitempty"`
	Multi       bool                   `json:"multi"`
	Name        string                 `json:"name"`
	Options     []TemplateOption       `json:"options"`
	Query       string                 `json:"query"`
	Refresh     int                    `json:"refresh"`
	Regex       string                 `json:"regex,omitempty"`
	SkipUrlSync bool                   `json:"skipUrlSync"`
	Sort        int                    `json:"sort"`
	TagValuesQuery string              `json:"tagValuesQuery,omitempty"`
	Tags        []interface{}          `json:"tags"`
	TagsQuery   string                 `json:"tagsQuery,omitempty"`
	Type        string                 `json:"type"`
	UseTags     bool                   `json:"useTags"`
}

// TemplateOption represents a template option
type TemplateOption struct {
	Selected bool   `json:"selected"`
	Text     string `json:"text"`
	Value    string `json:"value"`
}

// Annotations represents dashboard annotations
type Annotations struct {
	List []Annotation `json:"list"`
}

// Annotation represents a dashboard annotation
type Annotation struct {
	BuiltIn    int         `json:"builtIn"`
	Datasource *Datasource `json:"datasource"`
	Enable     bool        `json:"enable"`
	Hide       bool        `json:"hide"`
	IconColor  string      `json:"iconColor"`
	Name       string      `json:"name"`
	Type       string      `json:"type"`
}

// GetTracingDashboard returns a dashboard for distributed tracing
func GetTracingDashboard() *DashboardConfig {
	return &DashboardConfig{
		Dashboard: Dashboard{
			UID:         "gravity-tracing",
			Title:       "Gravity - Distributed Tracing",
			Description: "Distributed tracing overview for Gravity data pipeline",
			Tags:        []string{"gravity", "tracing", "observability"},
			Timezone:    "browser",
			Editable:    true,
			GraphTooltip: 0,
			Time: TimeRange{
				From: "now-1h",
				To:   "now",
			},
			Timepicker: Timepicker{
				RefreshIntervals: []string{"5s", "10s", "30s", "1m", "5m", "15m", "30m", "1h", "2h", "1d"},
				TimeOptions:      []string{"5m", "15m", "1h", "6h", "12h", "24h", "2d", "7d", "30d"},
			},
			Refresh:       "30s",
			SchemaVersion: 27,
			Version:       1,
			Panels: []Panel{
				{
					ID:    1,
					Title: "Trace Count by Pipeline",
					Type:  "stat",
					GridPos: GridPos{H: 8, W: 12, X: 0, Y: 0},
					Targets: []Target{
						{
							Expr:         "sum by (pipeline) (rate(gravity_traces_total[5m]))",
							LegendFormat: "{{pipeline}}",
							RefID:        "A",
						},
					},
					FieldConfig: FieldConfig{
						Defaults: FieldDefaults{
							Unit: "reqps",
						},
					},
				},
				{
					ID:    2,
					Title: "Average Trace Duration",
					Type:  "timeseries",
					GridPos: GridPos{H: 8, W: 12, X: 12, Y: 0},
					Targets: []Target{
						{
							Expr:         "histogram_quantile(0.50, rate(gravity_trace_duration_bucket[5m]))",
							LegendFormat: "p50",
							RefID:        "A",
						},
						{
							Expr:         "histogram_quantile(0.95, rate(gravity_trace_duration_bucket[5m]))",
							LegendFormat: "p95",
							RefID:        "B",
						},
						{
							Expr:         "histogram_quantile(0.99, rate(gravity_trace_duration_bucket[5m]))",
							LegendFormat: "p99",
							RefID:        "C",
						},
					},
					FieldConfig: FieldConfig{
						Defaults: FieldDefaults{
							Unit: "ms",
						},
					},
				},
				{
					ID:    3,
					Title: "Error Rate by Component",
					Type:  "timeseries",
					GridPos: GridPos{H: 8, W: 24, X: 0, Y: 8},
					Targets: []Target{
						{
							Expr:         "rate(gravity_error_total[5m])",
							LegendFormat: "{{pipeline}}.{{component}}",
							RefID:        "A",
						},
					},
					FieldConfig: FieldConfig{
						Defaults: FieldDefaults{
							Unit: "reqps",
						},
					},
				},
			},
			Templating: Templating{
				List: []Template{
					{
						Name:        "pipeline",
						Type:        "query",
						Query:       "label_values(gravity_traces_total, pipeline)",
						Refresh:     1,
						IncludeAll:  true,
						Multi:       true,
						Current:     map[string]interface{}{"text": "All", "value": "$__all"},
						Options:     []TemplateOption{},
						Hide:        0,
						SkipUrlSync: false,
						Sort:        1,
						Tags:        []interface{}{},
						UseTags:     false,
					},
				},
			},
			Annotations: Annotations{
				List: []Annotation{
					{
						BuiltIn:   1,
						Enable:    true,
						Hide:      true,
						IconColor: "rgba(0, 211, 255, 1)",
						Name:      "Annotations & Alerts",
						Type:      "dashboard",
					},
				},
			},
		},
		Overwrite: true,
	}
}

// GetErrorAnalysisDashboard returns a dashboard for error analysis
func GetErrorAnalysisDashboard() *DashboardConfig {
	return &DashboardConfig{
		Dashboard: Dashboard{
			UID:         "gravity-errors",
			Title:       "Gravity - Error Analysis",
			Description: "Error analysis and troubleshooting for Gravity data pipeline",
			Tags:        []string{"gravity", "errors", "troubleshooting"},
			Timezone:    "browser",
			Editable:    true,
			GraphTooltip: 0,
			Time: TimeRange{
				From: "now-1h",
				To:   "now",
			},
			Refresh:       "30s",
			SchemaVersion: 27,
			Version:       1,
			Panels: []Panel{
				{
					ID:    1,
					Title: "Error Rate Overview",
					Type:  "stat",
					GridPos: GridPos{H: 8, W: 6, X: 0, Y: 0},
					Targets: []Target{
						{
							Expr:         "sum(rate(gravity_error_total[5m]))",
							LegendFormat: "Total Errors/sec",
							RefID:        "A",
						},
					},
					FieldConfig: FieldConfig{
						Defaults: FieldDefaults{
							Unit: "reqps",
							Thresholds: Thresholds{
								Mode: "absolute",
								Steps: []Threshold{
									{Color: "green", Value: nil},
									{Color: "yellow", Value: func() *float64 { v := 0.1; return &v }()},
									{Color: "red", Value: func() *float64 { v := 1.0; return &v }()},
								},
							},
						},
					},
				},
				{
					ID:    2,
					Title: "Connection Failures",
					Type:  "stat",
					GridPos: GridPos{H: 8, W: 6, X: 6, Y: 0},
					Targets: []Target{
						{
							Expr:         "sum(rate(gravity_connection_failure_total[5m]))",
							LegendFormat: "Connection Failures/sec",
							RefID:        "A",
						},
					},
					FieldConfig: FieldConfig{
						Defaults: FieldDefaults{
							Unit: "reqps",
						},
					},
				},
				{
					ID:    3,
					Title: "Retry Rate",
					Type:  "stat",
					GridPos: GridPos{H: 8, W: 6, X: 12, Y: 0},
					Targets: []Target{
						{
							Expr:         "sum(rate(gravity_retry_total[5m]))",
							LegendFormat: "Retries/sec",
							RefID:        "A",
						},
					},
					FieldConfig: FieldConfig{
						Defaults: FieldDefaults{
							Unit: "reqps",
						},
					},
				},
				{
					ID:    4,
					Title: "Data Corruption Events",
					Type:  "stat",
					GridPos: GridPos{H: 8, W: 6, X: 18, Y: 0},
					Targets: []Target{
						{
							Expr:         "sum(increase(gravity_data_corruption_total[1h]))",
							LegendFormat: "Corruption Events (1h)",
							RefID:        "A",
						},
					},
					FieldConfig: FieldConfig{
						Defaults: FieldDefaults{
							Unit: "short",
							Thresholds: Thresholds{
								Mode: "absolute",
								Steps: []Threshold{
									{Color: "green", Value: func() *float64 { v := 0.0; return &v }()},
									{Color: "red", Value: func() *float64 { v := 1.0; return &v }()},
								},
							},
						},
					},
				},
				{
					ID:    5,
					Title: "Error Rate by Component",
					Type:  "timeseries",
					GridPos: GridPos{H: 8, W: 24, X: 0, Y: 8},
					Targets: []Target{
						{
							Expr:         "rate(gravity_error_total[5m])",
							LegendFormat: "{{pipeline}}.{{component}} - {{error_type}}",
							RefID:        "A",
						},
					},
				},
			},
		},
		Overwrite: true,
	}
}

// GetDataQualityDashboard returns a dashboard for data quality monitoring
func GetDataQualityDashboard() *DashboardConfig {
	return &DashboardConfig{
		Dashboard: Dashboard{
			UID:         "gravity-data-quality",
			Title:       "Gravity - Data Quality",
			Description: "Data quality and consistency monitoring for Gravity pipeline",
			Tags:        []string{"gravity", "data-quality", "consistency"},
			Timezone:    "browser",
			Editable:    true,
			GraphTooltip: 0,
			Time: TimeRange{
				From: "now-6h",
				To:   "now",
			},
			Refresh:       "1m",
			SchemaVersion: 27,
			Version:       1,
			Panels: []Panel{
				{
					ID:    1,
					Title: "Data Consistency Ratio",
					Type:  "gauge",
					GridPos: GridPos{H: 8, W: 12, X: 0, Y: 0},
					Targets: []Target{
						{
							Expr:         "gravity_data_consistency_ratio",
							LegendFormat: "{{schema}}.{{table}}",
							RefID:        "A",
						},
					},
					FieldConfig: FieldConfig{
						Defaults: FieldDefaults{
							Unit: "percentunit",
							Min:  func() *float64 { v := 0.0; return &v }(),
							Max:  func() *float64 { v := 1.0; return &v }(),
							Thresholds: Thresholds{
								Mode: "absolute",
								Steps: []Threshold{
									{Color: "red", Value: func() *float64 { v := 0.0; return &v }()},
									{Color: "yellow", Value: func() *float64 { v := 0.95; return &v }()},
									{Color: "green", Value: func() *float64 { v := 0.99; return &v }()},
								},
							},
						},
					},
				},
				{
					ID:    2,
					Title: "Data Lag Distribution",
					Type:  "timeseries",
					GridPos: GridPos{H: 8, W: 12, X: 12, Y: 0},
					Targets: []Target{
						{
							Expr:         "histogram_quantile(0.50, rate(gravity_data_lag_bucket[5m]))",
							LegendFormat: "p50",
							RefID:        "A",
						},
						{
							Expr:         "histogram_quantile(0.95, rate(gravity_data_lag_bucket[5m]))",
							LegendFormat: "p95",
							RefID:        "B",
						},
						{
							Expr:         "histogram_quantile(0.99, rate(gravity_data_lag_bucket[5m]))",
							LegendFormat: "p99",
							RefID:        "C",
						},
					},
					FieldConfig: FieldConfig{
						Defaults: FieldDefaults{
							Unit: "s",
						},
					},
				},
				{
					ID:    3,
					Title: "Throughput by Table",
					Type:  "timeseries",
					GridPos: GridPos{H: 8, W: 24, X: 0, Y: 8},
					Targets: []Target{
						{
							Expr:         "gravity_throughput_rows_per_second",
							LegendFormat: "{{schema}}.{{table}} - {{operation}}",
							RefID:        "A",
						},
					},
					FieldConfig: FieldConfig{
						Defaults: FieldDefaults{
							Unit: "rps",
						},
					},
				},
			},
		},
		Overwrite: true,
	}
}

// GetResourceMonitoringDashboard returns a dashboard for resource monitoring
func GetResourceMonitoringDashboard() *DashboardConfig {
	return &DashboardConfig{
		Dashboard: Dashboard{
			UID:         "gravity-resources",
			Title:       "Gravity - Resource Monitoring",
			Description: "Resource utilization monitoring for Gravity components",
			Tags:        []string{"gravity", "resources", "performance"},
			Timezone:    "browser",
			Editable:    true,
			GraphTooltip: 0,
			Time: TimeRange{
				From: "now-1h",
				To:   "now",
			},
			Refresh:       "30s",
			SchemaVersion: 27,
			Version:       1,
			Panels: []Panel{
				{
					ID:    1,
					Title: "Memory Usage by Component",
					Type:  "timeseries",
					GridPos: GridPos{H: 8, W: 12, X: 0, Y: 0},
					Targets: []Target{
						{
							Expr:         "gravity_memory_usage_bytes",
							LegendFormat: "{{component}}",
							RefID:        "A",
						},
					},
					FieldConfig: FieldConfig{
						Defaults: FieldDefaults{
							Unit: "bytes",
						},
					},
				},
				{
					ID:    2,
					Title: "CPU Usage by Component",
					Type:  "timeseries",
					GridPos: GridPos{H: 8, W: 12, X: 12, Y: 0},
					Targets: []Target{
						{
							Expr:         "gravity_cpu_usage_percent",
							LegendFormat: "{{component}}",
							RefID:        "A",
						},
					},
					FieldConfig: FieldConfig{
						Defaults: FieldDefaults{
							Unit: "percent",
							Max:  func() *float64 { v := 100.0; return &v }(),
						},
					},
				},
				{
					ID:    3,
					Title: "Database Connections",
					Type:  "timeseries",
					GridPos: GridPos{H: 8, W: 12, X: 0, Y: 8},
					Targets: []Target{
						{
							Expr:         "gravity_database_connections",
							LegendFormat: "{{database_type}} - {{endpoint}}",
							RefID:        "A",
						},
					},
					FieldConfig: FieldConfig{
						Defaults: FieldDefaults{
							Unit: "short",
						},
					},
				},
				{
					ID:    4,
					Title: "Kafka Partition Offsets",
					Type:  "timeseries",
					GridPos: GridPos{H: 8, W: 12, X: 12, Y: 8},
					Targets: []Target{
						{
							Expr:         "gravity_kafka_partition_offset",
							LegendFormat: "{{topic}}.{{partition}}",
							RefID:        "A",
						},
					},
					FieldConfig: FieldConfig{
						Defaults: FieldDefaults{
							Unit: "short",
						},
					},
				},
			},
		},
		Overwrite: true,
	}
}

// TopologyNode represents a node in the data flow topology
type TopologyNode struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Type        string                 `json:"type"` // input, filter, scheduler, output
	Status      string                 `json:"status"` // healthy, warning, error
	Metrics     map[string]interface{} `json:"metrics"`
	Position    Position               `json:"position"`
	Connections []string               `json:"connections"`
}

// Position represents node position in topology
type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// TopologyEdge represents a connection between nodes
type TopologyEdge struct {
	Source string                 `json:"source"`
	Target string                 `json:"target"`
	Label  string                 `json:"label,omitempty"`
	Metrics map[string]interface{} `json:"metrics,omitempty"`
}

// TopologyGraph represents the complete data flow topology
type TopologyGraph struct {
	Nodes []TopologyNode `json:"nodes"`
	Edges []TopologyEdge `json:"edges"`
	UpdatedAt time.Time   `json:"updated_at"`
}

// GetDataFlowTopology returns the current data flow topology
func GetDataFlowTopology(pipelineName string) *TopologyGraph {
	// This would be populated with real-time data from the pipeline
	return &TopologyGraph{
		Nodes: []TopologyNode{
			{
				ID:   "mysql-input",
				Name: "MySQL Input",
				Type: "input",
				Status: "healthy",
				Metrics: map[string]interface{}{
					"throughput": 1000.0,
					"latency":    50.0,
					"errors":     0,
				},
				Position: Position{X: 100, Y: 100},
				Connections: []string{"filter-1"},
			},
			{
				ID:   "mongo-input",
				Name: "MongoDB Input",
				Type: "input",
				Status: "healthy",
				Metrics: map[string]interface{}{
					"throughput": 800.0,
					"latency":    45.0,
					"errors":     0,
				},
				Position: Position{X: 100, Y: 300},
				Connections: []string{"filter-1"},
			},
			{
				ID:   "filter-1",
				Name: "Data Filter",
				Type: "filter",
				Status: "healthy",
				Metrics: map[string]interface{}{
					"throughput": 1800.0,
					"latency":    10.0,
					"errors":     0,
				},
				Position: Position{X: 300, Y: 200},
				Connections: []string{"scheduler-1"},
			},
			{
				ID:   "scheduler-1",
				Name: "Scheduler",
				Type: "scheduler",
				Status: "healthy",
				Metrics: map[string]interface{}{
					"throughput": 1800.0,
					"latency":    5.0,
					"errors":     0,
				},
				Position: Position{X: 500, Y: 200},
				Connections: []string{"kafka-output"},
			},
			{
				ID:   "kafka-output",
				Name: "Kafka Output",
				Type: "output",
				Status: "healthy",
				Metrics: map[string]interface{}{
					"throughput": 1800.0,
					"latency":    20.0,
					"errors":     0,
				},
				Position: Position{X: 700, Y: 200},
				Connections: []string{},
			},
		},
		Edges: []TopologyEdge{
			{
				Source: "mysql-input",
				Target: "filter-1",
				Label:  "MySQL Data",
				Metrics: map[string]interface{}{
					"throughput": 1000.0,
					"latency":    2.0,
				},
			},
			{
				Source: "mongo-input",
				Target: "filter-1",
				Label:  "MongoDB Data",
				Metrics: map[string]interface{}{
					"throughput": 800.0,
					"latency":    2.0,
				},
			},
			{
				Source: "filter-1",
				Target: "scheduler-1",
				Label:  "Filtered Data",
				Metrics: map[string]interface{}{
					"throughput": 1800.0,
					"latency":    1.0,
				},
			},
			{
				Source: "scheduler-1",
				Target: "kafka-output",
				Label:  "Scheduled Data",
				Metrics: map[string]interface{}{
					"throughput": 1800.0,
					"latency":    1.0,
				},
			},
		},
		UpdatedAt: time.Now(),
	}
}

// ExportDashboardsToJSON exports all dashboards to JSON format
func ExportDashboardsToJSON() (map[string]string, error) {
	dashboards := map[string]*DashboardConfig{
		"tracing":      GetTracingDashboard(),
		"errors":       GetErrorAnalysisDashboard(),
		"data-quality": GetDataQualityDashboard(),
		"resources":    GetResourceMonitoringDashboard(),
	}

	result := make(map[string]string)
	for name, dashboard := range dashboards {
		jsonData, err := json.MarshalIndent(dashboard, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("failed to marshal dashboard %s: %w", name, err)
		}
		result[name] = string(jsonData)
	}

	return result, nil
}