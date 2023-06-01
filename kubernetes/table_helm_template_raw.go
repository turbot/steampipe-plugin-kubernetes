package kubernetes

import (
	"context"
	"path"

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
			{Name: "chart_name", Type: proto.ColumnType_STRING, Description: "The name of the chart."},
			{Name: "path", Type: proto.ColumnType_STRING, Description: "The path to the template file."},
			{Name: "raw", Type: proto.ColumnType_STRING, Description: "Data is the template as byte data."},
		},
	}
}

type helmTemplateRaw struct {
	ChartName  string
	Path       string
	SourceType string
	Raw        string
}

//// LIST FUNCTION

func listHelmRawTemplates(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	charts, err := getUniqueHelmCharts(ctx, d)
	if err != nil {
		return nil, err
	}

	for _, chart := range charts {
		for _, template := range chart.Chart.Templates {
			d.StreamListItem(ctx, helmTemplateRaw{
				ChartName:  chart.Chart.Metadata.Name,
				SourceType: "helm",
				Raw:        string(template.Data),
				Path:       path.Join(chart.Path, template.Name),
			})
		}
	}

	return nil, nil
}
