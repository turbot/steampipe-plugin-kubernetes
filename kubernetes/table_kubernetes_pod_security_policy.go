package kubernetes

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableKubernetesPodSecurityPolicy(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		// Name:        "kubernetes_pod_security_policy",
		// Description: "A Pod Security Policy is a cluster-level resource that controls security sensitive aspects of the pod specification. The PodSecurityPolicy objects define a set of conditions that a pod must run with in order to be accepted into the system, as well as defaults for the related fields.",
		// Get: &plugin.GetConfig{
		// 	KeyColumns: plugin.SingleColumn("name"),
		// 	Hydrate:    getPodSecurityPolicy,
		// },
		// List: &plugin.ListConfig{
		// 	Hydrate: listPodSecurityPolicy,
		// },
		// // PodSecurityPolicy, is a non-namespaced resource.
		// Columns: k8sCommonGlobalColumns([]*plugin.Column{

		// 	// PodSecurityPolicySpec
		// 	{
		// 		Name:        "allow_privilege_escalation",
		// 		Type:        proto.ColumnType_BOOL,
		// 		Description: "Determines if a pod can request to allow privilege escalation. If unspecified, defaults to true.",
		// 		Transform:   transform.FromField("Spec.AllowPrivilegeEscalation"),
		// 	},
		// 	{
		// 		Name:        "default_allow_privilege_escalation",
		// 		Type:        proto.ColumnType_BOOL,
		// 		Description: "Controls the default setting for whether a process can gain more privileges than its parent process.",
		// 		Transform:   transform.FromField("Spec.DefaultAllowPrivilegeEscalation"),
		// 	},
		// 	{
		// 		Name:        "host_network",
		// 		Type:        proto.ColumnType_BOOL,
		// 		Description: "Determines if the policy allows the use of HostNetwork in the pod spec.",
		// 		Transform:   transform.FromField("Spec.HostNetwork"),
		// 	},
		// 	{
		// 		Name:        "host_ports",
		// 		Type:        proto.ColumnType_JSON,
		// 		Description: "Determines which host port ranges are allowed to be exposed.",
		// 		Transform:   transform.FromField("Spec.HostPorts"),
		// 	},
		// 	{
		// 		Name:        "host_pid",
		// 		Type:        proto.ColumnType_BOOL,
		// 		Description: "Determines if the policy allows the use of HostPID in the pod spec.",
		// 		Transform:   transform.FromField("Spec.HostPID"),
		// 	},
		// 	{
		// 		Name:        "host_ipc",
		// 		Type:        proto.ColumnType_BOOL,
		// 		Description: "Determines if the policy allows the use of HostIPC in the pod spec.",
		// 		Transform:   transform.FromField("Spec.HostIPC"),
		// 	},
		// 	{
		// 		Name:        "privileged",
		// 		Type:        proto.ColumnType_BOOL,
		// 		Description: "privileged determines if a pod can request to be run as privileged.",
		// 		Transform:   transform.FromField("Spec.Privileged"),
		// 	},
		// 	{
		// 		Name:        "read_only_root_filesystem",
		// 		Type:        proto.ColumnType_BOOL,
		// 		Description: "If set to true will force containers to run with a read only root file system. If set to false the container may run with a read only root file system if it wishes but it will not be forced to.",
		// 		Transform:   transform.FromField("Spec.ReadOnlyRootFilesystem"),
		// 	},

		// 	// JSON Fields
		// 	{
		// 		Name:        "allowed_csi_drivers",
		// 		Type:        proto.ColumnType_JSON,
		// 		Description: "An allowlist of inline CSI drivers that must be explicitly set to be embedded within a pod spec.",
		// 		Transform:   transform.FromField("Spec.allowedCSIDrivers"),
		// 	},
		// 	{
		// 		Name:        "allowed_host_paths",
		// 		Type:        proto.ColumnType_JSON,
		// 		Description: "An allowlist of host paths. Empty indicates that all host paths may be used.",
		// 		Transform:   transform.FromField("Spec.AllowedHostPaths"),
		// 	},
		// 	{
		// 		Name:        "allowed_flex_volumes",
		// 		Type:        proto.ColumnType_JSON,
		// 		Description: "An allowlist of Flexvolumes. Empty or nil indicates that all Flexvolumes may be used.",
		// 		Transform:   transform.FromField("Spec.AllowedFlexVolumes"),
		// 	},
		// 	{
		// 		Name:        "allowed_proc_mount_types",
		// 		Type:        proto.ColumnType_JSON,
		// 		Description: "An allowlist of allowed ProcMountTypes. Empty or nil indicates that only the DefaultProcMountType may be used.",
		// 		Transform:   transform.FromField("Spec.AllowedProcMountTypes"),
		// 	},
		// 	{
		// 		Name:        "allowed_unsafe_sysctls",
		// 		Type:        proto.ColumnType_JSON,
		// 		Description: "List of explicitly allowed unsafe sysctls, defaults to none.",
		// 		Transform:   transform.FromField("Spec.AllowedUnsafeSysctls"),
		// 	},
		// 	{
		// 		Name:        "default_add_capabilities",
		// 		Type:        proto.ColumnType_JSON,
		// 		Description: "List of the default set of capabilities that will be added to the container unless the pod spec specifically drops the capability.",
		// 		Transform:   transform.FromField("Spec.DefaultAddCapabilities"),
		// 	},
		// 	{
		// 		Name:        "forbidden_sysctls",
		// 		Type:        proto.ColumnType_JSON,
		// 		Description: "List of explicitly forbidden sysctls, defaults to none.",
		// 		Transform:   transform.FromField("Spec.ForbiddenSysctls"),
		// 	},
		// 	{
		// 		Name:        "fs_group",
		// 		Type:        proto.ColumnType_JSON,
		// 		Description: "The strategy that will dictate what fs group is used by the SecurityContext.",
		// 		Transform:   transform.FromField("Spec.FSGroup"),
		// 	},
		// 	{
		// 		Name:        "required_drop_capabilities",
		// 		Type:        proto.ColumnType_JSON,
		// 		Description: "List of the capabilities that will be dropped from the container. These are required to be dropped and cannot be added.",
		// 		Transform:   transform.FromField("Spec.RequiredDropCapabilities"),
		// 	},
		// 	{
		// 		Name:        "run_as_group",
		// 		Type:        proto.ColumnType_JSON,
		// 		Description: "The strategy that will dictate the allowable RunAsGroup values that may be set.",
		// 		Transform:   transform.FromField("Spec.RunAsGroup"),
		// 	},
		// 	{
		// 		Name:        "run_as_user",
		// 		Type:        proto.ColumnType_JSON,
		// 		Description: "The strategy that will dictate the allowable RunAsUser values that may be set.",
		// 		Transform:   transform.FromField("Spec.RunAsUser"),
		// 	},
		// 	{
		// 		Name:        "runtime_class",
		// 		Type:        proto.ColumnType_JSON,
		// 		Description: "The strategy that will dictate the allowable RuntimeClasses for a pod.",
		// 		Transform:   transform.FromField("Spec.RuntimeClass"),
		// 	},
		// 	{
		// 		Name:        "se_linux",
		// 		Type:        proto.ColumnType_JSON,
		// 		Description: "The strategy that will dictate the allowable labels that may be set.",
		// 		Transform:   transform.FromField("Spec.SELinux"),
		// 	},
		// 	{
		// 		Name:        "supplemental_groups",
		// 		Type:        proto.ColumnType_JSON,
		// 		Description: "The strategy that will dictate what supplemental groups are used by the SecurityContext.",
		// 		Transform:   transform.FromField("Spec.SupplementalGroups"),
		// 	},
		// 	{
		// 		Name:        "volumes",
		// 		Type:        proto.ColumnType_JSON,
		// 		Description: "An allowlist of volume plugins. Empty indicates that no volumes may be used.",
		// 		Transform:   transform.FromField("Spec.Volumes"),
		// 	},
		// 	{
		// 		Name:        "context_name",
		// 		Type:        proto.ColumnType_STRING,
		// 		Description: "Kubectl config context name.",
		// 		Hydrate:     getPodSecurityPolicyResourceContext,
		// 	},

		// 	{
		// 		Name:        "source_type",
		// 		Type:        proto.ColumnType_STRING,
		// 		Description: "The source of the resource. Possible values are: deployed and manifest. If the resource is fetched from the spec file the value will be manifest.",
		// 	},

		// 	//// Steampipe Standard Columns
		// 	{
		// 		Name:        "title",
		// 		Type:        proto.ColumnType_STRING,
		// 		Description: ColumnDescriptionTitle,
		// 		Transform:   transform.FromField("Name"),
		// 	},
		// 	{
		// 		Name:        "tags",
		// 		Type:        proto.ColumnType_JSON,
		// 		Description: ColumnDescriptionTags,
		// 		Transform:   transform.From(transformPodSecurityPolicyTags),
		// 	},
		// }),
	}
}

