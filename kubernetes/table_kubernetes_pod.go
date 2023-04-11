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

func tableKubernetesPod(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_pod",
		Description: "Kubernetes Pod is a collection of containers that can run on a host. This resource is created by clients and scheduled onto hosts.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "namespace"}),
			Hydrate:    getK8sPod,
		},
		List: &plugin.ListConfig{
			Hydrate: listK8sPods,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "selector_search", Require: plugin.Optional, CacheMatch: "exact"},
				{Name: "restart_policy", Require: plugin.Optional},       // spec.retryPolicy spec.serviceAccountName
				{Name: "service_account_name", Require: plugin.Optional}, // spec.serviceAccountName
				{Name: "scheduler_name", Require: plugin.Optional},       // spec.schedulerName
				{Name: "phase", Require: plugin.Optional},                // status.phase
				{Name: "nominated_node_name", Require: plugin.Optional},  // status.nominatedNodeName
				{Name: "pod_ip", Require: plugin.Optional},               // status.podIP
				{Name: "name", Require: plugin.Optional},
				{Name: "namespace", Require: plugin.Optional},
			},
		},
		Columns: k8sCommonColumns([]*plugin.Column{
			//// PodSpec Columns
			{
				Name:        "selector_search",
				Type:        proto.ColumnType_STRING,
				Description: "A label selector string to restrict the list of returned objects by their labels.",
				Transform:   transform.FromQual("selector_search"),
			},
			{
				Name:        "volumes",
				Type:        proto.ColumnType_JSON,
				Description: "List of volumes that can be mounted by containers belonging to the pod.",
				Transform:   transform.FromField("Spec.Volumes"),
			},
			{
				Name:        "containers",
				Type:        proto.ColumnType_JSON,
				Description: "List of containers belonging to the pod.",
				Transform:   transform.FromField("Spec.Containers"),
			},
			{
				Name: "ephemeral_containers",
				Type: proto.ColumnType_JSON,
				Description: "List of ephemeral containers run in this pod. Ephemeral containers may be run in an existing " +
					"pod to perform user-initiated actions such as debugging. This list cannot be specified when " +
					"creating a pod, and it cannot be modified by updating the pod spec. In order to add an " +
					"ephemeral container to an existing pod, use the pod's ephemeralcontainers subresource. " +
					"This field is alpha-level and is only honored by servers that enable the EphemeralContainers feature.",
				Transform: transform.FromField("Spec.EphemeralContainers"),
			},
			{
				Name: "init_containers",
				Type: proto.ColumnType_JSON,
				Description: "List of initialization containers belonging to the pod. Init containers " +
					"are executed in order prior to containers being started. If any " +
					"init container fails, the pod is considered to have failed and is handled according " +
					"to its restartPolicy. The name for an init container or normal container must be " +
					"unique among all containers.",
				Transform: transform.FromField("Spec.InitContainers"),
			},
			{
				Name:        "restart_policy",
				Type:        proto.ColumnType_STRING,
				Description: "Restart policy for all containers within the pod. One of Always, OnFailure, Never.",
				Transform:   transform.FromField("Spec.RestartPolicy"),
			},
			{
				Name: "termination_grace_period_seconds",
				Type: proto.ColumnType_INT,
				Description: "Optional duration in seconds the pod needs to terminate gracefully. May be decreased in delete request. " +
					"Value must be non-negative integer. The value zero indicates delete immediately. " +
					"If this value is nil, the default grace period will be used instead. " +
					"The grace period is the duration in seconds after the processes running in the pod are sent " +
					"a termination signal and the time when the processes are forcibly halted with a kill signal. " +
					"Set this value longer than the expected cleanup time for your process.",
				Transform: transform.FromField("Spec.TerminationGracePeriodSeconds"),
			},
			{
				Name: "active_deadline_seconds",
				Type: proto.ColumnType_STRING,
				Description: "Optional duration in seconds the pod may be active on the node relative to " +
					"StartTime before the system will actively try to mark it failed and kill associated containers.",
				Transform: transform.FromField("Spec.ActiveDeadlineSeconds"),
			},
			{
				Name:        "dns_policy",
				Type:        proto.ColumnType_STRING,
				Description: "DNS policy for pod.  Valid values are 'ClusterFirstWithHostNet', 'ClusterFirst', 'Default' or 'None'.",
				Transform:   transform.FromField("Spec.DNSPolicy"),
			},
			{
				Name:        "node_selector",
				Type:        proto.ColumnType_JSON,
				Description: "NodeSelector is a selector which must be true for the pod to fit on a node.",
				Transform:   transform.FromField("Spec.NodeSelector"),
			},
			{
				Name:        "service_account_name",
				Type:        proto.ColumnType_STRING,
				Description: "ServiceAccountName is the name of the ServiceAccount to use to run this pod.",
				Transform:   transform.FromField("Spec.ServiceAccountName"),
			},
			// // Dont include deprecated columns...
			// {
			// 	Name: "deprecated_service_account",
			// 	Type: proto.ColumnType_STRING,
			// 	Description: "DeprecatedServiceAccount is a depreciated alias for ServiceAccountName. " +
			// 		"Deprecated: Use serviceAccountName instead.",
			// 	Transform: transform.FromField("Spec.DeprecatedServiceAccount"),
			// },
			{
				Name:        "automount_service_account_token",
				Type:        proto.ColumnType_BOOL,
				Description: "AutomountServiceAccountToken indicates whether a service account token should be automatically mounted.",
				Transform:   transform.FromField("Spec.AutomountServiceAccountToken"),
			},
			{
				Name: "node_name",
				Type: proto.ColumnType_STRING,
				Description: "NodeName is a request to schedule this pod onto a specific node. If it is non-empty, " +
					"the scheduler simply schedules this pod onto that node, assuming that it fits resource " +
					"requirements.",
				Transform: transform.FromField("Spec.NodeName"),
			},
			{
				Name: "host_network",
				Type: proto.ColumnType_BOOL,
				Description: "Host networking requested for this pod. Use the host's network namespace. " +
					"If this option is set, the ports that will be used must be specified.",
				Transform: transform.FromField("Spec.HostNetwork"),
			},
			{
				Name:        "host_pid",
				Type:        proto.ColumnType_BOOL,
				Description: "Use the host's pid namespace.",
				Transform:   transform.FromField("Spec.HostPID"),
			},
			{
				Name:        "host_ipc",
				Type:        proto.ColumnType_BOOL,
				Description: "Use the host's ipc namespace.",
				Transform:   transform.FromField("Spec.HostIPC"),
			},
			{
				Name: "share_process_namespace",
				Type: proto.ColumnType_BOOL,
				Description: "Share a single process namespace between all of the containers in a pod. " +
					"When this is set containers will be able to view and signal processes from other containers " +
					"in the same pod, and the first process in each container will not be assigned PID 1. " +
					"HostPID and ShareProcessNamespace cannot both be set.",
				Transform: transform.FromField("Spec.ShareProcessNamespace"),
			},
			{
				Name:        "security_context",
				Type:        proto.ColumnType_JSON,
				Description: "SecurityContext holds pod-level security attributes and common container settings.",
				Transform:   transform.FromField("Spec.SecurityContext"),
			},

			{
				Name:        "image_pull_secrets",
				Type:        proto.ColumnType_JSON,
				Description: "ImagePullSecrets is an optional list of references to secrets in the same namespace to use for pulling any of the images used by this PodSpec.",
				Transform:   transform.FromField("Spec.ImagePullSecrets"),
			},
			{
				Name:        "hostname",
				Type:        proto.ColumnType_STRING,
				Description: "Specifies the hostname of the Pod. If not specified, the pod's hostname will be set to a system-defined value.",
				Transform:   transform.FromField("Spec.Hostname"),
			},
			{
				Name: "subdomain",
				Type: proto.ColumnType_STRING,
				Description: "If specified, the fully qualified Pod hostname will be '<hostname>.<subdomain>.<pod namespace>.svc.<cluster domain>'. " +
					"If not specified, the pod will not have a domainname at all.",
				Transform: transform.FromField("Spec.Subdomain"),
			},
			{
				Name:        "affinity",
				Type:        proto.ColumnType_JSON,
				Description: "If specified, the pod's scheduling constraints.",
				Transform:   transform.FromField("Spec.Affinity"),
			},
			{
				Name:        "scheduler_name",
				Type:        proto.ColumnType_STRING,
				Description: "If specified, the pod will be dispatched by specified scheduler.",
				Transform:   transform.FromField("Spec.SchedulerName"),
			},
			{
				Name:        "tolerations",
				Type:        proto.ColumnType_JSON,
				Description: "If specified, the pod's tolerations.",
				Transform:   transform.FromField("Spec.Tolerations"),
			},
			{
				Name: "host_aliases",
				Type: proto.ColumnType_JSON,
				Description: "HostAliases is an optional list of hosts and IPs that will be injected into the pod's hosts " +
					"file if specified. This is only valid for non-hostNetwork pods.",
				Transform: transform.FromField("Spec.HostAliases"),
			},
			{
				Name: "priority_class_name",
				Type: proto.ColumnType_STRING,
				Description: "If specified, indicates the pod's priority. 'system-node-critical' and " +
					"'system-cluster-critical' are two special keywords which indicate the " +
					"highest priorities with the former being the highest priority. Any other " +
					"name must be defined by creating a PriorityClass object with that name.",
				Transform: transform.FromField("Spec.PriorityClassName"),
			},
			{
				Name: "priority",
				Type: proto.ColumnType_INT,
				Description: "The priority value. Various system components use this field to find the " +
					"priority of the pod. When Priority Admission Controller is enabled, it " +
					"prevents users from setting this field. The admission controller populates " +
					"this field from PriorityClassName. " +
					"The higher the value, the higher the priority.",
				Transform: transform.FromField("Spec.Priority"),
			},
			{
				Name: "dns_config",
				Type: proto.ColumnType_JSON,
				Description: "Specifies the DNS parameters of a pod. " +
					"Parameters specified here will be merged to the generated DNS " +
					"configuration based on DNSPolicy.",
				Transform: transform.FromField("Spec.DNSConfig"),
			},

			{
				Name: "readiness_gates",
				Type: proto.ColumnType_JSON,
				Description: "If specified, all readiness gates will be evaluated for pod readiness. " +
					"A pod is ready when all its containers are ready AND " +
					"all conditions specified in the readiness gates have status equal to 'True'",
				Transform: transform.FromField("Spec.ReadinessGates"),
			},
			{
				Name: "runtime_class_name",
				Type: proto.ColumnType_STRING,
				Description: "RuntimeClassName refers to a RuntimeClass object in the node.k8s.io group, which should be used " +
					"to run this pod.  If no RuntimeClass resource matches the named class, the pod will not be run. " +
					"If unset or empty, the 'legacy' RuntimeClass will be used, which is an implicit class with an " +
					"empty definition that uses the default runtime handler.",
				Transform: transform.FromField("Spec.RuntimeClassName"),
			},
			{
				Name: "enable_service_links",
				Type: proto.ColumnType_BOOL,
				Description: "EnableServiceLinks indicates whether information about services should be injected into pod's " +
					"environment variables, matching the syntax of Docker links.",
				Transform: transform.FromField("Spec.EnableServiceLinks"),
			},
			{
				Name: "preemption_policy",
				Type: proto.ColumnType_STRING,
				Description: "PreemptionPolicy is the Policy for preempting pods with lower priority. " +
					"One of Never, PreemptLowerPriority.",
				Transform: transform.FromField("Spec.PreemptionPolicy"),
			},
			{
				Name:        "overhead",
				Type:        proto.ColumnType_JSON,
				Description: "Overhead represents the resource overhead associated with running a pod for a given RuntimeClass.",
				Transform:   transform.FromField("Spec.Overhead"),
			},
			{
				Name: "topology_spread_constraints",
				Type: proto.ColumnType_JSON,
				Description: "TopologySpreadConstraints describes how a group of pods ought to spread across topology " +
					"domains. Scheduler will schedule pods in a way which abides by the constraints. " +
					"All topologySpreadConstraints are ANDed.",
				Transform: transform.FromField("Spec.TopologySpreadConstraints"),
			},
			{
				Name: "set_hostname_as_fqdn",
				Type: proto.ColumnType_BOOL,
				Description: "If true the pod's hostname will be configured as the pod's FQDN, rather than the leaf name (the default). " +
					"In Linux containers, this means setting the FQDN in the hostname field of the kernel (the nodename field of struct utsname). " +
					"In Windows containers, this means setting the registry value of hostname for the registry key HKEY_LOCAL_MACHINE\\SYSTEM\\CurrentControlSet\\Services\\Tcpip\\Parameters to FQDN. " +
					"If a pod does not have FQDN, this has no effect.",
				Transform: transform.FromField("Spec.SetHostnameAsFQDN"),
			},

			//// PodStatus Columns
			{
				Name: "phase",
				Type: proto.ColumnType_STRING,
				Description: "The phase of a Pod is a simple, high-level summary of where the Pod is in its lifecycle. " +
					"The conditions array, the reason and message fields, and the individual container status " +
					"arrays contain more detail about the pod's status. There are five possible phase values: " +
					"Pending, Running, Succeeded, Failed, Unknown",
				Transform: transform.FromField("Status.Phase"),
			},
			{
				Name:        "conditions",
				Type:        proto.ColumnType_JSON,
				Description: "Current service state of pod.",
				Transform:   transform.FromField("Status.Conditions"),
			},

			{
				Name:        "status_message",
				Type:        proto.ColumnType_STRING,
				Description: "A human readable message indicating details about why the pod is in this condition.",
				Transform:   transform.FromField("Status.Message"),
			},
			{
				Name:        "status_reason",
				Type:        proto.ColumnType_STRING,
				Description: "A brief CamelCase message indicating details about why the pod is in this state. e.g. 'Evicted'",
				Transform:   transform.FromField("Status.Reason"),
			},
			{
				Name: "nominated_node_name",
				Type: proto.ColumnType_STRING,
				Description: "nominatedNodeName is set only when this pod preempts other pods on the node, but it cannot be " +
					"scheduled right away as preemption victims receive their graceful termination periods. " +
					"This field does not guarantee that the pod will be scheduled on this node. Scheduler may decide " +
					"to place the pod elsewhere if other nodes become available sooner. Scheduler may also decide to " +
					"give the resources on this node to a higher priority pod that is created after preemption. " +
					"As a result, this field may be different than PodSpec.nodeName when the pod is " +
					"scheduled.",
				Transform: transform.FromField("Status.NominatedNodeName"),
			},
			{
				Name:        "host_ip",
				Type:        proto.ColumnType_IPADDR,
				Description: "IP address of the host to which the pod is assigned. Empty if not yet scheduled.",
				Transform:   transform.FromField("Status.HostIP"),
			},

			{
				Name: "pod_ip",
				Type: proto.ColumnType_IPADDR,
				Description: "IP address allocated to the pod. Routable at least within the cluster. " +
					"Empty if not yet allocated.",
				Transform: transform.FromField("Status.PodIP"),
			},
			{
				Name: "pod_ips",
				Type: proto.ColumnType_JSON,
				Description: "podIPs holds the IP addresses allocated to the pod. If this field is specified, the 0th entry must " +
					"match the podIP field. Pods may be allocated at most 1 value for each of IPv4 and IPv6. This list " +
					"is empty if no IPs have been allocated yet.",
				Transform: transform.FromField("Status.PodIPs"),
			},

			{
				Name: "start_time",
				Type: proto.ColumnType_TIMESTAMP,
				Description: "Date and time at which the object was acknowledged by the Kubelet. " +
					"This is before the Kubelet pulled the container image(s) for the pod.",
				Transform: transform.FromField("Status.StartTime").Transform(v1TimeToRFC3339),
			},

			{
				Name: "init_container_statuses",
				Type: proto.ColumnType_JSON,
				Description: "The list has one entry per init container in the manifest. The most recent successful " +
					"init container will have ready = true, the most recently started container will have " +
					"startTime set.",
				Transform: transform.FromField("Status.InitContainerStatuses"),
			},
			{
				Name: "container_statuses",
				Type: proto.ColumnType_JSON,
				Description: "The list has one entry per container in the manifest. Each entry is currently the output " +
					"of `docker inspect`.",
				Transform: transform.FromField("Status.ContainerStatuses"),
			},
			{
				Name:        "qos_class",
				Type:        proto.ColumnType_STRING,
				Description: "The Quality of Service (QOS) classification assigned to the pod based on resource requirements.",
				Transform:   transform.FromField("Status.QOSClass"),
			},
			{
				Name: "ephemeral_container_statuses",
				Type: proto.ColumnType_JSON,
				Description: "Status for any ephemeral containers that have run in this pod. " +
					"This field is alpha-level and is only populated by servers that enable the EphemeralContainers feature.",
				Transform: transform.FromField("Status.EphemeralContainerStatuses"),
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
				Transform:   transform.From(transformPodTags),
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

type Pod struct {
	v1.Pod
	ManifestFilePath string
	StartLine        int
}

//// HYDRATE FUNCTIONS

func listK8sPods(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listK8sPods")

	// Get the client for querying the K8s APIs for the provided context.
	// If the connection is configured for the manifest files, the client will return nil.
	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		return nil, err
	}

	//
	// Check for manifest files
	//
	parsedContents, err := fetchResourceFromManifestFileByKind(ctx, d, "Pod")
	if err != nil {
		return nil, err
	}

	for _, content := range parsedContents {
		pod := content.Data.(*v1.Pod)

		d.StreamListItem(ctx, Pod{*pod, content.Path, content.Line})

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

	if d.EqualsQuals["selector_search"] != nil {
		input.LabelSelector = d.EqualsQuals["selector_search"].GetStringValue()
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

	fieldSelectors := buildKubernetsPodFieldSelectorFilter(ctx, d)
	commonFieldSelectorValue := getCommonOptionalKeyQualsValueForFieldSelector(d)
	fieldSelectors = append(fieldSelectors, commonFieldSelectorValue...)

	if len(fieldSelectors) > 0 {
		input.FieldSelector = strings.Join(fieldSelectors, ",")
	}

	var response *v1.PodList
	pageLeft := true

	for pageLeft {
		response, err = clientset.CoreV1().Pods("").List(ctx, input)
		if err != nil {
			return nil, err
		}

		if response.GetContinue() != "" {
			input.Continue = response.Continue
		} else {
			pageLeft = false
		}

		for _, pod := range response.Items {
			d.StreamListItem(ctx, Pod{pod, "", 0})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

func getK8sPod(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getK8sPod")

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
	parsedContents, err := fetchResourceFromManifestFileByKind(ctx, d, "Pod")
	if err != nil {
		return nil, err
	}

	for _, content := range parsedContents {
		pod := content.Data.(*v1.Pod)

		if pod.Name == name && pod.Namespace == namespace {
			return Pod{*pod, content.Path, content.Line}, nil
		}
	}

	//
	// Get the deployed resource
	//
	if clientset == nil {
		return nil, nil
	}

	pod, err := clientset.CoreV1().Pods(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil && !isNotFoundError(err) {
		return nil, err
	}

	return Pod{*pod, "", 0}, nil
}

//// TRANSFORM FUNCTIONS

func transformPodTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(Pod)
	return mergeTags(obj.Labels, obj.Annotations), nil
}

// // UTILITY FUNCTION
// Build kubernetes pod list call input field selector filter
func buildKubernetsPodFieldSelectorFilter(ctx context.Context, d *plugin.QueryData) []string {

	filterQuals := map[string]string{
		"restart_policy":       "spec.restartPolicy",
		"service_account_name": "spec.serviceAccountName",
		"scheduler_name":       "spec.schedulerName",
		"phase":                "status.phase",
		"nominated_node_name":  "status.nominatedNodeName",
		"pod_ip":               "status.podIP",
	}

	commonFieldSelectorValue := getCommonOptionalKeyQualsValueForFieldSelector(d)

	for columnName, filterName := range filterQuals {
		if columnName == "pod_ip" {
			if d.EqualsQuals["pod_ip"] != nil {
				value := d.EqualsQuals["pod_ip"].GetInetValue().GetAddr()
				commonFieldSelectorValue = append(commonFieldSelectorValue, filterName+"="+value)
			}
			continue
		}
		if d.EqualsQualString(columnName) != "" {
			commonFieldSelectorValue = append(commonFieldSelectorValue, filterName+"="+d.EqualsQualString(columnName))
		}
	}

	return commonFieldSelectorValue
}
