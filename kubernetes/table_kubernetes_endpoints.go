package kubernetes

import (
	"context"
	"strings"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
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
			Hydrate:    listK8sEnpoints,
			KeyColumns: getCommonOptionalKeyQuals(),
		},
		Columns: k8sCommonColumns([]*plugin.Column{
			{
				Name:        "subsets",
				Type:        proto.ColumnType_JSON,
				Description: "List of addresses and ports that comprise a service.",
			},
			{
				Name:        "context_name",
				Type:        proto.ColumnType_STRING,
				Description: "Kubectl config context name.",
				Hydrate:     endpointResourceContext,
			},
			{
				Name:        "source_type",
				Type:        proto.ColumnType_STRING,
				Description: "The source of the resource. Possible values are: deployed and manifest. If the resource is fetched from the spec file the value will be manifest.",
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

type Endpoints struct {
	v1.Endpoints
	parsedContent
}

//// HYDRATE FUNCTIONS

func listK8sEnpoints(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listK8sEnpoints")

	// Get the client for querying the K8s APIs for the provided context.
	// If the connection is configured for the manifest files, the client will return nil.
	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		return nil, err
	}

	// Check for manifest files
	parsedContents, err := fetchResourceFromManifestFileByKind(ctx, d, "Endpoints")
	if err != nil {
		return nil, err
	}

	for _, content := range parsedContents {
		endpoints := content.Data.(*v1.Endpoints)

		d.StreamListItem(ctx, Endpoints{*endpoints, content})

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

	var response *v1.EndpointsList
	pageLeft := true

	for pageLeft {
		response, err = clientset.CoreV1().Endpoints("").List(ctx, input)
		if err != nil {
			return nil, err
		}

		if response.GetContinue() != "" {
			input.Continue = response.Continue
		} else {
			pageLeft = false
		}

		for _, endpoint := range response.Items {
			d.StreamListItem(ctx, Endpoints{endpoint, parsedContent{SourceType: "deployed"}})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

func getK8sEndpoint(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getK8sEndpoint")

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

	// Get the manifest resource
	parsedContents, err := fetchResourceFromManifestFileByKind(ctx, d, "Endpoints")
	if err != nil {
		return nil, err
	}

	for _, content := range parsedContents {
		endpoints := content.Data.(*v1.Endpoints)

		if endpoints.Name == name && endpoints.Namespace == namespace {
			return Endpoints{*endpoints, content}, nil
		}
	}

	// Get the deployed resource
	if clientset == nil {
		return nil, nil
	}

	endpoint, err := clientset.CoreV1().Endpoints(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil && !isNotFoundError(err) {
		return nil, err
	}

	return Endpoints{*endpoint, parsedContent{SourceType: "deployed"}}, nil
}

func endpointResourceContext(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	obj := h.Item.(Endpoints)

	// Set the context_name as nil
	data := map[string]interface{}{}
	if obj.Path != "" {
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

//// TRANSFORM FUNCTIONS

func transformEndpointTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(Endpoints)
	return mergeTags(obj.Labels, obj.Annotations), nil
}
