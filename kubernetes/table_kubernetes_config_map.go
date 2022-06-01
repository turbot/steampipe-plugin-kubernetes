package kubernetes

import (
	"context"
	"strings"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
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

			//// Steampipe Standard Columns
			// {
			// 	Name:        "title",
			// 	Type:        proto.ColumnType_STRING,
			// 	Description: ColumnDescriptionTitle,
			// 	Transform:   transform.FromField("Name"),
			// },
			{
				Name:        "tags",
				Type:        proto.ColumnType_JSON,
				Description: ColumnDescriptionTags,
				Transform:   transform.From(transformConfigMapTags),
			},
		}),
	}
}

//// HYDRATE FUNCTIONS

func listK8sConfigMaps(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listK8sConfigMaps")

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
			d.StreamListItem(ctx, configMap)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

func getK8sConfigMap(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getK8sConfigMap")

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

	configMap, err := clientset.CoreV1().ConfigMaps(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil && !isNotFoundError(err) {
		return nil, err
	}

	return *configMap, nil
}

//// TRANSFORM FUNCTIONS

func transformConfigMapTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(v1.ConfigMap)
	return mergeTags(obj.Labels, obj.Annotations), nil
}
