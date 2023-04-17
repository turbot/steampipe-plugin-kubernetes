package kubernetes

import (
	"context"
	"strings"

	v1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableKubernetesStorageClass(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_storage_class",
		Description: "Storage class provides a way for administrators to describe the classes of storage they offer.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getK8sStorageClass,
		},
		List: &plugin.ListConfig{
			Hydrate:    listK8sStorageClasses,
			KeyColumns: getCommonOptionalKeyQuals(),
		},
		Columns: k8sCommonColumns([]*plugin.Column{
			{
				Name:        "provisioner",
				Type:        proto.ColumnType_STRING,
				Description: "Provisioner indicates the type of the provisioner.",
			},
			{
				Name:        "reclaim_policy",
				Type:        proto.ColumnType_STRING,
				Description: "Dynamically provisioned PersistentVolumes of this storage class are created with this reclaimPolicy. Defaults to Delete.",
			},
			{
				Name:        "allow_volume_expansion",
				Type:        proto.ColumnType_BOOL,
				Description: "AllowVolumeExpansion shows whether the storage class allows volume expand.",
				Transform:   transform.FromField("AllowVolumeExpansion"),
			},
			{
				Name:        "volume_binding_mode",
				Type:        proto.ColumnType_STRING,
				Description: "VolumeBindingMode indicates how PersistentVolumeClaims should be provisioned and bound. When unset, VolumeBindingImmediate is used. This field is only honored by servers that enable the VolumeScheduling feature.",
			},
			{
				Name:        "allowed_topologies",
				Type:        proto.ColumnType_JSON,
				Description: "Restrict the node topologies where volumes can be dynamically provisioned. Each volume plugin defines its own supported topology specifications. An empty TopologySelectorTerm list means there is no topology restriction.",
			},
			{
				Name:        "mount_options",
				Type:        proto.ColumnType_JSON,
				Description: "Dynamically provisioned PersistentVolumes of this storage class are created with these mountOptions, e.g. ['ro', 'soft']. Not validated - mount of the PVs will simply fail if one is invalid.",
			},
			{
				Name:        "parameters",
				Type:        proto.ColumnType_JSON,
				Description: "Parameters holds the parameters for the provisioner that should create volumes of this storage class.",
			},
			{
				Name:        "source_type",
				Type:        proto.ColumnType_STRING,
				Description: "The source of the resource. Possible values are: deployed and manifest. If the resource is fetched from the spec file the value will be manifest.",
				Transform:   transform.From(storageClassResourceSourceType),
			},

			//// Steampipe Standard Columns
			{
				Name:        "title",
				Type:        proto.ColumnType_STRING,
				Description: ColumnDescriptionTitle,
				Transform:   transform.FromField("Name"),
			},
		}),
	}
}

type StorageClass struct {
	v1.StorageClass
	Path      string
	StartLine int
	EndLine   int
}

//// HYDRATE FUNCTIONS

func listK8sStorageClasses(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Get the client for querying the K8s APIs for the provided context.
	// If the connection is configured for the manifest files, the client will return nil.
	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("kubernetes_storage_class.listK8sStorageClasses", "client_error", err)
		return nil, err
	}

	//
	// Check for manifest files
	//
	parsedContents, err := fetchResourceFromManifestFileByKind(ctx, d, "StorageClass")
	if err != nil {
		return nil, err
	}

	for _, content := range parsedContents {
		storageClass := content.Data.(*v1.StorageClass)

		d.StreamListItem(ctx, StorageClass{*storageClass, content.Path, content.StartLine, content.EndLine})

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
			input.Limit = *limit
		}
	}

	commonFieldSelectorValue := getCommonOptionalKeyQualsValueForFieldSelector(d)

	if len(commonFieldSelectorValue) > 0 {
		input.FieldSelector = strings.Join(commonFieldSelectorValue, ",")
	}

	var response *v1.StorageClassList
	pageLeft := true

	for pageLeft {
		response, err = clientset.StorageV1().StorageClasses().List(ctx, input)
		if err != nil {
			plugin.Logger(ctx).Error("kubernetes_storage_class.listK8sStorageClasses", "api_error", err)
			return nil, err
		}

		if response.GetContinue() != "" {
			input.Continue = response.Continue
		} else {
			pageLeft = false
		}

		for _, item := range response.Items {
			d.StreamListItem(ctx, StorageClass{item, "", 0, 0})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

func getK8sStorageClass(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Get the client for querying the K8s APIs for the provided context.
	// If the connection is configured for the manifest files, the client will return nil.
	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("kubernetes_storage_class.getK8sStorageClass", "client_error", err)
		return nil, err
	}

	name := d.EqualsQuals["name"].GetStringValue()

	// check if name is empty
	if name == "" {
		return nil, nil
	}

	//
	// Get the manifest resource
	//
	parsedContents, err := fetchResourceFromManifestFileByKind(ctx, d, "StorageClass")
	if err != nil {
		return nil, err
	}

	for _, content := range parsedContents {
		storageClass := content.Data.(*v1.StorageClass)

		if storageClass.Name == name {
			return StorageClass{*storageClass, content.Path, content.StartLine, content.EndLine}, nil
		}
	}

	//
	// Get the deployed resource
	//
	if clientset == nil {
		return nil, nil
	}

	rs, err := clientset.StorageV1().StorageClasses().Get(ctx, name, metav1.GetOptions{})
	if err != nil && !isNotFoundError(err) {
		plugin.Logger(ctx).Error("kubernetes_storage_class.getK8sStorageClass", "api_error", err)
		return nil, err
	}

	return StorageClass{*rs, "", 0, 0}, nil
}

//// TRANSFORM FUNCTIONS

func storageClassResourceSourceType(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(StorageClass)

	if obj.Path != "" {
		return "manifest", nil
	}
	return "deployed", nil
}
