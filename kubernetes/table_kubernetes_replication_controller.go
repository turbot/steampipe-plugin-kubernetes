package kubernetes

import (
	"context"
	"strings"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableKubernetesReplicaController(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_replication_controller",
		Description: "A Replication Controller makes sure that a pod or homogeneous set of pods are always up and available. If there are too many pods, it will kill some. If there are too few, the Replication Controller will start more.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "namespace"}),
			Hydrate:    getK8sReplicaController,
		},
		List: &plugin.ListConfig{
			Hydrate:    listK8sReplicaControllers,
			KeyColumns: getCommonOptionalKeyQuals(),
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
				Name:        "selector_query",
				Type:        proto.ColumnType_STRING,
				Description: "A query string representation of the selector.",
				Transform:   transform.FromField("Spec.Selector").Transform(labelSelectorToString),
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
				Name:        "ready_replicas",
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
				Description: "Reflects the generation of the most recently observed replication controller.",
				Transform:   transform.FromField("Status.ObservedGeneration"),
			},
			{
				Name:        "conditions",
				Type:        proto.ColumnType_JSON,
				Description: "Represents the latest available observations of a replication controller's current state.",
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
				Transform:   transform.From(transformReplicaControllerTags),
			},
		}),
	}
}

//// HYDRATE FUNCTIONS

func listK8sReplicaControllers(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listK8sReplicaControllers")

	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		return nil, err
	}

	input := metav1.ListOptions{
		Limit: 500,
	}

	// Limiting the results
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < input.Limit {
			if *limit < 1 {
				input.Limit = 1
			} else {
				input.Limit = *limit
			}
		}
	}

	commonFieldSelectorValue := getCommonOptionalKeyQualsValueForFieldSelector(d)

	if len(commonFieldSelectorValue) > 0 {
		input.FieldSelector = strings.Join(commonFieldSelectorValue, ",")
	}

	var response *v1.ReplicationControllerList
	pageLeft := true

	for pageLeft {
		response, err = clientset.CoreV1().ReplicationControllers("").List(ctx, input)
		if err != nil {
			return nil, err
		}

		if response.GetContinue() != "" {
			input.Continue = response.Continue
		} else {
			pageLeft = false
		}

		for _, replicaController := range response.Items {
			d.StreamListItem(ctx, replicaController)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

func getK8sReplicaController(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getK8sReplicaController")

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

	replicaController, err := clientset.CoreV1().ReplicationControllers(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil && !isNotFoundError(err) {
		return nil, err
	}

	return *replicaController, nil
}

//// TRANSFORM FUNCTIONS

func transformReplicaControllerTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(v1.ReplicationController)
	return mergeTags(obj.Labels, obj.Annotations), nil
}
