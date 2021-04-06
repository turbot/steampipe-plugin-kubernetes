/*
package kubernetes implements a steampipe plugin for kubernetes.

This plugin provides data that Steampipe uses to present foreign
tables that represent kubernetes resources.
*/
package kubernetes

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

const pluginName = "steampipe-plugin-kubernetes"

// Plugin creates this (k8s) plugin
func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name:             pluginName,
		DefaultTransform: transform.FromGo(),
		// DefaultGetConfig: &plugin.GetConfig{
		// 	ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFoundException", "NoSuchEntity"}),
		// },
		TableMap: map[string]*plugin.Table{
			"kubernetes_cluster_role":         tableKubernetesClusterRole(ctx),
			"kubernetes_cluster_role_binding": tableKubernetesClusterRoleBinding(ctx),
			"kubernetes_deployment":           tableKubernetesDeployment(ctx),
			"kubernetes_namespace":            tableKubernetesNamespace(ctx),
			"kubernetes_node":                 tableKubernetesNode(ctx),
			"kubernetes_pod":                  tableKubernetesPod(ctx),
			"kubernetes_replicaset":           tableKubernetesReplicaSet(ctx),
		},
	}

	return p
}
