package kubernetes

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableHelmTemplateRaw(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "helm_template_raw",
		Description: "Templates defines in a specific chart directory",
		List: &plugin.ListConfig{
			Hydrate: listHelmRawTemplates,
		},
		Columns: []*plugin.Column{
			{Name: "chart_path", Type: proto.ColumnType_STRING, Description: "The path to the chart directory."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name is the path-like name of the template."},
			{Name: "raw", Type: proto.ColumnType_STRING, Description: "Data is the template as byte data."},
			{Name: "source_type", Type: proto.ColumnType_STRING, Description: "The source of the template."},
		},
	}
}

type helmTemplateRaw struct {
	ChartPath  string
	Name       string
	SourceType string
	Raw        string
}

//// LIST FUNCTION

func listHelmRawTemplates(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	charts, err := getParsedHelmChart(ctx, d)
	if err != nil {
		return nil, err
	}

	for _, chart := range charts {
		for _, template := range chart.Chart.Templates {
			d.StreamListItem(ctx, helmTemplateRaw{
				ChartPath:  chart.Path,
				SourceType: "helm",
				Raw:        string(template.Data),
				Name:       template.Name,
			})
		}
	}

	return nil, nil
}
