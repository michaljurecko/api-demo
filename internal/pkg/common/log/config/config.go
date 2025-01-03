package config

type (
	ExporterType string
)

const (
	NoneExporter ExporterType = "none"
	HTTPExporter ExporterType = "http"
)

type Config struct {
	Exporter          ExporterType `default:"none"                          help:"Type of log exporter, one of [ none, http ]." validate:"oneof=none http"`
	HTTPEndpointURL   string       `default:"http://localhost:4318/v1/logs" help:"Endpoint URL for log exporter."               validate:"required_if=Exporter http"`
	HTTPAuthorization string       `default:"Basic ...."                    help:"Authorization for log exporter."              json:"-"                             validate:"required_if=Exporter http"`
}
