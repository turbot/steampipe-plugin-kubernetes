package kubernetes

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableKubernetesNamespace(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_namespace",
		Description: "Kubernetes Namespace provides a scope for Names.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getK8sNamespace,
		},
		List: &plugin.ListConfig{
			Hydrate: listK8sNamespaces,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "phase", Require: plugin.Optional},
			},
		},
		Columns: k8sCommonGlobalColumns([]*plugin.Column{

			//// NamespaceSpec Columns
			{
				Name:        "spec_finalizers",
				Type:        proto.ColumnType_JSON,
				Description: "Finalizers is an opaque list of values that must be empty to permanently remove object from storage.",
				Transform:   transform.FromField("Spec.Finalizers"),
			},

			//// NamespaceStatus Columns
			{
				Name:        "phase",
				Type:        proto.ColumnType_STRING,
				Description: "The current lifecycle phase of the namespace.",
				Transform:   transform.FromField("Status.Phase"),
			},
			{
				Name:        "conditions",
				Type:        proto.ColumnType_JSON,
				Description: "The latest available observations of namespace's current state.",
				Transform:   transform.FromField("Status.NamespaceCondition"),
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
				Transform:   transform.From(transformNamespaceTags),
			},
			{
				Name:        "manifest_file_path",
				Type:        proto.ColumnType_STRING,
				Description: "The path to the manifest file.",
				Transform:   transform.FromField("ManifestFilePath").Transform(transform.NullIfZeroValue),
			},
		}),
	}
}

type Namespace struct {
	v1.Namespace
	ManifestFilePath string
	StartLine        int
}

//// HYDRATE FUNCTIONS

func listK8sNamespaces(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listK8sNamespaces")

	// Get the client for querying the K8s APIs for the provided context.
	// If the connection is configured for the manifest files, the client will return nil.
	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		return nil, err
	}

	//
	// Check for manifest files
	//
	parsedContents, err := fetchResourceFromManifestFileByKind(ctx, d, "Namespace")
	if err != nil {
		return nil, err
	}

	for _, content := range parsedContents {
		namespace := content.Data.(*v1.Namespace)

		d.StreamListItem(ctx, Namespace{*namespace, content.Path, content.Line})

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	//
	// Check for deployed resources
	//
	if clientset == nil {
		return nil, nil
	}

	input := metav1.ListOptions{
		Limit: 500,
	}

	// Limiting the results
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < input.Limit {
			if *limit < 1 {
				input.Limit = 1
			} else {
				input.Limit = *limit
			}
		}
	}

	if d.EqualsQualString("phase") != "" {
		input.FieldSelector = fmt.Sprintf("status.phase=%v", d.EqualsQualString("phase"))
	}

	var response *v1.NamespaceList
	pageLeft := true

	for pageLeft {
		response, err = clientset.CoreV1().Namespaces().List(ctx, input)
		if err != nil {
			return nil, err
		}

		if response.GetContinue() != "" {
			input.Continue = response.Continue
		} else {
			pageLeft = false
		}

		for _, namespace := range response.Items {
			d.StreamListItem(ctx, Namespace{namespace, "", 0})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

func getK8sNamespace(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getK8sNamespace")

	// Get the client for querying the K8s APIs for the provided context.
	// If the connection is configured for the manifest files, the client will return nil.
	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		return nil, err
	}

	name := d.EqualsQuals["name"].GetStringValue()

	// return if name is empty
	if name == "" {
		return nil, nil
	}

	//
	// Get the manifest resource
	//
	parsedContents, err := fetchResourceFromManifestFileByKind(ctx, d, "Namespace")
	if err != nil {
		return nil, err
	}

	for _, content := range parsedContents {
		namespace := content.Data.(*v1.Namespace)

		if namespace.Name == name {
			return Namespace{*namespace, content.Path, content.Line}, nil
		}
	}

	//
	// Get the deployed resource
	//
	if clientset == nil {
		return nil, nil
	}

	namespace, err := clientset.CoreV1().Namespaces().Get(ctx, name, metav1.GetOptions{})
	if err != nil && !isNotFoundError(err) {
		return nil, err
	}

	return Namespace{*namespace, "", 0}, nil
}

//// TRANSFORM FUNCTIONS

func transformNamespaceTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(Namespace)
	return mergeTags(obj.Labels, obj.Annotations), nil
}
