package kubernetes

import (
	"context"
	"strings"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
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
			Hydrate:    listK8sResourceQuotas,
			KeyColumns: getCommonOptionalKeyQuals(),
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

	var response *v1.ResourceQuotaList
	pageLeft := true

	for pageLeft {
		response, err = clientset.CoreV1().ResourceQuotas("").List(ctx, input)
		if err != nil {
			return nil, err
		}

		if response.GetContinue() != "" {
			input.Continue = response.Continue
		} else {
			pageLeft = false
		}

		for _, resourceQuota := range response.Items {
			d.StreamListItem(ctx, resourceQuota)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
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
