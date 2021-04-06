package kubernetes

import (
	"context"

	v1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableKubernetesClusterRole(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_cluster_role",
		Description: "A role at the cluster level and in all namespaces.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getK8sClusterRole,
		},
		List: &plugin.ListConfig{
			Hydrate: listK8sClusterRoles,
		},
		Columns: k8sCommonGlobalColumns([]*plugin.Column{
			{
				Name:        "rules",
				Type:        proto.ColumnType_JSON,
				Description: "List of the PolicyRules for this Role.",
				// Transform:   transform.FromField("Rules"),
			},
			{
				Name:        "aggregation_rule",
				Type:        proto.ColumnType_JSON,
				Description: "An optional field that describes how to build the Rules for this ClusterRole",
				Transform:   transform.FromField("Spec.Containers"),
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
				Transform:   transform.From(transformRoleTags),
			},
		}),
	}
}

//// HYDRATE FUNCTIONS

func listK8sClusterRoles(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listK8sClusterRoles")

	clientset, err := GetNewClientset(ctx, d.ConnectionManager)
	if err != nil {
		return nil, err
	}

	response, err := clientset.RbacV1().ClusterRoles().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, clusterRole := range response.Items {
		d.StreamListItem(ctx, clusterRole)
	}

	return nil, nil
}

func getK8sClusterRole(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getK8sClusterRole")

	clientset, err := GetNewClientset(ctx, d.ConnectionManager)
	if err != nil {
		return nil, err
	}

	name := d.KeyColumnQuals["name"].GetStringValue()
	// namespace := d.KeyColumnQuals["namespace"].GetStringValue()

	clusterRole, err := clientset.RbacV1().ClusterRoles().Get(ctx, name, metav1.GetOptions{})
	if err != nil && !isNotFoundError(err) {
		return nil, err
	}

	return *clusterRole, nil
}

//// TRANSFORM FUNCTIONS

func transformRoleTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(v1.ClusterRole)
	return mergeTags(obj.Labels, obj.Annotations), nil
}
