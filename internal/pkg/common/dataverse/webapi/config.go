package webapi

// Config represents the configuration for the Dataverse OData client.
// It is a part of the application configuration.
type Config struct {
	TenantID      string `json:"-" mapstructure:"tenant_id"      validate:"required"`
	ClientID      string `json:"-" mapstructure:"client_id"      validate:"required"`
	ClientSecret  string `json:"-" mapstructure:"client_secret"  validate:"required"`
	APIHost       string `mapstructure:"api_host"       validate:"required"`
	DebugRequest  bool   `mapstructure:"debug_request"`
	DebugResponse bool   `mapstructure:"debug_response"`
}

func (c Config) IsEmpty() bool {
	return c.TenantID == "" || c.ClientID == "" || c.ClientSecret == "" || c.APIHost == ""
}
