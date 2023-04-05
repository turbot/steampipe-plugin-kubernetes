package kubernetes

import (
	"context"
	"strings"

	"k8s.io/api/discovery/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
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
			Hydrate:    listK8sEnpointSlices,
			KeyColumns: getCommonOptionalKeyQuals(),
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
			{
				Name:        "manifest_file_path",
				Type:        proto.ColumnType_STRING,
				Description: "The path to the manifest file.",
				Transform:   transform.FromField("ManifestFilePath").Transform(transform.NullIfZeroValue),
			},
		}),
	}
}

type EndpointSlice struct {
	v1beta1.EndpointSlice
	ManifestFilePath string
}

//// HYDRATE FUNCTIONS

func listK8sEnpointSlices(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listK8sEnpointSlices")

	// Get the client for querying the K8s APIs for the provided context.
	// If the connection is configured for the manifest files, the client will return nil.
	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		return nil, err
	}

	//
	// Check for manifest files
	//
	parsedContents, err := fetchResourceFromManifestFileByKind(ctx, d, "EndpointSlice")
	if err != nil {
		return nil, err
	}

	for _, content := range parsedContents {
		endpointSlice := content.Data.(*v1beta1.EndpointSlice)

		d.StreamListItem(ctx, EndpointSlice{*endpointSlice, content.Path})

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	//
	// Check for deployed resources
	//
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

	var response *v1beta1.EndpointSliceList
	pageLeft := true

	for pageLeft {
		response, err = clientset.DiscoveryV1beta1().EndpointSlices("").List(ctx, input)
		if err != nil {
			return nil, err
		}

		if response.GetContinue() != "" {
			input.Continue = response.Continue
		} else {
			pageLeft = false
		}

		for _, endpointSlice := range response.Items {
			d.StreamListItem(ctx, EndpointSlice{endpointSlice, ""})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

func getK8sEnpointSlice(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getK8sEnpointSlice")

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

	//
	// Get the manifest resource
	//
	parsedContents, err := fetchResourceFromManifestFileByKind(ctx, d, "EndpointSlice")
	if err != nil {
		return nil, err
	}

	for _, content := range parsedContents {
		endpointSlice := content.Data.(*v1beta1.EndpointSlice)

		if endpointSlice.Name == name && endpointSlice.Namespace == namespace {
			return EndpointSlice{*endpointSlice, content.Path}, nil
		}
	}

	//
	// Get the deployed resource
	//
	if clientset == nil {
		return nil, nil
	}

	endpointSlice, err := clientset.DiscoveryV1beta1().EndpointSlices(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil && !isNotFoundError(err) {
		return nil, err
	}

	return EndpointSlice{*endpointSlice, ""}, nil
}

//// TRANSFORM FUNCTIONS

func transformEndpointSliceTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(EndpointSlice)
	return mergeTags(obj.Labels, obj.Annotations), nil
}
