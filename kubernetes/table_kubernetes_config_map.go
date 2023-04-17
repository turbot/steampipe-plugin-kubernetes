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

func tableKubernetesConfigMap(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_config_map",
		Description: "Config Map can be used to store fine-grained information like individual properties or coarse-grained information like entire config files or JSON blobs.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "namespace"}),
			Hydrate:    getK8sConfigMap,
		},
		List: &plugin.ListConfig{
			Hydrate:    listK8sConfigMaps,
			KeyColumns: getCommonOptionalKeyQuals(),
		},
		// ClusterRole, is a non-namespaced resource.
		Columns: k8sCommonColumns([]*plugin.Column{
			{
				Name:        "immutable",
				Type:        proto.ColumnType_BOOL,
				Description: "If set to true, ensures that data stored in the ConfigMap cannot be updated (only object metadata can be modified). If not set to true, the field can be modified at any time. Defaulted to nil.",
			},
			{
				Name:        "data",
				Type:        proto.ColumnType_JSON,
				Description: "Contains the configuration data.",
			},
			{
				Name:        "binary_data",
				Type:        proto.ColumnType_JSON,
				Description: "Contains the configuration binary data.",
			},
			{
				Name:        "source_type",
				Type:        proto.ColumnType_STRING,
				Description: "The source of the resource. Possible values are: deployed and manifest. If the resource is fetched from the spec file the value will be manifest.",
				Transform:   transform.From(configMapResourceSourceType),
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
				Transform:   transform.From(transformConfigMapTags),
			},
		}),
	}
}

type ConfigMap struct {
	v1.ConfigMap
	Path      string
	StartLine int
	EndLine   int
}

//// HYDRATE FUNCTIONS

func listK8sConfigMaps(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listK8sConfigMaps")

	// Get the client for querying the K8s APIs for the provided context.
	// If the connection is configured for the manifest files, the client will return nil.
	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		return nil, err
	}

	//
	// Check for manifest files
	//
	parsedContents, err := fetchResourceFromManifestFileByKind(ctx, d, "ConfigMap")
	if err != nil {
		return nil, err
	}

	for _, content := range parsedContents {
		configMap := content.Data.(*v1.ConfigMap)

		d.StreamListItem(ctx, ConfigMap{*configMap, content.Path, content.StartLine, content.EndLine})

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

	commonFieldSelectorValue := getCommonOptionalKeyQualsValueForFieldSelector(d)

	if len(commonFieldSelectorValue) > 0 {
		input.FieldSelector = strings.Join(commonFieldSelectorValue, ",")
	}

	var response *v1.ConfigMapList
	pageLeft := true

	for pageLeft {
		response, err = clientset.CoreV1().ConfigMaps("").List(ctx, input)
		if err != nil {
			return nil, err
		}

		if response.GetContinue() != "" {
			input.Continue = response.Continue
		} else {
			pageLeft = false
		}

		for _, configMap := range response.Items {
			d.StreamListItem(ctx, ConfigMap{configMap, "", 0, 0})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

func getK8sConfigMap(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getK8sConfigMap")

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

	//
	// Get the manifest resource
	//
	parsedContents, err := fetchResourceFromManifestFileByKind(ctx, d, "ConfigMap")
	if err != nil {
		return nil, err
	}

	for _, content := range parsedContents {
		configMap := content.Data.(*v1.ConfigMap)

		if configMap.Name == name && configMap.Namespace == namespace {
			return ConfigMap{*configMap, content.Path, content.StartLine, content.EndLine}, nil
		}
	}

	//
	// Get the deployed resource
	//
	if clientset == nil {
		return nil, nil
	}

	configMap, err := clientset.CoreV1().ConfigMaps(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil && !isNotFoundError(err) {
		return nil, err
	}

	return ConfigMap{*configMap, "", 0, 0}, nil
}

//// TRANSFORM FUNCTIONS

func transformConfigMapTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(ConfigMap)
	return mergeTags(obj.Labels, obj.Annotations), nil
}

func configMapResourceSourceType(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(ConfigMap)

	if obj.Path != "" {
		return "manifest", nil
	}
	return "deployed", nil
}
