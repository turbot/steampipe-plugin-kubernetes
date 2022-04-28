package kubernetes

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

func tableKubernetesLimitRange(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_limit_range",
		Description: "Kubernetes Limit Range",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "namespace"}),
			Hydrate:    getK8sLimitRange,
		},
		List: &plugin.ListConfig{
			Hydrate: listK8sLimitRanges,
		},
		Columns: k8sCommonColumns([]*plugin.Column{

			//// LimitRangeSpec Columns
			{
				Name:        "spec_limits",
				Type:        proto.ColumnType_JSON,
				Description: "List of limit range item objects that are enforced.",
				Transform:   transform.FromField("Spec.Limits"),
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
				Transform:   transform.From(transformLimitRangeTags),
			},
		}),
	}
}

//// HYDRATE FUNCTIONS

func listK8sLimitRanges(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listK8sLimitRanges")

	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		return nil, err
	}

	limitRangeList, err := clientset.CoreV1().LimitRanges("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, limitRange := range limitRangeList.Items {
		d.StreamListItem(ctx, limitRange)
	}

	return nil, nil
}

func getK8sLimitRange(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getK8sLimitRange")

	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		return nil, err
	}

	name := d.KeyColumnQuals["name"].GetStringValue()
	namespace := d.KeyColumnQuals["namespace"].GetStringValue()

	// return nil if namespace or name is empty
	if namespace == "" || name == "" {
		return nil, nil
	}

	limitRange, err := clientset.CoreV1().LimitRanges(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil && !isNotFoundError(err) {
		return nil, err
	}

	return *limitRange, nil
}

//// TRANSFORM FUNCTIONS

func transformLimitRangeTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(v1.LimitRange)
	return mergeTags(obj.Labels, obj.Annotations), nil
}
