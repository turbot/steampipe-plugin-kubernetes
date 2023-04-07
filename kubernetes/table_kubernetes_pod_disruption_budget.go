package kubernetes

import (
	"context"
	"strings"

	v1 "k8s.io/api/policy/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableKubernetesPDB(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_pod_disruption_budget",
		Description: "A Pod Disruption Budget limits the number of Pods of a replicated application that are down simultaneously from voluntary disruptions.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "namespace"}),
			Hydrate:    getPDB,
		},
		List: &plugin.ListConfig{
			Hydrate:    listPDBs,
			KeyColumns: getCommonOptionalKeyQuals(),
		},
		Columns: k8sCommonColumns([]*plugin.Column{

			// PodDisruptionBudgetSpec
			{
				Name:        "min_available",
				Type:        proto.ColumnType_STRING,
				Description: "An eviction is allowed if at least 'minAvailable' pods selected by 'selector' will still be available after the eviction.",
				Transform:   transform.FromField("Spec.MinAvailable"),
			},
			{
				Name:        "selector",
				Type:        proto.ColumnType_JSON,
				Description: "Label query over pods whose evictions are managed by the disruption budget.",
				Transform:   transform.FromField("Spec.Selector"),
			},
			{
				Name:        "max_unavailable",
				Type:        proto.ColumnType_STRING,
				Description: "An eviction is allowed if at most 'maxAvailable' pods selected by 'selector' will still be unavailable after the eviction.",
				Transform:   transform.FromField("Spec.MaxUnavailable"),
			},

			// Steampipe Standard Columns
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
				Transform:   transform.From(transformPDBTags),
			},
		}),
	}
}

//// HYDRATE FUNCTIONS

func listPDBs(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listPDBs")

	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		return nil, err
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

	var response *v1.PodDisruptionBudgetList
	pageLeft := true

	for pageLeft {
		response, err = clientset.PolicyV1().PodDisruptionBudgets("").List(ctx, input)
		if err != nil {
			return nil, err
		}

		if response.GetContinue() != "" {
			input.Continue = response.Continue
		} else {
			pageLeft = false
		}

		for _, item := range response.Items {
			d.StreamListItem(ctx, item)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

func getPDB(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getPDB")

	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		return nil, err
	}

	name := d.EqualsQuals["name"].GetStringValue()
	namespace := d.EqualsQuals["namespace"].GetStringValue()

	// return if namespace or name is empty
	if namespace == "" || name == "" {
		return nil, nil
	}

	pdb, err := clientset.PolicyV1().PodDisruptionBudgets(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil && !isNotFoundError(err) {
		return nil, err
	}

	return *pdb, nil
}

//// TRANSFORM FUNCTIONS

func transformPDBTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(v1.PodDisruptionBudget)
	return mergeTags(obj.Labels, obj.Annotations), nil
}
