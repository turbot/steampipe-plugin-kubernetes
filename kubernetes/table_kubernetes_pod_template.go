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

func tableKubernetesPodTemplate(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_pod_template",
		Description: "Kubernetes Pod Template is a collection of templates for creating copies of a predefined pod.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "namespace"}),
			Hydrate:    getK8sPodTemplate,
		},
		List: &plugin.ListConfig{
			Hydrate: listK8sPodTemplates,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "name", Require: plugin.Optional},
				{Name: "namespace", Require: plugin.Optional},
			},
		},
		Columns: k8sCommonColumns([]*plugin.Column{
			{
				Name:        "template",
				Type:        proto.ColumnType_JSON,
				Description: "Template describes the pods that will be created.",
			},

			{
				Name:        "context_name",
				Type:        proto.ColumnType_STRING,
				Description: "Kubectl config context name.",
				Hydrate:     getPodTemplateResourceContext,
			},

			{
				Name:        "source_type",
				Type:        proto.ColumnType_STRING,
				Description: "The source of the resource. Possible values are: deployed and manifest. If the resource is fetched from the spec file the value will be manifest.",
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
				Transform:   transform.From(transformPodTemplateTags),
			},
		}),
	}
}

type PodTemplate struct {
	v1.PodTemplate
	parsedContent
}

//// HYDRATE FUNCTIONS

func listK8sPodTemplates(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Debug("listK8sPodTemplates")

	// Get the client for querying the K8s APIs for the provided context.
	// If the connection is configured for the manifest files, the client will return nil.
	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		return nil, err
	}

	// Check for manifest files
	parsedContents, err := fetchResourceFromManifestFileByKind(ctx, d, "PodTemplate")
	if err != nil {
		return nil, err
	}

	for _, content := range parsedContents {
		podTemplate := content.ParsedData.(*v1.PodTemplate)

		d.StreamListItem(ctx, PodTemplate{*podTemplate, content})

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
		if *limit < 500 {
			input.Limit = *limit
		}
	}

	commonFieldSelectorValue := getCommonOptionalKeyQualsValueForFieldSelector(d)
	if len(commonFieldSelectorValue) > 0 {
		input.FieldSelector = strings.Join(commonFieldSelectorValue, ",")
	}

	var response *v1.PodTemplateList
	pageLeft := true

	for pageLeft {
		response, err = clientset.CoreV1().PodTemplates("").List(ctx, input)
		if err != nil {
			return nil, err
		}

		if response.GetContinue() != "" {
			input.Continue = response.Continue
		} else {
			pageLeft = false
		}

		for _, podTemplate := range response.Items {
			d.StreamListItem(ctx, PodTemplate{podTemplate, parsedContent{SourceType: "deployed"}})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

func getK8sPodTemplate(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Debug("getK8sPodTemplate")

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
	parsedContents, err := fetchResourceFromManifestFileByKind(ctx, d, "PodTemplate")
	if err != nil {
		return nil, err
	}

	for _, content := range parsedContents {
		podTemplate := content.ParsedData.(*v1.PodTemplate)

		if podTemplate.Name == name && podTemplate.Namespace == namespace {
			return PodTemplate{*podTemplate, content}, nil
		}
	}

	// Get the deployed resource
	if clientset == nil {
		return nil, nil
	}

	podTemplate, err := clientset.CoreV1().PodTemplates(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil && !isNotFoundError(err) {
		return nil, err
	}

	return PodTemplate{*podTemplate, parsedContent{SourceType: "deployed"}}, nil
}

func getPodTemplateResourceContext(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	obj := h.Item.(PodTemplate)

	// Set the context_name as nil
	data := map[string]interface{}{}
	if obj.Path != "" {
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

func transformPodTemplateTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(PodTemplate)
	return mergeTags(obj.Labels, obj.Annotations), nil
}
