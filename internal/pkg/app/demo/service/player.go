package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"connectrpc.com/connect"

	api "github.com/michaljurecko/ch-demo/api/gen/go/demo/v1"
	"github.com/michaljurecko/ch-demo/internal/pkg/app/demo/ares"
	"github.com/michaljurecko/ch-demo/internal/pkg/app/demo/service/mapper"
	"github.com/michaljurecko/ch-demo/internal/pkg/common/dataverse/webapi"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Service) CreatePlayer(
	ctx context.Context,
	req *connect.Request[api.CreatePlayerRequest],
) (*connect.Response[api.Player], error) {
	entity := mapper.CreatePlayer(req.Msg)

	if err := s.playerSvc.AddPlayerInformation(ctx, entity); errors.As(err, &ares.EntityNotFoundError{}) {
		return nil, err
	} else if err != nil {
		return nil, err
	}

	if _, err := s.repo.Player().Create(entity).Do(ctx); err != nil {
		return nil, fmt.Errorf("failed to save player entity: %w", err)
	}

	s.invalidateCache(ctx, PlayersCacheTag)

	return connect.NewResponse(mapper.Player(entity)), nil
}

func (s *Service) UpdatePlayer(
	ctx context.Context,
	req *connect.Request[api.UpdatePlayerRequest],
) (*connect.Response[api.Player], error) {
	entity, err := s.repo.Player().ByID(req.Msg.GetId()).Do(ctx)
	if err != nil {
		var statusErr webapi.UnexpectedStatusError
		if errors.As(err, &statusErr) && statusErr.Actual == http.StatusNotFound {
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("player '%s' not found", req.Msg.GetId()))
		}
		return nil, fmt.Errorf("failed to load player entity: %w", err)
	}

	entity.TrackChanges()
	if err := mapper.UpdatePlayer(req.Msg, entity); err != nil {
		return nil, err
	}

	if err := s.playerSvc.AddPlayerInformation(ctx, entity); errors.As(err, &ares.EntityNotFoundError{}) {
		return nil, err
	} else if err != nil {
		return nil, err
	}

	if _, err := s.repo.Player().Update(entity).Do(ctx); err != nil {
		return nil, fmt.Errorf("failed to save player entity: %w", err)
	}

	s.invalidateCache(ctx, PlayersCacheTag)

	return connect.NewResponse(mapper.Player(entity)), nil
}

func (s *Service) DeletePlayer(
	ctx context.Context,
	req *connect.Request[api.DeletePlayerRequest],
) (*connect.Response[emptypb.Empty], error) {
	if _, err := s.repo.Player().Delete(req.Msg.GetId()).Do(ctx); err != nil {
		var statusErr webapi.UnexpectedStatusError
		if errors.As(err, &statusErr) && statusErr.Actual == http.StatusNotFound {
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("player '%s' not found", req.Msg.GetId()))
		}
		return nil, fmt.Errorf("failed to delete player entity: %w", err)
	}

	s.invalidateCache(ctx, PlayersCacheTag)

	return connect.NewResponse(&emptypb.Empty{}), nil
}

func (s *Service) GetPlayer(
	ctx context.Context,
	req *connect.Request[api.GetPlayerRequest],
) (*connect.Response[api.Player], error) {
	entity, err := s.repo.Player().ByID(req.Msg.GetId()).Do(ctx)
	if err != nil {
		var statusErr webapi.UnexpectedStatusError
		if errors.As(err, &statusErr) && statusErr.Actual == http.StatusNotFound {
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("player '%s' not found", req.Msg.GetId()))
		}
		return nil, fmt.Errorf("failed to delete player entity: %w", err)
	}
	return connect.NewResponse(mapper.Player(entity)), nil
}
