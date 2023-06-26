package kubernetes

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"k8s.io/apimachinery/pkg/version"
	// "k8s.io/component-base/version"
)

func tableKubernetesVersion(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_version",
		Description: "Obtain information about the version of the Kubernetes.",
		List: &plugin.ListConfig{
			Hydrate: listK8sVersion,
		},
		Columns: []*plugin.Column{
			{
				Name:        "component",
				Type:        proto.ColumnType_STRING,
				Description: "The type of the version information.",
			},
			{
				Name:        "git_version",
				Type:        proto.ColumnType_STRING,
				Description: "The full git version tag of the Kubernetes cluster, including additional information such as commit hash, build date, and build environment.",
			},
			{
				Name:        "major",
				Type:        proto.ColumnType_STRING,
				Description: "The major version number of the Kubernetes cluster.",
			},
			{
				Name:        "minor",
				Type:        proto.ColumnType_STRING,
				Description: "The minor version number of the Kubernetes cluster.",
			},
			{
				Name:        "git_commit",
				Type:        proto.ColumnType_STRING,
				Description: "The git commit hash of the Kubernetes cluster.",
			},
			{
				Name:        "git_tree_state",
				Type:        proto.ColumnType_STRING,
				Description: "The state of the git tree of the Kubernetes cluster, indicating whether it is clean or has uncommitted changes.",
			},
			{
				Name:        "build_date",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The date and time when the Kubernetes cluster was built.",
			},
			{
				Name:        "go_version",
				Type:        proto.ColumnType_STRING,
				Description: "The version of the Go programming language used to build the Kubernetes cluster.",
			},
			{
				Name:        "compiler",
				Type:        proto.ColumnType_STRING,
				Description: "The Go compiler used to build the Kubernetes cluster.",
			},
			{
				Name:        "platform",
				Type:        proto.ColumnType_STRING,
				Description: "The platform or operating system on which the Kubernetes cluster is running.",
			},
		},
	}
}

type VersionInfo struct {
	Component string
	version.Info
}

//// HYDRATE FUNCTIONS

func listK8sVersion(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listK8sVersion")

	// Get the client for querying the K8s APIs for the provided context.
	// If the connection is configured for the manifest files, the client will return nil.
	clientset, err := GetNewClientset(ctx, d)
	if err != nil {
		return nil, err
	}

	// Check for deployed resources
	if clientset == nil {
		return nil, nil
	}

	// Get the server's version
	response, err := clientset.Discovery().ServerVersion()
	if err != nil {
		return nil, err
	}
	d.StreamListItem(ctx, VersionInfo{"server", *response})

	// Context can be cancelled due to manual cancellation or the limit has been hit
	if d.RowsRemaining(ctx) == 0 {
		return nil, nil
	}

	// // Get the client version
	// // The method mentioned below not returning a valid client version.
	// // Issue link: https://github.com/kubernetes/client-go/issues/1274
	// clientVersion := version.Get()
	// d.StreamListItem(ctx, clientVersion)
	// //

	return nil, nil
}
