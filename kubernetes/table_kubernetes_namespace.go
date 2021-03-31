package kubernetes

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
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
		},
		Columns: k8sCommonColumns([]*plugin.Column{

			//// Spec Columns
			// To do - add Spec Columns...

			//// Status Columns
			// To do - add Status Columns...

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

	clientset, err := GetNewClientset(ctx, d.ConnectionManager)
	if err != nil {
		return nil, err
	}

	namespaces, err := clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, pod := range namespaces.Items {
		d.StreamListItem(ctx, pod)
	}

	return nil, nil
}

func getK8sNamespace(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getK8sNamespace")

	clientset, err := GetNewClientset(ctx, d.ConnectionManager)
	if err != nil {
		return nil, err
	}

	name := d.KeyColumnQuals["name"].GetStringValue()

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
