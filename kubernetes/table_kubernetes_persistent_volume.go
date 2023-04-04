package kubernetes

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableKubernetesPersistentVolume(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_persistent_volume",
		Description: "A PersistentVolume (PV) is a piece of storage in the cluster that has been provisioned by an administrator or dynamically provisioned using Storage Classes. PVs are volume plugins like Volumes, but have a lifecycle independent of any individual Pod that uses the PV.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getK8sPV,
		},
		List: &plugin.ListConfig{
			Hydrate: listK8sPVs,
		},
		Columns: k8sCommonGlobalColumns([]*plugin.Column{
			//// PersistentVolumeSpec columns
			{
				Name:        "storage_class",
				Type:        proto.ColumnType_STRING,
				Description: "Name of StorageClass to which this persistent volume belongs. Empty value means that this volume does not belong to any StorageClass.",
				Transform:   transform.FromField("Spec.StorageClassName"),
			},
			{
				Name:        "volume_mode",
				Type:        proto.ColumnType_STRING,
				Description: "Defines if a volume is intended to be used with a formatted filesystem or to remain in raw block state.",
				Transform:   transform.FromField("Spec.VolumeMode"),
			},
			{
				Name:        "persistent_volume_reclaim_policy",
				Type:        proto.ColumnType_STRING,
				Description: "What happens to a persistent volume when released from its claim. Valid options are Retain (default for manually created PersistentVolumes), Delete (default for dynamically provisioned PersistentVolumes), and Recycle (deprecated). Recycle must be supported by the volume plugin underlying this PersistentVolume.",
				Transform:   transform.FromField("Spec.PersistentVolumeReclaimPolicy"),
			},
			{
				Name:        "access_modes",
				Type:        proto.ColumnType_JSON,
				Description: "List of ways the volume can be mounted.",
				Transform:   transform.FromField("Spec.AccessModes"),
			},
			{
				Name:        "capacity",
				Type:        proto.ColumnType_JSON,
				Description: "A description of the persistent volume's resources and capacity.",
				Transform:   transform.FromField("Spec.Capacity"),
			},
			{
				Name:        "claim_ref",
				Type:        proto.ColumnType_JSON,
				Description: "ClaimRef is part of a bi-directional binding between PersistentVolume and PersistentVolumeClaim. Expected to be non-nil when bound.",
				Transform:   transform.FromField("Spec.ClaimRef"),
			},
			{
				Name:        "mount_options",
				Type:        proto.ColumnType_JSON,
				Description: "A list of mount options, e.g. [\"ro\", \"soft\"].",
				Transform:   transform.FromField("Spec.MountOptions"),
			},
			{
				Name:        "node_affinity",
				Type:        proto.ColumnType_JSON,
				Description: "Defines constraints that limit what nodes this volume can be accessed from.",
				Transform:   transform.FromField("Spec.NodeAffinity"),
			},
			{
				Name:        "persistent_volume_source",
				Type:        proto.ColumnType_JSON,
				Description: "The actual volume backing the persistent volume.",
				Transform:   transform.FromField("Spec.PersistentVolumeSource"),
			},

			//// PersistentVolumeStatus columns
			{
				Name:        "phase",
				Type:        proto.ColumnType_STRING,
				Description: "Phase indicates if a volume is available, bound to a claim, or released by a claim.",
				Transform:   transform.FromField("Status.Phase"),
			},
			{
				Name:        "message",
				Type:        proto.ColumnType_STRING,
				Description: "A human-readable message indicating details about why the volume is in this state.",
				Transform:   transform.FromField("Status.Message"),
			},
			{
				Name:        "reason",
				Type:        proto.ColumnType_STRING,
				Description: "Reason is a brief CamelCase string that describes any failure and is meant for machine parsing and tidy display in the CLI.",
				Transform:   transform.FromField("Status.Reason"),
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
				Transform:   transform.From(transformPVTags),
			},
			{
				Name:        "manifest_file_path",
				Type:        proto.ColumnType_STRING,
				Description: "The path to the manifest file.",
				Transform:   transform.FromField("ManifestFilePath").Transform(transform.NullIfZeroValue),
			},
		}),
	}
}

type PersistentVolume struct {
	v1.PersistentVolume
	ManifestFilePath string
}

//// HYDRATE FUNCTIONS

func listK8sPVs(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listK8sPVs")

	// Get the client for querying the K8s APIs for the provided context.
	// If the connection is configured for the manifest files, the client will return nil.
	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		return nil, err
	}

	//
	// Check for manifest files
	//
	parsedContents, err := fetchResourceFromManifestFileByKind(ctx, d, "PersistentVolume")
	if err != nil {
		return nil, err
	}

	for _, content := range parsedContents {
		persistentVolume := content.Data.(*v1.PersistentVolume)

		d.StreamListItem(ctx, PersistentVolume{*persistentVolume, content.Path})

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

	var response *v1.PersistentVolumeList
	pageLeft := true

	for pageLeft {
		response, err = clientset.CoreV1().PersistentVolumes().List(ctx, input)
		if err != nil {
			return nil, err
		}

		if response.GetContinue() != "" {
			input.Continue = response.Continue
		} else {
			pageLeft = false
		}

		for _, persistentVolume := range response.Items {
			d.StreamListItem(ctx, PersistentVolume{persistentVolume, ""})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

func getK8sPV(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getK8sPV")

	// Get the client for querying the K8s APIs for the provided context.
	// If the connection is configured for the manifest files, the client will return nil.
	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		return nil, err
	}

	name := d.EqualsQuals["name"].GetStringValue()

	// return if  name is empty
	if name == "" {
		return nil, nil
	}

	//
	// Get the manifest resource
	//
	parsedContents, err := fetchResourceFromManifestFileByKind(ctx, d, "PersistentVolume")
	if err != nil {
		return nil, err
	}

	for _, content := range parsedContents {
		persistentVolume := content.Data.(*v1.PersistentVolume)

		if persistentVolume.Name == name {
			return PersistentVolume{*persistentVolume, content.Path}, nil
		}
	}

	//
	// Get the deployed resource
	//
	if clientset == nil {
		return nil, nil
	}

	persistentVolume, err := clientset.CoreV1().PersistentVolumes().Get(ctx, name, metav1.GetOptions{})
	if err != nil && !isNotFoundError(err) {
		return nil, err
	}

	return PersistentVolume{*persistentVolume, ""}, nil
}

//// TRANSFORM FUNCTIONS

func transformPVTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(PersistentVolume)
	return mergeTags(obj.Labels, obj.Annotations), nil
}
