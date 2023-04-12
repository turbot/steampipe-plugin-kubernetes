package kubernetes

import (
	"context"
	"strings"

	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableKubernetesStatefulSet(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_stateful_set",
		Description: "A statefulSet is the workload API object used to manage stateful applications.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "namespace"}),
			Hydrate:    getK8sStatefulSet,
		},
		List: &plugin.ListConfig{
			Hydrate:    listK8sStatefulSets,
			KeyColumns: getCommonOptionalKeyQuals(),
		},
		// StatefulSet, is namespaced resource.
		Columns: k8sCommonColumns([]*plugin.Column{
			{
				Name:        "service_name",
				Type:        proto.ColumnType_STRING,
				Description: "The ame of the service that governs this StatefulSet.",
				Transform:   transform.FromField("Spec.ServiceName"),
			},
			{
				Name:        "replicas",
				Type:        proto.ColumnType_INT,
				Description: "The desired number of replicas of the given Template.",
				Transform:   transform.FromField("Spec.Replicas"),
			},
			{
				Name:        "collision_count",
				Type:        proto.ColumnType_INT,
				Description: "The count of hash collisions for the StatefulSet.",
				Transform:   transform.FromField("Status.CollisionCount"),
			},
			{
				Name:        "available_replicas",
				Type:        proto.ColumnType_INT,
				Description: "Total number of available pods (ready for at least minReadySeconds) targeted by this statefulset.",
				Transform:   transform.FromField("Status.AvailableReplicas"),
			},
			{
				Name:        "current_replicas",
				Type:        proto.ColumnType_INT,
				Description: "The number of Pods created by the StatefulSet controller from the StatefulSet version indicated by currentRevision.",
				Transform:   transform.FromField("Status.CurrentReplicas"),
			},
			{
				Name:        "current_revision",
				Type:        proto.ColumnType_STRING,
				Description: "Indicates the version of the StatefulSet used to generate Pods in the sequence [0,currentReplicas).",
				Transform:   transform.FromField("Status.CurrentRevision"),
			},
			{
				Name:        "observed_generation",
				Type:        proto.ColumnType_INT,
				Description: "The most recent generation observed for this StatefulSet.",
				Transform:   transform.FromField("Status.ObservedGeneration"),
			},
			{
				Name:        "pod_management_policy",
				Type:        proto.ColumnType_STRING,
				Description: "Policy that controls how pods are created during initial scale up, when replacing pods on nodes, or when scaling down.",
				Transform:   transform.FromField("Spec.PodManagementPolicy").Transform(transform.ToString),
			},
			{
				Name:        "ready_replicas",
				Type:        proto.ColumnType_INT,
				Description: "The number of Pods created by the StatefulSet controller that have a Ready Condition.",
				Transform:   transform.FromField("Status.ReadyReplicas"),
			},
			{
				Name:        "revision_history_limit",
				Type:        proto.ColumnType_INT,
				Description: "The maximum number of revisions that will be maintained in the StatefulSet's revision history.",
				Transform:   transform.FromField("Spec.RevisionHistoryLimit"),
			},
			{
				Name:        "updated_replicas",
				Type:        proto.ColumnType_INT,
				Description: "The number of Pods created by the StatefulSet controller from the StatefulSet version indicated by updateRevision.",
				Transform:   transform.FromField("Status.UpdatedReplicas"),
			},
			{
				Name:        "update_revision",
				Type:        proto.ColumnType_STRING,
				Description: "Indicates the version of the StatefulSet used to generate Pods in the sequence [replicas-updatedReplicas,replicas).",
				Transform:   transform.FromField("Status.UpdateRevision"),
			},
			{
				Name:        "conditions",
				Type:        proto.ColumnType_JSON,
				Description: "Represents the latest available observations of a stateful set's current state.",
				Transform:   transform.FromField("Status.Conditions"),
			},
			{
				Name:        "template",
				Type:        proto.ColumnType_JSON,
				Description: "Template is the object that describes the pod that will be created if insufficient replicas are detected.",
				Transform:   transform.FromField("Spec.Template"),
			},
			{
				Name:        "update_strategy",
				Type:        proto.ColumnType_JSON,
				Description: "Indicates the StatefulSetUpdateStrategy that will be employed to update Pods in the StatefulSet when a revision is made to Template.",
				Transform:   transform.FromField("Spec.UpdateStrategy"),
			},
			{
				Name:        "source",
				Type:        proto.ColumnType_STRING,
				Description: "The source of the resource. Possible values are: deployed and manifest. If the resource is fetched from the spec file the value will be manifest.",
				Transform:   transform.From(statefulSetResourceSource),
			},

			// Steampipe Standard Columns
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
				Transform:   transform.From(transformStatefulSetTags),
			},
		}),
	}
}

type StatefulSet struct {
	v1.StatefulSet
	Path      string
	StartLine int
	EndLine   int
}

//// HYDRATE FUNCTIONS

func listK8sStatefulSets(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listK8sStatefulSets")

	// Get the client for querying the K8s APIs for the provided context.
	// If the connection is configured for the manifest files, the client will return nil.
	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		return nil, err
	}

	//
	// Check for manifest files
	//
	parsedContents, err := fetchResourceFromManifestFileByKind(ctx, d, "StatefulSet")
	if err != nil {
		return nil, err
	}

	for _, content := range parsedContents {
		statefulSet := content.Data.(*v1.StatefulSet)

		d.StreamListItem(ctx, StatefulSet{*statefulSet, content.Path, content.StartLine, content.EndLine})

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	//
	// Check for deployed resources
	//
	if clientset == nil {
		return nil, nil
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

	var response *v1.StatefulSetList
	pageLeft := true

	for pageLeft {
		response, err = clientset.AppsV1().StatefulSets("").List(ctx, input)
		if err != nil {
			return nil, err
		}

		if response.GetContinue() != "" {
			input.Continue = response.Continue
		} else {
			pageLeft = false
		}

		for _, statefulSet := range response.Items {
			d.StreamListItem(ctx, StatefulSet{statefulSet, "", 0, 0})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

func getK8sStatefulSet(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getK8sStatefulSet")

	// Get the client for querying the K8s APIs for the provided context.
	// If the connection is configured for the manifest files, the client will return nil.
	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		return nil, err
	}

	name := d.EqualsQuals["name"].GetStringValue()
	namespace := d.EqualsQuals["namespace"].GetStringValue()

	// handle empty name and namespace value
	if name == "" || namespace == "" {
		return nil, nil
	}

	//
	// Get the manifest resource
	//
	parsedContents, err := fetchResourceFromManifestFileByKind(ctx, d, "StatefulSet")
	if err != nil {
		return nil, err
	}

	for _, content := range parsedContents {
		statefulSet := content.Data.(*v1.StatefulSet)

		if statefulSet.Name == name && statefulSet.Namespace == namespace {
			return StatefulSet{*statefulSet, content.Path, content.StartLine, content.EndLine}, nil
		}
	}

	//
	// Get the deployed resource
	//
	if clientset == nil {
		return nil, nil
	}

	statefulSet, err := clientset.AppsV1().StatefulSets(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil && !isNotFoundError(err) {
		logger.Debug("getK8sStatefulSet", "Error", err)
		return nil, err
	}

	return StatefulSet{*statefulSet, "", 0, 0}, nil
}

//// TRANSFORM FUNCTIONS

func transformStatefulSetTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(StatefulSet)
	return mergeTags(obj.Labels, obj.Annotations), nil
}

func statefulSetResourceSource(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(StatefulSet)

	if obj.Path != "" {
		return "manifest", nil
	}
	return "deployed", nil
}
