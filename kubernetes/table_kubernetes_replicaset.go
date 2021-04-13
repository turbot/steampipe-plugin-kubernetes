package kubernetes

import (
	"context"

	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableKubernetesReplicaSet(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_replicaset",
		Description: "Kubernetes ReplicaSet ensures that a specified number of pod replicas are running at any given time.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "namespace"}),
			Hydrate:    getK8sReplicaSet,
		},
		List: &plugin.ListConfig{
			Hydrate: listK8sReplicaSets,
		},
		Columns: k8sCommonColumns([]*plugin.Column{

			//// Spec Columns
			{
				Name:        "replicas",
				Type:        proto.ColumnType_INT,
				Description: "Replicas is the number of desired replicas. Defaults to 1.",
				Transform:   transform.FromField("Spec.Replicas"),
			},
			{
				Name:        "min_ready_seconds",
				Type:        proto.ColumnType_INT,
				Description: "Minimum number of seconds for which a newly created pod should be ready without any of its container crashing, for it to be considered available. Defaults to 0",
				Transform:   transform.FromField("Spec.MinReadySeconds"),
			},
			{
				Name:        "selector",
				Type:        proto.ColumnType_JSON,
				Description: "Selector is a label query over pods that should match the replica count. Label keys and values that must match in order to be controlled by this replica set.",
				Transform:   transform.FromField("Spec.Selector"),
			},
			{
				Name:        "template",
				Type:        proto.ColumnType_JSON,
				Description: "Template is the object that describes the pod that will be created if insufficient replicas are detected.",
				Transform:   transform.FromField("Spec.Template"),
			},

			//// Status Columns
			{
				Name:        "status_replicas",
				Type:        proto.ColumnType_INT,
				Description: "The most recently oberved number of replicas.",
				Transform:   transform.FromField("Status.Replicas"),
			},
			{
				Name:        "fully_labeled_replicas",
				Type:        proto.ColumnType_INT,
				Description: "The number of pods that have labels matching the labels of the pod template of the replicaset.",
				Transform:   transform.FromField("Status.FullyLabeledReplicas"),
			},
			{
				Name:        "ReadyReplicas",
				Type:        proto.ColumnType_INT,
				Description: "The number of ready replicas for this replica set.",
				Transform:   transform.FromField("Status.ReadyReplicas"),
			},
			{
				Name:        "available_replicas",
				Type:        proto.ColumnType_INT,
				Description: "The number of available replicas (ready for at least minReadySeconds) for this replica set.",
				Transform:   transform.FromField("Status.AvailableReplicas"),
			},
			{
				Name:        "observed_generation",
				Type:        proto.ColumnType_INT,
				Description: "ObservedGeneration reflects the generation of the most recently observed ReplicaSet.",
				Transform:   transform.FromField("Status.ObservedGeneration"),
			},
			{
				Name:        "conditions",
				Type:        proto.ColumnType_JSON,
				Description: "Represents the latest available observations of a replica set's current state.",
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
				Transform:   transform.From(transformReplicaSetTags),
			},
		}),
	}
}

//// HYDRATE FUNCTIONS

func listK8sReplicaSets(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listK8sReplicaSets")

	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		return nil, err
	}

	replicaSets, err := clientset.AppsV1().ReplicaSets("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, item := range replicaSets.Items {
		d.StreamListItem(ctx, item)
	}

	return nil, nil
}

func getK8sReplicaSet(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getK8sReplicaSet")

	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		return nil, err
	}

	name := d.KeyColumnQuals["name"].GetStringValue()
	namespace := d.KeyColumnQuals["namespace"].GetStringValue()

	rs, err := clientset.AppsV1().ReplicaSets(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil && !isNotFoundError(err) {
		return nil, err
	}

	return *rs, nil
}

//// TRANSFORM FUNCTIONS

func transformReplicaSetTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(v1.ReplicaSet)
	return mergeTags(obj.Labels, obj.Annotations), nil
}
