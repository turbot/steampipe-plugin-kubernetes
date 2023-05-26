package kubernetes

import (
	"context"
	"strings"
	"time"

	"helm.sh/helm/v3/pkg/action"
	helmTime "helm.sh/helm/v3/pkg/time"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableHelmRelease(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "helm_release",
		Description: "List of Helm releases for a specified namespace",
		List: &plugin.ListConfig{
			Hydrate: listHelmReleases,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "namespace", Require: plugin.Optional},
				{Name: "status", Require: plugin.Optional},
			},
		},
		Get: &plugin.GetConfig{
			Hydrate: getHelmRelease,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "name", Require: plugin.Required},
				{Name: "namespace", Require: plugin.Optional},
			},
		},
		Columns: []*plugin.Column{
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the release."},
			{Name: "namespace", Type: proto.ColumnType_STRING, Description: "The kubernetes namespace of the release."},
			{Name: "version", Type: proto.ColumnType_INT, Description: "The revision of the release."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "The current state of the release. Possible values: deployed, failed, pending-install, pending-rollback, pending-upgrade, superseded, uninstalled, uninstalling, unknown.", Transform: transform.FromField("Info.Status").Transform(transform.ToString)},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "A human-friendly description about the release.", Transform: transform.FromField("Info.Description")},
			{Name: "chart_name", Type: proto.ColumnType_STRING, Description: "The name of the chart that was released.", Transform: transform.FromField("Chart.Metadata.Name")},
			{Name: "first_deployed", Type: proto.ColumnType_TIMESTAMP, Description: "The time when the release was first deployed.", Transform: transform.FromField("Info.FirstDeployed").Transform(parseDateStringToTime)},
			{Name: "last_deployed", Type: proto.ColumnType_TIMESTAMP, Description: "The time when the release was last deployed.", Transform: transform.FromField("Info.LastDeployed").Transform(parseDateStringToTime)},
			{Name: "deleted", Type: proto.ColumnType_TIMESTAMP, Description: "The time when this object was deleted.", Transform: transform.FromField("Info.Deleted").Transform(parseDateStringToTime)},
			{Name: "notes", Type: proto.ColumnType_STRING, Description: "Contains the rendered templates/NOTES.txt if available.", Transform: transform.FromField("Info.Notes")},
			{Name: "config", Type: proto.ColumnType_JSON, Description: "The set of extra Values added to the chart. These values override the default values inside of the chart."},
			{Name: "labels", Type: proto.ColumnType_JSON, Description: "The labels of the release."},
			{Name: "manifest", Type: proto.ColumnType_STRING, Description: "The string representation of the rendered template."},
		},
	}
}

//// LIST FUNCTION

func listHelmReleases(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	namespace := d.EqualsQualString("namespace")

	// Create client
	client, err := getHelmClient(ctx, namespace)
	if err != nil {
		return nil, err
	}

	chart, err := getParsedHelmChart(ctx, d)
	if err != nil {
		return nil, err
	}

	// Return nil, if no chart configured in the connection
	if chart == nil {
		return nil, nil
	}

	for _, c := range chart {
		releaseState := action.ListAll
		if d.EqualsQuals["status"] != nil {
			givenState := d.EqualsQualString("status")
			releaseState = action.ListAll.FromName(givenState)
		}

		// Lists all releases for a specified namespace. By default it uses current namespace context, if nothing is set
		releases, err := client.ListReleasesByStateMask(releaseState)
		if err != nil {
			return nil, err
		}

		for _, release := range releases {
			if release.Chart.Metadata.Name != c.Chart.Metadata.Name {
				continue
			}
			d.StreamListItem(ctx, release)

			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getHelmRelease(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	namespace := d.EqualsQualString("namespace")
	releaseName := d.EqualsQualString("name")

	// Return nil, if empty
	if releaseName == "" {
		return nil, nil
	}

	// Create client
	client, err := getHelmClient(ctx, namespace)
	if err != nil {
		return nil, err
	}

	release, err := client.GetRelease(releaseName)
	if err != nil {
		// Return nil, if the requested resource is not present
		if strings.Contains(err.Error(), "not found") {
			return nil, nil
		}
		return nil, err
	}

	return release, nil
}

//// TRANSFORM FUNCTIONS

func parseDateStringToTime(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	if d.Value != nil {
		data := d.Value.(helmTime.Time)
		return data.Format(time.RFC3339), nil
	}
	return nil, nil
}
