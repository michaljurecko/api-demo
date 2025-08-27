package config

import (
	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	// Export partial configs
	wire.FieldsOf(new(Config), "Logger", "Telemetry", "Server", "Model"),
)
