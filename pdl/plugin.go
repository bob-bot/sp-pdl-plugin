package pdl

import (
	"context"
	"os"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// Plugin returns the steampipe plugin definition.
func Plugin(ctx context.Context) *plugin.Plugin {
	_ = ctx
	return &plugin.Plugin{
		Name: "steampipe-plugin-pdl",
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
			Schema:      ConfigSchema(),
		},
		TableMap: map[string]*plugin.Table{
			"pdl_person": tablePerson(),
		},
	}
}

// ConfigInstance returns a new instance of the plugin config.
func ConfigInstance() interface{} {
	return &Config{}
}

// ConfigSchema returns the schema for the plugin configuration.
func ConfigSchema() map[string]*plugin.Attribute {
	return map[string]*plugin.Attribute{
		"api_key": {Type: plugin.TypeString},
	}
}

// Config contains connection configuration for the plugin.
type Config struct {
	APIKey *string `hcl:"api_key"`
}

func (c *Config) apiKey(ctx context.Context, d *plugin.QueryData) (string, error) {
	if c != nil && c.APIKey != nil {
		return *c.APIKey, nil
	}

	if key := os.Getenv("PDL_API_KEY"); key != "" {
		return key, nil
	}

	return "", nil
}