// type PodSecurityPolicy struct {
// 	v1beta1.PodSecurityPolicy
// 	parsedContent
// }

// //// HYDRATE FUNCTIONS

// func listPodSecurityPolicy(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
// 	logger := plugin.Logger(ctx)
// 	logger.Trace("listPodSecurityPolicy")

// 	// Get the client for querying the K8s APIs for the provided context.
// 	// If the connection is configured for the manifest files, the client will return nil.
// 	clientset, err := GetNewClientset(ctx, d)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Check for manifest files
// 	parsedContents, err := fetchResourceFromManifestFileByKind(ctx, d, "PodSecurityPolicy")
// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, content := range parsedContents {
// 		podSecurityPolicy := content.ParsedData.(*v1beta1.PodSecurityPolicy)

// 		d.StreamListItem(ctx, PodSecurityPolicy{*podSecurityPolicy, content})

// 		// Context can be cancelled due to manual cancellation or the limit has been hit
// 		if d.RowsRemaining(ctx) == 0 {
// 			return nil, nil
// 		}
// 	}

// 	// Check for deployed resources
// 	if clientset == nil {
// 		return nil, nil
// 	}

// 	input := metav1.ListOptions{
// 		Limit: 500,
// 	}

// 	// Limiting the results
// 	limit := d.QueryContext.Limit
// 	if d.QueryContext.Limit != nil {
// 		if *limit < input.Limit {
// 			if *limit < 1 {
// 				input.Limit = 1
// 			} else {
// 				input.Limit = *limit
// 			}
// 		}
// 	}

