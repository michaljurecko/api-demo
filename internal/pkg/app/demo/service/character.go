package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"connectrpc.com/connect"

	api "github.com/michaljurecko/ch-demo/api/gen/go/demo/v1"
	"github.com/michaljurecko/ch-demo/internal/pkg/app/demo/biz/characterbiz"
	"github.com/michaljurecko/ch-demo/internal/pkg/app/demo/service/mapper"
	"github.com/michaljurecko/ch-demo/internal/pkg/common/dataverse/webapi"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Service) CreateCharacter(
	ctx context.Context,
	req *connect.Request[api.CreateCharacterRequest],
) (*connect.Response[api.Character], error) {
	entity := mapper.CreateCharacter(req.Msg)

	class, err := s.repo.Class().ByID(entity.Class.ID()).Do(ctx)
	if err != nil {
		var statusErr webapi.UnexpectedStatusError
		if errors.As(err, &statusErr) && statusErr.Actual == http.StatusNotFound {
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("class '%s' not found", entity.Class.ID()))
		}
		return nil, fmt.Errorf("failed to load class entity: %w", err)
	}

	race, err := s.repo.Race().ByID(entity.Race.ID()).Do(ctx)
	if err != nil {
		var statusErr webapi.UnexpectedStatusError
		if errors.As(err, &statusErr) && statusErr.Actual == http.StatusNotFound {
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("race '%s' not found", entity.Class.ID()))
		}
		return nil, fmt.Errorf("failed to load race entity: %w", err)
	}

	diceRolls := characterbiz.EnrichCharacterEntity("CreateCharacter", entity, class, race, nil)

	// Save atomically
	changes := s.repo.NewChangeSet()
	characterRef := changes.Add(s.repo.Character().Create(entity))
	for _, diceRoll := range diceRolls {
		diceRoll.Character.SetContentID(characterRef)
		changes.Add(s.repo.DiceRoll().Create(diceRoll))
	}
	if err := changes.Do(ctx); err != nil {
		return nil, fmt.Errorf("failed to save character entity: %w", err)
	}

	s.invalidateCache(ctx, CharactersCacheTag)

	return connect.NewResponse(mapper.Character(entity)), nil
}

func (s *Service) UpdateCharacter(
	ctx context.Context,
	req *connect.Request[api.UpdateCharacterRequest],
) (*connect.Response[api.Character], error) {
	entity, err := s.repo.Character().ByID(req.Msg.GetId()).Do(ctx)
	if err != nil {
		var statusErr webapi.UnexpectedStatusError
		if errors.As(err, &statusErr) && statusErr.Actual == http.StatusNotFound {
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("character '%s' not found", req.Msg.GetId()))
		}
		return nil, fmt.Errorf("failed to delete character entity: %w", err)
	}

	entity.TrackChanges()
	if err := mapper.UpdateCharacter(req.Msg, entity); err != nil {
		return nil, err
	}

	class, err := s.repo.Class().ByID(entity.Class.ID()).Do(ctx)
	if err != nil {
		var statusErr webapi.UnexpectedStatusError
		if errors.As(err, &statusErr) && statusErr.Actual == http.StatusNotFound {
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("class '%s' not found", entity.Class.ID()))
		}
		return nil, fmt.Errorf("failed to load class entity: %w", err)
	}

	race, err := s.repo.Race().ByID(entity.Race.ID()).Do(ctx)
	if err != nil {
		var statusErr webapi.UnexpectedStatusError
		if errors.As(err, &statusErr) && statusErr.Actual == http.StatusNotFound {
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("race '%s' not found", entity.Class.ID()))
		}
		return nil, fmt.Errorf("failed to load race entity: %w", err)
	}

	diceRolls := characterbiz.EnrichCharacterEntity("UpdateCharacter", entity, class, race, req.Msg.GetUpdateMask())

	// Save atomically
	changes := s.repo.NewChangeSet()
	characterRef := changes.Add(s.repo.Character().Update(entity))
	for _, diceRoll := range diceRolls {
		diceRoll.Character.SetContentID(characterRef)
		changes.Add(s.repo.DiceRoll().Create(diceRoll))
	}
	if err := changes.Do(ctx); err != nil {
		return nil, fmt.Errorf("failed to save character entity: %w", err)
	}

	s.invalidateCache(ctx, CharactersCacheTag)

	return connect.NewResponse(mapper.Character(entity)), nil
}

func (s *Service) DeleteCharacter(
	ctx context.Context,
	req *connect.Request[api.DeleteCharacterRequest],
) (*connect.Response[emptypb.Empty], error) {
	if _, err := s.repo.Character().Delete(req.Msg.GetId()).Do(ctx); err != nil {
		var statusErr webapi.UnexpectedStatusError
		if errors.As(err, &statusErr) && statusErr.Actual == http.StatusNotFound {
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("character '%s' not found", req.Msg.GetId()))
		}
		return nil, fmt.Errorf("failed to delete character entity: %w", err)
	}

	s.invalidateCache(ctx, CharactersCacheTag)

	return connect.NewResponse(&emptypb.Empty{}), nil
}

func (s *Service) GetCharacter(
	ctx context.Context,
	req *connect.Request[api.GetCharacterRequest],
) (*connect.Response[api.Character], error) {
	entity, err := s.repo.Character().ByID(req.Msg.GetId()).Do(ctx)
	if err != nil {
		var statusErr webapi.UnexpectedStatusError
		if errors.As(err, &statusErr) && statusErr.Actual == http.StatusNotFound {
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("character '%s' not found", req.Msg.GetId()))
		}
		return nil, fmt.Errorf("failed to delete character entity: %w", err)
	}
	return connect.NewResponse(mapper.Character(entity)), nil
}
