/*
package kubernetes implements a steampipe plugin for kubernetes.

This plugin provides data that Steampipe uses to present foreign
tables that represent kubernetes resources.
*/
package kubernetes

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

const pluginName = "steampipe-plugin-kubernetes"

// Uncomment once aggregator functionality available with dynamic tables
// type contextKey string

// Plugin creates this (k8s) plugin
func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name:             pluginName,
		DefaultTransform: transform.FromGo(),
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
			Schema:      ConfigSchema,
		},
		// TODO: Change to dynamic, once aggregator functionality available with dynamic tables
		SchemaMode:   plugin.SchemaModeStatic,
		TableMapFunc: pluginTableDefinitions,
	}

	return p
}

func pluginTableDefinitions(ctx context.Context, d *plugin.TableMapData) (map[string]*plugin.Table, error) {
	// Initialize tables
	tables := map[string]*plugin.Table{
		"kubernetes_cluster_role":               tableKubernetesClusterRole(ctx),
		"kubernetes_cluster_role_binding":       tableKubernetesClusterRoleBinding(ctx),
		"kubernetes_config_map":                 tableKubernetesConfigMap(ctx),
		"kubernetes_cronjob":                    tableKubernetesCronJob(ctx),
		"kubernetes_custom_resource_definition": tableKubernetesCustomResourceDefinition(ctx),
		"kubernetes_daemonset":                  tableKubernetesDaemonset(ctx),
		"kubernetes_deployment":                 tableKubernetesDeployment(ctx),
		"kubernetes_endpoint":                   tableKubernetesEndpoints(ctx),
		"kubernetes_endpoint_slice":             tableKubernetesEndpointSlice(ctx),
		"kubernetes_event":                      tableKubernetesEvent(ctx),
		"kubernetes_horizontal_pod_autoscaler":  tableKubernetesHorizontalPodAutoscaler(ctx),
		"kubernetes_ingress":                    tableKubernetesIngress(ctx),
		"kubernetes_job":                        tableKubernetesJob(ctx),
		"kubernetes_limit_range":                tableKubernetesLimitRange(ctx),
		"kubernetes_namespace":                  tableKubernetesNamespace(ctx),
		"kubernetes_network_policy":             tableKubernetesNetworkPolicy(ctx),
		"kubernetes_node":                       tableKubernetesNode(ctx),
		"kubernetes_persistent_volume":          tableKubernetesPersistentVolume(ctx),
		"kubernetes_persistent_volume_claim":    tableKubernetesPersistentVolumeClaim(ctx),
		"kubernetes_pod":                        tableKubernetesPod(ctx),
		"kubernetes_pod_disruption_budget":      tableKubernetesPDB(ctx),
		"kubernetes_pod_security_policy":        tableKubernetesPodSecurityPolicy(ctx),
		"kubernetes_replicaset":                 tableKubernetesReplicaSet(ctx),
		"kubernetes_replication_controller":     tableKubernetesReplicaController(ctx),
		"kubernetes_resource_quota":             tableKubernetesResourceQuota(ctx),
		"kubernetes_role":                       tableKubernetesRole(ctx),
		"kubernetes_role_binding":               tableKubernetesRoleBinding(ctx),
		"kubernetes_secret":                     tableKubernetesSecret(ctx),
		"kubernetes_service":                    tableKubernetesService(ctx),
		"kubernetes_service_account":            tableKubernetesServiceAccount(ctx),
		"kubernetes_stateful_set":               tableKubernetesStatefulSet(ctx),
	}

	// Fetch available CRDs
	// TODO: Re-enable once aggregator functionality works with dynamic tables
	// crds, err := listK8sDynamicCRDs(ctx, d.ConectionCache, d.Connection)
	// if err != nil {
	// 	plugin.Logger(ctx).Error("listK8sDynamicCRDs", "crds", err)
	// 	return nil, err
	// }

	// for _, crd := range crds {
	// 	ctx = context.WithValue(ctx, contextKey("CRDName"), crd.Name)
	// 	ctx = context.WithValue(ctx, contextKey("CustomResourceName"), crd.Spec.Names.Plural)
	// 	ctx = context.WithValue(ctx, contextKey("GroupName"), crd.Spec.Group)
	// 	for _, version := range crd.Spec.Versions {
	// 		if version.Served {
	// 			ctx = context.WithValue(ctx, contextKey("ActiveVersion"), version.Name)
	// 			ctx = context.WithValue(ctx, contextKey("VersionSchema"), version.Schema.OpenAPIV3Schema.Properties["spec"])
	// 		}
	// 	}
	// 	if tables[crd.Name] == nil {
	// 		tables[crd.Name] = tableKubernetesCustomResource(ctx)
	// 	}
	// }

	return tables, nil
}

// Uncomment once aggregator functionality available with dynamic tables

// func listK8sDynamicCRDs(ctx context.Context, cn *connection.ConnectionCache, c *plugin.Connection) ([]v1.CustomResourceDefinition, error) {
// 	clientset, err := GetNewClientCRDRaw(ctx, cn, c)
// 	if err != nil {
// 		plugin.Logger(ctx).Error("listK8sDynamicCRDs", "GetNewClientCRDRaw", err)
// 		return nil, err
// 	}

// 	input := metav1.ListOptions{
// 		Limit: 500,
// 	}

// 	crds := []v1.CustomResourceDefinition{}

// 	pageLeft := true
// 	for pageLeft {
// 		response, err := clientset.ApiextensionsV1().CustomResourceDefinitions().List(ctx, input)
// 		if err != nil {
// 			// Handle err at the plugin load time if config is not setup properly
// 			if strings.Contains(err.Error(), "/apis/apiextensions.k8s.io/v1/customresourcedefinitions?limit=500") {
// 				return nil, nil
// 			}
// 			plugin.Logger(ctx).Error("listK8sDynamicCRDs", "list_err", err)
// 			return nil, err
// 		}

// 		if response.GetContinue() != "" {
// 			input.Continue = response.Continue
// 		} else {
// 			pageLeft = false
// 		}

// 		crds = append(crds, response.Items...)
// 	}

// 	return crds, nil
// }
