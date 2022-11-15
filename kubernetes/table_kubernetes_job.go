package kubernetes

import (
	"context"
	"strings"

	v1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
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
			Hydrate:    listK8sJobs,
			KeyColumns: getCommonOptionalKeyQuals(),
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
				Name:        "selector_query",
				Type:        proto.ColumnType_STRING,
				Description: "A query string representation of the selector.",
				Transform:   transform.FromField("Spec.Selector").Transform(labelSelectorToString),
			},
			{
				Name:        "selector",
				Type:        proto.ColumnType_JSON,
				Description: "A label query over pods that should match the pod count.",
				Transform:   transform.FromField("Spec.Selector"),
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

	var response *v1.JobList
	pageLeft := true

	for pageLeft {
		response, err = clientset.BatchV1().Jobs("").List(ctx, input)
		if err != nil {
			return nil, err
		}

		if response.GetContinue() != "" {
			input.Continue = response.Continue
		} else {
			pageLeft = false
		}

		for _, job := range response.Items {
			d.StreamListItem(ctx, job)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
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

	name := d.EqualsQuals["name"].GetStringValue()
	namespace := d.EqualsQuals["namespace"].GetStringValue()

	// return if namespace or name is empty
	if namespace == "" || name == "" {
		return nil, nil
	}

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
