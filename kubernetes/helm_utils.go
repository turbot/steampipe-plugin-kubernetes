package kubernetes

import (
	"context"
	"bytes"
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func getUniqueHelmCharts(ctx context.Context, d *plugin.QueryData) ([]*parsedHelmChart, error) {
	var uniqueCharts []*parsedHelmChart
	var configuredChartPaths []string

	charts, err := getParsedHelmChart(ctx, d)
	if err != nil {
		return nil, err
	}

	for _, chart := range charts {
		if !helpers.StringSliceContains(configuredChartPaths, chart.Path) {
			uniqueCharts = append(uniqueCharts, chart)
		}
		configuredChartPaths = append(configuredChartPaths, chart.Path)
	}

	return uniqueCharts, nil
}

func getUniqueValueFilesFromConfig(ctx context.Context, d *plugin.QueryData) []string {
	var filePaths []string
	config := GetConfig(d.Connection)

	for _, chart := range config.HelmRenderedCharts {
		for _, path := range chart.ValuesFilePaths {
			if !helpers.StringSliceContains(filePaths, path) {
				filePaths = append(filePaths, path)
			}
		}
	}
	return filePaths
}

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