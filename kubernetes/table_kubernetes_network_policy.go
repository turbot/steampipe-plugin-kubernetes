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

func tableKubernetesNetworkPolicy(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_network_policy",
		Description: "Network policy specifiy how pods are allowed to communicate with each other and with other network endpoints.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "namespace"}),
			Hydrate:    getK8sNetworkPolicy,
		},
		List: &plugin.ListConfig{
			Hydrate:    listK8sNetworkPolicies,
			KeyColumns: getCommonOptionalKeyQuals(),
		},
		Columns: k8sCommonColumns([]*plugin.Column{
			//// NetworkPolicySpec
			{
				Name:        "pod_selector",
				Type:        proto.ColumnType_JSON,
				Description: "Selects the pods to which this NetworkPolicy object applies. The array of ingress rules is applied to any pods selected by this field. An empty podSelector matches all pods in this namespace.",
				Transform:   transform.FromField("Spec.PodSelector"),
			},
			{
				Name:        "ingress",
				Type:        proto.ColumnType_JSON,
				Description: "List of ingress rules to be applied to the selected pods. If this field is empty then this NetworkPolicy does not allow any traffic (and serves solely to ensure that the pods it selects are isolated by default)",
				Transform:   transform.FromField("Spec.Ingress"),
			},
			{
				Name:        "egress",
				Type:        proto.ColumnType_JSON,
				Description: "List of egress rules to be applied to the selected pods. If this field is empty then this NetworkPolicy limits all outgoing traffic (and serves solely to ensure that the pods it selects are isolated by default).",
				Transform:   transform.FromField("Spec.Egress"),
			},
			{
				Name:        "policy_types",
				Type:        proto.ColumnType_JSON,
				Description: "List of rule types that the NetworkPolicy relates to. Valid options are \"Ingress\", \"Egress\", or \"Ingress,Egress\". If this field is not specified, it will default based on the existence of Ingress or Egress rules.",
				Transform:   transform.FromField("Spec.PolicyTypes"),
			},
			{
				Name:        "context_name",
				Type:        proto.ColumnType_STRING,
				Description: "Kubectl config context name.",
				Hydrate:     getNetworkPolicyResourceContext,
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
				Transform:   transform.From(transformNetworkPolicyTags),
			},
		}),
	}
}

type NetworkPolicy struct {
	v1.NetworkPolicy
	parsedContent
}

//// HYDRATE FUNCTIONS

func listK8sNetworkPolicies(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listK8sNetworkPolicies")

	// Get the client for querying the K8s APIs for the provided context.
	// If the connection is configured for the manifest files, the client will return nil.
	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		return nil, err
	}

	// Check for manifest files
	parsedContents, err := fetchResourceFromManifestFileByKind(ctx, d, "NetworkPolicy")
	if err != nil {
		return nil, err
	}

	for _, content := range parsedContents {
		networkPolicy := content.Data.(*v1.NetworkPolicy)

		d.StreamListItem(ctx, NetworkPolicy{*networkPolicy, content})

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

	var response *v1.NetworkPolicyList
	pageLeft := true

	for pageLeft {
		response, err = clientset.NetworkingV1().NetworkPolicies("").List(ctx, input)
		if err != nil {
			return nil, err
		}

		if response.GetContinue() != "" {
			input.Continue = response.Continue
		} else {
			pageLeft = false
		}

		for _, networkPolicy := range response.Items {
			d.StreamListItem(ctx, NetworkPolicy{networkPolicy, parsedContent{SourceType: "deployed"}})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

func getK8sNetworkPolicy(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getK8sNetworkPolicy")

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
	parsedContents, err := fetchResourceFromManifestFileByKind(ctx, d, "NetworkPolicy")
	if err != nil {
		return nil, err
	}

	for _, content := range parsedContents {
		networkPolicy := content.Data.(*v1.NetworkPolicy)

		if networkPolicy.Name == name && networkPolicy.Namespace == namespace {
			return NetworkPolicy{*networkPolicy, content}, nil
		}
	}

	// Get the deployed resource
	if clientset == nil {
		return nil, nil
	}

	networkPolicy, err := clientset.NetworkingV1().NetworkPolicies(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil && !isNotFoundError(err) {
		return nil, err
	}

	return NetworkPolicy{*networkPolicy, parsedContent{SourceType: "deployed"}}, nil
}

func getNetworkPolicyResourceContext(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	obj := h.Item.(NetworkPolicy)

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

func transformNetworkPolicyTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(NetworkPolicy)
	return mergeTags(obj.Labels, obj.Annotations), nil
}
