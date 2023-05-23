package kubernetes

import (
	"context"
	"strings"

	v1 "k8s.io/api/autoscaling/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableKubernetesHorizontalPodAutoscaler(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_horizontal_pod_autoscaler",
		Description: "Kubernetes HorizontalPodAutoscaler is the configuration for a horizontal pod autoscaler, which automatically manages the replica count of any resource implementing the scale subresource based on the metrics specified.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "namespace"}),
			Hydrate:    getK8sHPA,
		},
		List: &plugin.ListConfig{
			Hydrate:    listK8sHPAs,
			KeyColumns: getCommonOptionalKeyQuals(),
		},
		Columns: k8sCommonColumns([]*plugin.Column{
			//// HpaSpec Columns
			{
				Name:        "scale_target_ref",
				Type:        proto.ColumnType_JSON,
				Description: "ScaleTargetRef points to the target resource to scale, and is used to the pods for which metrics should be collected, as well as to actually change the replica count.",
				Transform:   transform.FromField("Spec.ScaleTargetRef"),
			},
			{
				Name:        "min_replicas",
				Type:        proto.ColumnType_INT,
				Description: "MinReplicas is the lower limit for the number of replicas to which the autoscaler can scale down. It defaults to 1 pod. MinReplicas is allowed to be 0 if the alpha feature gate HPAScaleToZero is enabled and at least one Object or External metric is configured.",
				Transform:   transform.FromField("Spec.MinReplicas"),
			},
			{
				Name:        "max_replicas",
				Type:        proto.ColumnType_INT,
				Description: "The Upper limit for the number of pods that can be set by the autoscaler. It cannot be smaller than MinReplicas.",
				Transform:   transform.FromField("Spec.MaxReplicas"),
			},
			{
				Name:        "metrics",
				Type:        proto.ColumnType_JSON,
				Description: "Metrics contains the specifications for which to use to calculate the desired replica count (the maximum replica count across all metrics will be used). The desired replica count is calculated multiplying the ratio between the target value and the current value by the current number of pods.",
				Transform:   transform.FromField("Spec.Metrics"),
			},
			{
				Name:        "scale_up_behavior",
				Type:        proto.ColumnType_JSON,
				Description: "Behavior configures the scaling behavior of the target in both Up and Down directions (scaleUp and scaleDown fields respectively). If not set, the default value is the higher of: * increase no more than 4 pods per 60 seconds * double the number of pods per 60 seconds.",
				Transform:   transform.FromField("Spec.Behavior.ScaleUp"),
			},
			{
				Name:        "scale_down_behavior",
				Type:        proto.ColumnType_JSON,
				Description: "Behavior configures the scaling behavior of the target in both Up and Down directions (scaleUp and scaleDown fields respectively). If not set, the default value is to allow to scale down to minReplicas pods, with a 300 second stabilization window (i.e., the highest recommendation for the last 300sec is used).",
				Transform:   transform.FromField("Spec.Behavior.ScaleDown"),
			},

			//// HpaStatus Columns
			{
				Name:        "observed_generation",
				Type:        proto.ColumnType_INT,
				Description: "The most recent generation observed by this autoscaler.",
				Transform:   transform.FromField("Status.ObservedGeneration"),
			},
			{
				Name:        "last_scale_time",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The last time the HorizontalPodAutoscaler scaled the number of pods used by the autoscaler to control how often the number of pods is changed.",
				Transform:   transform.FromField("Status.LastScaleTime").Transform(v1TimeToRFC3339),
			},
			{
				Name:        "current_replicas",
				Type:        proto.ColumnType_INT,
				Description: "The current number of replicas of pods managed by this autoscaler.",
				Transform:   transform.FromField("Status.CurrentReplicas"),
			},
			{
				Name:        "desired_replicas",
				Type:        proto.ColumnType_INT,
				Description: "The desired number of replicas of pods managed by this autoscaler.",
				Transform:   transform.FromField("Status.DesiredReplicas"),
			},
			{
				Name:        "current_metrics",
				Type:        proto.ColumnType_JSON,
				Description: "CurrentMetrics is the last read state of the metrics used by this autoscaler.",
				Transform:   transform.FromField("Status.CurrentMetrics"),
			},
			{
				Name:        "conditions",
				Type:        proto.ColumnType_JSON,
				Description: "Conditions is the set of conditions required for this autoscaler to scale its target and indicates whether or not those conditions are met.",
				Transform:   transform.FromField("Status.Conditions"),
			},
			{
				Name:        "context_name",
				Type:        proto.ColumnType_STRING,
				Description: "Kubectl config context name.",
				Hydrate:     getHorizontalPodAutoscalarResourceAdditionalData,
			},
			{
				Name:        "source_type",
				Type:        proto.ColumnType_STRING,
				Description: "The source of the resource. Possible values are: deployed and manifest. If the resource is fetched from the spec file the value will be manifest.",
				Hydrate:     getHorizontalPodAutoscalarResourceAdditionalData,
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
				Transform:   transform.From(transformHpaTags),
			},
		}),
	}
}

