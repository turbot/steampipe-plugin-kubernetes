package kubernetes

import (
	"context"

	v1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

func tableKubernetesRole(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_role",
		Description: "Role contains rules that represent a set of permissions.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "namespace"}),
			Hydrate:    getK8sRole,
		},
		List: &plugin.ListConfig{
			Hydrate: listK8sRoles,
		},
		Columns: k8sCommonColumns([]*plugin.Column{
			{
				Name:        "rules",
				Type:        proto.ColumnType_JSON,
				Description: "List of the PolicyRules for this Role.",
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

func listK8sRoles(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listK8sRoles")

	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		return nil, err
	}

	response, err := clientset.RbacV1().Roles("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, role := range response.Items {
		d.StreamListItem(ctx, role)
	}

	return nil, nil
}

func getK8sRole(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getK8sRole")

	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		return nil, err
	}

	name := d.KeyColumnQuals["name"].GetStringValue()
	namespace := d.KeyColumnQuals["namespace"].GetStringValue()

	// return if namespace or name is empty
	if namespace == "" || name == "" {
		return nil, nil
	}

	role, err := clientset.RbacV1().Roles(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil && !isNotFoundError(err) {
		return nil, err
	}

	return *role, nil
}

//// TRANSFORM FUNCTIONS

func transformRoleTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(v1.Role)
	return mergeTags(obj.Labels, obj.Annotations), nil
}
