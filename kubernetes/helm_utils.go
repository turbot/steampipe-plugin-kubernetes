package kubernetes

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"slices"

	"gopkg.in/yaml.v3"

	helmClient "github.com/mittwald/go-helm-client"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// Utils functions from Helm charts

type parsedHelmChart struct {
	Chart *chart.Chart
	Path  string
}

// Get the parsed contents of the given Helm chart.
func getParsedHelmChart(ctx context.Context, d *plugin.QueryData) ([]*parsedHelmChart, error) {
	conn, err := parsedHelmChartCached(ctx, d, nil)
	if err != nil {
		return nil, err
	}

	if conn != nil {
		return conn.([]*parsedHelmChart), nil
	}
	return nil, nil
}

// Cached form of the parsed Helm chart.
var parsedHelmChartCached = plugin.HydrateFunc(parsedHelmChartUncached).Memoize()

// parsedHelmChartUncached is the actual implementation of getParsedHelmChart, which should
// be run only once per connection. Do not call this directly, use
// getParsedHelmChart instead.
func parsedHelmChartUncached(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	// Read the config
	kubernetesConfig := GetConfig(d.Connection)

	// Check for the sourceTypes argument in the config.
	// Default set to include values.
	var sources = All.ToSourceTypes()
	if kubernetesConfig.SourceTypes != nil {
		sources = kubernetesConfig.SourceTypes
	}
	// TODO: Remove once `SourceType` is obsolete
	if kubernetesConfig.SourceTypes == nil && kubernetesConfig.SourceType != nil {
		if *kubernetesConfig.SourceType != "all" { // if is all, sources is already set by default
			sources = []string{*kubernetesConfig.SourceType}
		}
	}

	if !slices.Contains(sources, "helm") {
		return nil, nil
	}

	var charts []*parsedHelmChart

	for _, v := range kubernetesConfig.HelmRenderedCharts {
		// Return error if source_types arg includes "helm" in the config, but
		// helm_chart_dir arg is not set.
		if v.ChartPath == "" {
			return nil, errors.New("helm_chart_dir must be set in the config while source_types includes 'helm'")
		}

		// Return empty parsedHelmChart object if no Helm chart directory path provided in the config
		chartDir := v.ChartPath
		if chartDir == "" {
			plugin.Logger(ctx).Debug("parsedHelmChartUncached", "helm_chart_dir not configured in the config", "connection", d.Connection.Name)
			return nil, nil
		}
		plugin.Logger(ctx).Debug("parsedHelmChartUncached", "Parsing Helm chart", chartDir, "connection", d.Connection.Name)

		// Load the given chart directory
		chart, err := loader.Load(chartDir)
		if err != nil {
			plugin.Logger(ctx).Error("parsedHelmChartUncached", "load_chart_error", err)
			return nil, err
		}

		charts = append(charts, &parsedHelmChart{
			Chart: chart,
			Path:  chartDir,
		})
	}

	return charts, nil
}

// getUniqueHelmCharts scans all the charts configured in the config and returns a list of unique charts
func getUniqueHelmCharts(ctx context.Context, d *plugin.QueryData) ([]*parsedHelmChart, error) {
	var uniqueCharts []*parsedHelmChart
	var configuredChartPaths []string

	charts, err := getParsedHelmChart(ctx, d)
	if err != nil {
		return nil, err
	}

	for _, chart := range charts {
		if !slices.Contains(configuredChartPaths, chart.Path) {
			uniqueCharts = append(uniqueCharts, chart)
		}
		configuredChartPaths = append(configuredChartPaths, chart.Path)
	}

	return uniqueCharts, nil
}

// getUniqueValueFilesFromConfig scans all the values files provided in the chart and returns a unique set of value files from it
func getUniqueValueFilesFromConfig(ctx context.Context, d *plugin.QueryData) []string {
	var filePaths []string
	config := GetConfig(d.Connection)

	for _, chart := range config.HelmRenderedCharts {
		for _, path := range chart.ValuesFilePaths {
			if !slices.Contains(filePaths, path) {
				filePaths = append(filePaths, path)
			}
		}
	}
	return filePaths
}

// getHelmClient creates the client for Helm
func getHelmClient(ctx context.Context, namespace string) (helmClient.Client, error) {
	// Return nil if no namespace provided
	if namespace == "" {
		return nil, nil
	}

	// Set the namespace if specified.
	// By default current namespace context is used.
	options := &helmClient.Options{
		Namespace: namespace,
	}

	// Create client
	client, err := helmClient.New(options)
	if err != nil {
		plugin.Logger(ctx).Error("getHelmClient", "client_error", err)
		return nil, err
	}

	return client, nil
}

// getRows takes the chart values values as input and returns all the keys and values in tree structure
func getRows(ctx context.Context, values map[string]interface{}) (Rows, error) {
	var root yaml.Node
	buf := new(bytes.Buffer)
	if err := yaml.NewEncoder(buf).Encode(values); err != nil {
		return nil, err
	}

	decoder := yaml.NewDecoder(buf)
	err := decoder.Decode(&root)
	if err != nil {
		return nil, fmt.Errorf("failed to decode content: %v", err)
	}

	var rows Rows
	treeToList(&root, []string{}, &rows, nil, nil, nil)

	return rows, nil
}
