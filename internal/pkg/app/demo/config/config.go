package config

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/go-playground/validator/v10"
	server "github.com/michaljurecko/ch-demo/internal/pkg/app/demo/server/config"
	"github.com/michaljurecko/ch-demo/internal/pkg/common/dataverse/webapi"
	logger "github.com/michaljurecko/ch-demo/internal/pkg/common/log/config"
	redis "github.com/michaljurecko/ch-demo/internal/pkg/common/redis/config"
	telemetry "github.com/michaljurecko/ch-demo/internal/pkg/common/telemetry/config"
)

type Config struct {
	Logger    logger.Config    `embed:"" prefix:"logger-"`
	Telemetry telemetry.Config `embed:"" prefix:"telemetry-"`
	Server    server.Config    `embed:"" prefix:"server-"`
	Model     webapi.Config    `embed:"" prefix:"model-"`
	Redis     redis.Config     `embed:"" prefix:"redis-"`
}

type Decorator func(Config) (Config, error)

func Load(validator *validator.Validate) (Config, error) {
	return doLoad(validator, os.Args[1:], nil)
}

func ForTest(validator *validator.Validate, fn Decorator) (Config, error) {
	return doLoad(validator, nil, fn)
}

func doLoad(validator *validator.Validate, args []string, decorator Decorator, options ...kong.Option) (Config, error) {
	cfg := Config{}

	options = append(
		options,
		kong.Description("Note: Each flag can be set as an ENV."),
		kong.DefaultEnvars("DEMO"),
	)

	parser, err := kong.New(&cfg, options...)
	if err != nil {
		return Config{}, err
	}

	_, err = parser.Parse(args)
	if err != nil {
		return Config{}, err
	}

	// Decorate
	if decorator != nil {
		cfg, err = decorator(cfg)
		if err != nil {
			return Config{}, err
		}
	}

	// Validate
	if err := validator.Struct(cfg); err != nil {
		return Config{}, fmt.Errorf("config is invalid: %w", err)
	}

	return cfg, nil
}
