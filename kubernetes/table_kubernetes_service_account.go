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

func tableKubernetesServiceAccount(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_service_account",
		Description: "A service account provides an identity for processes that run in a Pod.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "namespace"}),
			Hydrate:    getK8sServiceAccount,
		},
		List: &plugin.ListConfig{
			Hydrate:    listK8sServiceAccounts,
			KeyColumns: getCommonOptionalKeyQuals(),
		},
		// Service Account, is namespaced resource.
		Columns: k8sCommonColumns([]*plugin.Column{
			{
				Name:        "automount_service_account_token",
				Type:        proto.ColumnType_BOOL,
				Description: "Indicates whether pods running as this service account should have an API token automatically mounted. Can be overridden at the pod level.",
			},
			{
				Name:        "image_pull_secrets",
				Type:        proto.ColumnType_JSON,
				Description: "List of references to secrets in the same namespace to use for pulling any images in pods that reference this ServiceAccount. ImagePullSecrets are distinct from Secrets because Secrets can be mounted in the pod, but ImagePullSecrets are only accessed by the kubelet.",
			},
			{
				Name:        "secrets",
				Type:        proto.ColumnType_JSON,
				Description: "Secrets is the list of secrets allowed to be used by pods running using this ServiceAccount.",
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
				Transform:   transform.From(transformServiceAccountTags),
			},
		}),
	}
}

//// HYDRATE FUNCTIONS

func listK8sServiceAccounts(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listK8sServiceAccounts")

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

	var response *v1.ServiceAccountList
	pageLeft := true

	for pageLeft {
		response, err = clientset.CoreV1().ServiceAccounts("").List(ctx, input)
		if err != nil {
			return nil, err
		}

		if response.GetContinue() != "" {
			input.Continue = response.Continue
		} else {
			pageLeft = false
		}

		for _, serviceAccount := range response.Items {
			d.StreamListItem(ctx, serviceAccount)
		}
	}

	return nil, nil
}

func getK8sServiceAccount(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getK8sServiceAccount")

	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		return nil, err
	}

	name := d.EqualsQuals["name"].GetStringValue()
	namespace := d.EqualsQuals["namespace"].GetStringValue()

	// handle empty name and namespace value
	if name == "" || namespace == "" {
		return nil, nil
	}

	serviceAccount, err := clientset.CoreV1().ServiceAccounts(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil && !isNotFoundError(err) {
		logger.Debug("getK8sServiceAccount", "Error", err)
		return nil, err
	}

	return *serviceAccount, nil
}

//// TRANSFORM FUNCTIONS

func transformServiceAccountTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(v1.ServiceAccount)
	return mergeTags(obj.Labels, obj.Annotations), nil
}
