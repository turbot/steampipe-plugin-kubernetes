package kubernetes

import (
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/schema"
)

type kubernetesConfig struct {
	ConfigPaths   []string `cty:"config_paths"`
	ConfigPath    *string  `cty:"config_path"`
	ConfigContext *string  `cty:"config_context"`
}

var ConfigSchema = map[string]*schema.Attribute{
	"config_paths": {
		Type: schema.TypeList,
		Elem: &schema.Attribute{Type: schema.TypeString},
	},
	"config_path": {
		Type: schema.TypeString,
	},
	"config_context": {
		Type: schema.TypeString,
	},
}

func ConfigInstance() interface{} {
	return &kubernetesConfig{}
}

// GetConfig :: retrieve and cast connection config from query data
func GetConfig(connection *plugin.Connection) kubernetesConfig {
	if connection == nil || connection.Config == nil {
		return kubernetesConfig{}
	}
	config, _ := connection.Config.(kubernetesConfig)
	return config
}
