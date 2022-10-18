package kubernetes

import (
	"context"
	"strings"

	v1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableKubernetesRoleBinding(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_role_binding",
		Description: "A role binding grants the permissions defined in a role to a user or set of users. It holds a list of subjects (users, groups, or service accounts), and a reference to the role being granted.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "namespace"}),
			Hydrate:    getK8sRoleBinding,
		},
		List: &plugin.ListConfig{
			Hydrate:    listK8sRoleBindings,
			KeyColumns: getCommonOptionalKeyQuals(),
		},
		Columns: k8sCommonColumns([]*plugin.Column{
			{
				Name:        "subjects",
				Type:        proto.ColumnType_JSON,
				Description: "List of references to the objects the role applies to.",
			},

			//// RoleRef columns
			{
				Name:        "role_name",
				Type:        proto.ColumnType_STRING,
				Description: "Name of the role for which access is granted to subjects.",
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
				Description: "Type of the role referenced.",
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
				Transform:   transform.From(transformRoleBindingTags),
			},
		}),
	}
}

//// HYDRATE FUNCTIONS

func listK8sRoleBindings(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listK8sRoleBindings")

	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		return nil, err
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

	var response *v1.RoleBindingList
	pageLeft := true

	for pageLeft {
		response, err = clientset.RbacV1().RoleBindings("").List(ctx, input)
		if err != nil {
			return nil, err
		}

		if response.GetContinue() != "" {
			input.Continue = response.Continue
		} else {
			pageLeft = false
		}

		for _, roleBinding := range response.Items {
			d.StreamListItem(ctx, roleBinding)
		}
	}

	return nil, nil
}

func getK8sRoleBinding(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getK8sRoleBinding")

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

	roleBinding, err := clientset.RbacV1().RoleBindings(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil && !isNotFoundError(err) {
		return nil, err
	}

	return *roleBinding, nil
}

//// TRANSFORM FUNCTIONS

func transformRoleBindingTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(v1.RoleBinding)
	return mergeTags(obj.Labels, obj.Annotations), nil
}
