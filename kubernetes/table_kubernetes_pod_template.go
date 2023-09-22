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

func tableKubernetesPodTemplate(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_pod_template",
		Description: "Kubernetes Pod Template is a collection of templates for creating copies of a predefined pod.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "namespace"}),
			Hydrate:    getK8sPodTemplate,
		},
		List: &plugin.ListConfig{
			Hydrate: listK8sPodTemplates,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "name", Require: plugin.Optional},
				{Name: "namespace", Require: plugin.Optional},
			},
		},
		Columns: k8sCommonColumns([]*plugin.Column{
			//// PodSpec Columns
			{
				Name:        "volumes",
				Type:        proto.ColumnType_JSON,
				Description: "List of volumes that can be mounted by containers belonging to the pod.",
				Transform:   transform.FromField("Template.Template.Spec.Volumes"),
			},
			{
				Name:        "containers",
				Type:        proto.ColumnType_JSON,
				Description: "List of containers belonging to the pod.",
				Transform:   transform.FromField("Template.Spec.Containers"),
			},
			{
				Name: "ephemeral_containers",
				Type: proto.ColumnType_JSON,
				Description: "List of ephemeral containers run in this pod. Ephemeral containers may be run in an existing " +
					"pod to perform user-initiated actions such as debugging. This list cannot be specified when " +
					"creating a pod, and it cannot be modified by updating the pod spec. In order to add an " +
					"ephemeral container to an existing pod, use the pod's ephemeralcontainers subresource. " +
					"This field is alpha-level and is only honored by servers that enable the EphemeralContainers feature.",
				Transform: transform.FromField("Template.Spec.EphemeralContainers"),
			},
			{
				Name: "init_containers",
				Type: proto.ColumnType_JSON,
				Description: "List of initialization containers belonging to the pod. Init containers " +
					"are executed in order prior to containers being started. If any " +
					"init container fails, the pod is considered to have failed and is handled according " +
					"to its restartPolicy. The name for an init container or normal container must be " +
					"unique among all containers.",
				Transform: transform.FromField("Template.Spec.InitContainers"),
			},
			{
				Name:        "restart_policy",
				Type:        proto.ColumnType_STRING,
				Description: "Restart policy for all containers within the pod. One of Always, OnFailure, Never.",
				Transform:   transform.FromField("Template.Spec.RestartPolicy"),
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
				Transform: transform.FromField("Template.Spec.TerminationGracePeriodSeconds"),
			},
			{
				Name: "active_deadline_seconds",
				Type: proto.ColumnType_STRING,
				Description: "Optional duration in seconds the pod may be active on the node relative to " +
					"StartTime before the system will actively try to mark it failed and kill associated containers.",
				Transform: transform.FromField("Template.Spec.ActiveDeadlineSeconds"),
			},
			{
				Name:        "dns_policy",
				Type:        proto.ColumnType_STRING,
				Description: "DNS policy for pod.  Valid values are 'ClusterFirstWithHostNet', 'ClusterFirst', 'Default' or 'None'.",
				Transform:   transform.FromField("Template.Spec.DNSPolicy"),
			},
			{
				Name:        "node_selector",
				Type:        proto.ColumnType_JSON,
				Description: "NodeSelector is a selector which must be true for the pod to fit on a node.",
				Transform:   transform.FromField("Template.Spec.NodeSelector"),
			},
			{
				Name:        "service_account_name",
				Type:        proto.ColumnType_STRING,
				Description: "ServiceAccountName is the name of the ServiceAccount to use to run this pod.",
				Transform:   transform.FromField("Template.Spec.ServiceAccountName"),
			},
			{
				Name:        "automount_service_account_token",
				Type:        proto.ColumnType_BOOL,
				Description: "AutomountServiceAccountToken indicates whether a service account token should be automatically mounted.",
				Transform:   transform.FromField("Template.Spec.AutomountServiceAccountToken"),
			},
			{
				Name: "node_name",
				Type: proto.ColumnType_STRING,
				Description: "NodeName is a request to schedule this pod onto a specific node. If it is non-empty, " +
					"the scheduler simply schedules this pod onto that node, assuming that it fits resource " +
					"requirements.",
				Transform: transform.FromField("Template.Spec.NodeName"),
			},
			{
				Name: "host_network",
				Type: proto.ColumnType_BOOL,
				Description: "Host networking requested for this pod. Use the host's network namespace. " +
					"If this option is set, the ports that will be used must be specified.",
				Transform: transform.FromField("Template.Spec.HostNetwork"),
			},
			{
				Name:        "host_pid",
				Type:        proto.ColumnType_BOOL,
				Description: "Use the host's pid namespace.",
				Transform:   transform.FromField("Template.Spec.HostPID"),
			},
			{
				Name:        "host_ipc",
				Type:        proto.ColumnType_BOOL,
				Description: "Use the host's ipc namespace.",
				Transform:   transform.FromField("Template.Spec.HostIPC"),
			},
			{
				Name: "share_process_namespace",
				Type: proto.ColumnType_BOOL,
				Description: "Share a single process namespace between all of the containers in a pod. " +
					"When this is set containers will be able to view and signal processes from other containers " +
					"in the same pod, and the first process in each container will not be assigned PID 1. " +
					"HostPID and ShareProcessNamespace cannot both be set.",
				Transform: transform.FromField("Template.Spec.ShareProcessNamespace"),
			},
			{
				Name:        "security_context",
				Type:        proto.ColumnType_JSON,
				Description: "SecurityContext holds pod-level security attributes and common container settings.",
				Transform:   transform.FromField("Template.Spec.SecurityContext"),
			},

			{
				Name:        "image_pull_secrets",
				Type:        proto.ColumnType_JSON,
				Description: "ImagePullSecrets is an optional list of references to secrets in the same namespace to use for pulling any of the images used by this PodSpec.",
				Transform:   transform.FromField("Template.Spec.ImagePullSecrets"),
			},
			{
				Name:        "hostname",
				Type:        proto.ColumnType_STRING,
				Description: "Specifies the hostname of the Pod. If not specified, the pod's hostname will be set to a system-defined value.",
				Transform:   transform.FromField("Template.Spec.Hostname"),
			},
			{
				Name: "subdomain",
				Type: proto.ColumnType_STRING,
				Description: "If specified, the fully qualified Pod hostname will be '<hostname>.<subdomain>.<pod namespace>.svc.<cluster domain>'. " +
					"If not specified, the pod will not have a domainname at all.",
				Transform: transform.FromField("Template.Spec.Subdomain"),
			},
			{
				Name:        "affinity",
				Type:        proto.ColumnType_JSON,
				Description: "If specified, the pod's scheduling constraints.",
				Transform:   transform.FromField("Template.Spec.Affinity"),
			},
			{
				Name:        "scheduler_name",
				Type:        proto.ColumnType_STRING,
				Description: "If specified, the pod will be dispatched by specified scheduler.",
				Transform:   transform.FromField("Template.Spec.SchedulerName"),
			},
			{
				Name:        "tolerations",
				Type:        proto.ColumnType_JSON,
				Description: "If specified, the pod's tolerations.",
				Transform:   transform.FromField("Template.Spec.Tolerations"),
			},
			{
				Name: "host_aliases",
				Type: proto.ColumnType_JSON,
				Description: "HostAliases is an optional list of hosts and IPs that will be injected into the pod's hosts " +
					"file if specified. This is only valid for non-hostNetwork pods.",
				Transform: transform.FromField("Template.Spec.HostAliases"),
			},
			{
				Name: "priority_class_name",
				Type: proto.ColumnType_STRING,
				Description: "If specified, indicates the pod's priority. 'system-node-critical' and " +
					"'system-cluster-critical' are two special keywords which indicate the " +
					"highest priorities with the former being the highest priority. Any other " +
					"name must be defined by creating a PriorityClass object with that name.",
				Transform: transform.FromField("Template.Spec.PriorityClassName"),
			},
			{
				Name: "priority",
				Type: proto.ColumnType_INT,
				Description: "The priority value. Various system components use this field to find the " +
					"priority of the pod. When Priority Admission Controller is enabled, it " +
					"prevents users from setting this field. The admission controller populates " +
					"this field from PriorityClassName. " +
					"The higher the value, the higher the priority.",
				Transform: transform.FromField("Template.Spec.Priority"),
			},
			{
				Name: "dns_config",
				Type: proto.ColumnType_JSON,
				Description: "Specifies the DNS parameters of a pod. " +
					"Parameters specified here will be merged to the generated DNS " +
					"configuration based on DNSPolicy.",
				Transform: transform.FromField("Template.Spec.DNSConfig"),
			},

			{
				Name: "readiness_gates",
				Type: proto.ColumnType_JSON,
				Description: "If specified, all readiness gates will be evaluated for pod readiness. " +
					"A pod is ready when all its containers are ready AND " +
					"all conditions specified in the readiness gates have status equal to 'True'",
				Transform: transform.FromField("Template.Spec.ReadinessGates"),
			},
			{
				Name: "runtime_class_name",
				Type: proto.ColumnType_STRING,
				Description: "RuntimeClassName refers to a RuntimeClass object in the node.k8s.io group, which should be used " +
					"to run this pod.  If no RuntimeClass resource matches the named class, the pod will not be run. " +
					"If unset or empty, the 'legacy' RuntimeClass will be used, which is an implicit class with an " +
					"empty definition that uses the default runtime handler.",
				Transform: transform.FromField("Template.Spec.RuntimeClassName"),
			},
			{
				Name: "enable_service_links",
				Type: proto.ColumnType_BOOL,
				Description: "EnableServiceLinks indicates whether information about services should be injected into pod's " +
					"environment variables, matching the syntax of Docker links.",
				Transform: transform.FromField("Template.Spec.EnableServiceLinks"),
			},
			{
				Name: "preemption_policy",
				Type: proto.ColumnType_STRING,
				Description: "PreemptionPolicy is the Policy for preempting pods with lower priority. " +
					"One of Never, PreemptLowerPriority.",
				Transform: transform.FromField("Template.Spec.PreemptionPolicy"),
			},
			{
				Name:        "overhead",
				Type:        proto.ColumnType_JSON,
				Description: "Overhead represents the resource overhead associated with running a pod for a given RuntimeClass.",
				Transform:   transform.FromField("Template.Spec.Overhead"),
			},
			{
				Name: "topology_spread_constraints",
				Type: proto.ColumnType_JSON,
				Description: "TopologySpreadConstraints describes how a group of pods ought to spread across topology " +
					"domains. Scheduler will schedule pods in a way which abides by the constraints. " +
					"All topologySpreadConstraints are ANDed.",
				Transform: transform.FromField("Template.Spec.TopologySpreadConstraints"),
			},
			{
				Name: "set_hostname_as_fqdn",
				Type: proto.ColumnType_BOOL,
				Description: "If true the pod's hostname will be configured as the pod's FQDN, rather than the leaf name (the default). " +
					"In Linux containers, this means setting the FQDN in the hostname field of the kernel (the nodename field of struct utsname). " +
					"In Windows containers, this means setting the registry value of hostname for the registry key HKEY_LOCAL_MACHINE\\SYSTEM\\CurrentControlSet\\Services\\Tcpip\\Parameters to FQDN. " +
					"If a pod does not have FQDN, this has no effect.",
				Transform: transform.FromField("Template.Spec.SetHostnameAsFQDN"),
			},

			{
				Name:        "context_name",
				Type:        proto.ColumnType_STRING,
				Description: "Kubectl config context name.",
				Hydrate:     getPodTemplateResourceContext,
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
				Transform:   transform.From(transformPodTemplateTags),
			},
		}),
	}
}

