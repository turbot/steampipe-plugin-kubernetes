package kubernetes

import (
	"context"
	"strings"

	v2 "github.com/cilium/cilium/pkg/k8s/apis/cilium.io/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableKubernetesCiliumNetworkPolicy(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_cilium_network_policy",
		Description: "Cilium network policy specifiy how pods are allowed to communicate with each other and with other network endpoints.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "namespace"}),
			Hydrate:    getK8sCiliumNetworkPolicy,
		},
		List: &plugin.ListConfig{
			Hydrate:    listK8sCiliumNetworkPolicies,
			KeyColumns: getCommonOptionalKeyQuals(),
		},
		Columns: k8sCommonColumns([]*plugin.Column{
			//// CiliumNetworkPolicySpec
			{
				Name:        "endpoint_selector",
				Type:        proto.ColumnType_JSON,
				Description: "EndpointSelector selects all endpoints which should be subject to this rule. EndpointSelector and NodeSelector cannot be both empty and are mutually exclusive.",
				Transform:   transform.FromField("Spec.EndpointSelector"),
			},
			{
				Name:        "ingress_allow",
				Type:        proto.ColumnType_JSON,
				Description: "List of IngressRule which are enforced at ingress. If omitted or empty, this rule does not apply at ingress.",
				Transform:   transform.FromField("Spec.Ingress"),
			},
			{
				Name:        "ingress_deny",
				Type:        proto.ColumnType_JSON,
				Description: "List of IngressDenyRule which are enforced at ingress. Any rule inserted here will by denied regardless of the allowed ingress rules in the 'ingress' field. If omitted or empty, this rule does not apply at ingress.",
				Transform:   transform.FromField("Spec.IngressDeny"),
			},
			{
				Name:        "egress_allow",
				Type:        proto.ColumnType_JSON,
				Description: "List of EgressRule which are enforced at egress. If omitted or empty, this rule does not apply at egress.",
				Transform:   transform.FromField("Spec.Egress"),
			},
			{
				Name:        "egress_deny",
				Type:        proto.ColumnType_JSON,
				Description: "List of EgressDenyRule which are enforced at egress. Any rule inserted here will by denied regardless of the allowed egress rules in the 'egress' field. If omitted or empty, this rule does not apply at egress.",
				Transform:   transform.FromField("Spec.EgressDeny"),
			},
			{
				Name:        "description",
				Type:        proto.ColumnType_STRING,
				Description: "Description is a free form string, it can be used by the creator of the rule to store human readable explanation of the purpose of this rule. Rules cannot be identified by comment.",
				Transform:   transform.FromField("Spec.Description"),
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
				Transform:   transform.From(transformCiliumNetworkPolicyTags),
			},
		}),
	}
}

//// HYDRATE FUNCTIONS

func listK8sCiliumNetworkPolicies(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listK8sCilliumNetworkPolicies")

	clientset, err := GetNewClientCilium(ctx, d)
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

	var response *v2.CiliumNetworkPolicyList
	pageLeft := true

	for pageLeft {
		response, err = clientset.CiliumV2().CiliumNetworkPolicies("").List(ctx, input)
		if err != nil {
			return nil, err
		}

		if response.GetContinue() != "" {
			input.Continue = response.Continue
		} else {
			pageLeft = false
		}

		for _, networkPolicy := range response.Items {
			d.StreamListItem(ctx, networkPolicy)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

func getK8sCiliumNetworkPolicy(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getK8sNetworkPolicy")

	clientset, err := GetNewClientCilium(ctx, d)
	if err != nil {
		return nil, err
	}

	name := d.EqualsQuals["name"].GetStringValue()
	namespace := d.EqualsQuals["namespace"].GetStringValue()

	// return if namespace or name is empty
	if namespace == "" || name == "" {
		return nil, nil
	}

	networkPolicy, err := clientset.CiliumV2().CiliumNetworkPolicies(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil && !isNotFoundError(err) {
		return nil, err
	}

	return *networkPolicy, nil
}

//// TRANSFORM FUNCTIONS

func transformCiliumNetworkPolicyTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(v2.CiliumNetworkPolicy)
	return mergeTags(obj.Labels, obj.Annotations), nil
}
