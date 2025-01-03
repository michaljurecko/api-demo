package characterbiz

import (
	"math/rand"
	"time"

	"google.golang.org/protobuf/types/known/fieldmaskpb"

	model "github.com/michaljurecko/ch-demo/internal/pkg/app/demo/model/gen"
)

const DiceSideCount = 6

//nolint:revive // method cannot be split into smaller parts
func EnrichCharacterEntity(actionType string, entity *model.Character, class *model.Class, race *model.Race, updateMask *fieldmaskpb.FieldMask) (rolls []*model.DiceRoll) {
	isCreate := updateMask == nil
	updateMap := maskSliceToMap(updateMask)

	if isCreate || updateMap["strength"] || updateMap["class_id"] || updateMap["player_id"] {
		strengthRoll := rollDice("strength", actionType)
		rolls = append(rolls, strengthRoll)
		entity.Strength += race.StrengthBase + class.StrengthBase + strengthRoll.RollResult
	}

	if isCreate || updateMap["dexterity"] || updateMap["class_id"] || updateMap["player_id"] {
		dexterityRoll := rollDice("dexterity", actionType)
		rolls = append(rolls, dexterityRoll)
		entity.Dexterity += race.DexterityBase + class.DexterityBase + dexterityRoll.RollResult
	}

	if isCreate || updateMap["intelligence"] || updateMap["class_id"] || updateMap["player_id"] {
		intelligenceRoll := rollDice("intelligence", actionType)
		rolls = append(rolls, intelligenceRoll)
		entity.Intelligence += race.IntelligenceBase + class.IntelligenceBase + intelligenceRoll.RollResult
	}

	if isCreate || updateMap["charisma"] || updateMap["class_id"] || updateMap["player_id"] {
		charismaRoll := rollDice("charisma", actionType)
		rolls = append(rolls, charismaRoll)
		entity.Charisma += race.CharismaBase + class.CharismaBase + charismaRoll.RollResult
	}

	return rolls
}

func rollDice(name, actionType string) *model.DiceRoll {
	return &model.DiceRoll{
		DiceRoll:   name,
		RollResult: rand.Intn(DiceSideCount) + 1, //nolint:gosec // weak random number generator is acceptable here
		RollDate:   time.Now(),
		ActionType: actionType,
	}
}