// 	var response *v1beta1.PodSecurityPolicyList
// 	pageLeft := true

// 	for pageLeft {
// 		response, err = clientset.PolicyV1beta1().PodSecurityPolicies().List(ctx, input)
// 		if err != nil {
// 			return nil, err
// 		}

// 		if response.GetContinue() != "" {
// 			input.Continue = response.Continue
// 		} else {
// 			pageLeft = false
// 		}

// 		for _, podSecurityPolicy := range response.Items {
// 			d.StreamListItem(ctx, PodSecurityPolicy{podSecurityPolicy, parsedContent{SourceType: "deployed"}})

// 			// Context can be cancelled due to manual cancellation or the limit has been hit
// 			if d.RowsRemaining(ctx) == 0 {
// 				return nil, nil
// 			}
// 		}
// 	}

// 	return nil, nil
// }

// func getPodSecurityPolicy(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
// 	logger := plugin.Logger(ctx)
// 	logger.Trace("getPodSecurityPolicy")

// 	// Get the client for querying the K8s APIs for the provided context.
// 	// If the connection is configured for the manifest files, the client will return nil.
// 	clientset, err := GetNewClientset(ctx, d)
// 	if err != nil {
// 		return nil, err
// 	}

// 	name := d.EqualsQuals["name"].GetStringValue()

// 	// return if  name is empty
// 	if name == "" {
// 		return nil, nil
// 	}

// 	// Get the manifest resource
// 	parsedContents, err := fetchResourceFromManifestFileByKind(ctx, d, "PodSecurityPolicy")
// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, content := range parsedContents {
// 		podSecurityPolicy := content.ParsedData.(*v1beta1.PodSecurityPolicy)

// 		if podSecurityPolicy.Name == name {
// 			return PodSecurityPolicy{*podSecurityPolicy, content}, nil
// 		}
// 	}

// 	// Get the deployed resource
// 	if clientset == nil {
// 		return nil, nil
// 	}

// 	podSecurityPolicy, err := clientset.PolicyV1beta1().PodSecurityPolicies().Get(ctx, name, metav1.GetOptions{})
// 	if err != nil && !isNotFoundError(err) {
// 		return nil, err
// 	}

// 	return PodSecurityPolicy{*podSecurityPolicy, parsedContent{SourceType: "deployed"}}, nil
// }

// func getPodSecurityPolicyResourceContext(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
// 	obj := h.Item.(PodSecurityPolicy)

// 	// Set the context_name as nil
// 	data := map[string]interface{}{}
// 	if obj.Path != "" {
// 		return data, nil
// 	}

// 	// Else, set the current context as context_name
// 	currentContext, err := getKubectlContext(ctx, d, nil)
// 	if err != nil {
// 		return data, nil
// 	}
// 	data["ContextName"] = currentContext.(string)

// 	return data, nil
// }

// //// TRANSFORM FUNCTIONS

// func transformPodSecurityPolicyTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
// 	obj := d.HydrateItem.(PodSecurityPolicy)
// 	return mergeTags(obj.Labels, obj.Annotations), nil
// }
