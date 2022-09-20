/*
package kubernetes implements a steampipe plugin for kubernetes.

This plugin provides data that Steampipe uses to present foreign
tables that represent kubernetes resources.
*/
package kubernetes

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v3/connection"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	// Search for CRs to create as tables
	kubernetesTables := []string{}
	crds, err := listK8sDynamicCRDs(ctx, p.ConnectionManager, p.Connection)
	if err != nil {
		return nil, err
	}
	for _, crd := range crds {
		pluginTableName := "kubernetes_" + crd.Spec.Names.Plural
		if _, ok := tables[pluginTableName]; !ok {
			kubernetesTables = append(kubernetesTables, crd.Spec.Names.Plural)
		}
	}

	for _, kTable := range kubernetesTables {
		tableName := "kubernetes_" + kTable
		ctx = context.WithValue(ctx, contextKey("CustomResourceName"), kTable)
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

func listK8sDynamicCRDs(ctx context.Context, cm *connection.Manager, c *plugin.Connection) ([]v1.CustomResourceDefinition, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listK8sDynamicCRDs")

	clientset, err := GetNewClientCRDRaw(ctx, cm, c)
	if err != nil {
		return nil, err
	}

	input := metav1.ListOptions{
		Limit: 500,
	}

	crds := []v1.CustomResourceDefinition{}

	pageLeft := true
	for pageLeft {
		response, err := clientset.ApiextensionsV1().CustomResourceDefinitions().List(ctx, input)
		if err != nil {
			logger.Error("listK8sDynamicCRDs", "list_err", err)
			return nil, err
		}

		if response.GetContinue() != "" {
			input.Continue = response.Continue
		} else {
			pageLeft = false
		}

		for _, crd := range response.Items {
			crds = append(crds, crd)
		}
	}

	return crds, nil
}
