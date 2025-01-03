package service

import (
	"context"
	"fmt"

	"connectrpc.com/connect"

	api "github.com/michaljurecko/ch-demo/api/gen/go/demo/v1"
	"github.com/michaljurecko/ch-demo/internal/pkg/app/demo/service/mapper"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Service) ListRaces(
	ctx context.Context,
	_ *connect.Request[emptypb.Empty],
) (*connect.Response[api.ListRacesResponse], error) {
	entities, err := s.repo.Race().All().Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load race entities: %w", err)
	}

	return connect.NewResponse(&api.ListRacesResponse{Races: mapper.Races(entities)}), nil
}
