package httpclient

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

func New(traceProvider trace.TracerProvider, meterProvider metric.MeterProvider) *http.Client {
	return &http.Client{
		Transport: otelhttp.NewTransport(
			http.DefaultTransport,
			otelhttp.WithTracerProvider(traceProvider),
			otelhttp.WithMeterProvider(meterProvider),
		),
	}
}
