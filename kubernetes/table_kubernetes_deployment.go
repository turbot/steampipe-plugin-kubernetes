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

func tableKubernetesDeployment(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_deployment",
		Description: "Kubernetes Deployment enables declarative updates for Pods and ReplicaSets.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "namespace"}),
			Hydrate:    getK8sDeployment,
		},
		List: &plugin.ListConfig{
			Hydrate:    listK8sDeployments,
			KeyColumns: getCommonOptionalKeyQuals(),
		},
		Columns: k8sCommonColumns([]*plugin.Column{

			//// Spec Columns
			{
				Name:        "replicas",
				Type:        proto.ColumnType_INT,
				Description: "Number of desired pods. Defaults to 1.",
				Transform:   transform.FromField("Spec.Replicas"),
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
				Description: "Label selector for pods. A label selector is a label query over a set of resources.",
				Transform:   transform.FromField("Spec.Selector"),
			},
			{
				Name:        "template",
				Type:        proto.ColumnType_JSON,
				Description: "Template describes the pods that will be created.",
				Transform:   transform.FromField("Spec.Template"),
			},
			{
				Name:        "strategy",
				Type:        proto.ColumnType_JSON,
				Description: "The deployment strategy to use to replace existing pods with new ones.",
				Transform:   transform.FromField("Spec.Strategy"),
			},
			{
				Name:        "min_ready_seconds",
				Type:        proto.ColumnType_INT,
				Description: "Minimum number of seconds for which a newly created pod should be ready without any of its container crashing, for it to be considered available. Defaults to 0.",
				Transform:   transform.FromField("Spec.MinReadySeconds"),
			},
			{
				Name:        "revision_history_limit",
				Type:        proto.ColumnType_INT,
				Description: "The number of old ReplicaSets to retain to allow rollback.",
				Transform:   transform.FromField("Spec.RevisionHistoryLimit"),
			},
			{
				Name:        "paused",
				Type:        proto.ColumnType_BOOL,
				Description: "Indicates that the deployment is paused.",
				Transform:   transform.FromField("Spec.Paused"),
			},
			{
				Name:        "progress_deadline_seconds",
				Type:        proto.ColumnType_INT,
				Description: "The maximum time in seconds for a deployment to make progress before it is considered to be failed.",
				Transform:   transform.FromField("Spec.ProgressDeadlineSeconds"),
			},

			//// Status Columns
			{
				Name:        "observed_generation",
				Type:        proto.ColumnType_INT,
				Description: "The generation observed by the deployment controller.",
				Transform:   transform.FromField("Status.ObservedGeneration"),
			},
			{
				Name:        "status_replicas",
				Type:        proto.ColumnType_INT,
				Description: "Total number of non-terminated pods targeted by this deployment (their labels match the selector).",
				Transform:   transform.FromField("Status.Replicas"),
			},
			{
				Name:        "updated_replicas",
				Type:        proto.ColumnType_INT,
				Description: "Total number of non-terminated pods targeted by this deployment that have the desired template spec.",
				Transform:   transform.FromField("Status.UpdatedReplicas"),
			},
			{
				Name:        "ready_replicas",
				Type:        proto.ColumnType_INT,
				Description: "Total number of ready pods targeted by this deployment.",
				Transform:   transform.FromField("Status.ReadyReplicas"),
			},
			{
				Name:        "available_replicas",
				Type:        proto.ColumnType_INT,
				Description: "Total number of available pods (ready for at least minReadySeconds) targeted by this deployment.",
				Transform:   transform.FromField("Status.AvailableReplicas"),
			},
			{
				Name:        "unavailable_replicas",
				Type:        proto.ColumnType_INT,
				Description: "Total number of unavailable pods targeted by this deployment.",
				Transform:   transform.FromField("Status.UnavailableReplicas"),
			},
			{
				Name:        "conditions",
				Type:        proto.ColumnType_JSON,
				Description: "Represents the latest available observations of a deployment's current state.",
				Transform:   transform.FromField("Status.Conditions"),
			},
			{
				Name:        "collision_count",
				Type:        proto.ColumnType_INT,
				Description: "Count of hash collisions for the Deployment. The Deployment controller uses this field as a collision avoidance mechanism when it needs to create the name for the newest ReplicaSet.",
				Transform:   transform.FromField("Status.CollisionCount"),
			},
			{
				Name:        "context_name",
				Type:        proto.ColumnType_STRING,
				Description: "Kubectl config context name.",
				Hydrate:     getDeploymentResourceAdditionalData,
			},
			{
				Name:        "source_type",
				Type:        proto.ColumnType_STRING,
				Description: "The source of the resource. Possible values are: deployed and manifest. If the resource is fetched from the spec file the value will be manifest.",
				Hydrate:     getDeploymentResourceAdditionalData,
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
				Transform:   transform.From(transformDeploymentTags),
			},
		}),
	}
}

type Deployment struct {
	v1.Deployment
	Path      string
	StartLine int
	EndLine   int
}

//// HYDRATE FUNCTIONS

func listK8sDeployments(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listK8sDeployments")

	// Get the client for querying the K8s APIs for the provided context.
	// If the connection is configured for the manifest files, the client will return nil.
	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		return nil, err
	}

	// Check for manifest files
	parsedContents, err := fetchResourceFromManifestFileByKind(ctx, d, "Deployment")
	if err != nil {
		return nil, err
	}

	for _, content := range parsedContents {
		deployment := content.Data.(*v1.Deployment)

		d.StreamListItem(ctx, Deployment{*deployment, content.Path, content.StartLine, content.EndLine})

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	// Check for deployed resources
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

	var response *v1.DeploymentList
	pageLeft := true

	for pageLeft {

		response, err = clientset.AppsV1().Deployments("").List(ctx, input)
		if err != nil {
			return nil, err
		}

		if response.GetContinue() != "" {
			input.Continue = response.Continue
		} else {
			pageLeft = false
		}

		for _, item := range response.Items {
			d.StreamListItem(ctx, Deployment{item, "", 0, 0})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

func getK8sDeployment(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getK8sDeployment")

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

	// Get the manifest resource
	parsedContents, err := fetchResourceFromManifestFileByKind(ctx, d, "Deployment")
	if err != nil {
		return nil, err
	}

	for _, content := range parsedContents {
		deployment := content.Data.(*v1.Deployment)

		if deployment.Name == name && deployment.Namespace == namespace {
			return Deployment{*deployment, content.Path, content.StartLine, content.EndLine}, nil
		}
	}

	// Get the deployed resource
	if clientset == nil {
		return nil, nil
	}

	deployment, err := clientset.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil && !isNotFoundError(err) {
		return nil, err
	}

	return Deployment{*deployment, "", 0, 0}, nil
}

func getDeploymentResourceAdditionalData(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	obj := h.Item.(Deployment)

	data := map[string]interface{}{
		"SourceType": "deployed",
	}

	// Set the source_type as manifest, if path is not empty
	// also, set the context_name as nil
	if obj.Path != "" {
		data["SourceType"] = "manifest"
		return data, nil
	}

	// Else, set the current context as context_name
	currentContext, err := getKubectlContext(ctx, d, nil)
	if err != nil {
		return data, nil
	}
	data["ContextName"] = currentContext.(string)

	return data, nil
}

//// TRANSFORM FUNCTIONS

func transformDeploymentTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(Deployment)
	return mergeTags(obj.Labels, obj.Annotations), nil
}
