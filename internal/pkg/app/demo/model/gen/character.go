// DO NOT EDIT. Code generated by entitygen. For modifications use "make gen-model".

package model

import (
	"context"
	"errors"
	webapi "github.com/michaljurecko/ch-demo/internal/pkg/common/dataverse/webapi"
	"net/http"
)

// Character - This table contains records of Dungeons and Dragons characters.
type Character struct {
	// original field stores snapshot of the entity state for the Update operation, see TrackChanges method.
	original     *Character
	ID           string                `json:"cr568_characterid,omitempty"`
	Name         string                `json:"cr568_charactername"`
	Level        int                   `json:"cr568_level"`
	Strength     int                   `json:"cr568_strength"`
	Dexterity    int                   `json:"cr568_dexterity"`
	Intelligence int                   `json:"cr568_intelligence"`
	Charisma     int                   `json:"cr568_charisma"`
	Class        webapi.Lookup[Class]  `json:"_cr568_classid_value"`
	Race         webapi.Lookup[Race]   `json:"_cr568_raceid_value"`
	Player       webapi.Lookup[Player] `json:"_cr568_player_value"`
}

// TrackChanges internally stores entity actual state to track changes for the Update operation.
func (e *Character) TrackChanges() {
	clone := *e
	e.original = &clone
}

// resetChanges after the Update operation.
func (e *Character) resetChanges() {
	e.original = nil
}

type CharacterRepository struct {
	client *webapi.Client
}

func newCharacterRepository(client *webapi.Client) *CharacterRepository {
	return &CharacterRepository{client: client}
}

// Create entity. After successful operation, the new primary ID will be set to the original entity.
func (r *CharacterRepository) Create(entity *Character) *webapi.APIRequest[Character] {
	payload, err := r.createPayload(entity)
	if err != nil {
		return webapi.NewAPIRequestError[Character](entity, err)
	}

	path := "cr568_characters"
	result := entity
	httpReq := webapi.NewHTTPRequest(result, r.client, http.MethodPost, path, payload)
	httpReq.Header("Prefer", "return=representation")
	httpReq.ExpectStatus(http.StatusCreated)
	return webapi.NewAPIRequest(result, httpReq)
}

// Update entity. A diff of modifications is generated and saved via API.
// Before making changes, it is necessary to call the TrackChanges entity method.
func (r *CharacterRepository) Update(entity *Character) *webapi.APIRequest[Character] {

	if entity.original == nil {
		err := errors.New("changes are not tracked: use \"TrackChanges\" entity method to track changes and allow \"Update\" operation")
		return webapi.NewAPIRequestError[Character](entity, err)
	}

	payload, err := r.updatePayload(entity.original, entity)
	if err != nil {
		return webapi.NewAPIRequestError[Character](entity, err)
	}

	id := entity.ID
	path := "cr568_characters" + "(" + webapi.ID(id) + ")"
	result := entity
	httpReq := webapi.NewHTTPRequest(result, r.client, http.MethodPatch, path, payload)
	httpReq.Header("If-Match", "*") // prevent create if entity not exists
	httpReq.Header("Prefer", "return=representation")
	httpReq.OnSuccess(func(ctx context.Context, c *Character) error {
		c.resetChanges()
		return nil
	})
	return webapi.NewAPIRequest(result, httpReq)
}

// Delete entity by the ID.
func (r *CharacterRepository) Delete(id string) *webapi.APIRequest[webapi.NoResult] {
	path := "cr568_characters" + "(" + webapi.ID(id) + ")"
	httpReq := webapi.NewHTTPRequest(&webapi.NoResult{}, r.client, http.MethodDelete, path, nil)
	httpReq.ExpectStatus(http.StatusNoContent)
	return webapi.NewAPIRequest(&webapi.NoResult{}, httpReq)
}

func (r *CharacterRepository) All() *webapi.APIRequest[webapi.Collection[Character]] {
	path := "cr568_characters"
	result := &webapi.Collection[Character]{}
	httpReq := webapi.NewHTTPRequest(result, r.client, http.MethodGet, path, nil)
	return webapi.NewAPIRequest(result, httpReq)
}

