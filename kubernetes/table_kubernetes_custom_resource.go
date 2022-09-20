package kubernetes

import (
	"context"
	"fmt"
	"strings"
	"sync"

	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"

	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
)

func tableKubernetesCustomResource(ctx context.Context, p *plugin.Plugin) *plugin.Table {
	resourceName := ctx.Value(contextKey("CustomResourceName")).(string)
	version := ctx.Value(contextKey("ActiveVersion")).(string)
	resourceId := schema.GroupVersionResource{
		Group:    strings.Replace(resourceName, strings.Split(resourceName, ".")[0]+".", "", 1),
		Version:  version,
		Resource: strings.Split(resourceName, ".")[0],
	}

	clientset, err := GetNewClientDynamic(ctx, p.ConnectionManager, p.Connection)
	if err != nil {
		plugin.Logger(ctx).Error("GetNewClientDynamic", "connection_error", err)
		return nil
	}

	response, err := clientset.Resource(resourceId).List(ctx, metav1.ListOptions{})
	if err != nil {
		plugin.Logger(ctx).Error("GetNewClientDynamic", "connection_error", err)
		return nil
	}
	return &plugin.Table{
		Name:        fmt.Sprintf("\"" + resourceName + "\""),
		Description: fmt.Sprintf("Represents CRD object %s.", resourceName),
		List: &plugin.ListConfig{
			ParentHydrate: listK8sCustomResourceDefinitions,
			Hydrate:       listK8sCustomResources(resourceName),
		},
		Columns: k8sCRDResourceCommonColumns([]*plugin.Column{}),
	}
}

type CRDResourceInfo struct {
	Kind        string
	APIVersion  string
	Name        string
	Namespace   string
	Annotations interface{}
	Spec        interface{}
}

//// HYDRATE FUNCTIONS

func listK8sCustomResources(resourceName string) func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		version := h.Item.(v1.CustomResourceDefinition).Spec.Versions
		groupName := h.Item.(v1.CustomResourceDefinition).Spec.Group
		names := h.Item.(v1.CustomResourceDefinition).Spec.Names.Plural

		// check if the resourceNames is matching with the definition
		if names != resourceName {
			return nil, nil
		}

		var wg sync.WaitGroup
		errorCh := make(chan error, len(version))
		for _, v := range version {
			wg.Add(1)
			go getCustomResourceAsync(ctx, d, groupName, resourceName, v.Name, clientset, &wg, errorCh)
		}
		// wait for all inline policies to be processed
		wg.Wait()

		// NOTE: close channel before ranging over results
		close(errorCh)

		for err := range errorCh {
			// return the first error
			plugin.Logger(ctx).Error("getCustomResourceAsync", "channel_error", err)
			return nil, err
		}
		return nil, nil
	}

}

func getCustomResourceAsync(ctx context.Context, d *plugin.QueryData, groupName string, resourceName string, version string, clientset dynamic.Interface, wg *sync.WaitGroup, errorCh chan error) {
	defer wg.Done()

	err := getCustomResources(ctx, d, groupName, resourceName, version, clientset)
	if err != nil {
		errorCh <- err
	}
}

func getCustomResources(ctx context.Context, d *plugin.QueryData, groupName string, resourceName string, version string, clientset dynamic.Interface) error {
	resourceId := schema.GroupVersionResource{
		Group:    groupName,
		Version:  version,
		Resource: resourceName,
	}

	for _, crd := range response.Items {
		plugin.Logger(ctx).Error("getCustomResources", "crd", crd)
		ob := crd.Object
		d.StreamListItem(ctx, &CRDResourceInfo{
			Kind:        ob["kind"].(string),
			APIVersion:  ob["apiVersion"].(string),
			Name:        ob["metadata"].(map[string]interface{})["name"].(string),
			Annotations: ob["metadata"].(map[string]interface{})["annotations"],
			Namespace:   ob["metadata"].(map[string]interface{})["namespace"].(string),
			Spec:        ob["spec"],
		})

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil
		}
	}
	return nil
}