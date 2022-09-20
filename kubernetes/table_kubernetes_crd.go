package kubernetes

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

func tableKubernetesCRD(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_crd",
		Description: "Cron jobs are useful for creating periodic and recurring tasks, like running backups or sending emails.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.KeyColumnSlice{
				{
					Name: "name", Require: plugin.Required, Operators: []string{"="},
				},
			},
			Hydrate: getK8sCRD,
		},
		List: &plugin.ListConfig{
			Hydrate: listK8sCRDs,
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

func listK8sCRDs(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	clientset, err := GetNewClientCRD(ctx, d)
	if err != nil {
		logger.Error("kubernetes_crd.listK8sCRDs", "connection_error", err)
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

//// Hydrated Function

func getK8sCRD(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	clientset, err := GetNewClientCRD(ctx, d)
	if err != nil {
		logger.Error("kubernetes_crd.getK8sCRD", "connection_error", err)
		return nil, err
	}

	//queryCols := d.KeyColumnQuals

	// version := make(map[string]interface{})

	// versionString := queryCols["versions"].GetJsonbValue()
	// logger.Debug("kubernetes_crd.getK8sCRD", "versions", versionString)

	// if versionString != "" {
	// 	err := json.Unmarshal([]byte(versionString), &version)
	// 	if err != nil {
	// 		plugin.Logger(ctx).Error("kubernetes_crd.getK8sCRD", "unmarshal_error", err)
	// 		return nil, fmt.Errorf("failed to unmarshal versions: %v", err)
	// 	}
	// }

	// if version["version"] == nil {
	// 	panic("Version must to be pass")
	// }

	resourceName := d.KeyColumnQualString("name")
	//resourceVersion := version["version"].(string)

	//plugin.Logger(ctx).Debug("Resource name ======>>>", resourceName)
	//plugin.Logger(ctx).Debug("Version ======>>>>", resourceVersion)

	if resourceName == "" {
		return nil, nil
	}

	input := metav1.GetOptions{
		// ResourceVersion: resourceVersion,
	}

	response, err := clientset.ApiextensionsV1().CustomResourceDefinitions().Get(ctx, resourceName, input)
	plugin.Logger(ctx).Debug("Response Object=====>> ", response.ObjectMeta)

	// plugin.Logger(ctx).Debug("================================================================================================================================================================================================")

	// plugin.Logger(ctx).Debug("response.TypeMeta=====>> ", response.TypeMeta)

	// plugin.Logger(ctx).Debug("================================================================================================================================================================================================")

	// plugin.Logger(ctx).Debug("response.Spec=====>> ", response.Spec)

	// plugin.Logger(ctx).Debug("================================================================================================================================================================================================")

	if err != nil {
		logger.Error("kubernetes_crd.getK8sCRD", "api_error", err)
		return nil, err
	}
	plugin.Logger(ctx).Debug("HERE ======>>>>>")
	return response, nil
	// return v1.CustomResourceDefinition{
	// 	ObjectMeta: response.ObjectMeta,
	// 	TypeMeta:   response.TypeMeta,
	// 	Spec:       response.Spec,
	// 	Status:     response.Status,
	// }, nil
}
