package kubernetes

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
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
		},
		Columns: k8sCommonColumns([]*plugin.Column{
			//// Resource definition specification
			{
				Name:        "spec",
				Description: "Spec describes how the user wants the resources to appear.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "status",
				Description: "Status indicates the actual state of the CustomResourceDefinition.",
				Type:        proto.ColumnType_JSON,
			},
		}),
	}
}

//// HYDRATE FUNCTIONS

func listK8sCustomResourceDefinitions(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	clientset, err := GetNewClientCRD(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listK8sCustomResourceDefinitions", "connection_error", err)
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
			plugin.Logger(ctx).Error("listK8sCustomResourceDefinitions", "api_error", err)
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
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

func getK8sCustomResourceDefinition(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	name := d.EqualsQuals["name"].GetStringValue()

	// check if name is empty
	if name == "" {
		return nil, nil
	}

	clientset, err := GetNewClientCRD(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("getK8sCustomResourceDefinition", "connection_err", err)
		return nil, err
	}

	response, err := clientset.ApiextensionsV1().CustomResourceDefinitions().Get(ctx, name, metav1.GetOptions{})
	if err != nil && !isNotFoundError(err) {
		plugin.Logger(ctx).Error("getK8sCustomResourceDefinition", "api_err", err)
		return nil, err
	}

	return response, nil
}