func (r *CharacterRepository) ByID(id string) *webapi.APIRequest[Character] {
	path := "cr568_characters" + "(" + webapi.ID(id) + ")"
	result := &Character{}
	httpReq := webapi.NewHTTPRequest(result, r.client, http.MethodGet, path, nil)
	return webapi.NewAPIRequest(result, httpReq)
}

func (r *CharacterRepository) ByClass(id string) *webapi.APIRequest[webapi.Collection[Character]] {
	path := "cr568_characters"
	result := &webapi.Collection[Character]{}
	httpReq := webapi.NewHTTPRequest(result, r.client, http.MethodGet, path, nil)
	httpReq.Filter("cr568_classid eq '" + webapi.ID(id) + "'")
	return webapi.NewAPIRequest(result, httpReq)
}

func (r *CharacterRepository) ByRace(id string) *webapi.APIRequest[webapi.Collection[Character]] {
	path := "cr568_characters"
	result := &webapi.Collection[Character]{}
	httpReq := webapi.NewHTTPRequest(result, r.client, http.MethodGet, path, nil)
	httpReq.Filter("cr568_raceid eq '" + webapi.ID(id) + "'")
	return webapi.NewAPIRequest(result, httpReq)
}

func (r *CharacterRepository) ByPlayer(id string) *webapi.APIRequest[webapi.Collection[Character]] {
	path := "cr568_characters"
	result := &webapi.Collection[Character]{}
	httpReq := webapi.NewHTTPRequest(result, r.client, http.MethodGet, path, nil)
	httpReq.Filter("cr568_player eq '" + webapi.ID(id) + "'")
	return webapi.NewAPIRequest(result, httpReq)
}

// createPayload method generates payload for the Create operation.
func (r *CharacterRepository) createPayload(entity *Character) (map[string]any, error) {
	payload := map[string]any{}
	payload["cr568_charactername"] = entity.Name
	payload["cr568_level"] = entity.Level
	payload["cr568_strength"] = entity.Strength
	payload["cr568_dexterity"] = entity.Dexterity
	payload["cr568_intelligence"] = entity.Intelligence
	payload["cr568_charisma"] = entity.Charisma
	{
		idOrNil, err := lookupFullIDOrNil(entity.Class)
		if err != nil {
			return nil, err
		}
		payload["cr568_ClassID@odata.bind"] = idOrNil
	}
	{
		idOrNil, err := lookupFullIDOrNil(entity.Race)
		if err != nil {
			return nil, err
		}
		payload["cr568_RaceID@odata.bind"] = idOrNil
	}
	{
		idOrNil, err := lookupFullIDOrNil(entity.Player)
		if err != nil {
			return nil, err
		}
		payload["cr568_Player@odata.bind"] = idOrNil
	}
	return payload, nil
}

// updatePayload method generates diff of original and modified entity state for the Update operation.
func (r *CharacterRepository) updatePayload(original, modified *Character) (map[string]any, error) {
	payload := map[string]any{}
	if original.Name != modified.Name {
		payload["cr568_charactername"] = modified.Name
	}
	if original.Level != modified.Level {
		payload["cr568_level"] = modified.Level
	}
	if original.Strength != modified.Strength {
		payload["cr568_strength"] = modified.Strength
	}
	if original.Dexterity != modified.Dexterity {
		payload["cr568_dexterity"] = modified.Dexterity
	}
	if original.Intelligence != modified.Intelligence {
		payload["cr568_intelligence"] = modified.Intelligence
	}
	if original.Charisma != modified.Charisma {
		payload["cr568_charisma"] = modified.Charisma
	}
	if original.Class.ID() != modified.Class.ID() {
		idOrNil, err := lookupFullIDOrNil(modified.Class)
		if err != nil {
			return nil, err
		}
		payload["cr568_ClassID@odata.bind"] = idOrNil
	}
	if original.Race.ID() != modified.Race.ID() {
		idOrNil, err := lookupFullIDOrNil(modified.Race)
		if err != nil {
			return nil, err
		}
		payload["cr568_RaceID@odata.bind"] = idOrNil
	}
	if original.Player.ID() != modified.Player.ID() {
		idOrNil, err := lookupFullIDOrNil(modified.Player)
		if err != nil {
			return nil, err
		}
		payload["cr568_Player@odata.bind"] = idOrNil
	}
	return payload, nil
}
