package kubernetes

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableKubernetesKotsConfig(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_kots_config",
		Description: "KOTS application configuration. Retrieves config values from the KOTS admin console for a given application and sequence.",
		List: &plugin.ListConfig{
			Hydrate: listKotsConfig,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "app_slug", Require: plugin.Required},
				{Name: "namespace", Require: plugin.Optional},
				{Name: "sequence", Require: plugin.Optional},
			},
		},
		Columns: []*plugin.Column{
			// Key columns
			{Name: "app_slug", Type: proto.ColumnType_STRING, Description: "The slug identifier of the KOTS application.", Transform: transform.FromField("AppSlug")},
			{Name: "namespace", Type: proto.ColumnType_STRING, Description: "The Kubernetes namespace where kotsadm is running.", Transform: transform.FromField("Namespace")},
			{Name: "sequence", Type: proto.ColumnType_INT, Description: "The application sequence number for this config.", Transform: transform.FromField("Sequence")},

			// Config group fields
			{Name: "group_name", Type: proto.ColumnType_STRING, Description: "The name of the config group.", Transform: transform.FromField("GroupName")},
			{Name: "group_title", Type: proto.ColumnType_STRING, Description: "The display title of the config group.", Transform: transform.FromField("GroupTitle")},
			{Name: "group_description", Type: proto.ColumnType_STRING, Description: "The description of the config group.", Transform: transform.FromField("GroupDescription")},

			// Config item fields
			{Name: "item_name", Type: proto.ColumnType_STRING, Description: "The name of the config item.", Transform: transform.FromField("ItemName")},
			{Name: "item_type", Type: proto.ColumnType_STRING, Description: "The type of the config item (e.g., text, password, bool, select_one).", Transform: transform.FromField("ItemType")},
			{Name: "item_title", Type: proto.ColumnType_STRING, Description: "The display title of the config item.", Transform: transform.FromField("ItemTitle")},
			{Name: "item_value", Type: proto.ColumnType_STRING, Description: "The current value of the config item.", Transform: transform.FromField("ItemValue")},
			{Name: "item_default", Type: proto.ColumnType_STRING, Description: "The default value of the config item.", Transform: transform.FromField("ItemDefault")},
			{Name: "item_filename", Type: proto.ColumnType_STRING, Description: "The filename associated with the config item, if any.", Transform: transform.FromField("ItemFilename")},
			{Name: "item_hidden", Type: proto.ColumnType_BOOL, Description: "Whether the config item is hidden.", Transform: transform.FromField("ItemHidden")},
			{Name: "item_read_only", Type: proto.ColumnType_BOOL, Description: "Whether the config item is read-only.", Transform: transform.FromField("ItemReadOnly")},
			{Name: "item_help_text", Type: proto.ColumnType_STRING, Description: "Help text for the config item.", Transform: transform.FromField("ItemHelpText")},

			// Context
			{Name: "context_name", Type: proto.ColumnType_STRING, Description: "Kubectl config context name.", Hydrate: getKotsConfigContext},
		},
	}
}

type KotsConfigRow struct {
	AppSlug          string
	Namespace        string
	Sequence         int64
	GroupName        string
	GroupTitle       string
	GroupDescription string
	ItemName         string
	ItemType         string
	ItemTitle        string
	ItemValue        string
	ItemDefault      string
	ItemFilename     string
	ItemHidden       bool
	ItemReadOnly     bool
	ItemHelpText     string
}

func listKotsConfig(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listKotsConfig")

	appSlug := d.EqualsQualString("app_slug")
	if appSlug == "" {
		return nil, nil
	}

	namespaces, err := getKotsNamespaces(ctx, d)
	if err != nil {
		logger.Error("listKotsConfig", "namespace_discovery_error", err)
		return nil, err
	}

	for _, namespace := range namespaces {
		session, err := getKotsSession(ctx, d, namespace)
		if err != nil {
			logger.Warn("listKotsConfig", "namespace", namespace, "session_error", err)
			continue
		}

		// Determine the sequence to fetch config for
		var sequence int64 = -1
		if d.EqualsQuals["sequence"] != nil {
			sequence = d.EqualsQuals["sequence"].GetInt64Value()
		}

		// If no sequence specified, get the current sequence from the app
		if sequence == -1 {
			apps, err := getKotsApps(session)
			if err != nil {
				logger.Warn("listKotsConfig", "namespace", namespace, "apps_error", err)
				continue
			}

			for _, app := range apps.Apps {
				if app.Slug == appSlug {
					sequence = app.CurrentSequence
					break
				}
			}

			if sequence == -1 {
				continue
			}
		}

		config, err := getKotsConfig(session, appSlug, sequence)
		if err != nil {
			logger.Warn("listKotsConfig", "namespace", namespace, "api_error", err)
			continue
		}

		for _, group := range config.ConfigGroups {
			for _, item := range group.Items {
				d.StreamListItem(ctx, KotsConfigRow{
					AppSlug:          appSlug,
					Namespace:        namespace,
					Sequence:         sequence,
					GroupName:        group.Name,
					GroupTitle:       group.Title,
					GroupDescription: group.Description,
					ItemName:         item.Name,
					ItemType:         item.Type,
					ItemTitle:        item.Title,
					ItemValue:        item.Value,
					ItemDefault:      item.Default,
					ItemFilename:     item.Filename,
					ItemHidden:       item.Hidden,
					ItemReadOnly:     item.ReadOnly,
					ItemHelpText:     item.HelpText,
				})

				if d.RowsRemaining(ctx) == 0 {
					return nil, nil
				}
			}
		}
	}

	return nil, nil
}

func getKotsConfigContext(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	currentContext, err := getKubectlContext(ctx, d, nil)
	if err != nil {
		return nil, nil
	}
	return currentContext, nil
}
