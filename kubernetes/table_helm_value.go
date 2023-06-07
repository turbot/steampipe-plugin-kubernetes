package kubernetes

import (
	"context"
	"os"
	"path"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"gopkg.in/yaml.v3"
)

//// TABLE DEFINITION

func tableHelmValue(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "helm_value",
		Description: "Lists the values from chart's values.yaml file as well as the values listed in the configured values override files",
		List: &plugin.ListConfig{
			Hydrate: listHelmValues,
		},
		Columns: []*plugin.Column{
			{Name: "path", Type: proto.ColumnType_STRING, Description: "Name is the path-like name of the template."},
			{Name: "key_path", Type: proto.ColumnType_LTREE, Transform: transform.FromField("Key").Transform(keysToSnakeCase), Description: "Specifies full path of a key in YML file."},
			{Name: "value", Type: proto.ColumnType_STRING, Description: "Specifies the value of the corresponding key."},
			{Name: "keys", Type: proto.ColumnType_JSON, Transform: transform.FromField("Key"), Description: "The array representation of path of a key."},
			{Name: "start_line", Type: proto.ColumnType_INT, Description: "Specifies the line number where the value is located."},
			{Name: "start_column", Type: proto.ColumnType_INT, Description: "Specifies the starting column of the value."},
			{Name: "pre_comments", Type: proto.ColumnType_JSON, Description: "Specifies the comments added above a key."},
			{Name: "head_comment", Type: proto.ColumnType_STRING, Description: "Specifies the comment in the lines preceding the node and not separated by an empty line."},
			{Name: "line_comment", Type: proto.ColumnType_STRING, Description: "Specifies the comment at the end of the line where the node is in."},
			{Name: "foot_comment", Type: proto.ColumnType_STRING, Description: "Specifies the comment following the node and before empty lines."},
		},
	}
}

//// LIST FUNCTION

func listHelmValues(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Stream all values from chart's values.yaml file
	charts, err := getUniqueHelmCharts(ctx, d)
	if err != nil {
		return nil, err
	}

	var allValues Rows
	for _, chart := range charts {
		defaultValues, err := getRows(ctx, chart.Chart.Values)
		if err != nil {
			plugin.Logger(ctx).Error("helm_value.listHelmValues", "parse_error", err, "path", chart.Path)
			return nil, err
		}

		for _, r := range defaultValues {
			r.Path = path.Join(chart.Path, "values.yaml")
			allValues = append(allValues, r)
		}
	}

	// Stream values from the unique set of override value files provided in the config
	overrideValueFiles := getUniqueValueFilesFromConfig(ctx, d)
	for _, path := range overrideValueFiles {
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		var values map[string]interface{}
		err = yaml.Unmarshal(content, &values)
		if err != nil {
			return nil, err
		}

		overrideValues, err := getRows(ctx, values)
		if err != nil {
			return nil, err
		}

		for _, r := range overrideValues {
			r.Path = path
			allValues = append(allValues, r)
		}
	}

	// Stream all the values from chart's default values.yaml as well as the value override files configured in the config
	for _, v := range allValues {
		d.StreamListItem(ctx, v)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}
