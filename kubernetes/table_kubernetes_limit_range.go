package kubernetes

import (
	"context"
	"strings"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableKubernetesLimitRange(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_limit_range",
		Description: "Kubernetes Limit Range",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "namespace"}),
			Hydrate:    getK8sLimitRange,
		},
		List: &plugin.ListConfig{
			Hydrate:    listK8sLimitRanges,
			KeyColumns: getCommonOptionalKeyQuals(),
		},
		Columns: k8sCommonColumns([]*plugin.Column{

			//// LimitRangeSpec Columns
			{
				Name:        "spec_limits",
				Type:        proto.ColumnType_JSON,
				Description: "List of limit range item objects that are enforced.",
				Transform:   transform.FromField("Spec.Limits"),
			},
			{
				Name:        "source_type",
				Type:        proto.ColumnType_STRING,
				Description: "The source of the resource. Possible values are: deployed and manifest. If the resource is fetched from the spec file the value will be manifest.",
				Transform:   transform.From(limitRangeResourceSourceType),
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
				Transform:   transform.From(transformLimitRangeTags),
			},
		}),
	}
}

type LimitRange struct {
	v1.LimitRange
	Path      string
	StartLine int
	EndLine   int
}

//// HYDRATE FUNCTIONS

func listK8sLimitRanges(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listK8sLimitRanges")

	// Get the client for querying the K8s APIs for the provided context.
	// If the connection is configured for the manifest files, the client will return nil.
	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		return nil, err
	}

	//
	// Check for manifest files
	//
	parsedContents, err := fetchResourceFromManifestFileByKind(ctx, d, "LimitRange")
	if err != nil {
		return nil, err
	}

	for _, content := range parsedContents {
		limitRange := content.Data.(*v1.LimitRange)

		d.StreamListItem(ctx, LimitRange{*limitRange, content.Path, content.StartLine, content.EndLine})

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

	commonFieldSelectorValue := getCommonOptionalKeyQualsValueForFieldSelector(d)

	if len(commonFieldSelectorValue) > 0 {
		input.FieldSelector = strings.Join(commonFieldSelectorValue, ",")
	}

	var response *v1.LimitRangeList
	pageLeft := true

	for pageLeft {
		response, err = clientset.CoreV1().LimitRanges("").List(ctx, input)
		if err != nil {
			return nil, err
		}

		if response.GetContinue() != "" {
			input.Continue = response.Continue
		} else {
			pageLeft = false
		}

		for _, limitRange := range response.Items {
			d.StreamListItem(ctx, LimitRange{limitRange, "", 0, 0})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

func getK8sLimitRange(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getK8sLimitRange")

	// Get the client for querying the K8s APIs for the provided context.
	// If the connection is configured for the manifest files, the client will return nil.
	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		return nil, err
	}

	name := d.EqualsQuals["name"].GetStringValue()
	namespace := d.EqualsQuals["namespace"].GetStringValue()

	// return nil if namespace or name is empty
	if namespace == "" || name == "" {
		return nil, nil
	}

	//
	// Get the manifest resource
	//
	parsedContents, err := fetchResourceFromManifestFileByKind(ctx, d, "LimitRange")
	if err != nil {
		return nil, err
	}

	for _, content := range parsedContents {
		limitRange := content.Data.(*v1.LimitRange)

		if limitRange.Name == name && limitRange.Namespace == namespace {
			return LimitRange{*limitRange, content.Path, content.StartLine, content.EndLine}, nil
		}
	}

	//
	// Get the deployed resource
	//
	if clientset == nil {
		return nil, nil
	}

	limitRange, err := clientset.CoreV1().LimitRanges(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil && !isNotFoundError(err) {
		return nil, err
	}

	return LimitRange{*limitRange, "", 0, 0}, nil
}

//// TRANSFORM FUNCTIONS

func transformLimitRangeTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(LimitRange)
	return mergeTags(obj.Labels, obj.Annotations), nil
}

func limitRangeResourceSourceType(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(LimitRange)

	if obj.Path != "" {
		return "manifest", nil
	}
	return "deployed", nil
}
