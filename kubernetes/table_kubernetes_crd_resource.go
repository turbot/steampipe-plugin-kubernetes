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

func tableKubernetesCRDResource(ctx context.Context) *plugin.Table {
	resourceName := ctx.Value(contextKey("CustomResourceName")).(string)
	tableName := ctx.Value(contextKey("PluginTableName")).(string)
	return &plugin.Table{
		Name:        tableName,
		Description: fmt.Sprintf("Represents CRD object %s.", resourceName),
		List: &plugin.ListConfig{
			ParentHydrate: listK8sCRDs,
			Hydrate:       listK8sCRDResources(resourceName),
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

func listK8sCRDResources(resourceName string) func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		version := h.Item.(v1.CustomResourceDefinition).Spec.Versions
		groupName := h.Item.(v1.CustomResourceDefinition).Spec.Group
		names := h.Item.(v1.CustomResourceDefinition).Spec.Names.Plural
		plugin.Logger(ctx).Error("tableKubernetesCRDResource", "resourceName", resourceName, "names", names)
		if names != resourceName {
			return nil, nil
		}

		clientset, err := GetNewClientDynamic(ctx, d)
		if err != nil {
			return nil, err
		}

		var wg sync.WaitGroup
		errorCh := make(chan error, len(version))
		for _, v := range version {
			wg.Add(1)
			go getCRDResourceAsync(ctx, d, groupName, resourceName, v.Name, clientset, &wg, errorCh)
		}
		// wait for all inline policies to be processed
		wg.Wait()

		// NOTE: close channel before ranging over results
		close(errorCh)

		for err := range errorCh {
			// return the first error
			plugin.Logger(ctx).Error("getCRDResourceAsync", "channel_error", err)
			return nil, err
		}
		return nil, nil
	}

}

func getCRDResourceAsync(ctx context.Context, d *plugin.QueryData, groupName string, resourceName string, version string, clientset dynamic.Interface, wg *sync.WaitGroup, errorCh chan error) {
	defer wg.Done()

	err := getCRDResource(ctx, d, groupName, resourceName, version, clientset)
	if err != nil {
		errorCh <- err
	}
}

func getCRDResource(ctx context.Context, d *plugin.QueryData, groupName string, resourceName string, version string, clientset dynamic.Interface) error {
	resourceId := schema.GroupVersionResource{
		Group:    groupName,
		Version:  version,
		Resource: resourceName,
	}
	plugin.Logger(ctx).Error("tableKubernetesCRDResource", "resourceId", resourceId)
	response, err := clientset.Resource(resourceId).List(ctx, metav1.ListOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "could not find the requested resource") {
			return nil
		}
		return err
	}

	var annotations interface{}

	for _, crd := range response.Items {
		plugin.Logger(ctx).Error("tableKubernetesCRDResource", "crd", crd)
		ob := crd.Object
		for _, v := range ob["metadata"].(map[string]interface{})["annotations"].(map[string]interface{}) {
			annotations = strings.TrimLeft(strings.TrimRight(v.(string), "\""), "\"")
		}
		d.StreamListItem(ctx, &CRDResourceInfo{
			Kind:        ob["kind"].(string),
			APIVersion:  ob["apiVersion"].(string),
			Name:        ob["metadata"].(map[string]interface{})["name"].(string),
			Annotations: annotations,
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
