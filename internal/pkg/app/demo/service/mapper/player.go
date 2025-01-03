package mapper

import (
	"fmt"

	api "github.com/michaljurecko/ch-demo/api/gen/go/demo/v1"
	model "github.com/michaljurecko/ch-demo/internal/pkg/app/demo/model/gen"
)

func Player(entity *model.Player) *api.Player {
	return &api.Player{
		Id:        entity.ID,
		FirstName: entity.FirstName,
		LastName:  entity.LastName,
		Phone:     entity.Phone,
		Email:     entity.Email,
		Ic:        entity.IC,
		VatId:     entity.VATID,
		Address:   entity.Address,
	}
}

func CreatePlayer(req *api.CreatePlayerRequest) *model.Player {
	return &model.Player{
		FirstName: req.GetFirstName(),
		LastName:  req.GetLastName(),
		Phone:     req.GetPhone(),
		Email:     req.GetEmail(),
		IC:        req.GetIc(),
	}
}

func UpdatePlayer(req *api.UpdatePlayerRequest, entity *model.Player) error {
	for _, field := range req.GetUpdateMask().GetPaths() {
		switch field {
		case "first_name":
			entity.FirstName = req.GetFirstName()
		case "last_name":
			entity.LastName = req.GetLastName()
		case "phone":
			entity.Phone = req.GetPhone()
		case "email":
			entity.Email = req.GetEmail()
		case "ic":
			entity.IC = req.GetIc()
		default:
			return fmt.Errorf("unexpected mask field '%s'", field)
		}
	}
	return nil
}
