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

func tableKubernetesManifestJob(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_manifest_job",
		Description: "",
		List: &plugin.ListConfig{
			ParentHydrate: listKubernetesManifestFiles,
			Hydrate:       listKubernetesManifestJobs,
			KeyColumns:    plugin.OptionalColumns([]string{"manifest_file_path"}),
		},
		Columns: []*plugin.Column{

			// Metadata Columns
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of resource."},
			{Name: "namespace", Type: proto.ColumnType_STRING, Description: "Namespace defines the space within which each name must be unique."},
			{Name: "labels", Type: proto.ColumnType_JSON, Description: "Map of string keys and values that can be used to organize and categorize (scope and select) objects. May match selectors of replication controllers and services."},

			// Spec columns
			{Name: "parallelism", Type: proto.ColumnType_INT, Description: "The maximum desired number of pods the job should run at any given time. The actual number of pods running in steady state will be less than this number when ((.spec.completions - .status.successful) < .spec.parallelism), i.e. when the work left to do is less than max parallelism.", Transform: transform.FromField("Spec.Parallelism")},
			{Name: "completions", Type: proto.ColumnType_INT, Description: "The desired number of successfully finished pods the job should be run with.", Transform: transform.FromField("Spec.Completions")},
			{Name: "active_deadline_seconds", Type: proto.ColumnType_INT, Description: "The duration in seconds relative to the startTime that the job may be active before the system tries to terminate it.", Transform: transform.FromField("Spec.ActiveDeadlineSeconds")},
			{Name: "backoff_limit", Type: proto.ColumnType_INT, Description: "The number of retries before marking this job failed. Defaults to 6.", Transform: transform.FromField("Spec.BackoffLimit")},
			{Name: "manual_selector", Type: proto.ColumnType_BOOL, Description: "ManualSelector controls generation of pod labels and pod selectors. When false or unset, the system pick labels unique to this job and appends those labels to the pod template.  When true, the user is responsible for picking unique labels and specifying the selector.", Transform: transform.FromField("Spec.ManualSelector")},
			{Name: "ttl_seconds_after_finished", Type: proto.ColumnType_INT, Description: "limits the lifetime of a Job that has finished execution (either Complete or Failed). If this field is set, ttlSecondsAfterFinished after the Job finishes, it is eligible to be automatically deleted.", Transform: transform.FromField("Spec.TTLSecondsAfterFinished")},
			{Name: "selector_query", Type: proto.ColumnType_STRING, Description: "A query string representation of the selector.", Transform: transform.FromField("Spec.Selector").Transform(labelSelectorToString)},
			{Name: "selector", Type: proto.ColumnType_JSON, Description: "A label query over pods that should match the pod count.", Transform: transform.FromField("Spec.Selector")},
			{Name: "template", Type: proto.ColumnType_JSON, Description: "Describes the pod that will be created when executing a job.", Transform: transform.FromField("Spec.Template")},

			// Steampipe Standard Columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTitle, Transform: transform.FromField("Name")},
			{Name: "tags", Type: proto.ColumnType_JSON, Description: ColumnDescriptionTags, Transform: transform.From(transformManifestJobTags)},

			{Name: "manifest_file_path", Description: "Path to the file.", Type: proto.ColumnType_STRING},
		},
	}
}

type KubernetesManifestJob struct {
	v1.Job
	ManifestFilePath string
}

//// LIST FUNCTION

func listKubernetesManifestJobs(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// The path comes from a parent hydrate, defaulting to the config paths or
	// available by the optional key column
	path := h.Item.(filePath).Path

	// Load the file into a buffer
	content, err := os.ReadFile(path)
	if err != nil {
		plugin.Logger(ctx).Error("kubernetes_manifest_job.listKubernetesManifestJobs", "failed to read file", err, "path", path)
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
			plugin.Logger(ctx).Error("kubernetes_manifest_job.listKubernetesManifestJobs", "failed to decode the file", err, "path", path)
			return nil, err
		}

		// Return if the definition is not for the job resource
		if groupVersionKind.Kind == "Job" {
			job := obj.(*v1.Job)

			d.StreamListItem(ctx, KubernetesManifestJob{*job, path})

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

func transformManifestJobTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(KubernetesManifestJob)
	return mergeTags(obj.Labels, obj.Annotations), nil
}
