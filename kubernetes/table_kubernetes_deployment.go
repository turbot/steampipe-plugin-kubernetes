package kubernetes

import (
	"context"

	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableKubernetesDeployment(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_deployment",
		Description: "Kubernetes Deployment enables declarative updates for Pods and ReplicaSets.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "namespace"}),
			Hydrate:    getK8sDeployment,
		},
		List: &plugin.ListConfig{
			Hydrate: listK8sDeployments,
		},
		Columns: k8sCommonColumns([]*plugin.Column{

			//// Spec Columns
			// To do - add Spec Columns...

			//// Status Columns
			// To do - add Status Columns...

			//// Steampipe Standard Columns
			{
				Name:        "title",
				Type:        proto.ColumnType_STRING,
				Description: ColumnDescriptionTitle,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "tags",
				Type:        proto.ColumnType_JSON,
				Description: ColumnDescriptionTags,
				Transform:   transform.From(transformDeploymentTags),
			},
		}),
	}
}

//// HYDRATE FUNCTIONS

func listK8sDeployments(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listK8sDeployments")

	clientset, err := GetNewClientset(ctx, d.ConnectionManager)
	if err != nil {
		return nil, err
	}

	deployments, err := clientset.AppsV1().Deployments("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, item := range deployments.Items {
		d.StreamListItem(ctx, item)
	}

	return nil, nil
}

func getK8sDeployment(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getK8sDeployment")

	clientset, err := GetNewClientset(ctx, d.ConnectionManager)
	if err != nil {
		return nil, err
	}

	name := d.KeyColumnQuals["name"].GetStringValue()
	namespace := d.KeyColumnQuals["namespace"].GetStringValue()

	deployment, err := clientset.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil && !isNotFoundError(err) {
		return nil, err
	}

	return *deployment, nil
}

//// TRANSFORM FUNCTIONS

func transformDeploymentTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(v1.Deployment)
	return mergeTags(obj.Labels, obj.Annotations), nil
}
