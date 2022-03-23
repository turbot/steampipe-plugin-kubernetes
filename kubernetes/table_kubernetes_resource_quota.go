package kubernetes

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/turbot/steampipe-plugin-sdk/v2/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin/transform"
)

func tableKubernetesResourceQuota(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_resource_quota",
		Description: "Kubernetes Resource Quota",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "namespace"}),
			Hydrate:    getK8sResourceQuota,
		},
		List: &plugin.ListConfig{
			Hydrate: listK8sResourceQuotas,
		},
		Columns: k8sCommonColumns([]*plugin.Column{

			//// ResourceQuotaSpec Columns
			{
				Name:        "spec_hard",
				Type:        proto.ColumnType_JSON,
				Description: "Spec hard is the set of desired hard limits for each named resource.",
				Transform:   transform.FromField("Spec.Hard"),
			},
			{
				Name:        "spec_scopes",
				Type:        proto.ColumnType_JSON,
				Description: "A collection of filters that must match each object tracked by a quota.",
				Transform:   transform.FromField("Spec.Scopes"),
			},
			{
				Name:        "spec_scope_selector",
				Type:        proto.ColumnType_JSON,
				Description: "A collection of filters like scopes that must match each object tracked by a quota but expressed using ScopeSelectorOperator in combination with possible values.",
				Transform:   transform.FromField("Spec.ScopeSelector"),
			},

			//// ResourceQuotaStatus Columns
			{
				Name:        "status_hard",
				Type:        proto.ColumnType_JSON,
				Description: "Status hard is the set of enforced hard limits for each named resource.",
				Transform:   transform.FromField("Status.Hard"),
			},
			{
				Name:        "status_used",
				Type:        proto.ColumnType_JSON,
				Description: "Indicates current observed total usage of the resource in the namespace.",
				Transform:   transform.FromField("Status.Used"),
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
				Transform:   transform.From(transformResourceQuotaTags),
			},
		}),
	}
}

//// HYDRATE FUNCTIONS

func listK8sResourceQuotas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listK8sResourceQuotas")

	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		return nil, err
	}

	resourceQuotas, err := clientset.CoreV1().ResourceQuotas("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, resourceQuota := range resourceQuotas.Items {
		d.StreamListItem(ctx, resourceQuota)
	}

	return nil, nil
}

func getK8sResourceQuota(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getK8sResourceQuota")

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

	resourceQuota, err := clientset.CoreV1().ResourceQuotas(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil && !isNotFoundError(err) {
		return nil, err
	}

	return *resourceQuota, nil
}

//// TRANSFORM FUNCTIONS

func transformResourceQuotaTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(v1.ResourceQuota)
	return mergeTags(obj.Labels, obj.Annotations), nil
}
