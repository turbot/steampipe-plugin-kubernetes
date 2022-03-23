package kubernetes

import (
	"context"

	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/turbot/steampipe-plugin-sdk/v2/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin/transform"
)

func tableKubernetesIngress(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_ingress",
		Description: "Ingress exposes HTTP and HTTPS routes from outside the cluster to services within the cluster. Traffic routing is controlled by rules defined on the Ingress resource.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "namespace"}),
			Hydrate:    getK8sIngress,
		},
		List: &plugin.ListConfig{
			Hydrate: listK8sIngresses,
		},
		Columns: k8sCommonColumns([]*plugin.Column{
			//// IngressSpec columns
			{
				Name:        "ingress_class_name",
				Type:        proto.ColumnType_STRING,
				Description: "Name of the IngressClass cluster resource. The associated IngressClass defines which controller will implement the resource.",
				Transform:   transform.FromField("Spec.IngressClassName"),
			},
			{
				Name:        "backend",
				Type:        proto.ColumnType_JSON,
				Description: "A default backend capable of servicing requests that don't match any rule. At least one of 'backend' or 'rules' must be specified.",
				Transform:   transform.FromField("Spec.Backend"),
			},
			{
				Name:        "tls",
				Type:        proto.ColumnType_JSON,
				Description: "TLS configuration.",
				Transform:   transform.FromField("Spec.TLS"),
			},
			{
				Name:        "rules",
				Type:        proto.ColumnType_JSON,
				Description: "A list of host rules used to configure the Ingress.",
				Transform:   transform.FromField("Spec.Rules"),
			},

			//// IngressStatus columns
			{
				Name:        "load_balancer",
				Type:        proto.ColumnType_JSON,
				Description: "a list containing ingress points for the load-balancer. Traffic intended for the service should be sent to these ingress points.",
				Transform:   transform.FromField("Status.LoadBalancer.Ingress"),
			},

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
				Transform:   transform.From(transformIngressTags),
			},
		}),
	}
}

//// HYDRATE FUNCTIONS

func listK8sIngresses(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listK8sIngresses")

	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		return nil, err
	}

	ingresses, err := clientset.ExtensionsV1beta1().Ingresses("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, ingress := range ingresses.Items {
		d.StreamListItem(ctx, ingress)
	}

	return nil, nil
}

func getK8sIngress(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getK8sIngress")

	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		return nil, err
	}

	name := d.KeyColumnQuals["name"].GetStringValue()
	namespace := d.KeyColumnQuals["namespace"].GetStringValue()

	// return if namespace or name is empty
	if namespace == "" || name == "" {
		return nil, nil
	}

	ingress, err := clientset.ExtensionsV1beta1().Ingresses(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil && !isNotFoundError(err) {
		return nil, err
	}

	return *ingress, nil
}

//// TRANSFORM FUNCTIONS

func transformIngressTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(v1beta1.Ingress)
	return mergeTags(obj.Labels, obj.Annotations), nil
}
