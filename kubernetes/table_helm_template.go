package kubernetes

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableHelmTemplate(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "helm_template",
		Description: "Templates defines in a specific chart directory",
		List: &plugin.ListConfig{
			Hydrate: listHelmTemplates,
		},
		Columns: []*plugin.Column{
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name is the path-like name of the template."},
			{Name: "rendered", Type: proto.ColumnType_STRING, Description: "Data is the template as byte data."},
			{Name: "raw", Type: proto.ColumnType_STRING, Description: "Data is the template as byte data."},
			{Name: "chart_name", Type: proto.ColumnType_STRING, Description: "The name of the chart."},
		},
	}
}

type helmTemplate struct {
	// Path string
	ChartName string
	Name      string
	Rendered  string
	Raw       string
}

//// LIST FUNCTION

func listHelmTemplates(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	chart, err := getParsedHelmChart(ctx, d)
	if err != nil {
		return nil, err
	}

	renderedTemplates, err := getHelmRenderedTemplates(ctx, d, nil)
	if err != nil {
		return nil, err
	}

	for _, template := range chart.Chart.Templates {
		for k, v := range renderedTemplates {
			if strings.HasSuffix(k, template.Name) {
				d.StreamListItem(ctx, helmTemplate{
					ChartName: chart.Chart.Metadata.Name,
					Name:      k,
					Rendered:  v,
					Raw:       string(template.Data),
				})
			}
		}
	}

	return nil, nil
}