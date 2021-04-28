package kubernetes

import (
	"context"

	v1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableKubernetesJob(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_job",
		Description: "A Job creates one or more Pods and will continue to retry execution of the Pods until a specified number of them successfully terminate.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "namespace"}),
			Hydrate:    getK8sJob,
		},
		List: &plugin.ListConfig{
			Hydrate: listK8sJobs,
		},
		Columns: k8sCommonColumns([]*plugin.Column{
			//// JobSpec columns
			{
				Name:        "parallelism",
				Type:        proto.ColumnType_INT,
				Description: "The maximum desired number of pods the job should run at any given time. The actual number of pods running in steady state will be less than this number when ((.spec.completions - .status.successful) < .spec.parallelism), i.e. when the work left to do is less than max parallelism.",
				Transform:   transform.FromField("Spec.Parallelism"),
			},
			{
				Name:        "completions",
				Type:        proto.ColumnType_INT,
				Description: "The desired number of successfully finished pods the job should be run with.",
				Transform:   transform.FromField("Spec.Completions"),
			},
			{
				Name:        "active_deadline_seconds",
				Type:        proto.ColumnType_INT,
				Description: "The duration in seconds relative to the startTime that the job may be active before the system tries to terminate it.",
				Transform:   transform.FromField("Spec.ActiveDeadlineSeconds"),
			},
			{
				Name:        "backoff_limit",
				Type:        proto.ColumnType_INT,
				Description: "The number of retries before marking this job failed. Defaults to 6.",
				Transform:   transform.FromField("Spec.BackoffLimit"),
			},
			{
				Name:        "selector",
				Type:        proto.ColumnType_JSON,
				Description: "A label query over pods that should match the pod count.",
				Transform:   transform.FromField("Spec.Selector"),
			},
			{
				Name:        "manual_selector",
				Type:        proto.ColumnType_BOOL,
				Description: "ManualSelector controls generation of pod labels and pod selectors. When false or unset, the system pick labels unique to this job and appends those labels to the pod template.  When true, the user is responsible for picking unique labels and specifying the selector.",
				Transform:   transform.FromField("Spec.ManualSelector"),
			},
			{
				Name:        "ttl_seconds_after_finished",
				Type:        proto.ColumnType_INT,
				Description: "limits the lifetime of a Job that has finished execution (either Complete or Failed). If this field is set, ttlSecondsAfterFinished after the Job finishes, it is eligible to be automatically deleted.",
				Transform:   transform.FromField("Spec.TTLSecondsAfterFinished"),
			},
			{
				Name:        "template",
				Type:        proto.ColumnType_JSON,
				Description: "Describes the pod that will be created when executing a job.",
				Transform:   transform.FromField("Spec.Template"),
			},

			//// JobStatus columns
			{
				Name:        "start_time",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "Time when the job was acknowledged by the job controller.",
				Transform:   transform.FromField("Status.StartTime").Transform(v1TimeToRFC3339),
			},
			{
				Name:        "completion_time",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "Time when the job was completed.",
				Transform:   transform.FromField("Status.CompletionTime").Transform(v1TimeToRFC3339),
			},
			{
				Name:        "active",
				Type:        proto.ColumnType_INT,
				Description: "The number of actively running pods.",
				Transform:   transform.FromField("Status.Active"),
			},
			{
				Name:        "succeeded",
				Type:        proto.ColumnType_INT,
				Description: "The number of pods which reached phase Succeeded.",
				Transform:   transform.FromField("Status.Succeeded"),
			},
			{
				Name:        "failed",
				Type:        proto.ColumnType_INT,
				Description: "The number of pods which reached phase Failed.",
				Transform:   transform.FromField("Status.Failed"),
			},
			{
				Name:        "conditions",
				Type:        proto.ColumnType_JSON,
				Description: "The latest available observations of an object's current state.",
				Transform:   transform.FromField("Status.Conditions"),
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
				Transform:   transform.From(transformJobTags),
			},
		}),
	}
}

//// HYDRATE FUNCTIONS

func listK8sJobs(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listK8sJobs")

	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		return nil, err
	}

	jobs, err := clientset.BatchV1().Jobs("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, job := range jobs.Items {
		d.StreamListItem(ctx, job)
	}

	return nil, nil
}

func getK8sJob(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getK8sJob")

	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		return nil, err
	}

	name := d.KeyColumnQuals["name"].GetStringValue()
	namespace := d.KeyColumnQuals["namespace"].GetStringValue()

	job, err := clientset.BatchV1().Jobs(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil && !isNotFoundError(err) {
		return nil, err
	}

	return *job, nil
}

//// TRANSFORM FUNCTIONS

func transformJobTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(v1.Job)
	return mergeTags(obj.Labels, obj.Annotations), nil
}
