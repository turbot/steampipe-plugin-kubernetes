package kubernetes

import (
	"context"

	v1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/turbot/steampipe-plugin-sdk/v2/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin/transform"
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

//// HYDRATE FUNCTIONS

func listK8sClusterRoleBindings(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listK8sClusterRoleBindings")

	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		return nil, err
	}

	response, err := clientset.RbacV1().ClusterRoleBindings().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, clusterRoleBinding := range response.Items {
		d.StreamListItem(ctx, clusterRoleBinding)
	}

	return nil, nil
}

func getK8sClusterRoleBinding(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getK8sClusterRoleBinding")

	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		return nil, err
	}

	name := d.KeyColumnQuals["name"].GetStringValue()

	// return if name is empty
	if name == "" {
		return nil, nil
	}

	clusterRoleBinding, err := clientset.RbacV1().ClusterRoleBindings().Get(ctx, name, metav1.GetOptions{})
	if err != nil && !isNotFoundError(err) {
		return nil, err
	}

	return *clusterRoleBinding, nil
}

//// TRANSFORM FUNCTIONS

func transformClusterRoleBindingTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(v1.ClusterRoleBinding)
	return mergeTags(obj.Labels, obj.Annotations), nil
}
