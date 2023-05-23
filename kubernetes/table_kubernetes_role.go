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

func tableKubernetesRole(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_role",
		Description: "Role contains rules that represent a set of permissions.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "namespace"}),
			Hydrate:    getK8sRole,
		},
		List: &plugin.ListConfig{
			Hydrate:    listK8sRoles,
			KeyColumns: getCommonOptionalKeyQuals(),
		},
		Columns: k8sCommonColumns([]*plugin.Column{
			{
				Name:        "rules",
				Type:        proto.ColumnType_JSON,
				Description: "List of the PolicyRules for this Role.",
			},
			{
				Name:        "context_name",
				Type:        proto.ColumnType_STRING,
				Description: "Kubectl config context name.",
				Hydrate:     getRoleResourceAdditionalData,
			},
			{
				Name:        "source_type",
				Type:        proto.ColumnType_STRING,
				Description: "The source of the resource. Possible values are: deployed and manifest. If the resource is fetched from the spec file the value will be manifest.",
				Hydrate:     getRoleResourceAdditionalData,
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

type Role struct {
	v1.Role
	Path      string
	StartLine int
	EndLine   int
}

//// HYDRATE FUNCTIONS

func listK8sRoles(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listK8sRoles")

	// Get the client for querying the K8s APIs for the provided context.
	// If the connection is configured for the manifest files, the client will return nil.
	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		return nil, err
	}

	// Check for manifest files
	parsedContents, err := fetchResourceFromManifestFileByKind(ctx, d, "Role")
	if err != nil {
		return nil, err
	}

	for _, content := range parsedContents {
		role := content.Data.(*v1.Role)

		d.StreamListItem(ctx, Role{*role, content.Path, content.StartLine, content.EndLine})

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

	var response *v1.RoleList
	pageLeft := true

	for pageLeft {
		response, err = clientset.RbacV1().Roles("").List(ctx, input)
		if err != nil {
			return nil, err
		}

		if response.GetContinue() != "" {
			input.Continue = response.Continue
		} else {
			pageLeft = false
		}

		for _, role := range response.Items {
			d.StreamListItem(ctx, Role{role, "", 0, 0})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

func getK8sRole(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getK8sRole")

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
	parsedContents, err := fetchResourceFromManifestFileByKind(ctx, d, "Role")
	if err != nil {
		return nil, err
	}

	for _, content := range parsedContents {
		role := content.Data.(*v1.Role)

		if role.Name == name && role.Namespace == namespace {
			return Role{*role, content.Path, content.StartLine, content.EndLine}, nil
		}
	}

	// Get the deployed resource
	if clientset == nil {
		return nil, nil
	}

	role, err := clientset.RbacV1().Roles(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil && !isNotFoundError(err) {
		return nil, err
	}

	return Role{*role, "", 0, 0}, nil
}

func getRoleResourceAdditionalData(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	obj := h.Item.(Role)

	data := map[string]interface{}{
		"SourceType": "deployed",
	}

	// Set the source_type as manifest, if path is not empty
	// also, set the context_name as nil
	if obj.Path != "" {
		data["SourceType"] = "manifest"
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

func transformRoleTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(Role)
	return mergeTags(obj.Labels, obj.Annotations), nil
}
