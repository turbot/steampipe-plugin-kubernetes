package kubernetes

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableKubernetesKotsVersion(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_kots_version",
		Description: "KOTS application version history. Retrieves version information from the KOTS admin console running in the cluster.",
		List: &plugin.ListConfig{
			Hydrate: listKotsVersions,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "app_slug", Require: plugin.Required},
				{Name: "namespace", Require: plugin.Optional},
			},
		},
		Columns: []*plugin.Column{
			// Key columns
			{Name: "app_slug", Type: proto.ColumnType_STRING, Description: "The slug identifier of the KOTS application.", Transform: transform.FromField("AppSlug")},
			{Name: "namespace", Type: proto.ColumnType_STRING, Description: "The Kubernetes namespace where kotsadm is running.", Transform: transform.FromField("Namespace")},

			// Version fields
			{Name: "version_label", Type: proto.ColumnType_STRING, Description: "The version label of the application release.", Transform: transform.FromField("Version.VersionLabel")},
			{Name: "sequence", Type: proto.ColumnType_INT, Description: "The sequence number of this version.", Transform: transform.FromField("Version.Sequence")},
			{Name: "created_on", Type: proto.ColumnType_TIMESTAMP, Description: "The time when this version was created.", Transform: transform.FromField("Version.CreatedOn")},
			{Name: "deployed_at", Type: proto.ColumnType_TIMESTAMP, Description: "The time when this version was deployed.", Transform: transform.FromField("Version.DeployedAt")},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "The deployment status of this version (e.g., deployed, pending, failed).", Transform: transform.FromField("Version.Status")},
			{Name: "source", Type: proto.ColumnType_STRING, Description: "The source of this version.", Transform: transform.FromField("Version.Source")},
			{Name: "channel_id", Type: proto.ColumnType_STRING, Description: "The channel ID from which this version was obtained.", Transform: transform.FromField("Version.ChannelID")},
			{Name: "update_cursor", Type: proto.ColumnType_STRING, Description: "The update cursor (channel sequence) for this version.", Transform: transform.FromField("Version.UpdateCursor")},
			{Name: "is_required", Type: proto.ColumnType_BOOL, Description: "Whether this version is a required release.", Transform: transform.FromField("Version.IsRequired")},
			{Name: "is_deployable", Type: proto.ColumnType_BOOL, Description: "Whether this version can be deployed.", Transform: transform.FromField("Version.IsDeployable")},
			{Name: "non_deployable_cause", Type: proto.ColumnType_STRING, Description: "The reason this version cannot be deployed, if applicable.", Transform: transform.FromField("Version.NonDeployableCause")},
			{Name: "release_notes", Type: proto.ColumnType_STRING, Description: "Release notes for this version.", Transform: transform.FromField("Version.ReleaseNotes")},
			{Name: "preflight_skipped", Type: proto.ColumnType_BOOL, Description: "Whether preflight checks were skipped for this version.", Transform: transform.FromField("Version.PreflightSkipped")},
			{Name: "upstream_released_at", Type: proto.ColumnType_TIMESTAMP, Description: "The time when this version was released upstream.", Transform: transform.FromField("Version.UpstreamReleasedAt")},

			// Context
			{Name: "context_name", Type: proto.ColumnType_STRING, Description: "Kubectl config context name.", Transform: transform.FromField("ContextName")},
		},
	}
}

type KotsVersionRow struct {
	AppSlug     string
	Namespace   string
	ContextName string
	Version     *KotsDownstreamVersion
}

func listKotsVersions(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listKotsVersions")

	appSlug := d.EqualsQualString("app_slug")
	if appSlug == "" {
		return nil, nil
	}

	// Resolve context name once
	contextName := ""
	if currentContext, err := getKubectlContext(ctx, d, nil); err == nil && currentContext != nil {
		contextName = currentContext.(string)
	}

	namespaces, err := getKotsNamespaces(ctx, d)
	if err != nil {
		logger.Error("listKotsVersions", "namespace_discovery_error", err)
		return nil, err
	}

	for _, namespace := range namespaces {
		session, err := getKotsSession(ctx, d, namespace)
		if err != nil {
			logger.Warn("listKotsVersions", "namespace", namespace, "session_error", err)
			continue
		}

		versions, err := getKotsVersions(session, appSlug)
		if err != nil {
			logger.Warn("listKotsVersions", "namespace", namespace, "api_error", err)
			continue
		}

		for _, version := range versions.VersionHistory {
			d.StreamListItem(ctx, KotsVersionRow{
				AppSlug:     appSlug,
				Namespace:   namespace,
				ContextName: contextName,
				Version:     version,
			})

			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}