type HorizontalPodAutoscaler struct {
	v1.HorizontalPodAutoscaler
	Path      string
	StartLine int
	EndLine   int
}

//// HYDRATE FUNCTIONS

func listK8sHPAs(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Get the client for querying the K8s APIs for the provided context.
	// If the connection is configured for the manifest files, the client will return nil.
	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listK8sHPAs", "clientset_err", err)
		return nil, err
	}

	// Check for manifest files
	parsedContents, err := fetchResourceFromManifestFileByKind(ctx, d, "HorizontalPodAutoscaler")
	if err != nil {
		return nil, err
	}

	for _, content := range parsedContents {
		hpa := content.Data.(*v1.HorizontalPodAutoscaler)

		d.StreamListItem(ctx, HorizontalPodAutoscaler{*hpa, content.Path, content.StartLine, content.EndLine})

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

	var response *v1.HorizontalPodAutoscalerList
	pageLeft := true

	for pageLeft {
		response, err = clientset.AutoscalingV1().HorizontalPodAutoscalers("").List(ctx, input)
		if err != nil {
			plugin.Logger(ctx).Error("listK8sHPAs", "api_err", err)
			return nil, err
		}

		if response.GetContinue() != "" {
			input.Continue = response.Continue
		} else {
			pageLeft = false
		}

		for _, hpa := range response.Items {
			d.StreamListItem(ctx, HorizontalPodAutoscaler{hpa, "", 0, 0})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

func getK8sHPA(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Get the client for querying the K8s APIs for the provided context.
	// If the connection is configured for the manifest files, the client will return nil.
	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("getK8sHPA", "clientset_err", err)
		return nil, err
	}

	name := d.EqualsQuals["name"].GetStringValue()
	namespace := d.EqualsQuals["namespace"].GetStringValue()

	// return if namespace or name is empty
	if namespace == "" || name == "" {
		return nil, nil
	}

	// Get the manifest resource
	parsedContents, err := fetchResourceFromManifestFileByKind(ctx, d, "HorizontalPodAutoscaler")
	if err != nil {
		return nil, err
	}

	for _, content := range parsedContents {
		hpa := content.Data.(*v1.HorizontalPodAutoscaler)

		if hpa.Name == name && hpa.Namespace == namespace {
			return HorizontalPodAutoscaler{*hpa, content.Path, content.StartLine, content.EndLine}, nil
		}
	}

	// Get the deployed resource
	if clientset == nil {
		return nil, nil
	}

	hpa, err := clientset.AutoscalingV1().HorizontalPodAutoscalers(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil && !isNotFoundError(err) {
		plugin.Logger(ctx).Error("getK8sHPA", "api_err", err)
		return nil, err
	}

	return HorizontalPodAutoscaler{*hpa, "", 0, 0}, nil
}

func getHorizontalPodAutoscalarResourceAdditionalData(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	obj := h.Item.(HorizontalPodAutoscaler)

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

////// TRANSFORM FUNCTIONS

func transformHpaTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(HorizontalPodAutoscaler)
	return mergeTags(obj.Labels, obj.Annotations), nil
}
