package kubernetes

import (
	"context"

	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/turbot/steampipe-plugin-sdk/v2/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin/transform"
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
			Hydrate: listK8sNetworkPolicies,
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

//// HYDRATE FUNCTIONS

func listK8sNetworkPolicies(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listK8sNetworkPolicies")

	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		return nil, err
	}

	networkPolicyList, err := clientset.NetworkingV1().NetworkPolicies("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, networkPolicy := range networkPolicyList.Items {
		d.StreamListItem(ctx, networkPolicy)
	}

	return nil, nil
}

func getK8sNetworkPolicy(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getK8sNetworkPolicy")

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

	networkPolicy, err := clientset.NetworkingV1().NetworkPolicies(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil && !isNotFoundError(err) {
		return nil, err
	}

	return *networkPolicy, nil
}

//// TRANSFORM FUNCTIONS

func transformNetworkPolicyTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(v1.NetworkPolicy)
	return mergeTags(obj.Labels, obj.Annotations), nil
}
