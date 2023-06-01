package kubernetes

import (
	"context"

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
