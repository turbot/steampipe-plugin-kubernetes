package kubernetes

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

func tableKubernetesEndpoints(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_endpoint",
		Description: "Set of addresses and ports that comprise a service. More info: https://kubernetes.io/docs/concepts/services-networking/service/#services-without-selectors.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "namespace"}),
			Hydrate:    getK8sEndpoint,
		},
		List: &plugin.ListConfig{
			Hydrate: listK8sEnpoints,
		},
		Columns: k8sCommonColumns([]*plugin.Column{
			{
				Name:        "subsets",
				Type:        proto.ColumnType_JSON,
				Description: "List of addresses and ports that comprise a service.",
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
				Transform:   transform.From(transformEndpointTags),
			},
		}),
	}
}

//// HYDRATE FUNCTIONS

func listK8sEnpoints(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listK8sEnpoints")

	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		return nil, err
	}

	nodes, err := clientset.CoreV1().Endpoints("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, endpoint := range nodes.Items {
		d.StreamListItem(ctx, endpoint)
	}

	return nil, nil
}

func getK8sEndpoint(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getK8sEndpoint")

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

	endpoint, err := clientset.CoreV1().Endpoints(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil && !isNotFoundError(err) {
		return nil, err
	}

	return *endpoint, nil
}

//// TRANSFORM FUNCTIONS

func transformEndpointTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(v1.Endpoints)
	return mergeTags(obj.Labels, obj.Annotations), nil
}
