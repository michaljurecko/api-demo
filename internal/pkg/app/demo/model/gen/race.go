// DO NOT EDIT. Code generated by entitygen. For modifications use "make gen-model".

package model

import (
	"context"
	"errors"
	webapi "github.com/michaljurecko/ch-demo/internal/pkg/common/dataverse/webapi"
	"net/http"
)

// Race - This table contains records of character races.
type Race struct {
	// original field stores snapshot of the entity state for the Update operation, see TrackChanges method.
	original         *Race
	ID               string `json:"cr568_raceid,omitempty"`
	Name             string `json:"cr568_racename"`
	Description      string `json:"cr568_description"`
	StrengthBase     int    `json:"cr568_strengthbase"`
	CharismaBase     int    `json:"cr568_charismabase"`
	IntelligenceBase int    `json:"cr568_intelligencebase"`
	DexterityBase    int    `json:"cr568_dexteritybase"`
}

// TrackChanges internally stores entity actual state to track changes for the Update operation.
func (e *Race) TrackChanges() {
	clone := *e
	e.original = &clone
}

// resetChanges after the Update operation.
func (e *Race) resetChanges() {
	e.original = nil
}

type RaceRepository struct {
	client *webapi.Client
}

func newRaceRepository(client *webapi.Client) *RaceRepository {
	return &RaceRepository{client: client}
}

// Create entity. After successful operation, the new primary ID will be set to the original entity.
func (r *RaceRepository) Create(entity *Race) *webapi.APIRequest[Race] {
	payload, err := r.createPayload(entity)
	if err != nil {
		return webapi.NewAPIRequestError[Race](entity, err)
	}

	path := "cr568_races"
	result := entity
	httpReq := webapi.NewHTTPRequest(result, r.client, http.MethodPost, path, payload)
	httpReq.Header("Prefer", "return=representation")
	httpReq.ExpectStatus(http.StatusCreated)
	return webapi.NewAPIRequest(result, httpReq)
}

// Update entity. A diff of modifications is generated and saved via API.
// Before making changes, it is necessary to call the TrackChanges entity method.
func (r *RaceRepository) Update(entity *Race) *webapi.APIRequest[Race] {

	if entity.original == nil {
		err := errors.New("changes are not tracked: use \"TrackChanges\" entity method to track changes and allow \"Update\" operation")
		return webapi.NewAPIRequestError[Race](entity, err)
	}

	payload, err := r.updatePayload(entity.original, entity)
	if err != nil {
		return webapi.NewAPIRequestError[Race](entity, err)
	}

	id := entity.ID
	path := "cr568_races" + "(" + webapi.ID(id) + ")"
	result := entity
	httpReq := webapi.NewHTTPRequest(result, r.client, http.MethodPatch, path, payload)
	httpReq.Header("If-Match", "*") // prevent create if entity not exists
	httpReq.Header("Prefer", "return=representation")
	httpReq.OnSuccess(func(ctx context.Context, c *Race) error {
		c.resetChanges()
		return nil
	})
	return webapi.NewAPIRequest(result, httpReq)
}

// Delete entity by the ID.
func (r *RaceRepository) Delete(id string) *webapi.APIRequest[webapi.NoResult] {
	path := "cr568_races" + "(" + webapi.ID(id) + ")"
	httpReq := webapi.NewHTTPRequest(&webapi.NoResult{}, r.client, http.MethodDelete, path, nil)
	httpReq.ExpectStatus(http.StatusNoContent)
	return webapi.NewAPIRequest(&webapi.NoResult{}, httpReq)
}

func (r *RaceRepository) All() *webapi.APIRequest[webapi.Collection[Race]] {
	path := "cr568_races"
	result := &webapi.Collection[Race]{}
	httpReq := webapi.NewHTTPRequest(result, r.client, http.MethodGet, path, nil)
	return webapi.NewAPIRequest(result, httpReq)
}

func (r *RaceRepository) ByID(id string) *webapi.APIRequest[Race] {
	path := "cr568_races" + "(" + webapi.ID(id) + ")"
	result := &Race{}
	httpReq := webapi.NewHTTPRequest(result, r.client, http.MethodGet, path, nil)
	return webapi.NewAPIRequest(result, httpReq)
}

// createPayload method generates payload for the Create operation.
func (r *RaceRepository) createPayload(entity *Race) (map[string]any, error) {
	payload := map[string]any{}
	payload["cr568_racename"] = entity.Name
	payload["cr568_description"] = entity.Description
	payload["cr568_strengthbase"] = entity.StrengthBase
	payload["cr568_charismabase"] = entity.CharismaBase
	payload["cr568_intelligencebase"] = entity.IntelligenceBase
	payload["cr568_dexteritybase"] = entity.DexterityBase
	return payload, nil
}

// updatePayload method generates diff of original and modified entity state for the Update operation.
func (r *RaceRepository) updatePayload(original, modified *Race) (map[string]any, error) {
	payload := map[string]any{}
	if original.Name != modified.Name {
		payload["cr568_racename"] = modified.Name
	}
	if original.Description != modified.Description {
		payload["cr568_description"] = modified.Description
	}
	if original.StrengthBase != modified.StrengthBase {
		payload["cr568_strengthbase"] = modified.StrengthBase
	}
	if original.CharismaBase != modified.CharismaBase {
		payload["cr568_charismabase"] = modified.CharismaBase
	}
	if original.IntelligenceBase != modified.IntelligenceBase {
		payload["cr568_intelligencebase"] = modified.IntelligenceBase
	}
	if original.DexterityBase != modified.DexterityBase {
		payload["cr568_dexteritybase"] = modified.DexterityBase
	}
	return payload, nil
}
