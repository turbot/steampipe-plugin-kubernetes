package kubernetes

import (
	"context"
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func tableKubernetesCustomResource(ctx context.Context) *plugin.Table {
	crdName := ctx.Value(contextKey("CRDName")).(string)
	resourceName := ctx.Value(contextKey("CustomResourceName")).(string)
	resourceNameSingular := ctx.Value(contextKey("CustomResourceNameSingular")).(string)
	groupName := ctx.Value(contextKey("GroupName")).(string)
	activeVersion := ctx.Value(contextKey("ActiveVersion")).(string)
	versionSchemaSpec := ctx.Value(contextKey("VersionSchemaSpec"))
	versionSchemaStatus := ctx.Value(contextKey("VersionSchemaStatus"))
	tableName := ctx.Value(contextKey("TableName")).(string)

	var description string
	if ctx.Value(contextKey("VersionSchemaDescription")) == nil {
		description = "Custom resource for " + crdName + "."
	} else {
		description = ctx.Value(contextKey("VersionSchemaDescription")).(string) + " Custom resource for " + crdName + "."
	}
	return &plugin.Table{
		Name:        tableName,
		Description: description,
		List: &plugin.ListConfig{
			Hydrate: listK8sCustomResources(ctx, crdName, resourceName, resourceNameSingular, groupName, activeVersion),
		},
		Columns: k8sCRDResourceCommonColumns(getCustomResourcesDynamicColumns(ctx, versionSchemaSpec, versionSchemaStatus)),
	}
}

func getCustomResourcesDynamicColumns(ctx context.Context, versionSchemaSpec interface{}, versionSchemaStatus interface{}) []*plugin.Column {
	columns := []*plugin.Column{}

	// default metadata columns
	allColumns := []string{"name", "uid", "kind", "api_version", "namespace", "creation_timestamp", "labels", "start_line", "end_line", "path", "source_type", "annotations", "context_name"}

	// add the spec columns
	schemaSpec := versionSchemaSpec.(v1.JSONSchemaProps)
	for k, v := range schemaSpec.Properties {
		flag := 0
		for _, specColumn := range allColumns {
			if specColumn == strcase.ToSnake(k) {
				flag = 1
				column := &plugin.Column{
					Name:        "spec_" + strcase.ToSnake(k),
					Description: v.Description,
					Transform:   transform.FromP(extractSpecProperty, k),
				}
				setDynamicColumns(v, column)
				columns = append(columns, column)
			}
		}
		if flag == 0 {
			column := &plugin.Column{
				Name:        strcase.ToSnake(k),
				Description: v.Description,
				Transform:   transform.FromP(extractSpecProperty, k),
			}
			allColumns = append(allColumns, strcase.ToSnake(k))
			setDynamicColumns(v, column)
			columns = append(columns, column)
		}
	}

	// add the status columns
	schemaStatus := versionSchemaStatus.(v1.JSONSchemaProps)
	for k, v := range schemaStatus.Properties {
		flag := 0
		for _, statusColumn := range allColumns {
			if statusColumn == strcase.ToSnake(k) {
				flag = 1
				column := &plugin.Column{
					Name:        "status_" + strcase.ToSnake(k),
					Description: v.Description,
					Transform:   transform.FromP(extractStatusProperty, k),
				}
				setDynamicColumns(v, column)
				columns = append(columns, column)
			}
		}
		if flag == 0 {
			column := &plugin.Column{
				Name:        strcase.ToSnake(k),
				Description: v.Description,
				Transform:   transform.FromP(extractStatusProperty, k),
			}
			setDynamicColumns(v, column)
			columns = append(columns, column)
		}
	}

	return columns
}

type CRDResourceInfo struct {
	Name              interface{}
	UID               interface{}
	CreationTimestamp interface{}
	Kind              interface{}
	APIVersion        interface{}
	Namespace         interface{}
	Annotations       interface{}
	Spec              interface{}
	Labels            interface{}
	Status            interface{}
	Report            interface{}
	Path              string
	StartLine         int
	EndLine           int
	SourceType        string
	ContextName       string
}

// //// HYDRATE FUNCTIONS

func listK8sCustomResources(ctx context.Context, crdName string, resourceName string, resourceNameSingular string, groupName string, activeVersion string) func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		clientset, err := GetNewClientDynamic(ctx, d)
		if err != nil {
			return nil, err
		}

		//
		// Check for manifest files
		//

		// In general, the kind of the custom resource must be same as the singular name defined in the CRD
		// Convert the singular name into title format, e.g. if the name is `certificate`, the custom resource kind must be `Certificate`
		caser := cases.Title(language.English)
		parsedContents, err := fetchResourceFromManifestFileByKind(ctx, d, caser.String(resourceNameSingular))
		if err != nil {
			return nil, err
		}

		for _, content := range parsedContents {
			deployment := content.ParsedData.(*unstructured.Unstructured)

			// Also, the apiVersion of the custom resource must be in format of <groupName in CRD>/<spec version in CRD>
			if !(deployment.GetAPIVersion() == fmt.Sprintf("%s/%s", groupName, activeVersion)) {
				continue
			}

			data := deployment.Object
			d.StreamListItem(ctx, &CRDResourceInfo{
				Name:              deployment.GetName(),
				UID:               deployment.GetUID(),
				APIVersion:        deployment.GetAPIVersion(),
				Kind:              deployment.GetKind(),
				Namespace:         deployment.GetNamespace(),
				CreationTimestamp: deployment.GetCreationTimestamp(),
				Annotations:       deployment.GetAnnotations(),
				Labels:            deployment.GetLabels(),
				Spec:              data["spec"],
				Status:            data["status"],
				Report:            data["report"],
				StartLine:         content.StartLine,
				EndLine:           content.EndLine,
				Path:              content.Path,
				SourceType:        content.SourceType,
			})
		}

		//
		// Check for deployed resources
		//
		if clientset == nil {
			return nil, nil
		}

		resourceId := schema.GroupVersionResource{
			Group:    groupName,
			Version:  activeVersion,
			Resource: resourceName,
		}

		response, err := clientset.Resource(resourceId).List(ctx, metav1.ListOptions{})
		if err != nil {
			// Handle not found error code
			if strings.Contains(err.Error(), "could not find the requested resource") {
				return nil, nil
			}
			return nil, err
		}

		currentContext := getCurrentContext(ctx, d, nil)

		for _, crd := range response.Items {
			data := crd.Object
			d.StreamListItem(ctx, &CRDResourceInfo{
				Name:              crd.GetName(),
				UID:               crd.GetUID(),
				APIVersion:        crd.GetAPIVersion(),
				Kind:              crd.GetKind(),
				Namespace:         crd.GetNamespace(),
				Annotations:       crd.GetAnnotations(),
				CreationTimestamp: crd.GetCreationTimestamp(),
				Labels:            crd.GetLabels(),
				Spec:              data["spec"],
				Status:            data["status"],
				Report:            data["report"],
				SourceType:        "deployed",
				ContextName:       currentContext.(string),
			})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		return nil, nil
	}

}

func extractSpecProperty(_ context.Context, d *transform.TransformData) (interface{}, error) {
	ob := d.HydrateItem.(*CRDResourceInfo).Spec
	if ob == nil {
		return nil, nil
	}
	param := d.Param.(string)
	spec := ob.(map[string]interface{})
	if spec[param] != nil {
		return spec[param], nil
	}

	return nil, nil
}

func extractStatusProperty(_ context.Context, d *transform.TransformData) (interface{}, error) {
	ob := d.HydrateItem.(*CRDResourceInfo).Status
	if ob == nil {
		return nil, nil
	}
	param := d.Param.(string)
	status := ob.(map[string]interface{})
	if status[param] != nil {
		return status[param], nil
	}
	report := ob.(map[string]interface{})
	if report[param] != nil {
		return report[param], nil
	}

	return nil, nil
}

func setDynamicColumns(v v1.JSONSchemaProps, column *plugin.Column) {
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
}

func getCurrentContext(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) any {
	currentContext, err := getKubectlContext(ctx, d, nil)
	if err != nil {
		return nil
	}
	return currentContext
}
