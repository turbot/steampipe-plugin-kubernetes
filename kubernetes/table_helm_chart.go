package kubernetes

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"helm.sh/helm/v3/pkg/chart"
)

//// TABLE DEFINITION

func tableHelmChart(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "helm_chart",
		Description: "Lists the configuration settings from the configured charts",
		List: &plugin.ListConfig{
			Hydrate: listHelmCharts,
		},
		Columns: []*plugin.Column{
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the chart."},
			{Name: "api_version", Type: proto.ColumnType_STRING, Description: "The API Version of the chart.", Transform: transform.FromField("APIVersion")},
			{Name: "version", Type: proto.ColumnType_STRING, Description: "A SemVer 2 conformant version string of the chart."},
			{Name: "app_version", Type: proto.ColumnType_STRING, Description: "The version of the application enclosed inside of this chart."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "A one-sentence description of the chart."},
			{Name: "deprecated", Type: proto.ColumnType_BOOL, Description: "Indicates whether or not this chart is deprecated."},
			{Name: "home", Type: proto.ColumnType_STRING, Description: "The URL to a relevant project page, git repo, or contact person."},
			{Name: "icon", Type: proto.ColumnType_STRING, Description: "The URL to an icon file."},
			{Name: "condition", Type: proto.ColumnType_STRING, Description: "The condition to check to enable chart."},
			{Name: "tags", Type: proto.ColumnType_STRING, Description: "The tags to check to enable chart."},
			{Name: "kube_version", Type: proto.ColumnType_STRING, Description: "A SemVer constraint specifying the version of Kubernetes required."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "Specifies the chart type. Possible values: application, or library."},
			{Name: "sources", Type: proto.ColumnType_JSON, Description: "Source is the URL to the source code of this chart."},
			{Name: "keywords", Type: proto.ColumnType_JSON, Description: "A list of string keywords."},
			{Name: "maintainers", Type: proto.ColumnType_JSON, Description: "A list of name and URL/email address combinations for the maintainer(s)."},
			{Name: "annotations", Type: proto.ColumnType_JSON, Description: "Annotations are additional mappings uninterpreted by Helm, made available for inspection by other applications."},
			{Name: "dependencies", Type: proto.ColumnType_JSON, Description: "Dependencies are a list of dependencies for a chart."},
			{Name: "chart_path", Type: proto.ColumnType_STRING, Description: "The path to the directory where the chart is located."},
		},
	}
}

type HelmChartInfo struct {
	chart.Metadata
	ChartPath string
}

//// LIST FUNCTION

func listHelmCharts(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Get the list of unique helm charts from the charts provided in the config
	charts, err := getUniqueHelmCharts(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listHelmCharts", "failed to list charts", err)
		return nil, err
	}

	for _, chart := range charts {
		d.StreamListItem(ctx, HelmChartInfo{*chart.Chart.Metadata, chart.Path})

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}
