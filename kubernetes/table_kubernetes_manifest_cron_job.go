package kubernetes

import (
	"context"
	"os"
	"strings"

	v1 "k8s.io/api/batch/v1"
	"k8s.io/client-go/kubernetes/scheme"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableKubernetesManifestCronJob(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_manifest_cron_job",
		Description: "",
		List: &plugin.ListConfig{
			ParentHydrate: listKubernetesManifestFiles,
			Hydrate:       listKubernetesManifestCronJobs,
			KeyColumns:    plugin.OptionalColumns([]string{"manifest_file_path"}),
		},
		Columns: []*plugin.Column{
			// Metadata Columns
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of resource."},
			{Name: "namespace", Type: proto.ColumnType_STRING, Description: "Namespace defines the space within which each name must be unique."},
			{Name: "labels", Type: proto.ColumnType_JSON, Description: "Map of string keys and values that can be used to organize and categorize (scope and select) objects. May match selectors of replication controllers and services."},

			//// CronJobSpec columns
			{Name: "failed_jobs_history_limit", Type: proto.ColumnType_INT, Description: "The number of failed finished jobs to retain. Value must be non-negative integer.", Transform: transform.FromField("Spec.FailedJobsHistoryLimit")},
			{Name: "schedule", Type: proto.ColumnType_STRING, Description: "The schedule in Cron format.", Transform: transform.FromField("Spec.Schedule")},
			{Name: "starting_deadline_seconds", Type: proto.ColumnType_INT, Description: "Optional deadline in seconds for starting the job if it misses scheduled time for any reason.", Transform: transform.FromField("Spec.StartingDeadlineSeconds")},
			{Name: "successful_jobs_history_limit", Type: proto.ColumnType_INT, Description: "The number of successful finished jobs to retain. Value must be non-negative integer.", Transform: transform.FromField("Spec.SuccessfulJobsHistoryLimit")},
			{Name: "suspend", Type: proto.ColumnType_BOOL, Description: "This flag tells the controller to suspend subsequent executions, it does not apply to already started executions.  Defaults to false.", Transform: transform.FromField("Spec.Suspend")},
			{Name: "concurrency_policy", Type: proto.ColumnType_JSON, Description: "Specifies how to treat concurrent executions of a Job.", Transform: transform.FromField("Spec.ConcurrencyPolicy")},
			{Name: "job_template", Type: proto.ColumnType_JSON, Description: "Specifies the job that will be created when executing a CronJob.", Transform: transform.FromField("Spec.JobTemplate")},

			{Name: "manifest_file_path", Description: "Path to the file.", Type: proto.ColumnType_STRING},
		},
	}
}

type KubernetesManifestCronJob struct {
	v1.CronJob
	ManifestFilePath string
}

//// LIST FUNCTION

func listKubernetesManifestCronJobs(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// The path comes from a parent hydrate, defaulting to the config paths or
	// available by the optional key column
	path := h.Item.(filePath).Path

	// Load the file into a buffer
	content, err := os.ReadFile(path)
	if err != nil {
		plugin.Logger(ctx).Error("kubernetes_manifest_cron_job.listKubernetesManifestCronJobs", "failed to read file", err, "path", path)
		return nil, err
	}
	decoder := scheme.Codecs.UniversalDeserializer()

	// Check for the start of the document
	for _, resource := range strings.Split(string(content), "---") {
		// skip empty documents, `Decode` will fail on them
		if len(resource) == 0 {
			continue
		}

		// Decode the file content
		obj, groupVersionKind, err := decoder.Decode([]byte(resource), nil, nil)
		if err != nil {
			plugin.Logger(ctx).Error("kubernetes_manifest_cron_job.listKubernetesManifestCronJobs", "failed to decode the file", err, "path", path)
			return nil, err
		}

		// Return if the definition is not for the cron job resource
		if groupVersionKind.Kind == "CronJob" {
			cronJob := obj.(*v1.CronJob)

			d.StreamListItem(ctx, KubernetesManifestCronJob{*cronJob, path})

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}
