package playerbiz

import (
	"context"

	"github.com/michaljurecko/ch-demo/internal/pkg/app/demo/ares"
	model "github.com/michaljurecko/ch-demo/internal/pkg/app/demo/model/gen"
)

type Service struct {
	client *ares.Client
}

func NewService(client *ares.Client) *Service {
	return &Service{client: client}
}

func (s *Service) AddPlayerInformation(ctx context.Context, entity *model.Player) error {
	result, err := s.client.ByIC(ctx, entity.IC)
	if err != nil {
		return err
	}

	entity.VATID = result.VatID
	entity.Address = result.Address.String()
	return nil
}
