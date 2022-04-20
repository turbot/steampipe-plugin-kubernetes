package kubernetes

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
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

//// HYDRATE FUNCTIONS

func listK8sSecrets(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listK8sSecrets")

	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		return nil, err
	}

	response, err := clientset.CoreV1().Secrets("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, secret := range response.Items {
		d.StreamListItem(ctx, secret)
	}

	return nil, nil
}

func getK8sSecret(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getK8sSecret")

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

	secret, err := clientset.CoreV1().Secrets(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil && !isNotFoundError(err) {
		return nil, err
	}

	return *secret, nil
}

//// TRANSFORM FUNCTIONS

func transformSecretTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(v1.Secret)
	return mergeTags(obj.Labels, obj.Annotations), nil
}
