package kubernetes

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableKubernetesNamespace(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_namespace",
		Description: "Kubernetes Namespace provides a scope for Names.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getK8sNamespace,
		},
		List: &plugin.ListConfig{
			Hydrate: listK8sNamespaces,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "phase", Require: plugin.Optional},
			},
		},
		Columns: k8sCommonGlobalColumns([]*plugin.Column{

			//// NamespaceSpec Columns
			{
				Name:        "spec_finalizers",
				Type:        proto.ColumnType_JSON,
				Description: "Finalizers is an opaque list of values that must be empty to permanently remove object from storage.",
				Transform:   transform.FromField("Spec.Finalizers"),
			},

			//// NamespaceStatus Columns
			{
				Name:        "phase",
				Type:        proto.ColumnType_STRING,
				Description: "The current lifecycle phase of the namespace.",
				Transform:   transform.FromField("Status.Phase"),
			},
			{
				Name:        "conditions",
				Type:        proto.ColumnType_JSON,
				Description: "The latest available observations of namespace's current state.",
				Transform:   transform.FromField("Status.NamespaceCondition"),
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
				Transform:   transform.From(transformNamespaceTags),
			},
		}),
	}
}

//// HYDRATE FUNCTIONS

func listK8sNamespaces(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listK8sNamespaces")

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

	if d.EqualsQualString("phase") != "" {
		input.FieldSelector = fmt.Sprintf("status.phase=%v", d.EqualsQualString("phase"))
	}

	var response *v1.NamespaceList
	pageLeft := true

	for pageLeft {
		response, err = clientset.CoreV1().Namespaces().List(ctx, input)
		if err != nil {
			return nil, err
		}

		if response.GetContinue() != "" {
			input.Continue = response.Continue
		} else {
			pageLeft = false
		}

		for _, pod := range response.Items {
			d.StreamListItem(ctx, pod)
		}
	}

	return nil, nil
}

func getK8sNamespace(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getK8sNamespace")

	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		return nil, err
	}

	name := d.EqualsQuals["name"].GetStringValue()

	// return if name is empty
	if name == "" {
		return nil, nil
	}

	namespace, err := clientset.CoreV1().Namespaces().Get(ctx, name, metav1.GetOptions{})
	if err != nil && !isNotFoundError(err) {
		return nil, err
	}

	return *namespace, nil
}

//// TRANSFORM FUNCTIONS

func transformNamespaceTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(v1.Namespace)
	return mergeTags(obj.Labels, obj.Annotations), nil
}
