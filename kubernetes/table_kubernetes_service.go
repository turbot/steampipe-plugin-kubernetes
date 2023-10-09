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

func tableKubernetesService(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_service",
		Description: "A service provides an abstract way to expose an application running on a set of Pods as a network service.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "namespace"}),
			Hydrate:    getK8sService,
		},
		List: &plugin.ListConfig{
			Hydrate:    listK8sServices,
			KeyColumns: getCommonOptionalKeyQuals(),
		},
		// Service is namespaced resource.
		Columns: k8sCommonColumns([]*plugin.Column{
			{
				Name:        "type",
				Type:        proto.ColumnType_STRING,
				Description: "Type determines how the Service is exposed.",
				Transform:   transform.FromField("Spec.Type").Transform(transform.ToString),
			},
			{
				Name:        "allocate_load_balancer_node_ports",
				Type:        proto.ColumnType_BOOL,
				Description: "Indicates whether NodePorts will be automatically allocated for services with type LoadBalancer, or not.",
				Transform:   transform.FromField("Spec.AllocateLoadBalancerNodePorts"),
			},
			{
				Name:        "cluster_ip",
				Type:        proto.ColumnType_STRING,
				Description: "IP address of the service and is usually assigned randomly.",
				Transform:   transform.FromField("Spec.ClusterIP"),
			},
			{
				Name:        "external_name",
				Type:        proto.ColumnType_STRING,
				Description: "The external reference that discovery mechanisms will return as an alias for this service (e.g. a DNS CNAME record).",
				Transform:   transform.FromField("Spec.ExternalName"),
			},
			{
				Name:        "external_traffic_policy",
				Type:        proto.ColumnType_STRING,
				Description: "Denotes whether the service desires to route external traffic to node-local or cluster-wide endpoints.",
				Transform:   transform.FromField("Spec.ExternalTrafficPolicy").Transform(transform.ToString),
			},
			{
				Name:        "health_check_node_port",
				Type:        proto.ColumnType_INT,
				Description: "Specifies the healthcheck nodePort for the service.",
				Transform:   transform.FromField("Spec.HealthCheckNodePort"),
			},
			{
				Name:        "ip_family_policy",
				Type:        proto.ColumnType_STRING,
				Description: "Specifies the dual-stack-ness requested or required by this service, and is gated by the 'IPv6DualStack' feature gate.",
				Transform:   transform.FromField("Spec.IPFamilyPolicy").Transform(transform.ToString),
			},
			{
				Name:        "load_balancer_ip",
				Type:        proto.ColumnType_IPADDR,
				Description: "The IP specified when the load balancer was created.",
				Transform:   transform.FromField("Spec.LoadBalancerIP"),
			},
			{
				Name:        "publish_not_ready_addresses",
				Type:        proto.ColumnType_BOOL,
				Description: "Indicates that any agent which deals with endpoints for this service should disregard any indications of ready/not-ready.",
				Transform:   transform.FromField("Spec.PublishNotReadyAddresses"),
			},
			{
				Name:        "session_affinity",
				Type:        proto.ColumnType_STRING,
				Description: "Supports 'ClientIP' and 'None'. Used to maintain session affinity.",
				Transform:   transform.FromField("Spec.SessionAffinity").Transform(transform.ToString),
			},
			{
				Name:        "session_affinity_client_ip_timeout",
				Type:        proto.ColumnType_INT,
				Description: "Specifies the ClientIP type session sticky time in seconds.",
				Transform:   transform.FromField("Spec.SessionAffinityConfig.ClientIP.TimeoutSeconds"),
			},
			{
				Name:        "cluster_ips",
				Type:        proto.ColumnType_JSON,
				Description: "A list of IP addresses assigned to this service, and are usually assigned randomly.",
				Transform:   transform.FromField("Spec.ClusterIPs"),
			},
			{
				Name:        "external_ips",
				Type:        proto.ColumnType_JSON,
				Description: "A list of IP addresses for which nodes in the cluster will also accept traffic for this service.",
				Transform:   transform.FromField("Spec.ExternalIPs"),
			},
			{
				Name:        "ip_families",
				Type:        proto.ColumnType_JSON,
				Description: "A list of IP families (e.g. IPv4, IPv6) assigned to this service, and is gated by the 'IPv6DualStack' feature gate.",
				Transform:   transform.FromField("Spec.IPFamilies"),
			},
			{
				Name:        "load_balancer_ingress",
				Type:        proto.ColumnType_JSON,
				Description: "A list containing ingress points for the load-balancer.",
				Transform:   transform.FromField("Status.LoadBalancer.Ingress"),
			},
			{
				Name:        "load_balancer_source_ranges",
				Type:        proto.ColumnType_JSON,
				Description: "A list of source ranges that will restrict traffic through the cloud-provider load-balancer will be restricted to the specified client IPs.",
				Transform:   transform.FromField("Spec.LoadBalancerSourceRanges"),
			},
			{
				Name:        "ports",
				Type:        proto.ColumnType_JSON,
				Description: "A list of ports that are exposed by this service.",
				Transform:   transform.FromField("Spec.Ports"),
			},
			{
				Name:        "selector_query",
				Type:        proto.ColumnType_STRING,
				Description: "A query string representation of the selector.",
				Transform:   transform.FromField("Spec.Selector").Transform(selectorMapToString),
			},
			{
				Name:        "selector",
				Type:        proto.ColumnType_JSON,
				Description: "Route service traffic to pods with label keys and values matching this selector.",
				Transform:   transform.FromField("Spec.Selector"),
			},
			{
				Name:        "topology_keys",
				Type:        proto.ColumnType_JSON,
				Description: "A preference-order list of topology keys which implementations of services should use to preferentially sort endpoints when accessing this Service, it can not be used at the same time as externalTrafficPolicy=Local.",
				Transform:   transform.FromField("Spec.TopologyKeys"),
			},
			{
				Name:        "context_name",
				Type:        proto.ColumnType_STRING,
				Description: "Kubectl config context name.",
				Hydrate:     getServiceResourceContext,
			},
			{
				Name:        "source_type",
				Type:        proto.ColumnType_STRING,
				Description: "The source of the resource. Possible values are: deployed and manifest. If the resource is fetched from the spec file the value will be manifest.",
			},

			// Steampipe Standard Columns
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
				Transform:   transform.From(transformServiceTags),
			},
		}),
	}
}

