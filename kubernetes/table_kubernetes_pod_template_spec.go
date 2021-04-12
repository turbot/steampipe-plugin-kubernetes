package kubernetes

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableKubernetesPodTemplateSpec(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name: "kubernetes_pod_template_spec",
		// Description: "Role contains rules that represent a set of permissions.",
		// DefaultTransform: transform.FromJSONTag(),
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("template"),
			Hydrate:    listPodTemplateSpec,
		},
		Columns: []*plugin.Column{
			{
				Name:        "template",
				Type:        proto.ColumnType_STRING,
				Description: "List of the PolicyRules for this Role.",
				Transform:   transform.FromField("template"),
			},

			//// PodSpec Columns
			{
				Name:        "volumes",
				Type:        proto.ColumnType_JSON,
				Description: "List of volumes that can be mounted by containers belonging to the pod.",
				Transform:   transform.FromField("spec.volumes"),
			},
			{
				Name:        "containers",
				Type:        proto.ColumnType_JSON,
				Description: "List of containers belonging to the pod.",
				Transform:   transform.FromField("spec.containers"),
			},
			{
				Name: "ephemeral_containers",
				Type: proto.ColumnType_JSON,
				Description: "List of ephemeral containers run in this pod. Ephemeral containers may be run in an existing " +
					"pod to perform user-initiated actions such as debugging. This list cannot be specified when " +
					"creating a pod, and it cannot be modified by updating the pod spec. In order to add an " +
					"ephemeral container to an existing pod, use the pod's ephemeralcontainers subresource. " +
					"This field is alpha-level and is only honored by servers that enable the EphemeralContainers feature.",
				Transform: transform.FromField("spec.ephemeralContainers"),
			},
			{
				Name: "init_containers",
				Type: proto.ColumnType_JSON,
				Description: "List of initialization containers belonging to the pod. Init containers " +
					"are executed in order prior to containers being started. If any " +
					"init container fails, the pod is considered to have failed and is handled according " +
					"to its restartPolicy. The name for an init container or normal container must be " +
					"unique among all containers.",
				Transform: transform.FromField("spec.initContainers"),
			},
			{
				Name:        "restart_policy",
				Type:        proto.ColumnType_STRING,
				Description: "Restart policy for all containers within the pod. One of Always, OnFailure, Never.",
				Transform:   transform.FromField("spec.restartPolicy"),
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
				Transform: transform.FromField("spec.terminationGracePeriodSeconds"),
			},
			{
				Name: "active_deadline_seconds",
				Type: proto.ColumnType_STRING,
				Description: "Optional duration in seconds the pod may be active on the node relative to " +
					"StartTime before the system will actively try to mark it failed and kill associated containers.",
				Transform: transform.FromField("spec.activeDeadlineSeconds"),
			},
			{
				Name:        "dns_policy",
				Type:        proto.ColumnType_STRING,
				Description: "DNS policy for pod.  Valid values are 'ClusterFirstWithHostNet', 'ClusterFirst', 'Default' or 'None'.",
				Transform:   transform.FromField("spec.dnsPolicy"),
			},
			{
				Name:        "node_selector",
				Type:        proto.ColumnType_JSON,
				Description: "NodeSelector is a selector which must be true for the pod to fit on a node.",
				Transform:   transform.FromField("spec.nodeSelector"),
			},
			{
				Name:        "service_account_name",
				Type:        proto.ColumnType_STRING,
				Description: "ServiceAccountName is the name of the ServiceAccount to use to run this pod.",
				Transform:   transform.FromField("spec.serviceAccountName"),
			},
			// // Dont include deprecated columns...
			// {
			// 	Name: "deprecated_service_account",
			// 	Type: proto.ColumnType_STRING,
			// 	Description: "DeprecatedServiceAccount is a depreciated alias for ServiceAccountName. " +
			// 		"Deprecated: Use serviceAccountName instead.",
			// 	Transform: transform.FromField("DeprecatedServiceAccount"),
			// },
			{
				Name:        "automount_service_account_token",
				Type:        proto.ColumnType_BOOL,
				Description: "AutomountServiceAccountToken indicates whether a service account token should be automatically mounted.",
				Transform:   transform.FromField("spec.automountServiceAccountToken"),
			},
			{
				Name: "node_name",
				Type: proto.ColumnType_STRING,
				Description: "NodeName is a request to schedule this pod onto a specific node. If it is non-empty, " +
					"the scheduler simply schedules this pod onto that node, assuming that it fits resource " +
					"requirements.",
				Transform: transform.FromField("spec.nodeName"),
			},
			{
				Name: "host_network",
				Type: proto.ColumnType_BOOL,
				Description: "Host networking requested for this pod. Use the host's network namespace. " +
					"If this option is set, the ports that will be used must be specified.",
				Transform: transform.FromField("spec.hostNetwork"),
			},
			{
				Name:        "host_pid",
				Type:        proto.ColumnType_BOOL,
				Description: "Use the host's pid namespace.",
				Transform:   transform.FromField("spec.hostPID"),
			},
			{
				Name:        "host_ipc",
				Type:        proto.ColumnType_BOOL,
				Description: "Use the host's ipc namespace.",
				Transform:   transform.FromField("spec.hostIPC"),
			},
			{
				Name: "share_process_namespace",
				Type: proto.ColumnType_BOOL,
				Description: "Share a single process namespace between all of the containers in a pod. " +
					"When this is set containers will be able to view and signal processes from other containers " +
					"in the same pod, and the first process in each container will not be assigned PID 1. " +
					"HostPID and ShareProcessNamespace cannot both be set.",
				Transform: transform.FromField("spec.shareProcessNamespace"),
			},
			{
				Name:        "security_context",
				Type:        proto.ColumnType_JSON,
				Description: "SecurityContext holds pod-level security attributes and common container settings.",
				Transform:   transform.FromField("spec.securityContext"),
			},

			{
				Name:        "image_pull_secrets",
				Type:        proto.ColumnType_JSON,
				Description: "ImagePullSecrets is an optional list of references to secrets in the same namespace to use for pulling any of the images used by this PodSpec.",
				Transform:   transform.FromField("spec.imagePullSecrets"),
			},
			{
				Name:        "hostname",
				Type:        proto.ColumnType_STRING,
				Description: "Specifies the hostname of the Pod. If not specified, the pod's hostname will be set to a system-defined value.",
				Transform:   transform.FromField("spec.hostname"),
			},
			{
				Name: "subdomain",
				Type: proto.ColumnType_STRING,
				Description: "If specified, the fully qualified Pod hostname will be '<hostname>.<subdomain>.<pod namespace>.svc.<cluster domain>'. " +
					"If not specified, the pod will not have a domainname at all.",
				Transform: transform.FromField("spec.subdomain"),
			},
			{
				Name:        "affinity",
				Type:        proto.ColumnType_JSON,
				Description: "If specified, the pod's scheduling constraints.",
				Transform:   transform.FromField("spec.affinity"),
			},
			{
				Name:        "scheduler_name",
				Type:        proto.ColumnType_STRING,
				Description: "If specified, the pod will be dispatched by specified scheduler.",
				Transform:   transform.FromField("spec.schedulerName"),
			},
			{
				Name:        "tolerations",
				Type:        proto.ColumnType_JSON,
				Description: "If specified, the pod's tolerations.",
				Transform:   transform.FromField("spec.tolerations"),
			},
			{
				Name: "host_aliases",
				Type: proto.ColumnType_JSON,
				Description: "HostAliases is an optional list of hosts and IPs that will be injected into the pod's hosts " +
					"file if specified. This is only valid for non-hostNetwork pods.",
				Transform: transform.FromField("spec.hostAliases"),
			},
			{
				Name: "priority_class_name",
				Type: proto.ColumnType_STRING,
				Description: "If specified, indicates the pod's priority. 'system-node-critical' and " +
					"'system-cluster-critical' are two special keywords which indicate the " +
					"highest priorities with the former being the highest priority. Any other " +
					"name must be defined by creating a PriorityClass object with that name.",
				Transform: transform.FromField("spec.priorityClassName"),
			},
			{
				Name: "priority",
				Type: proto.ColumnType_INT,
				Description: "The priority value. Various system components use this field to find the " +
					"priority of the pod. When Priority Admission Controller is enabled, it " +
					"prevents users from setting this field. The admission controller populates " +
					"this field from PriorityClassName. " +
					"The higher the value, the higher the priority.",
				Transform: transform.FromField("spec.priority"),
			},
			{
				Name: "dns_config",
				Type: proto.ColumnType_JSON,
				Description: "Specifies the DNS parameters of a pod. " +
					"Parameters specified here will be merged to the generated DNS " +
					"configuration based on DNSPolicy.",
				Transform: transform.FromField("spec.dnsConfig"),
			},

			{
				Name: "readiness_gates",
				Type: proto.ColumnType_JSON,
				Description: "If specified, all readiness gates will be evaluated for pod readiness. " +
					"A pod is ready when all its containers are ready AND " +
					"all conditions specified in the readiness gates have status equal to 'True'",
				Transform: transform.FromField("spec.readinessGates"),
			},
			{
				Name: "runtime_class_name",
				Type: proto.ColumnType_STRING,
				Description: "RuntimeClassName refers to a RuntimeClass object in the node.k8s.io group, which should be used " +
					"to run this pod.  If no RuntimeClass resource matches the named class, the pod will not be run. " +
					"If unset or empty, the 'legacy' RuntimeClass will be used, which is an implicit class with an " +
					"empty definition that uses the default runtime handler.",
				Transform: transform.FromField("spec.runtimeClassName"),
			},
			{
				Name: "enable_service_links",
				Type: proto.ColumnType_BOOL,
				Description: "EnableServiceLinks indicates whether information about services should be injected into pod's " +
					"environment variables, matching the syntax of Docker links.",
				Transform: transform.FromField("spec.enableServiceLinks"),
			},
			{
				Name: "preemption_policy",
				Type: proto.ColumnType_STRING,
				Description: "PreemptionPolicy is the Policy for preempting pods with lower priority. " +
					"One of Never, PreemptLowerPriority.",
				Transform: transform.FromField("spec.preemptionPolicy"),
			},
			{
				Name:        "overhead",
				Type:        proto.ColumnType_JSON,
				Description: "Overhead represents the resource overhead associated with running a pod for a given RuntimeClass.",
				Transform:   transform.FromField("spec.overhead"),
			},
			{
				Name: "topology_spread_constraints",
				Type: proto.ColumnType_JSON,
				Description: "TopologySpreadConstraints describes how a group of pods ought to spread across topology " +
					"domains. Scheduler will schedule pods in a way which abides by the constraints. " +
					"All topologySpreadConstraints are ANDed.",
				Transform: transform.FromField("spec.topologySpreadConstraints"),
			},
			{
				Name: "set_hostname_as_fqdn",
				Type: proto.ColumnType_BOOL,
				Description: "If true the pod's hostname will be configured as the pod's FQDN, rather than the leaf name (the default). " +
					"In Linux containers, this means setting the FQDN in the hostname field of the kernel (the nodename field of struct utsname). " +
					"In Windows containers, this means setting the registry value of hostname for the registry key HKEY_LOCAL_MACHINE\\SYSTEM\\CurrentControlSet\\Services\\Tcpip\\Parameters to FQDN. " +
					"If a pod does not have FQDN, this has no effect.",
				Transform: transform.FromField("spec.setHostnameAsFQDN"),
			},
		},
	}
}

//// HYDRATE FUNCTIONS

func listPodTemplateSpec(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Warn("listPodTemplate")

	template := d.KeyColumnQuals["template"].GetStringValue()
	logger.Warn("listPodTemplate", "template", template)

	if template == "" {
		return nil, nil
	}

	var rawPolicyMap map[string]interface{}
	err := json.Unmarshal([]byte(template), &rawPolicyMap)
	if err != nil {
		fmt.Println("error:", err)
		return nil, err
	}

	rawPolicyMap["template"] = template

	logger.Warn("listPodTemplate", "rawPolicyMap", rawPolicyMap)
	d.StreamListItem(ctx, rawPolicyMap)

	return nil, nil
}
