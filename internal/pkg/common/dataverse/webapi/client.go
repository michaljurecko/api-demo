package webapi

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

const (
	APIVersion      = "9.2"
	ContentTypeJSON = "application/json"
)

// Client is an HTTP client for Dataverse OData API with OAuth authentication middleware.
type Client struct {
	config  Config
	baseURL *url.URL
	client  *http.Client
}

func ConfigFromENV() Config {
	return Config{
		TenantID:     os.Getenv("DEMO_MODEL_TENANT_ID"),
		ClientID:     os.Getenv("DEMO_MODEL_CLIENT_ID"),
		ClientSecret: os.Getenv("DEMO_MODEL_CLIENT_SECRET"),
		APIHost:      os.Getenv("DEMO_MODEL_API_HOST"),
	}
}

func NewClient(ctx context.Context, config Config, baseClient *http.Client) (*Client, error) {
	// Create configuration for OAuth client credentials flow.
	oauthConfig := clientcredentials.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		TokenURL:     "https://login.microsoftonline.com/" + config.TenantID + "/oauth2/v2.0/token",
		Scopes:       []string{"https://" + config.APIHost + "/.default"},
	}

	baseURL, err := url.Parse("https://" + config.APIHost + "/api/data/v" + APIVersion + "/")
	if err != nil {
		return nil, fmt.Errorf(`invalid api host "%s": %w`, config.APIHost, err)
	}

	// Inject base HTTP client, so the http.Transport/http.RoundTripper can be modified in tests.
	ctx = context.WithValue(ctx, oauth2.HTTPClient, baseClient)

	return &Client{
		config:  config,
		baseURL: baseURL,
		client:  oauthConfig.Client(ctx),
	}, nil
}
