package kubernetes

import (
	"context"
	"os"
	"strings"

	v1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/kubernetes/scheme"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableKubernetesManifestDeployment(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_manifest_deployment",
		Description: "",
		List: &plugin.ListConfig{
			ParentHydrate: listKubernetesManifestFiles,
			Hydrate:       listKubernetesManifestDeployments,
			KeyColumns:    plugin.OptionalColumns([]string{"manifest_file_path"}),
		},
		Columns: []*plugin.Column{
			// Metadata Columns
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of resource."},
			{Name: "namespace", Type: proto.ColumnType_STRING, Description: "Namespace defines the space within which each name must be unique."},
			{Name: "labels", Type: proto.ColumnType_JSON, Description: "Map of string keys and values that can be used to organize and categorize (scope and select) objects. May match selectors of replication controllers and services."},

			// Spec Columns
			{Name: "replicas", Type: proto.ColumnType_INT, Description: "Number of desired pods. Defaults to 1.", Transform: transform.FromField("Spec.Replicas")},
			{Name: "selector_query", Type: proto.ColumnType_STRING, Description: "A query string representation of the selector.", Transform: transform.FromField("Spec.Selector").Transform(labelSelectorToString)},
			{Name: "selector", Type: proto.ColumnType_JSON, Description: " Label selector for pods. A label selector is a label query over a set of resources.", Transform: transform.FromField("Spec.Selector")},
			{Name: "template", Type: proto.ColumnType_JSON, Description: "Template describes the pods that will be created.", Transform: transform.FromField("Spec.Template")},
			{Name: "strategy", Type: proto.ColumnType_JSON, Description: "The deployment strategy to use to replace existing pods with new ones.", Transform: transform.FromField("Spec.Strategy")},
			{Name: "min_ready_seconds", Type: proto.ColumnType_INT, Description: "Minimum number of seconds for which a newly created pod should be ready without any of its container crashing, for it to be considered available. Defaults to 0.", Transform: transform.FromField("Spec.MinReadySeconds")},
			{Name: "revision_history_limit", Type: proto.ColumnType_INT, Description: "The number of old ReplicaSets to retain to allow rollback.", Transform: transform.FromField("Spec.RevisionHistoryLimit")},
			{Name: "paused", Type: proto.ColumnType_BOOL, Description: "Indicates that the deployment is paused.", Transform: transform.FromField("Spec.Paused")},
			{Name: "progress_deadline_seconds", Type: proto.ColumnType_INT, Description: "The maximum time in seconds for a deployment to make progress before it is considered to be failed.", Transform: transform.FromField("Spec.ProgressDeadlineSeconds")},

			{Name: "manifest_file_path", Description: "Path to the file.", Type: proto.ColumnType_STRING},
		},
	}
}

type KubernetesManifestDeployment struct {
	v1.Deployment
	ManifestFilePath string
}

//// LIST FUNCTION

func listKubernetesManifestDeployments(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// The path comes from a parent hydrate, defaulting to the config paths or
	// available by the optional key column
	path := h.Item.(filePath).Path

	// Load the file into a buffer
	content, err := os.ReadFile(path)
	if err != nil {
		plugin.Logger(ctx).Error("kubernetes_manifest_deployment.listKubernetesManifestDeployments", "failed to read file", err, "path", path)
		return nil, err
	}
	decoder := scheme.Codecs.UniversalDeserializer()

	for _, resourceYAML := range strings.Split(string(content), "---") {
		// skip empty documents, `Decode` will fail on them
		if len(resourceYAML) == 0 {
			continue
		}

		// Decode the file content
		obj, groupVersionKind, err := decoder.Decode([]byte(resourceYAML), nil, nil)
		if err != nil {
			plugin.Logger(ctx).Error("kubernetes_manifest_deployment.listKubernetesManifestDeployments", "failed to decode the file", err, "path", path)
			return nil, err
		}

		if groupVersionKind.Kind == "Deployment" {
			deployment := obj.(*v1.Deployment)

			d.StreamListItem(ctx, KubernetesManifestDeployment{*deployment, path})

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}
