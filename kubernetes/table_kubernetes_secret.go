package kubernetes

import (
	"context"
	"fmt"
	"strings"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableKubernetesSecret(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_secret",
		Description: "Secrets can be used to store sensitive information either as individual properties or coarse-grained entries like entire files or JSON blobs.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "namespace"}),
			Hydrate:    getK8sSecret,
		},
		List: &plugin.ListConfig{
			Hydrate: listK8sSecrets,
			KeyColumns: getOptionalKeyQualWithCommonKeyQuals([]*plugin.KeyColumn{
				{Name: "type", Require: plugin.Optional},
			}),
		},
		// ClusterRole, is a non-namespaced resource.
		Columns: k8sCommonColumns([]*plugin.Column{
			{
				Name:        "immutable",
				Type:        proto.ColumnType_BOOL,
				Description: "If set to true, ensures that data stored in the Secret cannot be updated (only object metadata can be modified). If not set to true, the field can be modified at any time. Defaulted to nil.",
			},
			{
				Name:        "type",
				Type:        proto.ColumnType_STRING,
				Description: "Type of the secret data.",
			},
			{
				Name:        "data",
				Type:        proto.ColumnType_JSON,
				Description: "Contains the secret data.",
			},
			{
				Name:        "string_data",
				Type:        proto.ColumnType_JSON,
				Description: "Contains the configuration binary data.",
			},
			{
				Name:        "context_name",
				Type:        proto.ColumnType_STRING,
				Description: "Kubectl config context name.",
				Hydrate:     getSecretResourceContext,
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
				Transform:   transform.From(transformSecretTags),
			},
		}),
	}
}

type Secret struct {
	v1.Secret
	parsedContent
}

//// HYDRATE FUNCTIONS

func listK8sSecrets(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listK8sSecrets")

	// Get the client for querying the K8s APIs for the provided context.
	// If the connection is configured for the manifest files, the client will return nil.
	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		return nil, err
	}

	// Check for manifest files
	parsedContents, err := fetchResourceFromManifestFileByKind(ctx, d, "Secret")
	if err != nil {
		return nil, err
	}

	for _, content := range parsedContents {
		secret := content.ParsedData.(*v1.Secret)

		d.StreamListItem(ctx, Secret{*secret, content})

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

	if d.EqualsQualString("type") != "" {
		commonFieldSelectorValue = append(commonFieldSelectorValue, fmt.Sprintf("type=%v", d.EqualsQualString("type")))
	}

	if len(commonFieldSelectorValue) > 0 {
		input.FieldSelector = strings.Join(commonFieldSelectorValue, ",")
	}

	var response *v1.SecretList
	pageLeft := true

	for pageLeft {
		response, err = clientset.CoreV1().Secrets(d.EqualsQualString("namespace")).List(ctx, input)
		if err != nil {
			return nil, err
		}

		if response.GetContinue() != "" {
			input.Continue = response.Continue
		} else {
			pageLeft = false
		}

		for _, secret := range response.Items {
			d.StreamListItem(ctx, Secret{secret, parsedContent{SourceType: "deployed"}})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

func getK8sSecret(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getK8sSecret")

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
	parsedContents, err := fetchResourceFromManifestFileByKind(ctx, d, "Secret")
	if err != nil {
		return nil, err
	}

	for _, content := range parsedContents {
		secret := content.ParsedData.(*v1.Secret)

		if secret.Name == name && secret.Namespace == namespace {
			return Secret{*secret, content}, nil
		}
	}

	// Get the deployed resource
	if clientset == nil {
		return nil, nil
	}

	secret, err := clientset.CoreV1().Secrets(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil && !isNotFoundError(err) {
		return nil, err
	}

	return Secret{*secret, parsedContent{SourceType: "deployed"}}, nil
}

func getSecretResourceContext(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	obj := h.Item.(Secret)

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

func transformSecretTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(Secret)
	return mergeTags(obj.Labels, obj.Annotations), nil
}
