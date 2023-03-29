package kubernetes

import (
	"context"
	"os"
	"strings"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableKubernetesManifestService(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_manifest_service",
		Description: "",
		List: &plugin.ListConfig{
			ParentHydrate: listKubernetesManifestFiles,
			Hydrate:       listKubernetesManifestServices,
			KeyColumns:    plugin.OptionalColumns([]string{"manifest_file_path"}),
		},
		Columns: []*plugin.Column{

			// Metadata Columns
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of resource."},
			{Name: "namespace", Type: proto.ColumnType_STRING, Description: "Namespace defines the space within which each name must be unique."},
			{Name: "labels", Type: proto.ColumnType_JSON, Description: "Map of string keys and values that can be used to organize and categorize (scope and select) objects. May match selectors of replication controllers and services."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "Type determines how the Service is exposed.", Transform: transform.FromField("Spec.Type").Transform(transform.ToString)},
			{Name: "allocate_load_balancer_node_ports", Type: proto.ColumnType_BOOL, Description: "Indicates whether NodePorts will be automatically allocated for services with type LoadBalancer, or not.", Transform: transform.FromField("Spec.AllocateLoadBalancerNodePorts")},
			{Name: "cluster_ip", Type: proto.ColumnType_STRING, Description: "IP address of the service and is usually assigned randomly.", Transform: transform.FromField("Spec.ClusterIP")},
			{Name: "external_name", Type: proto.ColumnType_STRING, Description: "The external reference that discovery mechanisms will return as an alias for this service (e.g. a DNS CNAME record).", Transform: transform.FromField("Spec.ExternalName")},
			{Name: "external_traffic_policy", Type: proto.ColumnType_STRING, Description: "Denotes whether the service desires to route external traffic to node-local or cluster-wide endpoints.", Transform: transform.FromField("Spec.ExternalTrafficPolicy").Transform(transform.ToString)},
			{Name: "health_check_node_port", Type: proto.ColumnType_INT, Description: "Specifies the healthcheck nodePort for the service.", Transform: transform.FromField("Spec.HealthCheckNodePort")},
			{Name: "ip_family_policy", Type: proto.ColumnType_STRING, Description: "Specifies the dual-stack-ness requested or required by this service, and is gated by the 'IPv6DualStack' feature gate.", Transform: transform.FromField("Spec.IPFamilyPolicy").Transform(transform.ToString)},
			{Name: "load_balancer_ip", Type: proto.ColumnType_IPADDR, Description: "The IP specified when the load balancer was created.", Transform: transform.FromField("Spec.LoadBalancerIP")},
			{Name: "publish_not_ready_addresses", Type: proto.ColumnType_BOOL, Description: "Indicates that any agent which deals with endpoints for this service should disregard any indications of ready/not-ready.", Transform: transform.FromField("Spec.PublishNotReadyAddresses")},
			{Name: "session_affinity", Type: proto.ColumnType_STRING, Description: "Supports 'ClientIP' and 'None'. Used to maintain session affinity.", Transform: transform.FromField("Spec.SessionAffinity").Transform(transform.ToString)},
			{Name: "session_affinity_client_ip_timeout", Type: proto.ColumnType_INT, Description: "Specifies the ClientIP type session sticky time in seconds.", Transform: transform.FromField("Spec.SessionAffinityConfig.ClientIP.TimeoutSeconds")},
			{Name: "cluster_ips", Type: proto.ColumnType_JSON, Description: "A list of IP addresses assigned to this service, and are usually assigned randomly.", Transform: transform.FromField("Spec.ClusterIPs")},
			{Name: "external_ips", Type: proto.ColumnType_JSON, Description: "A list of IP addresses for which nodes in the cluster will also accept traffic for this service.", Transform: transform.FromField("Spec.ExternalIPs")},
			{Name: "ip_families", Type: proto.ColumnType_JSON, Description: "A list of IP families (e.g. IPv4, IPv6) assigned to this service, and is gated by the 'IPv6DualStack' feature gate.", Transform: transform.FromField("Spec.IPFamilies")},
			{Name: "load_balancer_ingress", Type: proto.ColumnType_JSON, Description: "A list containing ingress points for the load-balancer.", Transform: transform.FromField("Status.LoadBalancer.Ingress")},
			{Name: "load_balancer_source_ranges", Type: proto.ColumnType_JSON, Description: "A list of source ranges that will restrict traffic through the cloud-provider load-balancer will be restricted to the specified client IPs.", Transform: transform.FromField("Spec.LoadBalancerSourceRanges")},
			{Name: "ports", Type: proto.ColumnType_JSON, Description: "A list of ports that are exposed by this service.", Transform: transform.FromField("Spec.Ports")},
			{Name: "selector_query", Type: proto.ColumnType_STRING, Description: "A query string representation of the selector.", Transform: transform.FromField("Spec.Selector").Transform(selectorMapToString)},
			{Name: "selector", Type: proto.ColumnType_JSON, Description: "Route service traffic to pods with label keys and values matching this selector.", Transform: transform.FromField("Spec.Selector")},
			{Name: "topology_keys", Type: proto.ColumnType_JSON, Description: "A preference-order list of topology keys which implementations of services should use to preferentially sort endpoints when accessing this Service, it can not be used at the same time as externalTrafficPolicy=Local.", Transform: transform.FromField("Spec.TopologyKeys")},

			// Steampipe Standard Columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTitle, Transform: transform.FromField("Name")},
			{Name: "tags", Type: proto.ColumnType_JSON, Description: ColumnDescriptionTags, Transform: transform.From(transformManifestServiceTags)},

			{Name: "manifest_file_path", Description: "Path to the file.", Type: proto.ColumnType_STRING},
		},
	}
}

type KubernetesManifestService struct {
	v1.Service
	ManifestFilePath string
}

//// LIST FUNCTION

func listKubernetesManifestServices(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// The path comes from a parent hydrate, defaulting to the config paths or
	// available by the optional key column
	path := h.Item.(filePath).Path

	// Load the file into a buffer
	content, err := os.ReadFile(path)
	if err != nil {
		plugin.Logger(ctx).Error("kubernetes_manifest_service.listKubernetesManifestServices", "failed to read file", err, "path", path)
		return nil, err
	}
	decoder := scheme.Codecs.UniversalDeserializer()

	// Check for the start of the document
	for _, resource := range strings.Split(string(content), "---") {
		// skip empty documents, `Decode` will fail on them
		if len(resource) == 0 {
			continue
		}

		// Decode the file content
		obj, groupVersionKind, err := decoder.Decode([]byte(resource), nil, nil)
		if err != nil {
			plugin.Logger(ctx).Error("kubernetes_manifest_service.listKubernetesManifestServices", "failed to decode the file", err, "path", path)
			return nil, err
		}

		// Return if the definition is not for the service resource
		if groupVersionKind.Kind == "Service" {
			service := obj.(*v1.Service)

			d.StreamListItem(ctx, KubernetesManifestService{*service, path})

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

func transformManifestServiceTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(KubernetesManifestService)
	return mergeTags(obj.Labels, obj.Annotations), nil
}
