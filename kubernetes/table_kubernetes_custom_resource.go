package kubernetes

// Uncomment once aggregator functionality available with dynamic tables

// func tableKubernetesCustomResource(ctx context.Context) *plugin.Table {
// 	crdName := ctx.Value(contextKey("CRDName")).(string)
// 	resourceName := ctx.Value(contextKey("CustomResourceName")).(string)
// 	groupName := ctx.Value(contextKey("GroupName")).(string)
// 	activeVersion := ctx.Value(contextKey("ActiveVersion")).(string)
// 	versionSchema := ctx.Value(contextKey("VersionSchema"))
// 	return &plugin.Table{
// 		Name:        crdName,
// 		Description: fmt.Sprintf("Custom resource for %s.", crdName),
// 		List: &plugin.ListConfig{
// 			Hydrate: listK8sCustomResources(ctx, crdName, resourceName, groupName, activeVersion),
// 		},
// 		Columns: k8sCRDResourceCommonColumns(getCustomResourcesDynamicColumns(ctx, versionSchema)),
// 	}
// }

// func getCustomResourcesDynamicColumns(ctx context.Context, versionSchema interface{}) []*plugin.Column {
// 	var columns []*plugin.Column

// 	schema := versionSchema.(v1.JSONSchemaProps)
// 	for k, v := range schema.Properties {
// 		column := &plugin.Column{
// 			Name:        k,
// 			Description: v.Description,
// 			Transform:   transform.FromP(extractSpecProperty, k),
// 		}
// 		switch v.Type {
// 		case "string":
// 			column.Type = proto.ColumnType_STRING
// 		case "integer":
// 			column.Type = proto.ColumnType_INT
// 		case "boolean":
// 			column.Type = proto.ColumnType_BOOL
// 		case "date", "dateTime":
// 			column.Type = proto.ColumnType_TIMESTAMP
// 		case "double":
// 			column.Type = proto.ColumnType_DOUBLE
// 		default:
// 			column.Type = proto.ColumnType_JSON
// 		}

// 		columns = append(columns, column)
// 	}

// 	return columns
// }

// type CRDResourceInfo struct {
// 	Name              interface{}
// 	UID               interface{}
// 	CreationTimestamp interface{}
// 	Kind              interface{}
// 	APIVersion        interface{}
// 	Namespace         interface{}
// 	Annotations       interface{}
// 	Spec              interface{}
// 	Labels            interface{}
// }

// //// HYDRATE FUNCTIONS

// func listK8sCustomResources(ctx context.Context, crdName string, resourceName string, groupName string, activeVersion string) func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
// 	return func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
// 		clientset, err := GetNewClientDynamic(ctx, d)
// 		if err != nil {
// 			return nil, err
// 		}

// 		resourceId := schema.GroupVersionResource{
// 			Group:    groupName,
// 			Version:  activeVersion,
// 			Resource: resourceName,
// 		}

// 		response, err := clientset.Resource(resourceId).List(ctx, metav1.ListOptions{})
// 		if err != nil {
// 			if strings.Contains(err.Error(), "could not find the requested resource") {
// 				return nil, nil
// 			}
// 			return nil, err
// 		}

// 		for _, crd := range response.Items {
// 			data := crd.Object
// 			d.StreamListItem(ctx, &CRDResourceInfo{
// 				Name:              crd.GetName(),
// 				UID:               crd.GetUID(),
// 				APIVersion:        crd.GetAPIVersion(),
// 				Kind:              crd.GetKind(),
// 				Namespace:         crd.GetNamespace(),
// 				CreationTimestamp: crd.GetCreationTimestamp(),
// 				Labels:            crd.GetLabels(),
// 				Spec:              data["spec"],
// 			})

// 			// Context can be cancelled due to manual cancellation or the limit has been hit
// 			if d.RowsRemaining(ctx) == 0 {
// 				return nil, nil
// 			}
// 		}

// 		return nil, nil
// 	}

// }

// func extractSpecProperty(_ context.Context, d *transform.TransformData) (interface{}, error) {
// 	ob := d.HydrateItem.(*CRDResourceInfo).Spec
// 	param := d.Param.(string)
// 	spec := ob.(map[string]interface{})
// 	if spec[param] != nil {
// 		return spec[param], nil
// 	}

// 	return nil, nil
// }
