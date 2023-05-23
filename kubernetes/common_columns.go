package kubernetes

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

const (
	ColumnDescriptionTitle = "Title of the resource."
	ColumnDescriptionAkas  = "Array of globally unique identifier strings (also known as) for the resource."
	ColumnDescriptionTags  = "A map of tags for the resource. This includes both labels and annotations."
)

func manifestResourceColumns() []*plugin.Column {
	return []*plugin.Column{
		{Name: "path", Type: proto.ColumnType_STRING, Description: "The path to the manifest file.", Transform: transform.FromField("Path").Transform(transform.NullIfZeroValue)},
		{Name: "start_line", Type: proto.ColumnType_INT, Description: "The path to the manifest file.", Transform: transform.FromField("StartLine").NullIfZero()},
		{Name: "end_line", Type: proto.ColumnType_INT, Description: "The path to the manifest file.", Transform: transform.FromField("EndLine").NullIfZero()},
	}
}

func objectMetadataPrimaryColumnsWithoutNamespace() []*plugin.Column {
	return []*plugin.Column{
		//{Name: "raw", Type: proto.ColumnType_JSON, Transform: transform.FromValue()},
		{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the object.  Name must be unique within a namespace."},
		// {Name: "namespace", Type: proto.ColumnType_STRING, Description: "Namespace defines the space within which each name must be unique."},
		{Name: "uid", Type: proto.ColumnType_STRING, Description: "UID is the unique in time and space value for this object.", Transform: transform.FromField("UID").Transform(transform.NullIfZeroValue)},
	}
}

func objectMetadataPrimaryColumns() []*plugin.Column {
	return []*plugin.Column{
		//{Name: "raw", Type: proto.ColumnType_JSON, Transform: transform.FromValue()},
		{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the object.  Name must be unique within a namespace."},
		{Name: "namespace", Type: proto.ColumnType_STRING, Description: "Namespace defines the space within which each name must be unique."},
		{Name: "uid", Type: proto.ColumnType_STRING, Description: "UID is the unique in time and space value for this object.", Transform: transform.FromField("UID").Transform(transform.NullIfZeroValue)},
	}
}

func objectMetadataSecondaryColumns() []*plugin.Column {
	return []*plugin.Column{
		{Name: "generate_name", Type: proto.ColumnType_STRING, Description: "GenerateName is an optional prefix, used by the server, to generate a unique name ONLY IF the Name field has not been provided."},
		{Name: "resource_version", Type: proto.ColumnType_STRING, Description: "An opaque value that represents the internal version of this object that can be used by clients to determine when objects have changed."},
		{Name: "generation", Type: proto.ColumnType_INT, Description: "A sequence number representing a specific generation of the desired state."},
		{Name: "creation_timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromGo().Transform(v1TimeToRFC3339), Description: "CreationTimestamp is a timestamp representing the server time when this object was created."},
		{Name: "deletion_timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromGo().Transform(v1TimeToRFC3339), Description: "DeletionTimestamp is RFC 3339 date and time at which this resource will be deleted."},
		{Name: "deletion_grace_period_seconds", Type: proto.ColumnType_INT, Description: "Number of seconds allowed for this object to gracefully terminate before it will be removed from the system.  Only set when deletionTimestamp is also set."},
		{Name: "labels", Type: proto.ColumnType_JSON, Description: "Map of string keys and values that can be used to organize and categorize (scope and select) objects. May match selectors of replication controllers and services."},
		{Name: "annotations", Type: proto.ColumnType_JSON, Description: "Annotations is an unstructured key value map stored with a resource that may be set by external tools to store and retrieve arbitrary metadata."},
		{Name: "owner_references", Type: proto.ColumnType_JSON, Description: "List of objects depended by this object. If ALL objects in the list have been deleted, this object will be garbage collected. If this object is managed by a controller, then an entry in this list will point to this controller, with the controller field set to true. There cannot be more than one managing controller."},
		{Name: "finalizers", Type: proto.ColumnType_JSON, Description: "Must be empty before the object is deleted from the registry. Each entry is an identifier for the responsible component that will remove the entry from the list. If the deletionTimestamp of the object is non-nil, entries in this list can only be removed."},

		// DEPRECATED Kubernetes will stop propagating this field in 1.20 release and the field is planned to be removed in 1.21 release.
		// {Name: "self_link", Type: proto.ColumnType_STRING, Description: "SelfLink is a URL representing this object."},

		// Per https://v1-18.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#objectmeta-v1-meta
		//   This field is not set anywhere right now and apiserver is going to ignore it if set in create or update request.
		//{Name: "cluster_name", Type: proto.ColumnType_STRING, Description: "The name of the cluster which the object belongs to."},

		// Since 'users typically shouldn't need to set or understand this field.' we will omit it until/unless
		// we need it, or someone requests it.
		//{Name: "managed_fields", Type: proto.ColumnType_JSON, Description: "ManagedFields maps workflow-id and version to the set of fields that are managed by that workflow. This is mostly for internal housekeeping, and users typically shouldn't need to set or understand this field."},
	}
}

// append the common kubernetes columns for REGIONAL resources onto the column list
func k8sCommonColumns(columns []*plugin.Column) []*plugin.Column {
	allColumns := objectMetadataPrimaryColumns()
	allColumns = append(allColumns, columns...)
	//allColumns = append(allColumns, typeMetaColumns...)
	//allColumns = append(allColumns, specStatusColumns...)
	allColumns = append(allColumns, objectMetadataSecondaryColumns()...)
	allColumns = append(allColumns, manifestResourceColumns()...)

	return allColumns
}

// append the common kubernetes columns for CRD resources
func k8sCRDResourceCommonColumns(columns []*plugin.Column) []*plugin.Column {
	allColumns := []*plugin.Column{
		{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of resource."},
		{Name: "uid", Type: proto.ColumnType_STRING, Description: "UID is the unique in time and space value for this object.", Transform: transform.FromField("UID")},
		{Name: "kind", Type: proto.ColumnType_STRING, Description: "Type of resource."},
		{Name: "api_version", Type: proto.ColumnType_STRING, Description: "The API version of the resource.", Transform: transform.FromField("APIVersion")},
		{Name: "namespace", Type: proto.ColumnType_STRING, Description: "Namespace defines the space within which each name must be unique."},
		{Name: "creation_timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromGo().Transform(v1TimeToRFC3339), Description: "CreationTimestamp is a timestamp representing the server time when this object was created."},
		{Name: "labels", Type: proto.ColumnType_JSON, Description: "Map of string keys and values that can be used to organize and categorize (scope and select) objects. May match selectors of replication controllers and services."},
		{Name: "context_name", Type: proto.ColumnType_STRING, Description: "Kubectl config context name.", Transform: transform.FromField("ContextName").Transform(transform.NullIfZeroValue)},
		{Name: "source_type", Type: proto.ColumnType_STRING, Description: "The source of the resource. Possible values are: deployed and manifest. If the resource is fetched from the spec file the value will be manifest."},
	}
	allColumns = append(allColumns, columns...)
	allColumns = append(allColumns, manifestResourceColumns()...)

	return allColumns
}

// append the common kubernetes columns for non-namespaced resources onto the column list
func k8sCommonGlobalColumns(columns []*plugin.Column) []*plugin.Column {
	allColumns := objectMetadataPrimaryColumnsWithoutNamespace()
	allColumns = append(allColumns, columns...)
	//allColumns = append(allColumns, typeMetaColumns...)
	//allColumns = append(allColumns, specStatusColumns...)
	allColumns = append(allColumns, objectMetadataSecondaryColumns()...)
	allColumns = append(allColumns, manifestResourceColumns()...)

	return allColumns
}