type Service struct {
	v1.Service
	parsedContent
}

//// HYDRATE FUNCTIONS

func listK8sServices(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listK8sServices")

	// Get the client for querying the K8s APIs for the provided context.
	// If the connection is configured for the manifest files, the client will return nil.
	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		return nil, err
	}

	// Check for manifest files
	parsedContents, err := fetchResourceFromManifestFileByKind(ctx, d, "Service")
	if err != nil {
		return nil, err
	}

	for _, content := range parsedContents {
		service := content.ParsedData.(*v1.Service)

		d.StreamListItem(ctx, Service{*service, content})

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

	var response *v1.ServiceList
	pageLeft := true

	for pageLeft {
		response, err = clientset.CoreV1().Services(d.EqualsQualString("namespace")).List(ctx, input)
		if err != nil {
			return nil, err
		}

		if response.GetContinue() != "" {
			input.Continue = response.Continue
		} else {
			pageLeft = false
		}

		for _, service := range response.Items {
			d.StreamListItem(ctx, Service{service, parsedContent{SourceType: "deployed"}})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

func getK8sService(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getK8sService")

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
	parsedContents, err := fetchResourceFromManifestFileByKind(ctx, d, "Service")
	if err != nil {
		return nil, err
	}

	for _, content := range parsedContents {
		service := content.ParsedData.(*v1.Service)

		if service.Name == name && service.Namespace == namespace {
			return Service{*service, content}, nil
		}
	}

	// Get the deployed resource
	if clientset == nil {
		return nil, nil
	}

	service, err := clientset.CoreV1().Services(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil && !isNotFoundError(err) {
		plugin.Logger(ctx).Debug("getK8sService", "Error", err)
		return nil, err
	}

	return Service{*service, parsedContent{SourceType: "deployed"}}, nil
}

func getServiceResourceContext(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	obj := h.Item.(Service)

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

func transformServiceTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(Service)
	return mergeTags(obj.Labels, obj.Annotations), nil
}
