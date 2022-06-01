package kubernetes

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

func tableKubernetesNode(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_node",
		Description: "Kubernetes Node is a worker node in Kubernetes.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getK8sNode,
		},
		List: &plugin.ListConfig{
			Hydrate: listK8sNodes,
		},
		Columns: k8sCommonGlobalColumns([]*plugin.Column{
			//// NodeSpec
			{
				Name:        "pod_cidr",
				Type:        proto.ColumnType_CIDR,
				Description: "Pod IP range assigned to the node.",
				Transform:   transform.FromField("Spec.PodCIDR"),
			},
			{
				Name:        "pod_cidrs",
				Type:        proto.ColumnType_JSON,
				Description: "List of the IP ranges assigned to the node for usage by Pods.",
				Transform:   transform.FromField("Spec.PodCIDRs"),
			},
			{
				Name:        "provider_id",
				Type:        proto.ColumnType_STRING,
				Description: "ID of the node assigned by the cloud provider in the format: <ProviderName>://<ProviderSpecificNodeID>.",
				Transform:   transform.FromField("Spec.ProviderID"),
			},
			{
				Name:        "unschedulable",
				Type:        proto.ColumnType_BOOL,
				Description: "Unschedulable controls node schedulability of new pods. By default, node is schedulable.",
				Transform:   transform.FromField("Spec.Unschedulable"),
			},
			{
				Name:        "taints",
				Type:        proto.ColumnType_JSON,
				Description: "List of the taints attached to the node to has the \"effect\" on pod that does not tolerate the Taint",
				Transform:   transform.FromField("Spec.Taints"),
			},
			{
				Name:        "config_source",
				Type:        proto.ColumnType_JSON,
				Description: "The source to get node configuration from.",
				Transform:   transform.FromField("Spec.ConfigSource"),
			},

			//// NodeStatus Columns
			{
				Name:        "capacity",
				Type:        proto.ColumnType_JSON,
				Description: "Capacity represents the total resources of a node.",
				Transform:   transform.FromField("Status.Capacity"),
			},
			{
				Name:        "allocatable",
				Type:        proto.ColumnType_JSON,
				Description: "Allocatable represents the resources of a node that are available for scheduling. Defaults to capacity.",
				Transform:   transform.FromField("Status.Allocatable"),
			},
			{
				Name:        "phase",
				Type:        proto.ColumnType_STRING,
				Description: "The recently observed lifecycle phase of the node.",
				Transform:   transform.FromField("Status.Phase"),
			},
			{
				Name:        "conditions",
				Type:        proto.ColumnType_JSON,
				Description: "List of current observed node conditions.",
				Transform:   transform.FromField("Status.Conditions"),
			},
			{
				Name:        "addresses",
				Type:        proto.ColumnType_JSON,
				Description: "Endpoints of daemons running on the Node.",
				Transform:   transform.FromField("Status.Addresses"),
			},
			{
				Name:        "daemon_endpoints",
				Type:        proto.ColumnType_JSON,
				Description: "Set of ids/uuids to uniquely identify the node.",
				Transform:   transform.FromField("Status.DaemonEndpoints"),
			},
			{
				Name:        "node_info",
				Type:        proto.ColumnType_JSON,
				Description: "List of container images on this node.",
				Transform:   transform.FromField("Status.NodeInfo"),
			},
			{
				Name:        "images",
				Type:        proto.ColumnType_JSON,
				Description: "List of container images on this node.",
				Transform:   transform.FromField("Status.Images"),
			},
			{
				Name:        "volumes_in_use",
				Type:        proto.ColumnType_JSON,
				Description: "List of attachable volumes in use (mounted) by the node.",
				Transform:   transform.FromField("Status.VolumesInUse"),
			},
			{
				Name:        "volumes_attached",
				Type:        proto.ColumnType_JSON,
				Description: "List of volumes that are attached to the node.",
				Transform:   transform.FromField("Status.VolumesAttached"),
			},
			{
				Name:        "config",
				Type:        proto.ColumnType_JSON,
				Description: "Status of the config assigned to the node via the dynamic Kubelet config feature.",
				Transform:   transform.FromField("Status.Config"),
			},
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
				Transform:   transform.From(transformNodeTags),
			},
		}),
	}
}

//// HYDRATE FUNCTIONS

func listK8sNodes(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listK8sNodes")

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

	var response *v1.NodeList
	pageLeft := true

	for pageLeft {
		response, err = clientset.CoreV1().Nodes().List(ctx, input)
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

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

func getK8sNode(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getK8sNode")

	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		return nil, err
	}

	name := d.KeyColumnQuals["name"].GetStringValue()

	// return if name is empty
	if name == "" {
		return nil, nil
	}

	node, err := clientset.CoreV1().Nodes().Get(ctx, name, metav1.GetOptions{})
	if err != nil && !isNotFoundError(err) {
		return nil, err
	}

	return *node, nil
}

//// TRANSFORM FUNCTIONS

func transformNodeTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(v1.Node)
	return mergeTags(obj.Labels, obj.Annotations), nil
}
