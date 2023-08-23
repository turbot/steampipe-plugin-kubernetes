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
			{
				Name:        "context_name",
				Type:        proto.ColumnType_STRING,
				Description: "Kubectl config context name.",
				Hydrate:     getServiceAccountResourceContext,
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
				Transform:   transform.From(transformServiceAccountTags),
			},
		}),
	}
}

type ServiceAccount struct {
	v1.ServiceAccount
	parsedContent
}

//// HYDRATE FUNCTIONS

func listK8sServiceAccounts(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listK8sServiceAccounts")

	// Get the client for querying the K8s APIs for the provided context.
	// If the connection is configured for the manifest files, the client will return nil.
	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		return nil, err
	}

	// Check for manifest files
	parsedContents, err := fetchResourceFromManifestFileByKind(ctx, d, "ServiceAccount")
	if err != nil {
		return nil, err
	}

	for _, content := range parsedContents {
		serviceAccount := content.ParsedData.(*v1.ServiceAccount)

		d.StreamListItem(ctx, ServiceAccount{*serviceAccount, content})

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
			d.StreamListItem(ctx, ServiceAccount{serviceAccount, parsedContent{SourceType: "deployed"}})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

func getK8sServiceAccount(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getK8sServiceAccount")

	// Get the client for querying the K8s APIs for the provided context.
	// If the connection is configured for the manifest files, the client will return nil.
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

	// Get the manifest resource
	parsedContents, err := fetchResourceFromManifestFileByKind(ctx, d, "ServiceAccount")
	if err != nil {
		return nil, err
	}

	for _, content := range parsedContents {
		serviceAccount := content.ParsedData.(*v1.ServiceAccount)

		if serviceAccount.Name == name && serviceAccount.Namespace == namespace {
			return ServiceAccount{*serviceAccount, content}, nil
		}
	}

	// Get the deployed resource
	if clientset == nil {
		return nil, nil
	}

	serviceAccount, err := clientset.CoreV1().ServiceAccounts(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil && !isNotFoundError(err) {
		logger.Debug("getK8sServiceAccount", "Error", err)
		return nil, err
	}

	return ServiceAccount{*serviceAccount, parsedContent{SourceType: "deployed"}}, nil
}

func getServiceAccountResourceContext(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	obj := h.Item.(ServiceAccount)

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

func transformServiceAccountTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(ServiceAccount)
	return mergeTags(obj.Labels, obj.Annotations), nil
}
