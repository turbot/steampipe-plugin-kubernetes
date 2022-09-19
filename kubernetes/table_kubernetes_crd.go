package kubernetes

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
)

func tableKubernetesCRD(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_crd",
		Description: "Cron jobs are useful for creating periodic and recurring tasks, like running backups or sending emails.",
		// Get: &plugin.GetConfig{
		// 	KeyColumns: plugin.AllColumns([]string{"name", "namespace"}),
		// 	Hydrate:    getK8sCronJob,
		// },
		List: &plugin.ListConfig{
			Hydrate: listK8sCRDs,
			//KeyColumns: getCommonOptionalKeyQuals(),
		},
		Columns: k8sCommonColumns([]*plugin.Column{
			//// CronJobSpec columns
		}),
	}
}

//// HYDRATE FUNCTIONS

func listK8sCRDs(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listK8sCRDs")

	clientset, err := GetNewClientCRD(ctx, d)
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

	pageLeft := true
	for pageLeft {
		response, err := clientset.ApiextensionsV1().CustomResourceDefinitions().List(ctx, input)
		if err != nil {
			logger.Error("listK8sCronJobs", "list_err", err)
			return nil, err
		}

		if response.GetContinue() != "" {
			input.Continue = response.Continue
		} else {
			pageLeft = false
		}

		for _, crd := range response.Items {
			d.StreamListItem(ctx, crd)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}