type PodTemplate struct {
	v1.PodTemplate
	parsedContent
}

//// HYDRATE FUNCTIONS

func listK8sPodTemplates(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listK8sPodTemplates")

	// Get the client for querying the K8s APIs for the provided context.
	// If the connection is configured for the manifest files, the client will return nil.
	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		return nil, err
	}

	// Check for manifest files
	parsedContents, err := fetchResourceFromManifestFileByKind(ctx, d, "PodTemplate")
	if err != nil {
		return nil, err
	}

	for _, content := range parsedContents {
		podTemplate := content.ParsedData.(*v1.PodTemplate)

		d.StreamListItem(ctx, PodTemplate{*podTemplate, content})

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

	var response *v1.PodTemplateList
	pageLeft := true

	for pageLeft {
		response, err = clientset.CoreV1().PodTemplates("").List(ctx, input)
		if err != nil {
			return nil, err
		}

		if response.GetContinue() != "" {
			input.Continue = response.Continue
		} else {
			pageLeft = false
		}

		for _, podTemplate := range response.Items {
			d.StreamListItem(ctx, PodTemplate{podTemplate, parsedContent{SourceType: "deployed"}})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

func getK8sPodTemplate(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getK8sPodTemplate")

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
	parsedContents, err := fetchResourceFromManifestFileByKind(ctx, d, "PodTemplate")
	if err != nil {
		return nil, err
	}

	for _, content := range parsedContents {
		podTemplate := content.ParsedData.(*v1.PodTemplate)

		if podTemplate.Name == name && podTemplate.Namespace == namespace {
			return PodTemplate{*podTemplate, content}, nil
		}
	}

	// Get the deployed resource
	if clientset == nil {
		return nil, nil
	}

	podTemplate, err := clientset.CoreV1().PodTemplates(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil && !isNotFoundError(err) {
		return nil, err
	}

	return PodTemplate{*podTemplate, parsedContent{SourceType: "deployed"}}, nil
}

func getPodTemplateResourceContext(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	obj := h.Item.(PodTemplate)

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

func transformPodTemplateTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(PodTemplate)
	return mergeTags(obj.Labels, obj.Annotations), nil
}
