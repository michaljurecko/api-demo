package telemetry

import (
	"github.com/google/wire"
	"github.com/michaljurecko/ch-demo/internal/pkg/common/telemetry/config"
)

var WireSet = wire.NewSet(config.WireSet, NewTracerProvider, NewMeterProvider)
