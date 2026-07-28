package kubernetes

import (
	"context"
	"strings"

	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableKubernetesIngress(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_ingress",
		Description: "Ingress exposes HTTP and HTTPS routes from outside the cluster to services within the cluster. Traffic routing is controlled by rules defined on the Ingress resource.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "namespace"}),
			Hydrate:    getK8sIngress,
		},
		List: &plugin.ListConfig{
			Hydrate:    listK8sIngresses,
			KeyColumns: getCommonOptionalKeyQuals(),
		},
		Columns: k8sCommonColumns([]*plugin.Column{
			//// IngressSpec columns
			{
				Name:        "ingress_class_name",
				Type:        proto.ColumnType_STRING,
				Description: "Name of the IngressClass cluster resource. The associated IngressClass defines which controller will implement the resource.",
				Transform:   transform.FromField("Spec.IngressClassName"),
			},
			{
				Name:        "default_backend",
				Type:        proto.ColumnType_JSON,
				Description: "A default backend capable of servicing requests that don't match any rule. At least one of 'backend' or 'rules' must be specified.",
				Transform:   transform.FromField("Spec.DefaultBackend"),
			},
			{
				Name:        "tls",
				Type:        proto.ColumnType_JSON,
				Description: "TLS configuration.",
				Transform:   transform.FromField("Spec.TLS"),
			},
			{
				Name:        "rules",
				Type:        proto.ColumnType_JSON,
				Description: "A list of host rules used to configure the Ingress.",
				Transform:   transform.FromField("Spec.Rules"),
			},
			{
				Name:        "context_name",
				Type:        proto.ColumnType_STRING,
				Description: "Kubectl config context name.",
				Hydrate:     getIngressResourceContext,
			},
			{
				Name:        "source_type",
				Type:        proto.ColumnType_STRING,
				Description: "The source of the resource. Possible values are: deployed and manifest. If the resource is fetched from the spec file the value will be manifest.",
			},

			//// IngressStatus columns
			{
				Name:        "load_balancer",
				Type:        proto.ColumnType_JSON,
				Description: "a list containing ingress points for the load-balancer. Traffic intended for the service should be sent to these ingress points.",
				Transform:   transform.FromField("Status.LoadBalancer.Ingress"),
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
				Transform:   transform.From(transformIngressTags),
			},
		}),
	}
}

type Ingress struct {
	v1.Ingress
	parsedContent
}

//// HYDRATE FUNCTIONS

func listK8sIngresses(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listK8sIngresses")

	// Get the client for querying the K8s APIs for the provided context.
	// If the connection is configured for the manifest files, the client will return nil.
	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		return nil, err
	}

	// Check for manifest files
	parsedContents, err := fetchResourceFromManifestFileByKind(ctx, d, "Ingress")
	if err != nil {
		return nil, err
	}

	for _, content := range parsedContents {
		ingress := content.ParsedData.(*v1.Ingress)

		d.StreamListItem(ctx, Ingress{*ingress, content})

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

	commonFieldSelectorValue := getCommonOptionalKeyQualsValueForFieldSelector(d)

	if len(commonFieldSelectorValue) > 0 {
		input.FieldSelector = strings.Join(commonFieldSelectorValue, ",")
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

	var response *v1.IngressList
	pageLeft := true

	for pageLeft {
		response, err = clientset.NetworkingV1().Ingresses(d.EqualsQualString("namespace")).List(ctx, input)
		if err != nil {
			return nil, err
		}

		if response.GetContinue() != "" {
			input.Continue = response.Continue
		} else {
			pageLeft = false
		}

		for _, ingress := range response.Items {
			d.StreamListItem(ctx, Ingress{ingress, parsedContent{SourceType: "deployed"}})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

func getK8sIngress(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getK8sIngress")

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
	parsedContents, err := fetchResourceFromManifestFileByKind(ctx, d, "Ingress")
	if err != nil {
		return nil, err
	}

	for _, content := range parsedContents {
		ingress := content.ParsedData.(*v1.Ingress)

		if ingress.Name == name && ingress.Namespace == namespace {
			return Ingress{*ingress, content}, nil
		}
	}

	// Get the deployed resource
	if clientset == nil {
		return nil, nil
	}

	ingress, err := clientset.NetworkingV1().Ingresses(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil && !isNotFoundError(err) {
		return nil, err
	}

	return Ingress{*ingress, parsedContent{SourceType: "deployed"}}, nil
}

func getIngressResourceContext(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	obj := h.Item.(Ingress)

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

func transformIngressTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(Ingress)
	return mergeTags(obj.Labels, obj.Annotations), nil
}
