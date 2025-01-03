package mapper

import (
	api "github.com/michaljurecko/ch-demo/api/gen/go/demo/v1"
	model "github.com/michaljurecko/ch-demo/internal/pkg/app/demo/model/gen"
	"github.com/michaljurecko/ch-demo/internal/pkg/common/dataverse/webapi"
)

func Class(entity *model.Class) *api.Class {
	return &api.Class{
		Id:               entity.ID,
		Name:             entity.Name,
		Description:      entity.Description,
		StrengthBase:     int32(entity.StrengthBase),     //nolint:gosec //conversion is safe, numbers are small and validated
		IntelligenceBase: int32(entity.IntelligenceBase), //nolint:gosec //conversion is safe, numbers are small and validated
		CharismaBase:     int32(entity.CharismaBase),     //nolint:gosec //conversion is safe, numbers are small and validated
		DexterityBase:    int32(entity.DexterityBase),    //nolint:gosec //conversion is safe, numbers are small and validated
	}
}

func Classes(entities *webapi.Collection[model.Class]) (out []*api.Class) {
	for _, entity := range entities.Items {
		out = append(out, Class(entity))
	}
	return out
}
