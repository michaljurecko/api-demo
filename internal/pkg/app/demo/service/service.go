package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/michaljurecko/ch-demo/internal/pkg/app/demo/biz/playerbiz"
	"github.com/michaljurecko/ch-demo/internal/pkg/common/cache"
	"github.com/michaljurecko/ch-demo/internal/pkg/common/cachestore"

	"github.com/michaljurecko/ch-demo/api/gen/go/demo/v1"

	"github.com/bufbuild/protovalidate-go"
	"github.com/michaljurecko/ch-demo/internal/pkg/app/demo/config"
	model "github.com/michaljurecko/ch-demo/internal/pkg/app/demo/model/gen"
	"github.com/michaljurecko/ch-demo/internal/pkg/common/log"
)

const (
	PlayersAndCharactersCacheKey = "players.v1"
	DistLockTTL                  = 30 * time.Second
	CacheTTL                     = 1 * time.Hour
	PlayersCacheTag              = "players"
	CharactersCacheTag           = "characters"
)

type Service struct {
	logger         *log.Logger
	repo           *model.Repository
	playerSvc      *playerbiz.Service
	protoValidator *protovalidate.Validator
	// locker         *distlock.Locker
	cacheStore   *cachestore.Store
	playersCache *cache.Marshaler[*api.ListPlayersAndCharactersResponse]
}

func New(
	ctx context.Context,
	logger *log.Logger,
	cfg config.Config,
	repo *model.Repository,
	cacheStore *cachestore.Store,
	// locker *distlock.Locker,
	playerSvc *playerbiz.Service,
) (*Service, error) {
	// Init service dependencies
	svc := &Service{
		logger:    logger.With(slog.String(log.LoggerKey, "service")),
		repo:      repo,
		playerSvc: playerSvc,
		// locker:       locker,
		cacheStore:   cacheStore,
		playersCache: cache.NewMarshaler[*api.ListPlayersAndCharactersResponse](cacheStore.Store()),
	}

	// Log configuration (without sensitives values)
	bytes, err := json.Marshal(cfg) //nolint:musttag
	if err != nil {
		return nil, fmt.Errorf("failed to encode configuration: %w", err)
	}
	logger.Info(ctx, "loaded configuration", slog.String(log.LoggerKey, "config"), slog.Any("config", string(bytes)))

	// Create proto validator
	svc.protoValidator, err = protovalidate.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create proto validator: %w", err)
	}

	return svc, nil
}
