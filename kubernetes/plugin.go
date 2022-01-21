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
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
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
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
			Schema:      ConfigSchema,
		},
		TableMapFunc: pluginTableDefinitions,
	}

	return p
}

func pluginTableDefinitions(ctx context.Context, p *plugin.Plugin) (map[string]*plugin.Table, error) {

	// Initialize tables
	tables := map[string]*plugin.Table{
		"kubernetes_cluster_role":            tableKubernetesClusterRole(ctx),
		"kubernetes_cluster_role_binding":    tableKubernetesClusterRoleBinding(ctx),
		"kubernetes_config_map":              tableKubernetesConfigMap(ctx),
		"kubernetes_daemonset":               tableKubernetesDaemonset(ctx),
		"kubernetes_deployment":              tableKubernetesDeployment(ctx),
		"kubernetes_endpoint":                tableKubernetesEndpoints(ctx),
		"kubernetes_endpoint_slice":          tableKubernetesEndpointSlice(ctx),
		"kubernetes_ingress":                 tableKubernetesIngress(ctx),
		"kubernetes_job":                     tableKubernetesJob(ctx),
		"kubernetes_limit_range":             tableKubernetesLimitRange(ctx),
		"kubernetes_namespace":               tableKubernetesNamespace(ctx),
		"kubernetes_network_policy":          tableKubernetesNetworkPolicy(ctx),
		"kubernetes_node":                    tableKubernetesNode(ctx),
		"kubernetes_persistent_volume":       tableKubernetesPersistentVolume(ctx),
		"kubernetes_persistent_volume_claim": tableKubernetesPersistentVolumeClaim(ctx),
		"kubernetes_pod":                     tableKubernetesPod(ctx),
		"kubernetes_pod_security_policy":     tableKubernetesPodSecurityPolicy(ctx),
		"kubernetes_replicaset":              tableKubernetesReplicaSet(ctx),
		"kubernetes_replication_controller":  tableKubernetesReplicaController(ctx),
		"kubernetes_resource_quota":          tableKubernetesResourceQuota(ctx),
		"kubernetes_role":                    tableKubernetesRole(ctx),
		"kubernetes_role_binding":            tableKubernetesRoleBinding(ctx),
		"kubernetes_secret":                  tableKubernetesSecret(ctx),
		"kubernetes_service":                 tableKubernetesService(ctx),
		"kubernetes_service_account":         tableKubernetesServiceAccount(ctx),
		"kubernetes_stateful_set":            tableKubernetesStatefulSet(ctx),
	}

	// Get a list of custom resource definitions configured in the server
	crdList, err := listCustomResourceDefinitions(ctx, p)
	if err != nil {
		return nil, err
	}
	for _, crd := range crdList.Items {
		tableCtx := context.WithValue(ctx, "crd", crd)
		// base := filepath.Base(i)
		// tableName := base[0 : len(base)-len(filepath.Ext(base))]
		// Add the table if it does not already exist, ensuring standard tables win
		tableObj := tableDynamicCRD(tableCtx, p)
		tables[tableObj.Name] = tableObj
	}

	return tables, nil
}

func listCustomResourceDefinitions(ctx context.Context, p *plugin.Plugin) (*v1beta1.CustomResourceDefinitionList, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listCustomResourceDefinitions")

	crdClientSet, err := GetNewCrdClientSet(ctx, &plugin.QueryData{Connection: p.Connection, ConnectionManager: p.ConnectionManager})
	if err != nil {
		logger.Error("kubernetes_dynamic_crd.listCustomResourceDefinitions", "get_client_set_error", err)
		return nil, err
	}

	crdList, err := crdClientSet.CustomResourceDefinitions().ListCustomResourceDefinition(ctx)
	if err != nil {
		logger.Error("kubernetes_dynamic_crd.listCustomResourceDefinitions", "list_crd_error", err)
		return nil, err
	}

	return crdList, nil
}
