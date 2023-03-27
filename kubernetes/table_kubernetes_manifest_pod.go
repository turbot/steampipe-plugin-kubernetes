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

func tableKubernetesManifestPod(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_manifest_pod",
		Description: "",
		List: &plugin.ListConfig{
			ParentHydrate: listKubernetesManifestFiles,
			Hydrate:       listKubernetesManifestPods,
			KeyColumns:    plugin.OptionalColumns([]string{"manifest_file_path"}),
		},
		Columns: []*plugin.Column{
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the object.  Name must be unique within a namespace."},
			{Name: "namespace", Type: proto.ColumnType_STRING, Description: "Namespace defines the space within which each name must be unique."},
			{Name: "labels", Type: proto.ColumnType_JSON, Description: "Map of string keys and values that can be used to organize and categorize (scope and select) objects. May match selectors of replication controllers and services."},
			{Name: "annotations", Type: proto.ColumnType_JSON, Description: "Annotations is an unstructured key value map stored with a resource that may be set by external tools to store and retrieve arbitrary metadata."},

			// PodSpec Columns
			{Name: "selector_search", Type: proto.ColumnType_STRING, Description: "A label selector string to restrict the list of returned objects by their labels.", Transform: transform.FromQual("selector_search")},
			{Name: "volumes", Type: proto.ColumnType_JSON, Description: "List of volumes that can be mounted by containers belonging to the pod.", Transform: transform.FromField("Spec.Volumes")},
			{Name: "containers", Type: proto.ColumnType_JSON, Description: "List of containers belonging to the pod.", Transform: transform.FromField("Spec.Containers")},
			{Name: "ephemeral_containers", Type: proto.ColumnType_JSON, Description: "List of ephemeral containers run in this pod. Ephemeral containers may be run in an existing pod to perform user-initiated actions such as debugging. This list cannot be specified when creating a pod, and it cannot be modified by updating the pod spec. In order to add an ephemeral container to an existing pod, use the pod's ephemeralcontainers subresource. This field is alpha-level and is only honored by servers that enable the EphemeralContainers feature.", Transform: transform.FromField("Spec.EphemeralContainers")},
			{Name: "init_containers", Type: proto.ColumnType_JSON, Description: "List of initialization containers belonging to the pod. Init containers are executed in order prior to containers being started. If any init container fails, the pod is considered to have failed and is handled according to its restartPolicy. The name for an init container or normal container must be unique among all containers.", Transform: transform.FromField("Spec.InitContainers")},
			{Name: "restart_policy", Type: proto.ColumnType_STRING, Description: "Restart policy for all containers within the pod. One of Always, OnFailure, Never.", Transform: transform.FromField("Spec.RestartPolicy")},
			{Name: "termination_grace_period_seconds", Type: proto.ColumnType_INT, Description: "Optional duration in seconds the pod needs to terminate gracefully. May be decreased in delete request. Value must be non-negative integer. The value zero indicates delete immediately. If this value is nil, the default grace period will be used instead. The grace period is the duration in seconds after the processes running in the pod are sent a termination signal and the time when the processes are forcibly halted with a kill signal. Set this value longer than the expected cleanup time for your process.", Transform: transform.FromField("Spec.TerminationGracePeriodSeconds")},
			{Name: "active_deadline_seconds", Type: proto.ColumnType_STRING, Description: "Optional duration in seconds the pod may be active on the node relative to StartTime before the system will actively try to mark it failed and kill associated containers.", Transform: transform.FromField("Spec.ActiveDeadlineSeconds")},
			{Name: "dns_policy", Type: proto.ColumnType_STRING, Description: "DNS policy for pod.  Valid values are 'ClusterFirstWithHostNet', 'ClusterFirst', 'Default' or 'None'.", Transform: transform.FromField("Spec.DNSPolicy")},
			{Name: "node_selector", Type: proto.ColumnType_JSON, Description: "NodeSelector is a selector which must be true for the pod to fit on a node.", Transform: transform.FromField("Spec.NodeSelector")},
			{Name: "service_account_name", Type: proto.ColumnType_STRING, Description: "ServiceAccountName is the name of the ServiceAccount to use to run this pod.", Transform: transform.FromField("Spec.ServiceAccountName")},
			{Name: "automount_service_account_token", Type: proto.ColumnType_BOOL, Description: "AutomountServiceAccountToken indicates whether a service account token should be automatically mounted.", Transform: transform.FromField("Spec.AutomountServiceAccountToken")},
			{Name: "node_name", Type: proto.ColumnType_STRING, Description: "NodeName is a request to schedule this pod onto a specific node. If it is non-empty, the scheduler simply schedules this pod onto that node, assuming that it fits resource requirements.", Transform: transform.FromField("Spec.NodeName")},
			{Name: "host_network", Type: proto.ColumnType_BOOL, Description: "Host networking requested for this pod. Use the host's network namespace. If this option is set, the ports that will be used must be specified.", Transform: transform.FromField("Spec.HostNetwork")},
			{Name: "host_pid", Type: proto.ColumnType_BOOL, Description: "Use the host's pid namespace.", Transform: transform.FromField("Spec.HostPID")},
			{Name: "host_ipc", Type: proto.ColumnType_BOOL, Description: "Use the host's ipc namespace.", Transform: transform.FromField("Spec.HostIPC")},
			{Name: "share_process_namespace", Type: proto.ColumnType_BOOL, Description: "Share a single process namespace between all of the containers in a pod. When this is set containers will be able to view and signal processes from other containers in the same pod, and the first process in each container will not be assigned PID 1. HostPID and ShareProcessNamespace cannot both be set.", Transform: transform.FromField("Spec.ShareProcessNamespace")},
			{Name: "security_context", Type: proto.ColumnType_JSON, Description: "SecurityContext holds pod-level security attributes and common container settings.", Transform: transform.FromField("Spec.SecurityContext")},
			{Name: "image_pull_secrets", Type: proto.ColumnType_JSON, Description: "ImagePullSecrets is an optional list of references to secrets in the same namespace to use for pulling any of the images used by this PodSpec.", Transform: transform.FromField("Spec.ImagePullSecrets")},
			{Name: "hostname", Type: proto.ColumnType_STRING, Description: "Specifies the hostname of the Pod. If not specified, the pod's hostname will be set to a system-defined value.", Transform: transform.FromField("Spec.Hostname")},
			{Name: "subdomain", Type: proto.ColumnType_STRING, Description: "If specified, the fully qualified Pod hostname will be '<hostname>.<subdomain>.<pod namespace>.svc.<cluster domain>'. If not specified, the pod will not have a domainname at all.", Transform: transform.FromField("Spec.Subdomain")},
			{Name: "affinity", Type: proto.ColumnType_JSON, Description: "If specified, the pod's scheduling constraints.", Transform: transform.FromField("Spec.Affinity")},
			{Name: "scheduler_name", Type: proto.ColumnType_STRING, Description: "If specified, the pod will be dispatched by specified scheduler.", Transform: transform.FromField("Spec.SchedulerName")},
			{Name: "tolerations", Type: proto.ColumnType_JSON, Description: "If specified, the pod's tolerations.", Transform: transform.FromField("Spec.Tolerations")},
			{Name: "host_aliases", Type: proto.ColumnType_JSON, Description: "HostAliases is an optional list of hosts and IPs that will be injected into the pod's hosts file if specified. This is only valid for non-hostNetwork pods.", Transform: transform.FromField("Spec.HostAliases")},
			{Name: "priority_class_name", Type: proto.ColumnType_STRING, Description: "If specified, indicates the pod's priority. 'system-node-critical' and 'system-cluster-critical' are two special keywords which indicate the highest priorities with the former being the highest priority. Any other name must be defined by creating a PriorityClass object with that name.", Transform: transform.FromField("Spec.PriorityClassName")},
			{Name: "priority", Type: proto.ColumnType_INT, Description: "The priority value. Various system components use this field to find the priority of the pod. When Priority Admission Controller is enabled, it prevents users from setting this field. The admission controller populates this field from PriorityClassName. The higher the value, the higher the priority.", Transform: transform.FromField("Spec.Priority")},
			{Name: "dns_config", Type: proto.ColumnType_JSON, Description: "Specifies the DNS parameters of a pod. Parameters specified here will be merged to the generated DNS configuration based on DNSPolicy.", Transform: transform.FromField("Spec.DNSConfig")},
			{Name: "readiness_gates", Type: proto.ColumnType_JSON, Description: "If specified, all readiness gates will be evaluated for pod readiness. A pod is ready when all its containers are ready AND all conditions specified in the readiness gates have status equal to 'True'", Transform: transform.FromField("Spec.ReadinessGates")},
			{Name: "runtime_class_name", Type: proto.ColumnType_STRING, Description: "RuntimeClassName refers to a RuntimeClass object in the node.k8s.io group, which should be used to run this pod.  If no RuntimeClass resource matches the named class, the pod will not be run. If unset or empty, the 'legacy' RuntimeClass will be used, which is an implicit class with an empty definition that uses the default runtime handler.", Transform: transform.FromField("Spec.RuntimeClassName")},
			{Name: "enable_service_links", Type: proto.ColumnType_BOOL, Description: "EnableServiceLinks indicates whether information about services should be injected into pod's environment variables, matching the syntax of Docker links.", Transform: transform.FromField("Spec.EnableServiceLinks")},
			{Name: "preemption_policy", Type: proto.ColumnType_STRING, Description: "PreemptionPolicy is the Policy for preempting pods with lower priority. One of Never, PreemptLowerPriority.", Transform: transform.FromField("Spec.PreemptionPolicy")},
			{Name: "overhead", Type: proto.ColumnType_JSON, Description: "Overhead represents the resource overhead associated with running a pod for a given RuntimeClass.", Transform: transform.FromField("Spec.Overhead")},
			{Name: "topology_spread_constraints", Type: proto.ColumnType_JSON, Description: "TopologySpreadConstraints describes how a group of pods ought to spread across topology domains. Scheduler will schedule pods in a way which abides by the constraints. All topologySpreadConstraints are ANDed.", Transform: transform.FromField("Spec.TopologySpreadConstraints")},
			{Name: "set_hostname_as_fqdn", Type: proto.ColumnType_BOOL, Description: "If true the pod's hostname will be configured as the pod's FQDN, rather than the leaf name (the default). In Linux containers, this means setting the FQDN in the hostname field of the kernel (the nodename field of struct utsname). In Windows containers, this means setting the registry value of hostname for the registry key HKEY_LOCAL_MACHINE\\SYSTEM\\CurrentControlSet\\Services\\Tcpip\\Parameters to FQDN. If a pod does not have FQDN, this has no effect.", Transform: transform.FromField("Spec.SetHostnameAsFQDN")},

			{Name: "manifest_file_path", Description: "Path to the file.", Type: proto.ColumnType_STRING},
		},
	}
}

type KubernetesManifestPod struct {
	v1.Pod
	ManifestFilePath string
}

//// LIST FUNCTION

func listKubernetesManifestPods(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// The path comes from a parent hydrate, defaulting to the config paths or
	// available by the optional key column
	path := h.Item.(filePath).Path

	// Load the file into a buffer
	content, err := os.ReadFile(path)
	if err != nil {
		plugin.Logger(ctx).Error("kubernetes_manifest_pod.listKubernetesManifestPods", "failed to read file", err, "path", path)
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
			plugin.Logger(ctx).Error("kubernetes_manifest_pod.listKubernetesManifestPods", "failed to decode the file", err, "path", path)
			return nil, err
		}

		// Return if the definition is not for the pod resource
		if groupVersionKind.Kind == "Pod" {
			pod := obj.(*v1.Pod)

			d.StreamListItem(ctx, KubernetesManifestPod{*pod, path})

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}
