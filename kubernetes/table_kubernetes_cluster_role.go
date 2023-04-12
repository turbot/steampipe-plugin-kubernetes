package kubernetes

import (
	"context"

	v1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableKubernetesClusterRole(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_cluster_role",
		Description: "ClusterRole contains rules that represent a set of permissions.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getK8sClusterRole,
		},
		List: &plugin.ListConfig{
			Hydrate: listK8sClusterRoles,
		},
		// ClusterRole, is a non-namespaced resource.
		Columns: k8sCommonGlobalColumns([]*plugin.Column{
			{
				Name:        "rules",
				Type:        proto.ColumnType_JSON,
				Description: "List of the PolicyRules for this Role.",
			},
			{
				Name:        "aggregation_rule",
				Type:        proto.ColumnType_JSON,
				Description: "An optional field that describes how to build the Rules for this ClusterRole",
			},
			{
				Name:        "source",
				Type:        proto.ColumnType_STRING,
				Description: "The source of the resource. Possible values are: deployed and manifest. If the resource is fetched from the spec file the value will be manifest.",
				Transform:   transform.From(clusterRoleResourceSource),
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
				Transform:   transform.From(transformClusterRoleTags),
			},
		}),
	}
}

type ClusterRole struct {
	v1.ClusterRole
	Path      string
	StartLine int
	EndLine   int
}

//// HYDRATE FUNCTIONS

func listK8sClusterRoles(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listK8sClusterRoles")

	// Get the client for querying the K8s APIs for the provided context.
	// If the connection is configured for the manifest files, the client will return nil.
	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		return nil, err
	}

	//
	// Check for manifest files
	//
	parsedContents, err := fetchResourceFromManifestFileByKind(ctx, d, "ClusterRole")
	if err != nil {
		return nil, err
	}

	for _, content := range parsedContents {
		clusterRole := content.Data.(*v1.ClusterRole)

		d.StreamListItem(ctx, ClusterRole{*clusterRole, content.Path, content.StartLine, content.EndLine})

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
	var response *v1.ClusterRoleList
	pageLeft := true
	for pageLeft {

		response, err = clientset.RbacV1().ClusterRoles().List(ctx, input)
		if err != nil {
			return nil, err
		}

		if response.GetContinue() != "" {
			input.Continue = response.Continue
		} else {
			pageLeft = false
		}

		for _, clusterRole := range response.Items {
			d.StreamListItem(ctx, ClusterRole{clusterRole, "", 0, 0})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

func getK8sClusterRole(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getK8sClusterRole")

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
	parsedContents, err := fetchResourceFromManifestFileByKind(ctx, d, "ClusterRole")
	if err != nil {
		return nil, err
	}

	for _, content := range parsedContents {
		clusterRole := content.Data.(*v1.ClusterRole)

		if clusterRole.Name == name {
			return ClusterRole{*clusterRole, content.Path, content.StartLine, content.EndLine}, nil
		}
	}

	//
	// Get the deployed resource
	//
	if clientset == nil {
		return nil, nil
	}

	clusterRole, err := clientset.RbacV1().ClusterRoles().Get(ctx, name, metav1.GetOptions{})
	if err != nil && !isNotFoundError(err) {
		return nil, err
	}

	return ClusterRole{*clusterRole, "", 0, 0}, nil
}

//// TRANSFORM FUNCTIONS

func transformClusterRoleTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(ClusterRole)
	return mergeTags(obj.Labels, obj.Annotations), nil
}

func clusterRoleResourceSource(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(ClusterRole)

	if obj.Path != "" {
		return "manifest", nil
	}
	return "deployed", nil
}
