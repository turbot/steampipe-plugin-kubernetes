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

func tableKubernetesDaemonset(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_daemonset",
		Description: "A DaemonSet ensures that all (or some) Nodes run a copy of a Pod.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "namespace"}),
			Hydrate:    getK8sDaemonSet,
		},
		List: &plugin.ListConfig{
			Hydrate:    listK8sDaemonSets,
			KeyColumns: getCommonOptionalKeyQuals(),
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
			{
				Name:        "source_type",
				Type:        proto.ColumnType_STRING,
				Description: "The source of the resource. Possible values are: deployed and manifest. If the resource is fetched from the spec file the value will be manifest.",
				Transform:   transform.From(daemonSetResourceSourceType),
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

type DaemonSet struct {
	v1.DaemonSet
	Path      string
	StartLine int
	EndLine   int
}

//// HYDRATE FUNCTIONS

func listK8sDaemonSets(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listK8sDaemonSets")

	// Get the client for querying the K8s APIs for the provided context.
	// If the connection is configured for the manifest files, the client will return nil.
	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		return nil, err
	}

	//
	// Check for manifest files
	//
	parsedContents, err := fetchResourceFromManifestFileByKind(ctx, d, "DaemonSet")
	if err != nil {
		return nil, err
	}

	for _, content := range parsedContents {
		daemonSet := content.Data.(*v1.DaemonSet)

		d.StreamListItem(ctx, DaemonSet{*daemonSet, content.Path, content.StartLine, content.EndLine})

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

	var response *v1.DaemonSetList
	pageLeft := true

	for pageLeft {
		response, err = clientset.AppsV1().DaemonSets("").List(ctx, input)
		if err != nil {
			return nil, err
		}

		if response.GetContinue() != "" {
			input.Continue = response.Continue
		} else {
			pageLeft = false
		}

		for _, daemonSet := range response.Items {
			d.StreamListItem(ctx, DaemonSet{daemonSet, "", 0, 0})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

func getK8sDaemonSet(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getK8sDaemonSet")

	// Get the client for querying the K8s APIs for the provided context.
	// If the connection is configured for the manifest files, the client will return nil.
	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		return nil, err
	}

	name := d.EqualsQuals["name"].GetStringValue()
	namespace := d.EqualsQuals["namespace"].GetStringValue()

	// return if namespace or name is empty
	if namespace == "" || name == "" {
		return nil, nil
	}

	//
	// Get the manifest resource
	//
	parsedContents, err := fetchResourceFromManifestFileByKind(ctx, d, "DaemonSet")
	if err != nil {
		return nil, err
	}

	for _, content := range parsedContents {
		daemonSet := content.Data.(*v1.DaemonSet)

		if daemonSet.Name == name && daemonSet.Namespace == namespace {
			return DaemonSet{*daemonSet, content.Path, content.StartLine, content.EndLine}, nil
		}
	}

	//
	// Get the deployed resource
	//
	if clientset == nil {
		return nil, nil
	}

	daemonSet, err := clientset.AppsV1().DaemonSets(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil && !isNotFoundError(err) {
		return nil, err
	}

	return DaemonSet{*daemonSet, "", 0, 0}, nil
}

//// TRANSFORM FUNCTIONS

func transformDaemonSetTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(DaemonSet)
	return mergeTags(obj.Labels, obj.Annotations), nil
}

func daemonSetResourceSourceType(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(DaemonSet)

	if obj.Path != "" {
		return "manifest", nil
	}
	return "deployed", nil
}
