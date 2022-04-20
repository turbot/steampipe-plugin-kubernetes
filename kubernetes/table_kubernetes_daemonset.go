package kubernetes

import (
	"context"

	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

func tableKubernetesDaemonset(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_daemonset",
		Description: "A DaemonSet ensures that all (or some) Nodes run a copy of a Pod.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "namespace"}),
			Hydrate:    getK8sDaemonSet,
		},
		List: &plugin.ListConfig{
			Hydrate: listK8sDaemonSets,
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
				Transform:   transform.From(transformDaemonSetTags),
			},
		}),
	}
}

//// HYDRATE FUNCTIONS

func listK8sDaemonSets(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listK8sDaemonSets")

	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		return nil, err
	}

	pods, err := clientset.AppsV1().DaemonSets("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, daemonSet := range pods.Items {
		d.StreamListItem(ctx, daemonSet)
	}

	return nil, nil
}

func getK8sDaemonSet(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getK8sDaemonSet")

	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		return nil, err
	}

	name := d.KeyColumnQuals["name"].GetStringValue()
	namespace := d.KeyColumnQuals["namespace"].GetStringValue()

	// return if namespace or name is empty
	if namespace == "" || name == "" {
		return nil, nil
	}

	daemonSet, err := clientset.AppsV1().DaemonSets(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil && !isNotFoundError(err) {
		return nil, err
	}

	return *daemonSet, nil
}

//// TRANSFORM FUNCTIONS

func transformDaemonSetTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(v1.DaemonSet)
	return mergeTags(obj.Labels, obj.Annotations), nil
}
