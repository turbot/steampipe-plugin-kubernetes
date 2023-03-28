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

func tableKubernetesManifestDaemonSet(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_manifest_daemonset",
		Description: "",
		List: &plugin.ListConfig{
			ParentHydrate: listKubernetesManifestFiles,
			Hydrate:       listKubernetesManifestDaemonsets,
			KeyColumns:    plugin.OptionalColumns([]string{"manifest_file_path"}),
		},
		Columns: k8sCommonColumns([]*plugin.Column{

			//// DaemonSetSpec Columns
			{
				Name:        "min_ready_seconds",
				Type:        proto.ColumnType_INT,
				Description: "The minimum number of seconds for which a newly created DaemonSet pod should be ready without any of its container crashing, for it to be considered available. Defaults to 0",
				Transform:   transform.FromField("Spec.MinReadySeconds"),
			},
			{
				Name:        "revision_history_limit",
				Type:        proto.ColumnType_INT,
				Description: "The number of old history to retain to allow rollback. This is a pointer to distinguish between explicit zero and not specified. Defaults to 10.",
				Transform:   transform.FromField("Spec.RevisionHistoryLimit"),
			},
			{
				Name:        "selector_query",
				Type:        proto.ColumnType_STRING,
				Description: "A query string representation of the selector.",
				Transform:   transform.FromField("Spec.Selector").Transform(labelSelectorToString),
			},
			{
				Name:        "selector",
				Type:        proto.ColumnType_JSON,
				Description: "A label query over pods that are managed by the daemon set.",
				Transform:   transform.FromField("Spec.Volumes"),
			},
			{
				Name:        "template",
				Type:        proto.ColumnType_JSON,
				Description: "An object that describes the pod that will be created.",
				Transform:   transform.FromField("Spec.Template"),
			},
			{
				Name:        "update_strategy",
				Type:        proto.ColumnType_JSON,
				Description: "An update strategy to replace existing DaemonSet pods with new pods.",
				Transform:   transform.FromField("Spec.UpdateStrategy"),
			},

			//// DaemonSetStatus Columns
			{
				Name:        "current_number_scheduled",
				Type:        proto.ColumnType_INT,
				Description: "The number of nodes that are running at least 1 daemon pod and are supposed to run the daemon pod.",
				Transform:   transform.FromField("Status.CurrentNumberScheduled"),
			},
			{
				Name:        "number_misscheduled",
				Type:        proto.ColumnType_INT,
				Description: "The number of nodes that are running the daemon pod, but are not supposed to run the daemon pod.",
				Transform:   transform.FromField("Status.NumberMisscheduled"),
			},
			{
				Name:        "desired_number_scheduled",
				Type:        proto.ColumnType_INT,
				Description: "The total number of nodes that should be running the daemon pod (including nodes correctly running the daemon pod).",
				Transform:   transform.FromField("Status.DesiredNumberScheduled"),
			},
			{
				Name:        "number_ready",
				Type:        proto.ColumnType_INT,
				Description: "The number of nodes that should be running the daemon pod and have one or more of the daemon pod running and ready.",
				Transform:   transform.FromField("Status.NumberReady"),
			},
			{
				Name:        "observed_generation",
				Type:        proto.ColumnType_INT,
				Description: "The most recent generation observed by the daemon set controller.",
				Transform:   transform.FromField("Status.ObservedGeneration"),
			},
			{
				Name:        "updated_number_scheduled",
				Type:        proto.ColumnType_INT,
				Description: "The total number of nodes that are running updated daemon pod.",
				Transform:   transform.FromField("Status.UpdatedNumberScheduled"),
			},
			{
				Name:        "number_available",
				Type:        proto.ColumnType_INT,
				Description: "The number of nodes that should be running the daemon pod and have one or more of the daemon pod running and available (ready for at least spec.minReadySeconds).",
				Transform:   transform.FromField("Status.NumberAvailable"),
			},
			{
				Name:        "number_unavailable",
				Type:        proto.ColumnType_INT,
				Description: "The number of nodes that should be running the daemon pod and have none of the daemon pod running and available (ready for at least spec.minReadySeconds).",
				Transform:   transform.FromField("Status.NumberUnavailable"),
			},
			{
				Name:        "collision_count",
				Type:        proto.ColumnType_INT,
				Description: "Count of hash collisions for the DaemonSet. The DaemonSet controller uses this field as a collision avoidance mechanism when it needs to create the name for the newest ControllerRevision.",
				Transform:   transform.FromField("Status.CollisionCount"),
			},
			{
				Name:        "conditions",
				Type:        proto.ColumnType_JSON,
				Description: "Represents the latest available observations of a DaemonSet's current state.",
				Transform:   transform.FromField("Status.Conditions"),
			},

			//// Steampipe Standard Columns
			{
				Name:        "title",
				Type:        proto.ColumnType_STRING,
				Description: ColumnDescriptionTitle,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "tags",
				Type:        proto.ColumnType_JSON,
				Description: ColumnDescriptionTags,
				Transform:   transform.From(transformDaemonTags),
			},

			{Name: "manifest_file_path", Description: "Path to the file.", Type: proto.ColumnType_STRING},
		}),
	}
}

type KubernetesManifestDaemonSet struct {
	v1.DaemonSet
	ManifestFilePath string
}

//// LIST FUNCTION

func listKubernetesManifestDaemonsets(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// The path comes from a parent hydrate, defaulting to the config paths or
	// available by the optional key column
	path := h.Item.(filePath).Path

	// Load the file into a buffer
	content, err := os.ReadFile(path)
	if err != nil {
		plugin.Logger(ctx).Error("kubernetes_manifest_daemonset.listKubernetesManifestDaemonsets", "failed to read file", err, "path", path)
		return nil, err
	}
	decoder := scheme.Codecs.UniversalDeserializer()

	// Check for the start of the document
	for _, resource := range strings.Split(string(content), "---") {
		// skip empty documents, `Decode` will fail on them
		if len(resource) == 0 {
			continue
		}

		// Decode the file content
		obj, groupVersionKind, err := decoder.Decode([]byte(resource), nil, nil)
		if err != nil {
			plugin.Logger(ctx).Error("kubernetes_manifest_daemonset.listKubernetesManifestDaemonsets", "failed to decode the file", err, "path", path)
			return nil, err
		}

		// Return if the definition is not for the daemonSet resource
		if groupVersionKind.Kind == "DaemonSet" {
			daemonSet := obj.(*v1.DaemonSet)

			d.StreamListItem(ctx, KubernetesManifestDaemonSet{*daemonSet, path})

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

func transformDaemonTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(KubernetesManifestDaemonSet)
	return mergeTags(obj.Labels, obj.Annotations), nil
}
