package kubernetes

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

type kubernetesConfig struct {
	ConfigPaths          []string               `hcl:"config_paths,optional"`
	ConfigPath           *string                `hcl:"config_path"`
	ConfigContext        *string                `hcl:"config_context"`
	CustomResourceTables []string               `hcl:"custom_resource_tables,optional"`
	ManifestFilePaths    []string               `hcl:"manifest_file_paths,optional" steampipe:"watch"`
	SourceType           *string                `hcl:"source_type"`
	SourceTypes          []string               `hcl:"source_types,optional"`
	HelmRenderedCharts   map[string]chartConfig `hcl:"helm_rendered_charts,optional"`
}

type chartConfig struct {
	ChartPath       string   `hcl:"chart_path" cty:"chart_path"`
	ValuesFilePaths []string `hcl:"values_file_paths,optional" cty:"values_file_paths"`
}

// GetConfig :: retrieve and cast connection config from query data
func GetConfig(connection *plugin.Connection) kubernetesConfig {
	if connection == nil || connection.Config == nil {
		return kubernetesConfig{}
	}
	config, _ := connection.Config.(kubernetesConfig)
	return config
}
