package kubernetes

import (
	"context"

	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

func tableKubernetesCRDResource(ctx context.Context) *plugin.Table {
	//plugin.Logger(ctx).Error("tableKubernetesCRDResource", "resourceName", resourceName)
	//resourceName := ctx.Value("custom_resource_name").(string)
	resourceName := ctx.Value(contextKey("CustomResourceName")).(string)
	tableName := ctx.Value(contextKey("PluginTableName")).(string)
	return &plugin.Table{
		Name:        tableName,
		Description: "Cron jobs are useful for creating periodic and recurring tasks, like running backups or sending emails.",
		List: &plugin.ListConfig{
			ParentHydrate: listK8sCRDs,
			Hydrate:       listK8sCRDResources(resourceName),
		},
		Columns: []*plugin.Column{
			{
				Name:        "kind",
				Type:        proto.ColumnType_STRING,
				Description: "The number of failed finished jobs to retain. Value must be non-negative integer.",
			},
			{
				Name:        "api_resources",
				Type:        proto.ColumnType_JSON,
				Description: "The number of failed finished jobs to retain. Value must be non-negative integer.",
				Transform:   transform.FromField("APIResources"),
			},
		},
	}
}

//// HYDRATE FUNCTIONS

func listK8sCRDResources(resourceName string) func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		version := h.Item.(v1.CustomResourceDefinition).Spec.Versions[0].Name
		groupName := h.Item.(v1.CustomResourceDefinition).Spec.Group

		clientset, err := GetNewClientDynamic(ctx, d)
		if err != nil {
			return nil, err
		}

		resourceId := schema.GroupVersionResource{
			Group:    groupName,
			Version:  version,
			Resource: resourceName,
		}
		plugin.Logger(ctx).Error("tableKubernetesCRDResource", "resourceId", resourceId)
		response, err := clientset.Resource(resourceId).List(ctx, metav1.ListOptions{})
		if err != nil {
			return nil, err
		}
		for _, crd := range response.Items {
			d.StreamListItem(ctx, crd)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		return nil, nil
	}

}
