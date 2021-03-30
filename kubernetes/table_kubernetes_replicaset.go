package kubernetes

import (
	"context"

	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableKubernetesReplicaSet(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_replicaset",
		Description: "Kubernetes ReplicaSet ensures that a specified number of pod replicas are running at any given time.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "namespace"}),
			Hydrate:    getK8sReplicaSet,
		},
		List: &plugin.ListConfig{
			Hydrate: listK8sReplicaSets,
		},
		Columns: k8sCommonColumns([]*plugin.Column{
			//// Steampipe Standard Columns
			{
				Name:        "title",
				Type:        proto.ColumnType_STRING,
				Description: ColumnDescriptionTitle,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "akas",
				Type:        proto.ColumnType_JSON,
				Description: ColumnDescriptionAkas,
				Transform:   transform.FromField("UID").Transform(ensureStringArray),
			},
			{
				Name:        "tags",
				Type:        proto.ColumnType_JSON,
				Description: ColumnDescriptionTags,
				Transform:   transform.From(transformReplicaSetTags),
			},
		}),
	}
}

//// HYDRATE FUNCTIONS

func listK8sReplicaSets(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listK8sReplicaSets")

	clientset, err := GetNewClientset(ctx, d.ConnectionManager)
	if err != nil {
		return nil, err
	}

	replicaSets, err := clientset.AppsV1().ReplicaSets("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, item := range replicaSets.Items {
		d.StreamListItem(ctx, item)
	}

	return nil, nil
}

func getK8sReplicaSet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getK8sReplicaSet")

	clientset, err := GetNewClientset(ctx, d.ConnectionManager)
	if err != nil {
		return nil, err
	}

	name := d.KeyColumnQuals["name"].GetStringValue()
	namespace := d.KeyColumnQuals["namespace"].GetStringValue()

	rs, err := clientset.AppsV1().ReplicaSets(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil && !isNotFoundError(err) {
		return nil, err
	}

	return rs, nil
}

//// TRANSFORM FUNCTIONS

func transformReplicaSetTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(v1.ReplicaSet)
	return mergeTags(obj.Labels, obj.Annotations), nil
}
