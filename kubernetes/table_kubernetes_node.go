package kubernetes

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableKubernetesNode(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_node",
		Description: "Kubernetes Node is a worker node in Kubernetes.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getK8sNode,
		},
		List: &plugin.ListConfig{
			Hydrate: listK8sNodes,
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
				Transform:   transform.From(transformNodeTags),
			},
		}),
	}
}

//// HYDRATE FUNCTIONS

func listK8sNodes(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listK8sNodes")

	clientset, err := GetNewClientset(ctx, d.ConnectionManager)
	if err != nil {
		return nil, err
	}

	nodes, err := clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, pod := range nodes.Items {
		d.StreamListItem(ctx, pod)
	}

	return nil, nil
}

func getK8sNode(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getK8sNode")

	clientset, err := GetNewClientset(ctx, d.ConnectionManager)
	if err != nil {
		return nil, err
	}

	name := d.KeyColumnQuals["name"].GetStringValue()

	node, err := clientset.CoreV1().Nodes().Get(ctx, name, metav1.GetOptions{})
	if err != nil && !isNotFoundError(err) {
		return nil, err
	}

	return node, nil
}

//// TRANSFORM FUNCTIONS

func transformNodeTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(v1.Node)
	return mergeTags(obj.Labels, obj.Annotations), nil
}
