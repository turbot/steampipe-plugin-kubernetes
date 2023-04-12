package kubernetes

import (
	"context"

	v1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableKubernetesClusterRoleBinding(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_cluster_role_binding",
		Description: "A ClusterRoleBinding grants the permissions defined in a cluster role to a user or set of users. Access granted by ClusterRoleBinding is cluster-wide.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getK8sClusterRoleBinding,
		},
		List: &plugin.ListConfig{
			Hydrate: listK8sClusterRoleBindings,
		},
		Columns: k8sCommonGlobalColumns([]*plugin.Column{
			{
				Name:        "subjects",
				Type:        proto.ColumnType_JSON,
				Description: "List of references to the objects the role applies to.",
			},

			//// RoleRef columns
			{
				Name:        "role_name",
				Type:        proto.ColumnType_STRING,
				Description: "Name of the cluster role for which access is granted to subjects.",
				Transform:   transform.FromField("RoleRef.Name"),
			},
			{
				Name:        "role_api_group",
				Type:        proto.ColumnType_STRING,
				Description: "The group for the referenced role.",
				Transform:   transform.FromField("RoleRef.APIGroup"),
			},
			{
				Name:        "role_kind",
				Type:        proto.ColumnType_STRING,
				Description: "Type of the role refrenced must be one of ClusterRole or Role.",
				Transform:   transform.FromField("RoleRef.Kind"),
			},
			{
				Name:        "source",
				Type:        proto.ColumnType_STRING,
				Description: "The source of the resource. Possible values are: deployed and manifest. If the resource is fetched from the spec file the value will be manifest.",
				Transform:   transform.From(clusterRoleBindingResourceSource),
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
				Transform:   transform.From(transformClusterRoleBindingTags),
			},
		}),
	}
}

type ClusterRoleBinding struct {
	v1.ClusterRoleBinding
	Path      string
	StartLine int
	EndLine   int
}

//// HYDRATE FUNCTIONS

func listK8sClusterRoleBindings(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listK8sClusterRoleBindings")

	// Get the client for querying the K8s APIs for the provided context.
	// If the connection is configured for the manifest files, the client will return nil.
	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		return nil, err
	}

	//
	// Check for manifest files
	//
	parsedContents, err := fetchResourceFromManifestFileByKind(ctx, d, "ClusterRoleBinding")
	if err != nil {
		return nil, err
	}

	for _, content := range parsedContents {
		clusterRoleBinding := content.Data.(*v1.ClusterRoleBinding)

		d.StreamListItem(ctx, ClusterRoleBinding{*clusterRoleBinding, content.Path, content.StartLine, content.EndLine})

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

	var response *v1.ClusterRoleBindingList
	pageLeft := true

	for pageLeft {
		response, err = clientset.RbacV1().ClusterRoleBindings().List(ctx, input)
		if err != nil {
			return nil, err
		}

		if response.GetContinue() != "" {
			input.Continue = response.Continue
		} else {
			pageLeft = false
		}

		for _, clusterRoleBinding := range response.Items {
			d.StreamListItem(ctx, ClusterRoleBinding{clusterRoleBinding, "", 0, 0})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

func getK8sClusterRoleBinding(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getK8sClusterRoleBinding")

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

	//
	// Get the manifest resource
	//
	parsedContents, err := fetchResourceFromManifestFileByKind(ctx, d, "ClusterRoleBinding")
	if err != nil {
		return nil, err
	}

	for _, content := range parsedContents {
		clusterRoleBinding := content.Data.(*v1.ClusterRoleBinding)

		if clusterRoleBinding.Name == name {
			return ClusterRoleBinding{*clusterRoleBinding, content.Path, content.StartLine, content.EndLine}, nil
		}
	}

	//
	// Get the deployed resource
	//
	if clientset == nil {
		return nil, nil
	}

	clusterRoleBinding, err := clientset.RbacV1().ClusterRoleBindings().Get(ctx, name, metav1.GetOptions{})
	if err != nil && !isNotFoundError(err) {
		return nil, err
	}

	return ClusterRoleBinding{*clusterRoleBinding, "", 0, 0}, nil
}

//// TRANSFORM FUNCTIONS

func transformClusterRoleBindingTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(ClusterRoleBinding)
	return mergeTags(obj.Labels, obj.Annotations), nil
}

func clusterRoleBindingResourceSource(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(ClusterRoleBinding)

	if obj.Path != "" {
		return "manifest", nil
	}
	return "deployed", nil
}
