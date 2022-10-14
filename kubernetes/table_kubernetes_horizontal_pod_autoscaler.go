package kubernetes

import (
	"context"
	"strings"

	"k8s.io/api/autoscaling/v2beta2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
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

//// HYDRATE FUNCTIONS

func listK8sHPAs(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listK8sHPAs", "clientset_err", err)
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

	var response *v2beta2.HorizontalPodAutoscalerList
	pageLeft := true

	for pageLeft {
		response, err = clientset.AutoscalingV2beta2().HorizontalPodAutoscalers("").List(ctx, input)
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
			d.StreamListItem(ctx, hpa)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

func getK8sHPA(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("getK8sHPA", "clientset_err", err)
		return nil, err
	}

	name := d.KeyColumnQuals["name"].GetStringValue()
	namespace := d.KeyColumnQuals["namespace"].GetStringValue()

	// return if namespace or name is empty
	if namespace == "" || name == "" {
		return nil, nil
	}

	hpa, err := clientset.AutoscalingV2beta2().HorizontalPodAutoscalers(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil && !isNotFoundError(err) {
		plugin.Logger(ctx).Error("getK8sHPA", "api_err", err)
		return nil, err
	}

	return *hpa, nil
}

////// TRANSFORM FUNCTIONS

func transformHpaTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(v2beta2.HorizontalPodAutoscaler)
	return mergeTags(obj.Labels, obj.Annotations), nil
}
