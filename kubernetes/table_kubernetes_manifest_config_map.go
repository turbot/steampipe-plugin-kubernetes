package kubernetes

import (
	"context"
	"os"
	"strings"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableKubernetesManifestConfigMap(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_manifest_config_map",
		Description: "",
		List: &plugin.ListConfig{
			ParentHydrate: listKubernetesManifestFiles,
			Hydrate:       listKubernetesManifestConfigMaps,
			KeyColumns:    plugin.OptionalColumns([]string{"manifest_file_path"}),
		},
		Columns: []*plugin.Column{

			// Metadata Columns
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of resource."},
			{Name: "namespace", Type: proto.ColumnType_STRING, Description: "Namespace defines the space within which each name must be unique."},
			{Name: "labels", Type: proto.ColumnType_JSON, Description: "Map of string keys and values that can be used to organize and categorize (scope and select) objects. May match selectors of replication controllers and services."},
			{Name: "immutable", Type: proto.ColumnType_BOOL, Description: "If set to true, ensures that data stored in the ConfigMap cannot be updated (only object metadata can be modified). If not set to true, the field can be modified at any time. Defaulted to nil."},
			{Name: "data", Type: proto.ColumnType_JSON, Description: "Contains the configuration data."},
			{Name: "binary_data", Type: proto.ColumnType_JSON, Description: "Contains the configuration binary data."},

			// Steampipe Standard Columns
			{Name: "title", Type: proto.ColumnType_STRING, Description: ColumnDescriptionTitle, Transform: transform.FromField("Name")},
			{Name: "tags", Type: proto.ColumnType_JSON, Description: ColumnDescriptionTags, Transform: transform.From(transformManifestConfigMapTags)},

			{Name: "manifest_file_path", Description: "Path to the file.", Type: proto.ColumnType_STRING},
		},
	}
}

type KubernetesManifestConfigMap struct {
	v1.ConfigMap
	ManifestFilePath string
}

//// LIST FUNCTION

func listKubernetesManifestConfigMaps(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// The path comes from a parent hydrate, defaulting to the config paths or
	// available by the optional key column
	path := h.Item.(filePath).Path

	// Load the file into a buffer
	content, err := os.ReadFile(path)
	if err != nil {
		plugin.Logger(ctx).Error("kubernetes_manifest_config_map.listKubernetesManifestConfigMaps", "failed to read file", err, "path", path)
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
			plugin.Logger(ctx).Error("kubernetes_manifest_config_map.listKubernetesManifestConfigMaps", "failed to decode the file", err, "path", path)
			return nil, err
		}

		// Return if the definition is not for the config map resource
		if groupVersionKind.Kind == "ConfigMap" {
			configMap := obj.(*v1.ConfigMap)

			d.StreamListItem(ctx, KubernetesManifestConfigMap{*configMap, path})

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

func transformManifestConfigMapTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(KubernetesManifestConfigMap)
	return mergeTags(obj.Labels, obj.Annotations), nil
}
