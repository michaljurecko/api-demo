package mapper

import (
	"sort"

	api "github.com/michaljurecko/ch-demo/api/gen/go/demo/v1"
	model "github.com/michaljurecko/ch-demo/internal/pkg/app/demo/model/gen"
	"github.com/michaljurecko/ch-demo/internal/pkg/common/dataverse/webapi"
)

func PlayersAndCharacters(
	playersCol *webapi.Collection[model.Player],
	charactersCol *webapi.Collection[model.Character],
) []*api.CharactersPerPlayer {
	players := playersCol.Items
	sort.SliceStable(players, func(i, j int) bool {
		return players[i].ID < players[j].ID
	})

	characters := charactersCol.Items
	sort.SliceStable(characters, func(i, j int) bool {
		return characters[i].ID < characters[j].ID
	})

	charactersMap := make(map[string][]*api.Character)
	for _, character := range characters {
		playerID := character.Player.ID()
		charactersMap[playerID] = append(charactersMap[playerID], Character(character))
	}

	var out []*api.CharactersPerPlayer
	for _, player := range players {
		out = append(out, &api.CharactersPerPlayer{
			Player:     Player(player),
			Characters: charactersMap[player.ID],
		})
	}

	return out
}
