package ares

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const EconomicEntityURL = "https://ares.gov.cz/ekonomicke-subjekty-v-be/rest/ekonomicke-subjekty/"

type Client struct {
	client *http.Client
}

type EconomicEntity struct {
	Name    string  `json:"obchodniJmeno"`
	IC      string  `json:"ico"`
	VatID   string  `json:"dic"`
	Address Address `json:"adresaDorucovaci"`
}

type Address struct {
	Line1 string `json:"radekAdresy1"`
	Line2 string `json:"radekAdresy2"`
	Line3 string `json:"radekAdresy3"`
}

type EntityNotFoundError struct {
	IC string
}

func (e EntityNotFoundError) Error() string {
	return fmt.Sprintf("economic entity with IC '%s' not found", e.IC)
}

func NewClient(client *http.Client) *Client {
	return &Client{client: client}
}

func (c *Client) ByIC(ctx context.Context, id string) (result *EconomicEntity, err error) {
	endpoint := EconomicEntityURL + id

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request to ARES: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform HTTP request to ARES: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	switch resp.StatusCode {
	case http.StatusOK:
		// ok
	case http.StatusNotFound:
		return nil, fmt.Errorf("economic entity with IC '%s' not found", id)
	default:
		return nil, EntityNotFoundError{IC: id}
	}

	result = &EconomicEntity{}
	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return nil, fmt.Errorf("failed to decode response from ARES: %w", err)
	}

	return result, nil
}

func (a Address) String() string {
	var parts []string
	if v := strings.TrimSpace(a.Line1); v != "" {
		parts = append(parts, v)
	}
	if v := strings.TrimSpace(a.Line2); v != "" {
		parts = append(parts, v)
	}
	if v := strings.TrimSpace(a.Line3); v != "" {
		parts = append(parts, v)
	}
	return strings.Join(parts, "\n")
}
