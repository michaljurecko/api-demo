// Package webapi provides HTTP client for Dataverse OData Web API with OAuth authentication.
package metadata

import (
	"fmt"
	"net/http"

	"github.com/michaljurecko/ch-demo/internal/pkg/common/dataverse/webapi"
)

// API contains part of DataVerse Web API related to metadata.
type API struct {
	client *webapi.Client
}

// Entities represents the response from the Entities endpoint.
type Entities struct {
	Definitions []Entity `json:"value"`
}

// Entity represents the metadata of an entity.
type Entity struct {
	MetadataID    string `json:"MetadataId"`
	LogicalName   string `json:"LogicalName"`
	EntitySetName string `json:"EntitySetName"` // <- API endpoint name
	DisplayName   Labels `json:"DisplayName"`
	Description   Labels `json:"Description"`
}

// Labels represents a collection of localized labels.
type Labels struct {
	LocalizedLabels    []Label `json:"LocalizedLabels"`
	UserLocalizedLabel Label   `json:"UserLocalizedLabel"`
}

// Label represents a localized label.
type Label struct {
	Label        string `json:"Label"`
	LanguageCode int    `json:"LanguageCode"`
}

// Attributes represents the metadata of an entity attributes.
type Attributes struct {
	Attributes []Attribute `json:"value"`
}

// Attribute represents the metadata of an entity attribute.
type Attribute struct {
	MetadataID    string   `json:"MetadataId"`
	LogicalName   string   `json:"LogicalName"`
	AttributeType string   `json:"AttributeType"`
	DisplayName   Labels   `json:"DisplayName"`
	SchemaName    string   `json:"SchemaName"`
	Description   Labels   `json:"Description"`
	ColumnNumber  int      `json:"ColumnNumber"`
	IsPrimaryID   bool     `json:"IsPrimaryId"`
	IsPrimaryName bool     `json:"IsPrimaryName"`
	Targets       []string `json:"Targets"` // AttributeType == "Lookup"
}

func NewAPI(client *webapi.Client) *API {
	return &API{client: client}
}

func (a *API) EntityDefinitions() *webapi.APIRequest[Entities] {
	path := "EntityDefinitions"
	result := &Entities{}
	return webapi.NewAPIRequest(result, webapi.NewHTTPRequest(result, a.client, http.MethodGet, path, nil))
}

func (a *API) CustomEntityDefinitions() *webapi.APIRequest[Entities] {
	path := "EntityDefinitions"
	result := &Entities{}
	return webapi.NewAPIRequest(result,
		webapi.NewHTTPRequest(result, a.client, http.MethodGet, path, nil).
			Filter(
				webapi.And(
					"OwnershipType eq 'UserOwned'",
					"IsCustomEntity eq true",
					"IsManaged eq false",
					"IsMappable/Value eq true",
					"IsRenameable/Value eq true",
				),
			),
	)
}

func (a *API) EntityDefinition(logicalName string) *webapi.APIRequest[Entity] {
	path := fmt.Sprintf("EntityDefinitions(LogicalName='%s')", logicalName)
	result := &Entity{}
	return webapi.NewAPIRequest(result, webapi.NewHTTPRequest(result, a.client, http.MethodGet, path, nil))
}

func (a *API) EntityAttributes(logicalName string) *webapi.APIRequest[Attributes] {
	path := fmt.Sprintf("EntityDefinitions(LogicalName='%s')/Attributes", logicalName)
	result := &Attributes{}
	return webapi.NewAPIRequest(result, webapi.NewHTTPRequest(result, a.client, http.MethodGet, path, nil))
}

func (a *API) CustomEntityAttributes(logicalName string) *webapi.APIRequest[Attributes] {
	path := fmt.Sprintf("EntityDefinitions(LogicalName='%s')/Attributes", logicalName)
	result := &Attributes{}
	return webapi.NewAPIRequest(result,
		webapi.NewHTTPRequest(result, a.client, http.MethodGet, path, nil).
			Filter(
				webapi.Or(
					"IsPrimaryId eq true",
					webapi.And(
						"IsCustomAttribute eq true",
						"IsValidODataAttribute eq true",
						"IsLogical eq false",
					),
				),
			),
	)
}
