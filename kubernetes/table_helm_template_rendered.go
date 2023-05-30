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
		Description: "Templates defines in a specific chart directory",
		List: &plugin.ListConfig{
			Hydrate: listHelmRenderedTemplates,
		},
		Columns: []*plugin.Column{
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name is the path-like name of the template."},
			{Name: "source_type", Type: proto.ColumnType_STRING, Description: "The source of the template."},
			{Name: "rendered", Type: proto.ColumnType_STRING, Description: "Data is the template as byte data."},
			{Name: "chart_name", Type: proto.ColumnType_STRING, Description: "The name of the chart."},
		},
	}
}

type helmTemplate struct {
	// Path string
	ChartName  string
	Name       string
	Rendered   string
	SourceType string
}

//// LIST FUNCTION

func listHelmRenderedTemplates(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// renderedTemplates, err := getHelmRenderedTemplates(ctx, d, nil)
	// if err != nil {
	// 	return nil, err
	// }

	// for _, template := range renderedTemplates {
	// 	for k, v := range template.Data {
	// 		d.StreamListItem(ctx, helmTemplate{
	// 			ChartName:  template.Chart.Metadata.Name,
	// 			Name:       k,
	// 			Rendered:   v,
	// 			SourceType: fmt.Sprintf("helm_rendered:%s", template.Name),
	// 		})
	// 	}
	// }

	renderedTemplates, err := getHelmTemplatesUsingKics(ctx, d, nil)
	if err != nil {
		return nil, err
	}

	for _, template := range renderedTemplates {
		d.StreamListItem(ctx, helmTemplate{
			ChartName:  template.Chart.Metadata.Name,
			Name:       template.TemplateName,
			Rendered:   template.Data,
			SourceType: fmt.Sprintf("helm_rendered:%s", template.Name),
		})
	}

	return nil, nil
}
