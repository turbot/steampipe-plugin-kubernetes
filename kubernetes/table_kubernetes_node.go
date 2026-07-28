package kubernetes

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
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
				Name:        "capacity_cpu",
				Type:        proto.ColumnType_STRING,
				Description: "Raw capacity CPU value as provided by the system.",
				Transform:   transform.FromP(transformNodeCpuAndMemory, "Capacity.CPU"),
			},
			{
				Name:        "capacity_memory",
				Type:        proto.ColumnType_STRING,
				Description: "Raw capacity memory value as provided by the system.",
				Transform:   transform.FromP(transformNodeCpuAndMemory, "Capacity.Memory"),
			},
			{
				Name:        "allocatable_cpu",
				Type:        proto.ColumnType_STRING,
				Description: "Raw allocatable CPU value as provided by the system.",
				Transform:   transform.FromP(transformNodeCpuAndMemory, "Allocatable.CPU"),
			},
			{
				Name:        "allocatable_memory",
				Type:        proto.ColumnType_STRING,
				Description: "Raw allocatable memory value as provided by the system.",
				Transform:   transform.FromP(transformNodeCpuAndMemory, "Allocatable.Memory"),
			},
			{
				Name:        "capacity_cpu_std",
				Type:        proto.ColumnType_INT,
				Description: "Standardized capacity CPU value in millicores (m).",
				Transform:   transform.FromP(transformNodeCpuAndMemoryUnit, "Capacity.CPU"),
			},
			{
				Name:        "capacity_memory_std",
				Type:        proto.ColumnType_INT,
				Description: "Standardized capacity memory value in bytes.",
				Transform:   transform.FromP(transformNodeCpuAndMemoryUnit, "Capacity.Memory"),
			},
			{
				Name:        "allocatable_cpu_std",
				Type:        proto.ColumnType_INT,
				Description: "Standardized allocatable CPU value in millicores (m).",
				Transform:   transform.FromP(transformNodeCpuAndMemoryUnit, "Allocatable.CPU"),
			},
			{
				Name:        "allocatable_memory_std",
				Type:        proto.ColumnType_INT,
				Description: "Standardized allocatable memory value in bytes.",
				Transform:   transform.FromP(transformNodeCpuAndMemoryUnit, "Allocatable.Memory"),
			},
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
			{
				Name:        "context_name",
				Type:        proto.ColumnType_STRING,
				Description: "Kubectl config context name.",
				Hydrate:     getNodeResourceContext,
			},
			{
				Name:        "source_type",
				Type:        proto.ColumnType_STRING,
				Description: "The source of the resource. Possible values are: deployed and manifest. If the resource is fetched from the spec file the value will be manifest.",
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

type Node struct {
	v1.Node
	parsedContent
}

//// HYDRATE FUNCTIONS

func listK8sNodes(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listK8sNodes")

	// Get the client for querying the K8s APIs for the provided context.
	// If the connection is configured for the manifest files, the client will return nil.
	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		return nil, err
	}

	// Check for manifest files
	parsedContents, err := fetchResourceFromManifestFileByKind(ctx, d, "Node")
	if err != nil {
		return nil, err
	}

	for _, content := range parsedContents {
		node := content.ParsedData.(*v1.Node)

		d.StreamListItem(ctx, Node{*node, content})

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

		for _, node := range response.Items {
			d.StreamListItem(ctx, Node{node, parsedContent{SourceType: "deployed"}})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

func getK8sNode(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getK8sNode")

	// Get the client for querying the K8s APIs for the provided context.
	// If the connection is configured for the manifest files, the client will return nil.
	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		return nil, err
	}

	name := d.EqualsQuals["name"].GetStringValue()

	// return if name is empty
	if name == "" {
		return nil, nil
	}

	// Get the manifest resource
	parsedContents, err := fetchResourceFromManifestFileByKind(ctx, d, "Node")
	if err != nil {
		return nil, err
	}

	for _, content := range parsedContents {
		node := content.ParsedData.(*v1.Node)

		if node.Name == name {
			return Node{*node, content}, nil
		}
	}

	// Get the deployed resource
	if clientset == nil {
		return nil, nil
	}

	node, err := clientset.CoreV1().Nodes().Get(ctx, name, metav1.GetOptions{})
	if err != nil && !isNotFoundError(err) {
		return nil, err
	}

	return Node{*node, parsedContent{SourceType: "deployed"}}, nil
}

func getNodeResourceContext(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	obj := h.Item.(Node)

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

func transformNodeTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(Node)
	return mergeTags(obj.Labels, obj.Annotations), nil
}

func transformNodeCpuAndMemoryUnit(_ context.Context, d *transform.TransformData) (interface{}, error) {
	param := d.Param.(string)

	node := d.HydrateItem.(Node)

	switch param {
	case "Capacity.CPU":
		return normalizeCPUToMilliCores(node.Status.Capacity.Cpu().String())
	case "Capacity.Memory":
		return normalizeMemoryToBytes(node.Status.Capacity.Memory().String())
	case "Allocatable.CPU":
		return normalizeCPUToMilliCores(node.Status.Allocatable.Cpu().String())
	case "Allocatable.Memory":
		return normalizeMemoryToBytes(node.Status.Allocatable.Memory().String())
	}

	return nil, nil
}

func transformNodeCpuAndMemory(_ context.Context, d *transform.TransformData) (interface{}, error) {
	param := d.Param.(string)

	node := d.HydrateItem.(Node)

	switch param {
	case "Capacity.CPU":
		return node.Status.Capacity.Cpu().String(), nil
	case "Capacity.Memory":
		return node.Status.Capacity.Memory().String(), nil
	case "Allocatable.CPU":
		return node.Status.Allocatable.Cpu().String(), nil
	case "Allocatable.Memory":
		return node.Status.Allocatable.Memory().String(), nil
	}

	return nil, nil
}
