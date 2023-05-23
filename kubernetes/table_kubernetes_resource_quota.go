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

func tableKubernetesResourceQuota(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_resource_quota",
		Description: "Kubernetes Resource Quota",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "namespace"}),
			Hydrate:    getK8sResourceQuota,
		},
		List: &plugin.ListConfig{
			Hydrate:    listK8sResourceQuotas,
			KeyColumns: getCommonOptionalKeyQuals(),
		},
		Columns: k8sCommonColumns([]*plugin.Column{

			//// ResourceQuotaSpec Columns
			{
				Name:        "spec_hard",
				Type:        proto.ColumnType_JSON,
				Description: "Spec hard is the set of desired hard limits for each named resource.",
				Transform:   transform.FromField("Spec.Hard"),
			},
			{
				Name:        "spec_scopes",
				Type:        proto.ColumnType_JSON,
				Description: "A collection of filters that must match each object tracked by a quota.",
				Transform:   transform.FromField("Spec.Scopes"),
			},
			{
				Name:        "spec_scope_selector",
				Type:        proto.ColumnType_JSON,
				Description: "A collection of filters like scopes that must match each object tracked by a quota but expressed using ScopeSelectorOperator in combination with possible values.",
				Transform:   transform.FromField("Spec.ScopeSelector"),
			},

			//// ResourceQuotaStatus Columns
			{
				Name:        "status_hard",
				Type:        proto.ColumnType_JSON,
				Description: "Status hard is the set of enforced hard limits for each named resource.",
				Transform:   transform.FromField("Status.Hard"),
			},
			{
				Name:        "status_used",
				Type:        proto.ColumnType_JSON,
				Description: "Indicates current observed total usage of the resource in the namespace.",
				Transform:   transform.FromField("Status.Used"),
			},
			{
				Name:        "context_name",
				Type:        proto.ColumnType_STRING,
				Description: "Kubectl config context name.",
				Hydrate:     getResourceQuotaResourceAdditionalData,
			},
			{
				Name:        "source_type",
				Type:        proto.ColumnType_STRING,
				Description: "The source of the resource. Possible values are: deployed and manifest. If the resource is fetched from the spec file the value will be manifest.",
				Hydrate:     getResourceQuotaResourceAdditionalData,
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
				Transform:   transform.From(transformResourceQuotaTags),
			},
		}),
	}
}

type ResourceQuota struct {
	v1.ResourceQuota
	Path      string
	StartLine int
	EndLine   int
}

//// HYDRATE FUNCTIONS

func listK8sResourceQuotas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listK8sResourceQuotas")

	// Get the client for querying the K8s APIs for the provided context.
	// If the connection is configured for the manifest files, the client will return nil.
	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		return nil, err
	}

	// Check for manifest files
	parsedContents, err := fetchResourceFromManifestFileByKind(ctx, d, "ResourceQuota")
	if err != nil {
		return nil, err
	}

	for _, content := range parsedContents {
		resourceQuota := content.Data.(*v1.ResourceQuota)

		d.StreamListItem(ctx, ResourceQuota{*resourceQuota, content.Path, content.StartLine, content.EndLine})

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

	var response *v1.ResourceQuotaList
	pageLeft := true

	for pageLeft {
		response, err = clientset.CoreV1().ResourceQuotas("").List(ctx, input)
		if err != nil {
			return nil, err
		}

		if response.GetContinue() != "" {
			input.Continue = response.Continue
		} else {
			pageLeft = false
		}

		for _, resourceQuota := range response.Items {
			d.StreamListItem(ctx, ResourceQuota{resourceQuota, "", 0, 0})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

func getK8sResourceQuota(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getK8sResourceQuota")

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
	parsedContents, err := fetchResourceFromManifestFileByKind(ctx, d, "ResourceQuota")
	if err != nil {
		return nil, err
	}

	for _, content := range parsedContents {
		resourceQuota := content.Data.(*v1.ResourceQuota)

		if resourceQuota.Name == name && resourceQuota.Namespace == namespace {
			return ResourceQuota{*resourceQuota, content.Path, content.StartLine, content.EndLine}, nil
		}
	}

	// Get the deployed resource
	if clientset == nil {
		return nil, nil
	}

	resourceQuota, err := clientset.CoreV1().ResourceQuotas(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil && !isNotFoundError(err) {
		return nil, err
	}

	return ResourceQuota{*resourceQuota, "", 0, 0}, nil
}

func getResourceQuotaResourceAdditionalData(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	obj := h.Item.(ResourceQuota)

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

func transformResourceQuotaTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(ResourceQuota)
	return mergeTags(obj.Labels, obj.Annotations), nil
}
