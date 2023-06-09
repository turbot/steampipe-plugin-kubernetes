package kubernetes

import (
	"context"
	"path"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableHelmTemplates(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "helm_template",
		Description: "Lists the raw templates defined in the configured charts",
		List: &plugin.ListConfig{
			Hydrate: listHelmTemplates,
		},
		Columns: []*plugin.Column{
			{Name: "chart_name", Type: proto.ColumnType_STRING, Description: "The name of the chart."},
			{Name: "path", Type: proto.ColumnType_STRING, Description: "The path to the template file."},
			{Name: "raw", Type: proto.ColumnType_STRING, Description: "Data is the template as byte data."},
		},
	}
}

type helmTemplateRaw struct {
	ChartName string
	Path      string
	Raw       string
}

//// LIST FUNCTION

func listHelmTemplates(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	charts, err := getUniqueHelmCharts(ctx, d)
	if err != nil {
		return nil, err
	}

	for _, chart := range charts {
		for _, template := range chart.Chart.Templates {
			d.StreamListItem(ctx, helmTemplateRaw{
				ChartName: chart.Chart.Metadata.Name,
				Raw:       string(template.Data),
				Path:      path.Join(chart.Path, template.Name),
			})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}
