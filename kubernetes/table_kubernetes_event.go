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

func tableKubernetesEvent(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_event",
		Description: "Kubernetes Event is a report of an event somewhere in the cluster.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "namespace"}),
			Hydrate:    getK8sEvent,
		},
		List: &plugin.ListConfig{
			Hydrate:    listK8sEvents,
			KeyColumns: getCommonOptionalKeyQuals(),
		},
		Columns: k8sCommonColumns([]*plugin.Column{
			{
				Name:        "last_timestamp",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "Time when this event was last observed.",
				Transform:   transform.FromField("LastTimestamp").Transform(v1TimeToRFC3339),
			},
			{
				Name:        "type",
				Type:        proto.ColumnType_STRING,
				Description: "Type of this event (Normal, Warning), new types could be added in the future.",
			},
			{
				Name:        "reason",
				Type:        proto.ColumnType_STRING,
				Description: "The reason the transition into the object's current status.",
			},
			{
				Name:        "message",
				Type:        proto.ColumnType_STRING,
				Description: "A description of the status of this operation.",
			},
			{
				Name:        "action",
				Type:        proto.ColumnType_STRING,
				Description: "What action was taken/failed with the regarding object.",
			},
			{
				Name:        "count",
				Type:        proto.ColumnType_INT,
				Description: "The number of times this event has occurred.",
			},
			{
				Name:        "event_time",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "Time when this event was first observed.",
				Transform:   transform.FromField("EventTime").Transform(v1MicroTimeToRFC3339),
			},
			{
				Name:        "first_timestamp",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The time at which the event was first recorded.",
				Transform:   transform.FromField("FirstTimestamp").Transform(v1TimeToRFC3339),
			},
			{
				Name:        "reporting_component",
				Type:        proto.ColumnType_STRING,
				Description: "Name of the controller that emitted this event.",
				Transform:   transform.FromField("ReportingComponent"),
			},
			{
				Name:        "reporting_instance",
				Type:        proto.ColumnType_STRING,
				Description: "ID of the controller instance.",
				Transform:   transform.FromField("ReportingInstance"),
			},
			{
				Name:        "involved_object",
				Type:        proto.ColumnType_JSON,
				Description: "The object that this event is about.",
				Transform:   transform.FromField("InvolvedObject"),
			},
			{
				Name:        "related",
				Type:        proto.ColumnType_JSON,
				Description: "Optional secondary object for more complex actions.",
			},
			{
				Name:        "series",
				Type:        proto.ColumnType_JSON,
				Description: "Data about the event series this event represents.",
			},
			{
				Name:        "source_type",
				Type:        proto.ColumnType_JSON,
				Description: "The component reporting this event.",
			},
			{
				Name:        "config_source",
				Type:        proto.ColumnType_STRING,
				Description: "The source of the resource. Possible values are: deployed and manifest. If the resource is fetched from the spec file the value will be manifest.",
				Transform:   transform.From(eventResourceSourceType),
			},
		}),
	}
}

type Event struct {
	v1.Event
	Path      string
	StartLine int
	EndLine   int
}

//// HYDRATE FUNCTIONS

func listK8sEvents(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Get the client for querying the K8s APIs for the provided context.
	// If the connection is configured for the manifest files, the client will return nil.
	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listK8sEvents", "client_err", err)
		return nil, err
	}

	//
	// Check for manifest files
	//
	parsedContents, err := fetchResourceFromManifestFileByKind(ctx, d, "Event")
	if err != nil {
		return nil, err
	}

	for _, content := range parsedContents {
		event := content.Data.(*v1.Event)

		d.StreamListItem(ctx, Event{*event, content.Path, content.StartLine, content.EndLine})

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

	var response *v1.EventList
	pageLeft := true

	for pageLeft {
		response, err = clientset.CoreV1().Events("").List(ctx, input)
		if err != nil {
			plugin.Logger(ctx).Error("listK8sEvents", "api_err", err)
			return nil, err
		}

		if response.GetContinue() != "" {
			input.Continue = response.Continue
		} else {
			pageLeft = false
		}

		for _, event := range response.Items {
			d.StreamListItem(ctx, Event{event, "", 0, 0})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

func getK8sEvent(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Get the client for querying the K8s APIs for the provided context.
	// If the connection is configured for the manifest files, the client will return nil.
	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("getK8sEvent", "client_err", err)
		return nil, err
	}

	name := d.EqualsQuals["name"].GetStringValue()
	namespace := d.EqualsQuals["namespace"].GetStringValue()

	// return if namespace or name is empty
	if namespace == "" || name == "" {
		return nil, nil
	}

	//
	// Get the manifest resource
	//
	parsedContents, err := fetchResourceFromManifestFileByKind(ctx, d, "Event")
	if err != nil {
		return nil, err
	}

	for _, content := range parsedContents {
		event := content.Data.(*v1.Event)

		if event.Name == name && event.Namespace == namespace {
			return Event{*event, content.Path, content.StartLine, content.EndLine}, nil
		}
	}

	//
	// Get the deployed resource
	//
	if clientset == nil {
		return nil, nil
	}

	event, err := clientset.CoreV1().Events(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil && !isNotFoundError(err) {
		plugin.Logger(ctx).Error("getK8sEvent", "api_err", err)
		return nil, err
	}

	return Event{*event, "", 0, 0}, nil
}

//// TRANSFORM FUNCTIONS

func eventResourceSourceType(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(Event)

	if obj.Path != "" {
		return "manifest", nil
	}
	return "deployed", nil
}
