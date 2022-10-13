package kubernetes

import (
	"context"
	"strings"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableKubernetesPersistentVolumeClaim(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_persistent_volume_claim",
		Description: "A PersistentVolumeClaim (PVC) is a request for storage by a user.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "namespace"}),
			Hydrate:    getK8sPVC,
		},
		List: &plugin.ListConfig{
			Hydrate:    listK8sPVCs,
			KeyColumns: getCommonOptionalKeyQuals(),
		},
		Columns: k8sCommonColumns([]*plugin.Column{
			//// PersistentVolumeClaimSpec columns
			{
				Name:        "volume_name",
				Type:        proto.ColumnType_STRING,
				Description: "The binding reference to the PersistentVolume backing this claim.",
				Transform:   transform.FromField("Spec.VolumeName"),
			},
			{
				Name:        "volume_mode",
				Type:        proto.ColumnType_STRING,
				Description: "Defines if a volume is intended to be used with a formatted filesystem or to remain in raw block state.",
				Transform:   transform.FromField("Spec.VolumeMode"),
			},
			{
				Name:        "storage_class",
				Type:        proto.ColumnType_STRING,
				Description: "Name of StorageClass to which this persistent volume belongs. Empty value means that this volume does not belong to any StorageClass.",
				Transform:   transform.FromField("Spec.StorageClassName"),
			},
			{
				Name:        "access_modes",
				Type:        proto.ColumnType_JSON,
				Description: "List of ways the volume can be mounted.",
				Transform:   transform.FromField("Spec.AccessModes"),
			},
			{
				Name:        "data_source",
				Type:        proto.ColumnType_JSON,
				Description: "The source of the volume. This can be used to specify either: an existing VolumeSnapshot object (snapshot.storage.k8s.io/VolumeSnapshot), an existing PVC (PersistentVolumeClaim) or an existing custom resource that implements data population (Alpha).",
				Transform:   transform.FromField("Spec.DataSource"),
			},
			{
				Name:        "resources",
				Type:        proto.ColumnType_JSON,
				Description: "Represents the minimum resources the volume should have.",
				Transform:   transform.FromField("Spec.Resources"),
			},
			{
				Name:        "selector",
				Type:        proto.ColumnType_JSON,
				Description: "The actual volume backing the persistent volume.",
				Transform:   transform.FromField("Spec.Selector"),
			},

			//// PersistentVolumeClaimStatus columns
			{
				Name:        "phase",
				Type:        proto.ColumnType_STRING,
				Description: "Phase indicates the current phase of PersistentVolumeClaim.",
				Transform:   transform.FromField("Status.Phase"),
			},
			{
				Name:        "status_access_modes",
				Type:        proto.ColumnType_JSON,
				Description: "The actual access modes the volume backing the PVC has.",
				Transform:   transform.FromField("Status.AccessModes"),
			},
			{
				Name:        "capacity",
				Type:        proto.ColumnType_JSON,
				Description: "The actual resources of the underlying volume.",
				Transform:   transform.FromField("Status.Capacity"),
			},
			{
				Name:        "conditions",
				Type:        proto.ColumnType_JSON,
				Description: "The Condition of persistent volume claim.",
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
				Transform:   transform.From(transformPVCTags),
			},
		}),
	}
}

//// HYDRATE FUNCTIONS

func listK8sPVCs(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listK8sPVCs")

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

	var response *v1.PersistentVolumeClaimList
	pageLeft := true

	for pageLeft {
		response, err = clientset.CoreV1().PersistentVolumeClaims("").List(ctx, input)
		if err != nil {
			return nil, err
		}

		if response.GetContinue() != "" {
			input.Continue = response.Continue
		} else {
			pageLeft = false
		}

		for _, persistentVolumeClaim := range response.Items {
			d.StreamListItem(ctx, persistentVolumeClaim)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

func getK8sPVC(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getK8sPVC")

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

	persistentVolumeClaim, err := clientset.CoreV1().PersistentVolumeClaims(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil && !isNotFoundError(err) {
		return nil, err
	}

	return *persistentVolumeClaim, nil
}

//// TRANSFORM FUNCTIONS

func transformPVCTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(v1.PersistentVolumeClaim)
	return mergeTags(obj.Labels, obj.Annotations), nil
}
