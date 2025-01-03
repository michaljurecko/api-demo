//go:build wireinject
// +build wireinject

package cmd

import (
	"context"
	"github.com/google/wire"
	"github.com/michaljurecko/ch-demo/internal/pkg/app/demo/ares"
	"github.com/michaljurecko/ch-demo/internal/pkg/app/demo/biz/playerbiz"
	"github.com/michaljurecko/ch-demo/internal/pkg/app/demo/config"
	model "github.com/michaljurecko/ch-demo/internal/pkg/app/demo/model/gen"
	"github.com/michaljurecko/ch-demo/internal/pkg/app/demo/server"
	"github.com/michaljurecko/ch-demo/internal/pkg/app/demo/service"
	"github.com/michaljurecko/ch-demo/internal/pkg/common/cachestore"
	"github.com/michaljurecko/ch-demo/internal/pkg/common/dataverse/webapi"
	"github.com/michaljurecko/ch-demo/internal/pkg/common/distlock"
	"github.com/michaljurecko/ch-demo/internal/pkg/common/httpclient"
	"github.com/michaljurecko/ch-demo/internal/pkg/common/log"
	"github.com/michaljurecko/ch-demo/internal/pkg/common/redis"
	"github.com/michaljurecko/ch-demo/internal/pkg/common/shutdown"
	"github.com/michaljurecko/ch-demo/internal/pkg/common/telemetry"
	"github.com/michaljurecko/ch-demo/internal/pkg/common/validate"
)

var set = wire.NewSet(
	log.WireSet,
	shutdown.WireSet,
	telemetry.WireSet,
	validate.WireSet,
	config.WireSet,
	httpclient.WireSet,
	redis.WireSet,
	cachestore.WireSet,
	distlock.WireSet,
	webapi.NewClient,
	ares.NewClient,
	model.NewRepository,
	playerbiz.WireSet,
	service.New,
	server.New,
)

func NewServer(ctx context.Context) (*server.Server, error) {
	wire.Build(set, config.Load)
	return &server.Server{}, nil
}

func NewServerForTest(ctx context.Context, cfgFn config.Decorator) (*server.Server, error) {
	wire.Build(set, config.ForTest)
	return &server.Server{}, nil
}
