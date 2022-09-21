package kubernetes

import (
	"context"
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
)

func tableKubernetesCustomResource(ctx context.Context) *plugin.Table {
	crdName := ctx.Value(contextKey("CRDName")).(string)
	resourceName := ctx.Value(contextKey("CustomResourceName")).(string)
	groupName := ctx.Value(contextKey("GroupName")).(string)
	activeVersion := ctx.Value(contextKey("ActiveVersion")).(string)
	return &plugin.Table{
		Name:        crdName,
		Description: fmt.Sprintf("Represents Custom resource %s.", crdName),
		List: &plugin.ListConfig{
			ParentHydrate: listK8sCustomResourceDefinitions,
			Hydrate:       listK8sCustomResources(resourceName, groupName, activeVersion),
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

func listK8sCustomResources(resourceName string, groupName string, activeVersion string) func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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
			d.StreamListItem(ctx, crd)
		}
		return nil, nil
	}

}
