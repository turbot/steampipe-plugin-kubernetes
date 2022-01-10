package kubernetes

import (
	"context"

	v1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableKubernetesCronJob(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_cron_job",
		Description: "A Job creates one or more Pods and will continue to retry execution of the Pods until a specified number of them successfully terminate.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "namespace"}),
			Hydrate:    getK8sCronJob,
		},
		List: &plugin.ListConfig{
			Hydrate: listK8sCronJobs,
		},
		Columns: k8sCommonColumns([]*plugin.Column{
			//// CronJobSpec columns
			{
				Name:        "failed_jobs_history_limit",
				Type:        proto.ColumnType_INT,
				Description: "The number of failed finished jobs to retain. Value must be non-negative integer.",
				Transform:   transform.FromField("Spec.FailedJobsHistoryLimit"),
			},
			{
				Name:        "schedule",
				Type:        proto.ColumnType_STRING,
				Description: "The schedule in Cron format.",
				Transform:   transform.FromField("Spec.Schedule"),
			},
			{
				Name:        "starting_deadline_seconds",
				Type:        proto.ColumnType_INT,
				Description: "Optional deadline in seconds for starting the job if it misses scheduledtime for any reason.",
				Transform:   transform.FromField("Spec.StartingDeadlineSeconds"),
			},
			{
				Name:        "successful_jobs_history_limit",
				Type:        proto.ColumnType_INT,
				Description: "The number of successful finished jobs to retain. Value must be non-negative integer.",
				Transform:   transform.FromField("Spec.SuccessfulJobsHistoryLimit"),
			},
			{
				Name:        "suspend",
				Type:        proto.ColumnType_BOOL,
				Description: "This flag tells the controller to suspend subsequent executions, it does not apply to already started executions.  Defaults to false.",
				Transform:   transform.FromField("Spec.Suspend"),
			},
			{
				Name:        "concurrency_policy",
				Type:        proto.ColumnType_JSON,
				Description: "Specifies how to treat concurrent executions of a Job.",
				Transform:   transform.FromField("Spec.ConcurrencyPolicy"),
			},
			{
				Name:        "job_template",
				Type:        proto.ColumnType_JSON,
				Description: "Specifies the job that will be created when executing a CronJob.",
				Transform:   transform.FromField("Spec.JobTemplate"),
			},

			//// CronJobStatus columns
			{
				Name:        "last_schedule_time",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "Information when was the last time the job was successfully scheduled.",
				Transform:   transform.FromField("Status.LastScheduleTime").Transform(v1TimeToRFC3339),
			},
			{
				Name:        "last_successful_time",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "Information when was the last time the job successfully completed.",
				Transform:   transform.FromField("Status.LastSuccessfulTime").Transform(v1TimeToRFC3339),
			},
			{
				Name:        "active",
				Type:        proto.ColumnType_JSON,
				Description: "A list of pointers to currently running jobs.",
				Transform:   transform.FromField("Status.Active"),
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
				Transform:   transform.From(transformCronJobTags),
			},
		}),
	}
}

//// HYDRATE FUNCTIONS

func listK8sCronJobs(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listK8sCronJobs")

	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		return nil, err
	}

	cronJobs, err := clientset.BatchV1().CronJobs("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, cronJob := range cronJobs.Items {
		d.StreamListItem(ctx, cronJob)
	}

	return nil, nil
}

func getK8sCronJob(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getK8sCronJob")

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

	cronJob, err := clientset.BatchV1().CronJobs(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil && !isNotFoundError(err) {
		return nil, err
	}

	return *cronJob, nil
}

//// TRANSFORM FUNCTIONS

func transformCronJobTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(v1.CronJob)
	return mergeTags(obj.Labels, obj.Annotations), nil
}
