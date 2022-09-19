/*
package kubernetes implements a steampipe plugin for kubernetes.

This plugin provides data that Steampipe uses to present foreign
tables that represent kubernetes resources.
*/
package kubernetes

import (
	"context"
	"regexp"

	"github.com/iancoleman/strcase"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

const pluginName = "steampipe-plugin-kubernetes"

type contextKey string

// Plugin creates this (k8s) plugin
func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name:             pluginName,
		DefaultTransform: transform.FromGo(),
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
			Schema:      ConfigSchema,
		},
		SchemaMode:   plugin.SchemaModeDynamic,
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
		"kubernetes_cronjob":                 tableKubernetesCronJob(ctx),
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
		"kubernetes_pod_disruption_budget":   tableKubernetesPDB(ctx),
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
		"kubernetes_crd":                     tableKubernetesCRD(ctx),
	}

	// Search for metrics to create as tables
	var re = regexp.MustCompile(`\d+`)
	var substitution = ``
	kubernetesTables := []string{}
	config := GetConfig(p.Connection)
	if config.CustomResources != nil && len(*config.CustomResources) > 0 {
		for _, tableName := range *config.CustomResources {
			pluginTableName := "kubernetes_crd_" + strcase.ToSnake(re.ReplaceAllString(tableName, substitution))
			if _, ok := tables[pluginTableName]; !ok {
				kubernetesTables = append(kubernetesTables, tableName)
			}
		}
	}

	plugin.Logger(ctx).Error("tableKubernetesCRDResource", "kubernetesTables", kubernetesTables[0])
	for _, i := range kubernetesTables {
		tableName := "kubernetes_crd_" + strcase.ToSnake(re.ReplaceAllString(i, substitution))
		ctx = context.WithValue(ctx, contextKey("CustomResourceName"), i)
		ctx = context.WithValue(ctx, contextKey("PluginTableName"), tableName)
		plugin.Logger(ctx).Error("tableKubernetesCRDResource", "tableName", tableName)
		// Add the table if it does not already exist, ensuring standard tables win
		if tables[tableName] == nil {
			tables[tableName] = tableKubernetesCRDResource(ctx)
		} else {
			plugin.Logger(ctx).Error("tableKubernetesCRDResource", "table_already_exists", tableName)
		}
	}

	return tables, nil
}
