package kubernetes

import (
	"context"
	"fmt"
	"strings"

	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/turbot/steampipe-plugin-sdk/v3/connection"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

func tableKubernetesCustomResource(ctx context.Context, p *plugin.Plugin) *plugin.Table {
	crdName := ctx.Value(contextKey("CRDName")).(string)
	resourceName := ctx.Value(contextKey("CustomResourceName")).(string)
	groupName := ctx.Value(contextKey("GroupName")).(string)
	activeVersion := ctx.Value(contextKey("ActiveVersion")).(string)
	versionSchema := ctx.Value(contextKey("VersionSchema"))
	return &plugin.Table{
		Name:        crdName,
		Description: fmt.Sprintf("Represents Custom resource %s.", crdName),
		List: &plugin.ListConfig{
			Hydrate: listK8sCustomResources(ctx, crdName, resourceName, groupName, activeVersion),
		},
		Columns: k8sCRDResourceCommonColumns(getCustomResourcesDynamincColumns(ctx, p.ConnectionManager, p.Connection, versionSchema)),
	}
}

func getCustomResourcesDynamincColumns(ctx context.Context, cm *connection.Manager, c *plugin.Connection, versionSchema interface{}) []*plugin.Column {
	var columns []*plugin.Column

	schema := versionSchema.(v1.JSONSchemaProps)
	for k, v := range schema.Properties {
		column := &plugin.Column{
			Name:        "spec_" + k,
			Description: v.Description,
			Transform:   transform.FromP(extractSpecProperty, k),
		}
		switch v.Type {
		case "string":
			column.Type = proto.ColumnType_STRING
		case "integer":
			column.Type = proto.ColumnType_INT
		case "boolean":
			column.Type = proto.ColumnType_BOOL
		case "date", "dateTime":
			column.Type = proto.ColumnType_TIMESTAMP
		case "double":
			column.Type = proto.ColumnType_DOUBLE
		default:
			column.Type = proto.ColumnType_JSON
		}

		columns = append(columns, column)
	}

	return columns
}

type CRDResourceInfo struct {
	Name        interface{}
	Kind        interface{}
	APIVersion  interface{}
	Namespace   interface{}
	Annotations interface{}
	Spec        interface{}
}

//// HYDRATE FUNCTIONS

func listK8sCustomResources(ctx context.Context, crdName string, resourceName string, groupName string, activeVersion string) func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		clientset, err := GetNewClientDynamic(ctx, d)
		if err != nil {
			return nil, err
		}

		resourceId := schema.GroupVersionResource{
			Group:    groupName,
			Version:  activeVersion,
			Resource: resourceName,
		}

		response, err := clientset.Resource(resourceId).List(ctx, metav1.ListOptions{})
		if err != nil {
			if strings.Contains(err.Error(), "could not find the requested resource") {
				return nil, nil
			}
			return nil, err
		}
		for _, crd := range response.Items {
			data := crd.Object
			var annotations interface{}
			for _, v := range crd.GetAnnotations(){
			annotations = strings.TrimLeft(strings.TrimRight(v, "\""), "\"")
		}
			d.StreamListItem(ctx, &CRDResourceInfo{
				Name:        crd.GetName(),
				APIVersion:  crd.GetAPIVersion(),
				Kind:        crd.GetKind(),
				Namespace:   crd.GetNamespace(),
				Annotations: annotations,
				Spec:        data["spec"],
			})

			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		return nil, nil
	}

}

func extractSpecProperty(_ context.Context, d *transform.TransformData) (interface{}, error) {
	ob := d.HydrateItem.(*CRDResourceInfo).Spec
	param := d.Param.(string)
	spec := ob.(map[string]interface{})
	if spec[param] != nil {
		return spec[param], nil
	}

	return nil, nil
}
