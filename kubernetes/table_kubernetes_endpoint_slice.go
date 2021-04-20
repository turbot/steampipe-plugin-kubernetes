package kubernetes

import (
	"context"

	"k8s.io/api/discovery/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableKubernetesEndpointSlice(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_endpoint_slice",
		Description: "EndpointSlice represents a subset of the endpoints that implement a service.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "namespace"}),
			Hydrate:    getK8sEnpointSlice,
		},
		List: &plugin.ListConfig{
			Hydrate: listK8sEnpointSlices,
		},
		Columns: k8sCommonColumns([]*plugin.Column{
			{
				Name:        "address_type",
				Type:        proto.ColumnType_STRING,
				Description: "Type of address carried by this EndpointSlice. All addresses in the slice are of the same type. Supported types are IPv4, IPv6, and FQDN.",
			},
			{
				Name:        "endpoints",
				Type:        proto.ColumnType_JSON,
				Description: "List of unique endpoints in this slice.",
			},
			{
				Name:        "ports",
				Type:        proto.ColumnType_JSON,
				Description: "List of network ports exposed by each endpoint in this slice. Each port must have a unique name. When ports is empty, it indicates that there are no defined ports.",
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
				Transform:   transform.From(transformEndpointSliceTags),
			},
		}),
	}
}

//// HYDRATE FUNCTIONS

func listK8sEnpointSlices(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listK8sEnpointSlices")

	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		return nil, err
	}

	nodes, err := clientset.DiscoveryV1beta1().EndpointSlices("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, endpointSlice := range nodes.Items {
		d.StreamListItem(ctx, endpointSlice)
	}

	return nil, nil
}

func getK8sEnpointSlice(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getK8sEnpointSlice")

	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		return nil, err
	}

	name := d.KeyColumnQuals["name"].GetStringValue()
	namespace := d.KeyColumnQuals["namespace"].GetStringValue()

	endpointSlice, err := clientset.DiscoveryV1beta1().EndpointSlices(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil && !isNotFoundError(err) {
		return nil, err
	}

	return *endpointSlice, nil
}

//// TRANSFORM FUNCTIONS

func transformEndpointSliceTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(v1beta1.EndpointSlice)
	return mergeTags(obj.Labels, obj.Annotations), nil
}
