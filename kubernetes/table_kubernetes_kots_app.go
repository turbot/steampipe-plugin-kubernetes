package kubernetes

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableKubernetesKotsApp(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_kots_app",
		Description: "KOTS applications installed in the cluster. Retrieves app information and status from the KOTS admin console.",
		List: &plugin.ListConfig{
			Hydrate: listKotsApps,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "namespace", Require: plugin.Optional},
			},
		},
		Columns: []*plugin.Column{
			{Name: "namespace", Type: proto.ColumnType_STRING, Description: "The Kubernetes namespace where kotsadm is running.", Transform: transform.FromField("Namespace")},
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The unique identifier of the application.", Transform: transform.FromField("App.ID")},
			{Name: "slug", Type: proto.ColumnType_STRING, Description: "The slug identifier of the application.", Transform: transform.FromField("App.Slug")},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The display name of the application.", Transform: transform.FromField("App.Name")},
			{Name: "state", Type: proto.ColumnType_STRING, Description: "The runtime state of the application (e.g., ready, degraded, unavailable).", Transform: transform.FromField("State")},
			{Name: "is_airgap", Type: proto.ColumnType_BOOL, Description: "Whether the application is installed in airgap mode.", Transform: transform.FromField("App.IsAirgap")},
			{Name: "current_sequence", Type: proto.ColumnType_INT, Description: "The sequence number of the latest available version.", Transform: transform.FromField("App.CurrentSequence")},
			{Name: "upstream_uri", Type: proto.ColumnType_STRING, Description: "The upstream URI for the application.", Transform: transform.FromField("App.UpstreamURI")},
			{Name: "icon_uri", Type: proto.ColumnType_STRING, Description: "The URI of the application icon.", Transform: transform.FromField("App.IconURI")},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "The time when the application was created.", Transform: transform.FromField("App.CreatedAt")},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "The time when the application was last updated.", Transform: transform.FromField("App.UpdatedAt")},
			{Name: "last_update_check_at", Type: proto.ColumnType_TIMESTAMP, Description: "The time when updates were last checked.", Transform: transform.FromField("App.LastUpdateCheckAt")},
			{Name: "has_preflight", Type: proto.ColumnType_BOOL, Description: "Whether the application has preflight checks.", Transform: transform.FromField("App.HasPreflight")},
			{Name: "is_configurable", Type: proto.ColumnType_BOOL, Description: "Whether the application has configurable settings.", Transform: transform.FromField("App.IsConfigurable")},
			{Name: "update_checker_spec", Type: proto.ColumnType_STRING, Description: "The cron spec for automatic update checking.", Transform: transform.FromField("App.UpdateCheckerSpec")},
			{Name: "auto_deploy", Type: proto.ColumnType_STRING, Description: "The auto-deploy policy for the application.", Transform: transform.FromField("App.AutoDeploy")},
			{Name: "license_type", Type: proto.ColumnType_STRING, Description: "The license type of the application.", Transform: transform.FromField("App.LicenseType")},
			{Name: "allow_rollback", Type: proto.ColumnType_BOOL, Description: "Whether rollback is allowed for this application.", Transform: transform.FromField("App.AllowRollback")},
			{Name: "allow_snapshots", Type: proto.ColumnType_BOOL, Description: "Whether snapshots are allowed for this application.", Transform: transform.FromField("App.AllowSnapshots")},
			{Name: "target_kots_version", Type: proto.ColumnType_STRING, Description: "The target KOTS version for this application.", Transform: transform.FromField("App.TargetKotsVersion")},
			{Name: "is_semver_required", Type: proto.ColumnType_BOOL, Description: "Whether semantic versioning is required.", Transform: transform.FromField("App.IsSemverRequired")},
			{Name: "current_version_label", Type: proto.ColumnType_STRING, Description: "The version label of the currently deployed version.", Transform: transform.FromField("CurrentVersionLabel")},
			{Name: "downstream_name", Type: proto.ColumnType_STRING, Description: "The name of the downstream cluster.", Transform: transform.FromField("App.Downstream.Name")},
			{Name: "context_name", Type: proto.ColumnType_STRING, Description: "Kubectl config context name.", Transform: transform.FromField("ContextName")},
		},
	}
}

type KotsAppRow struct {
	Namespace           string
	ContextName         string
	State               string
	CurrentVersionLabel string
	App                 *KotsApp
}

func listKotsApps(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listKotsApps")

	// Resolve context name once
	contextName := ""
	if currentContext, err := getKubectlContext(ctx, d, nil); err == nil && currentContext != nil {
		contextName = currentContext.(string)
	}

	namespaces, err := getKotsNamespaces(ctx, d)
	if err != nil {
		logger.Error("listKotsApps", "namespace_discovery_error", err)
		return nil, err
	}

	for _, namespace := range namespaces {
		session, err := getKotsSession(ctx, d, namespace)
		if err != nil {
			logger.Warn("listKotsApps", "namespace", namespace, "session_error", err)
			continue
		}

		apps, err := getKotsApps(session)
		if err != nil {
			logger.Warn("listKotsApps", "namespace", namespace, "api_error", err)
			continue
		}

		for i := range apps.Apps {
			app := &apps.Apps[i]

			// Get the runtime state via the status API
			state := ""
			status, err := getKotsAppStatus(session, app.Slug)
			if err != nil {
				logger.Warn("listKotsApps", "namespace", namespace, "app", app.Slug, "status_error", err)
			} else if status.AppStatus != nil {
				state = status.AppStatus.State
			}

			// Extract current version label from downstream
			currentVersionLabel := ""
			if app.Downstream.CurrentVersion != nil {
				currentVersionLabel = app.Downstream.CurrentVersion.VersionLabel
			}

			d.StreamListItem(ctx, KotsAppRow{
				Namespace:           namespace,
				ContextName:         contextName,
				State:               state,
				CurrentVersionLabel: currentVersionLabel,
				App:                 app,
			})

			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}
