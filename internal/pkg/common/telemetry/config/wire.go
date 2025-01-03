package config

import (
	"github.com/google/wire"
)

var WireSet = wire.NewSet(GetTraceConfig, GetMetricConfig)

func GetTraceConfig(c Config) TraceConfig {
	return c.Trace
}

func GetMetricConfig(c Config) MetricConfig {
	return c.Metric
}
