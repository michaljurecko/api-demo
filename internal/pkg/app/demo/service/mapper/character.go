package mapper

import (
	"fmt"

	api "github.com/michaljurecko/ch-demo/api/gen/go/demo/v1"
	model "github.com/michaljurecko/ch-demo/internal/pkg/app/demo/model/gen"
	"github.com/michaljurecko/ch-demo/internal/pkg/common/dataverse/webapi"
)

func Character(entity *model.Character) *api.Character {
	return &api.Character{
		Id:   entity.ID,
		Name: entity.Name,
		// Level: entity.Level,
		Strength:     int32(entity.Strength),     //nolint:gosec // conversion is safe, numbers are small and validated
		Dexterity:    int32(entity.Dexterity),    //nolint:gosec // conversion is safe, numbers are small and validated
		Intelligence: int32(entity.Intelligence), //nolint:gosec // conversion is safe, numbers are small and validated
		Charisma:     int32(entity.Charisma),     //nolint:gosec // conversion is safe, numbers are small and validated
		ClassId:      entity.Class.ID(),
		RaceId:       entity.Race.ID(),
		PlayerId:     entity.Player.ID(),
	}
}

func CreateCharacter(req *api.CreateCharacterRequest) *model.Character {
	return &model.Character{
		Name:         req.GetName(),
		Strength:     int(req.GetStrength()),
		Dexterity:    int(req.GetDexterity()),
		Intelligence: int(req.GetIntelligence()),
		Charisma:     int(req.GetCharisma()),
		Class:        webapi.LookupValue[model.Class](req.GetClassId()),
		Race:         webapi.LookupValue[model.Race](req.GetRaceId()),
		Player:       webapi.LookupValue[model.Player](req.GetPlayerId()),
	}
}

func UpdateCharacter(req *api.UpdateCharacterRequest, entity *model.Character) error {
	for _, field := range req.GetUpdateMask().GetPaths() {
		switch field {
		case "name":
			entity.Name = req.GetName()
		case "strength":
			entity.Strength = int(req.GetStrength())
		case "dexterity":
			entity.Dexterity = int(req.GetDexterity())
		case "intelligence":
			entity.Intelligence = int(req.GetIntelligence())
		case "charisma":
			entity.Charisma = int(req.GetCharisma())
		case "class":
			entity.Class = webapi.LookupValue[model.Class](req.GetClassId())
		case "race":
			entity.Race = webapi.LookupValue[model.Race](req.GetRaceId())
		case "player":
			entity.Player = webapi.LookupValue[model.Player](req.GetPlayerId())
		default:
			return fmt.Errorf("unexpected mask field '%s'", field)
		}
	}
	return nil
}
