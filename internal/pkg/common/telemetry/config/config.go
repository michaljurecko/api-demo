package config

type (
	TraceExporterType  string
	MetricExporterType string
)

const (
	NoneMetricExporter MetricExporterType = "none"
	HTTPMetricExporter MetricExporterType = "http"
	NoneTraceExporter  TraceExporterType  = "none"
	HTTPTraceExporter  TraceExporterType  = "http"
)

type Config struct {
	Trace  TraceConfig  `embed:"" prefix:"trace-"`
	Metric MetricConfig `embed:"" prefix:"metric-"`
}

type TraceConfig struct {
	Exporter          TraceExporterType `default:"none"                            help:"Type of telemetry trace exporter, one of [ none, http ]." validate:"oneof=none http"`
	HTTPEndpointURL   string            `default:"http://localhost:4318/v1/traces" help:"Endpoint URL for HTTP traces exporter."                   validate:"required_if=Exporter http"`
	HTTPAuthorization string            `default:"Basic ...."                      help:"Authorization for traces exporter."                       json:"-"                             validate:"required_if=Exporter http"`
}

type MetricConfig struct {
	Exporter          MetricExporterType `default:"none"                             help:"Type of telemetry metric exporter, one of [ none, http ]." validate:"oneof=none http"`
	HTTPEndpointURL   string             `default:"http://localhost:4318/v1/metrics" help:"Endpoint URL for HTTP metrics exporter."                   validate:"required_if=Exporter http"`
	HTTPAuthorization string             `default:"Basic ...."                       help:"Authorization for metrics exporter."                       json:"-"                             validate:"required_if=Exporter http"`
}
