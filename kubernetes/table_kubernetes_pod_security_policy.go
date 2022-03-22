package kubernetes

import (
	"context"

	v1beta1 "k8s.io/api/policy/v1beta1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/turbot/steampipe-plugin-sdk/v2/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin/transform"
)

func tableKubernetesPodSecurityPolicy(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_pod_security_policy",
		Description: "A Pod Security Policy is a cluster-level resource that controls security sensitive aspects of the pod specification. The PodSecurityPolicy objects define a set of conditions that a pod must run with in order to be accepted into the system, as well as defaults for the related fields.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getPodSecurityPolicy,
		},
		List: &plugin.ListConfig{
			Hydrate: listPodSecurityPolicy,
		},
		// PodSecurityPolicy, is a non-namespaced resource.
		Columns: k8sCommonGlobalColumns([]*plugin.Column{

			// PodSecurityPolicySpec
			{
				Name:        "allow_privilege_escalation",
				Type:        proto.ColumnType_BOOL,
				Description: "Determines if a pod can request to allow privilege escalation. If unspecified, defaults to true.",
				Transform:   transform.FromField("Spec.AllowPrivilegeEscalation"),
			},
			{
				Name:        "default_allow_privilege_escalation",
				Type:        proto.ColumnType_BOOL,
				Description: "Controls the default setting for whether a process can gain more privileges than its parent process.",
				Transform:   transform.FromField("Spec.DefaultAllowPrivilegeEscalation"),
			},
			{
				Name:        "host_network",
				Type:        proto.ColumnType_BOOL,
				Description: "Determines if the policy allows the use of HostNetwork in the pod spec.",
				Transform:   transform.FromField("Spec.HostNetwork"),
			},
			{
				Name:        "host_ports",
				Type:        proto.ColumnType_JSON,
				Description: "Determines which host port ranges are allowed to be exposed.",
				Transform:   transform.FromField("Spec.HostPorts"),
			},
			{
				Name:        "host_pid",
				Type:        proto.ColumnType_BOOL,
				Description: "Determines if the policy allows the use of HostPID in the pod spec.",
				Transform:   transform.FromField("Spec.HostPID"),
			},
			{
				Name:        "host_ipc",
				Type:        proto.ColumnType_BOOL,
				Description: "Determines if the policy allows the use of HostIPC in the pod spec.",
				Transform:   transform.FromField("Spec.HostIPC"),
			},
			{
				Name:        "privileged",
				Type:        proto.ColumnType_BOOL,
				Description: "privileged determines if a pod can request to be run as privileged.",
				Transform:   transform.FromField("Spec.Privileged"),
			},
			{
				Name:        "read_only_root_filesystem",
				Type:        proto.ColumnType_BOOL,
				Description: "If set to true will force containers to run with a read only root file system. If set to false the container may run with a read only root file system if it wishes but it will not be forced to.",
				Transform:   transform.FromField("Spec.ReadOnlyRootFilesystem"),
			},

			// JSON Fields
			{
				Name:        "allowed_csi_drivers",
				Type:        proto.ColumnType_JSON,
				Description: "An allowlist of inline CSI drivers that must be explicitly set to be embedded within a pod spec.",
				Transform:   transform.FromField("Spec.allowedCSIDrivers"),
			},
			{
				Name:        "allowed_host_paths",
				Type:        proto.ColumnType_JSON,
				Description: "An allowlist of host paths. Empty indicates that all host paths may be used.",
				Transform:   transform.FromField("Spec.AllowedHostPaths"),
			},
			{
				Name:        "allowed_flex_volumes",
				Type:        proto.ColumnType_JSON,
				Description: "An allowlist of Flexvolumes. Empty or nil indicates that all Flexvolumes may be used.",
				Transform:   transform.FromField("Spec.AllowedFlexVolumes"),
			},
			{
				Name:        "allowed_proc_mount_types",
				Type:        proto.ColumnType_JSON,
				Description: "An allowlist of allowed ProcMountTypes. Empty or nil indicates that only the DefaultProcMountType may be used.",
				Transform:   transform.FromField("Spec.AllowedProcMountTypes"),
			},
			{
				Name:        "allowed_unsafe_sysctls",
				Type:        proto.ColumnType_JSON,
				Description: "List of explicitly allowed unsafe sysctls, defaults to none.",
				Transform:   transform.FromField("Spec.AllowedUnsafeSysctls"),
			},
			{
				Name:        "default_add_capabilities",
				Type:        proto.ColumnType_JSON,
				Description: "List of the default set of capabilities that will be added to the container unless the pod spec specifically drops the capability.",
				Transform:   transform.FromField("Spec.DefaultAddCapabilities"),
			},
			{
				Name:        "forbidden_sysctls",
				Type:        proto.ColumnType_JSON,
				Description: "List of explicitly forbidden sysctls, defaults to none.",
				Transform:   transform.FromField("Spec.ForbiddenSysctls"),
			},
			{
				Name:        "fs_group",
				Type:        proto.ColumnType_JSON,
				Description: "The strategy that will dictate what fs group is used by the SecurityContext.",
				Transform:   transform.FromField("Spec.FSGroup"),
			},
			{
				Name:        "required_drop_capabilities",
				Type:        proto.ColumnType_JSON,
				Description: "List of the capabilities that will be dropped from the container. These are required to be dropped and cannot be added.",
				Transform:   transform.FromField("Spec.RequiredDropCapabilities"),
			},
			{
				Name:        "run_as_group",
				Type:        proto.ColumnType_JSON,
				Description: "The strategy that will dictate the allowable RunAsGroup values that may be set.",
				Transform:   transform.FromField("Spec.RunAsGroup"),
			},
			{
				Name:        "run_as_user",
				Type:        proto.ColumnType_JSON,
				Description: "The strategy that will dictate the allowable RunAsUser values that may be set.",
				Transform:   transform.FromField("Spec.RunAsUser"),
			},
			{
				Name:        "runtime_class",
				Type:        proto.ColumnType_JSON,
				Description: "The strategy that will dictate the allowable RuntimeClasses for a pod.",
				Transform:   transform.FromField("Spec.RuntimeClass"),
			},
			{
				Name:        "se_linux",
				Type:        proto.ColumnType_JSON,
				Description: "The strategy that will dictate the allowable labels that may be set.",
				Transform:   transform.FromField("Spec.SELinux"),
			},
			{
				Name:        "supplemental_groups",
				Type:        proto.ColumnType_JSON,
				Description: "The strategy that will dictate what supplemental groups are used by the SecurityContext.",
				Transform:   transform.FromField("Spec.SupplementalGroups"),
			},
			{
				Name:        "volumes",
				Type:        proto.ColumnType_JSON,
				Description: "An allowlist of volume plugins. Empty indicates that no volumes may be used.",
				Transform:   transform.FromField("Spec.Volumes"),
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
				Transform:   transform.From(transformPodSecurityPolicyTags),
			},
		}),
	}
}

//// HYDRATE FUNCTIONS

func listPodSecurityPolicy(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listPodSecurityPolicy")

	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		return nil, err
	}

	response, err := clientset.PolicyV1beta1().PodSecurityPolicies().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, podSecurityPolicy := range response.Items {
		d.StreamListItem(ctx, podSecurityPolicy)
	}

	return nil, nil
}

func getPodSecurityPolicy(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getK8sSecret")

	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		return nil, err
	}

	name := d.KeyColumnQuals["name"].GetStringValue()

	// return if  name is empty
	if name == "" {
		return nil, nil
	}

	podSecurityPolicy, err := clientset.PolicyV1beta1().PodSecurityPolicies().Get(ctx, name, metav1.GetOptions{})
	if err != nil && !isNotFoundError(err) {
		return nil, err
	}

	return *podSecurityPolicy, nil
}

//// TRANSFORM FUNCTIONS

func transformPodSecurityPolicyTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(v1beta1.PodSecurityPolicy)
	return mergeTags(obj.Labels, obj.Annotations), nil
}
