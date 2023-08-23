package kubernetes

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableHelmTemplateRendered(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "helm_template_rendered",
		Description: "Lists the fully rendered templates using the values provided in the config",
		List: &plugin.ListConfig{
			Hydrate: listHelmRenderedTemplates,
		},
		Columns: []*plugin.Column{
			{Name: "path", Type: proto.ColumnType_STRING, Description: "The path to the template file."},
			{Name: "chart_name", Type: proto.ColumnType_STRING, Description: "The name of the chart."},
			{Name: "source_type", Type: proto.ColumnType_STRING, Description: "The source of the template."},
			{Name: "rendered", Type: proto.ColumnType_STRING, Description: "Rendered is the rendered template as byte data."},
		},
	}
}

type helmTemplate struct {
	ChartName  string
	Path       string
	Rendered   string
	SourceType string
}

//// LIST FUNCTION

func listHelmRenderedTemplates(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	renderedTemplates, err := getHelmRenderedTemplates(ctx, d, nil)
	if err != nil {
		return nil, err
	}

	for _, template := range renderedTemplates {
		d.StreamListItem(ctx, helmTemplate{
			ChartName:  template.Chart.Metadata.Name,
			Path:       template.Path,
			Rendered:   template.Data,
			SourceType: fmt.Sprintf("helm_rendered:%s", template.ConfigKey),
		})

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}
