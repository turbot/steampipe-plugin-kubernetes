package kubernetes

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

func tableKubernetesCustomResourceDefinition(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_custom_resource_definition",
		Description: "Kubernetes Custom Resource Definition.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getK8sCustomResourceDefinition,
		},
		List: &plugin.ListConfig{
			Hydrate: listK8sCustomResourceDefinitions,
			//KeyColumns: getCommonOptionalKeyQuals(),
		},
		Columns: []*plugin.Column{
			//// Resource definition specification
			{
				Name:        "name",
				Description: "Group is the API group of the defined custom resource.",
				Type:        proto.ColumnType_STRING,
				// Transform:   transform.FromField("Spec.Group"),
			},
			{
				Name:        "plural",
				Description: "Plural is the plural name of the resource to serve.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Spec.Names.Plural"),
			},
			{
				Name:        "custom_resource_conversion_strategy",
				Description: "Strategy specifies how custom resources are converted between versions.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Spec.Conversion.Strategy"),
			},
			{
				Name:        "preserve_unknown_fields",
				Description: "PreserveUnknownFields indicates that object fields which are not specified in the OpenAPI schema should be preserved when persisting to storage.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Spec.PreserveUnknownFields"),
			},
			{
				Name:        "scope",
				Description: "Group is the API group of the defined custom resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Spec.Scope"),
			},
			{
				Name:        "custom_resource_conversion_webhook",
				Description: "webhook describes how to call the conversion webhook. Required when `strategy` is set to `Webhook`.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Spec.Conversion.Webhook"),
			},
			{
				Name:        "names",
				Description: "Names specify the resource and kind names for the custom resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Spec.Names"),
			},
			{
				Name:        "versions",
				Description: "Versions is the list of all API versions of the defined custom resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Spec.Versions"),
			},
		},
	}
}

//// HYDRATE FUNCTIONS

func listK8sCustomResourceDefinitions(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listK8sCustomResourceDefinitions")

	clientset, err := GetNewClientCRD(ctx, d)
	if err != nil {
		logger.Error("kubernetes_crd.listK8sCustomResourceDefinitions", "connection_error", err)
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

	pageLeft := true
	for pageLeft {
		response, err := clientset.ApiextensionsV1().CustomResourceDefinitions().List(ctx, input)
		if err != nil {
			logger.Error("kubernetes_crd.listK8sCRDs", "api_error", err)
			return nil, err
		}

		if response.GetContinue() != "" {
			input.Continue = response.Continue
		} else {
			pageLeft = false
		}

		for _, crd := range response.Items {
			d.StreamListItem(ctx, crd)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

func getK8sCustomResourceDefinition(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getK8sCustomResourceDefinition")

	clientset, err := GetNewClientCRD(ctx, d)
	if err != nil {
		return nil, err
	}
	name := d.KeyColumnQuals["name"].GetStringValue()

	response, err := clientset.ApiextensionsV1().CustomResourceDefinitions().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		logger.Error("listK8sCustomResourceDefinitions", "list_err", err)
		return nil, err
	}

	return response, nil
}
