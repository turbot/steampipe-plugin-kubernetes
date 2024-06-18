/*
package kubernetes implements a steampipe plugin for kubernetes.

This plugin provides data that Steampipe uses to present foreign
tables that represent kubernetes resources.
*/
package kubernetes

import (
	"context"
	"path"
	"regexp"
	"sort"
	"strings"

	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/connection"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
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
		ConnectionKeyColumns: []plugin.ConnectionKeyColumn{
			{
				Name:    "context_name",
				Hydrate: getKubectlContext,
			},
		},
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: func() interface{} { return &kubernetesConfig{} },
		},
		SchemaMode:   plugin.SchemaModeDynamic,
		TableMapFunc: pluginTableDefinitions,
	}

	return p
}

func pluginTableDefinitions(ctx context.Context, d *plugin.TableMapData) (map[string]*plugin.Table, error) {
	// Initialize tables
	tables := map[string]*plugin.Table{
		"helm_chart":                            tableHelmChart(ctx),
		"helm_release":                          tableHelmRelease(ctx),
		"helm_template":                         tableHelmTemplates(ctx),
		"helm_template_rendered":                tableHelmTemplateRendered(ctx),
		"helm_value":                            tableHelmValue(ctx),
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
		"kubernetes_pod_template":               tableKubernetesPodTemplate(ctx),
		"kubernetes_replicaset":                 tableKubernetesReplicaSet(ctx),
		"kubernetes_replication_controller":     tableKubernetesReplicaController(ctx),
		"kubernetes_resource_quota":             tableKubernetesResourceQuota(ctx),
		"kubernetes_role":                       tableKubernetesRole(ctx),
		"kubernetes_role_binding":               tableKubernetesRoleBinding(ctx),
		"kubernetes_secret":                     tableKubernetesSecret(ctx),
		"kubernetes_service":                    tableKubernetesService(ctx),
		"kubernetes_service_account":            tableKubernetesServiceAccount(ctx),
		"kubernetes_stateful_set":               tableKubernetesStatefulSet(ctx),
		"kubernetes_storage_class":              tableKubernetesStorageClass(ctx),
	}

	// Fetch available CRDs
	crds, err := listK8sDynamicCRDs(ctx, d.ConnectionCache, d.Connection)
	if err != nil {
		plugin.Logger(ctx).Error("listK8sDynamicCRDs", "crds", err)
		return tables, nil
	}

	for _, crd := range crds {
		ctx = context.WithValue(ctx, contextKey("CRDName"), crd.Name)
		ctx = context.WithValue(ctx, contextKey("CustomResourceName"), crd.Spec.Names.Plural)
		ctx = context.WithValue(ctx, contextKey("CustomResourceNameSingular"), crd.Spec.Names.Singular)
		ctx = context.WithValue(ctx, contextKey("GroupName"), crd.Spec.Group)
		for _, version := range crd.Spec.Versions {
			if version.Served {
				ctx = context.WithValue(ctx, contextKey("ActiveVersion"), version.Name)
				if version.Schema != nil && version.Schema.OpenAPIV3Schema != nil {
					ctx = context.WithValue(ctx, contextKey("VersionSchemaSpec"), version.Schema.OpenAPIV3Schema.Properties["spec"])
					ctx = context.WithValue(ctx, contextKey("VersionSchemaStatus"), version.Schema.OpenAPIV3Schema.Properties["status"])
					if len(version.Schema.OpenAPIV3Schema.Description) > 0 {
						ctx = context.WithValue(ctx, contextKey("VersionSchemaDescription"), strings.TrimSuffix(version.Schema.OpenAPIV3Schema.Description, ".")+".")
					}
				}
			}
		}

		// add the tables in snake case
		// if there is any name collision, plugin will create the dynamic tables in below order:
		// plugin will use singular name for the first one, e.g. kubernetes_certificate
		// plugin will use fully qualified names for the subsequent ones, e.g. kubernetes_certificate_cert_manager_io
		re := regexp.MustCompile(`[-.]`)
		if tables["kubernetes_"+crd.Spec.Names.Singular] == nil {
			ctx = context.WithValue(ctx, contextKey("TableName"), "kubernetes_"+crd.Spec.Names.Singular)
			tables["kubernetes_"+crd.Spec.Names.Singular] = tableKubernetesCustomResource(ctx)
		}
		ctx = context.WithValue(ctx, contextKey("TableName"), "kubernetes_"+crd.Spec.Names.Singular+"_"+re.ReplaceAllString(crd.Spec.Group, "_"))
		tables["kubernetes_"+crd.Spec.Names.Singular+"_"+re.ReplaceAllString(crd.Spec.Group, "_")] = tableKubernetesCustomResource(ctx)
	}

	return tables, nil
}

func listK8sDynamicCRDs(ctx context.Context, cn *connection.ConnectionCache, c *plugin.Connection) ([]v1.CustomResourceDefinition, error) {
	// get the crds from config if any
	kubernetesConfig := GetConfig(c)
	filterCrds := kubernetesConfig.CustomResourceTables
	if len(filterCrds) == 0 {
		return nil, nil
	}

	clientset, err := GetNewClientCRDRaw(ctx, cn, c)
	if err != nil {
		plugin.Logger(ctx).Error("listK8sDynamicCRDs", "GetNewClientCRDRaw", err)

		// At the plugin load time, if the config file does not contain valid properties, return nil
		return nil, nil
	}

	crds := []v1.CustomResourceDefinition{}
	temp_crd_names := []string{}

	// Build the queryData
	queryData := &plugin.QueryData{
		Connection:      c,
		ConnectionCache: cn,
	}

	// Parse the manifest file content based on the kind, e.g. CustomResourceDefinition
	parsedContents, err := fetchResourceFromManifestFileByKind(ctx, queryData, "CustomResourceDefinition")
	if err != nil {
		return nil, err
	}

	// Match the CRDs based on the pattern/list of custom resource name provided in the config.
	// Return only those CRDs which are matched with the pattern provided.
	// Also, skip the duplicate CRDs to avoid the conflicts.
	for _, pattern := range filterCrds {
		for _, item := range parsedContents {
			crd := item.ParsedData.(*v1.CustomResourceDefinition)

			if helpers.StringSliceContains(temp_crd_names, crd.Name) {
				continue
			} else if ok, _ := path.Match(pattern, crd.Name); ok {
				crds = append(crds, *crd)
				temp_crd_names = append(temp_crd_names, crd.Name)
			} else if ok, _ := path.Match(pattern, crd.Spec.Names.Singular); ok {
				crds = append(crds, *crd)
				temp_crd_names = append(temp_crd_names, crd.Name)
			}
		}
	}

	// clientset provides the kubernetes API client to list the deployed resources.
	// If the source_type is set to "manifest" in the config file, clientset value will be nil, since
	// the deployed resources are not expected when the source_type is "manifest".
	if clientset != nil {
		input := metav1.ListOptions{
			Limit:          500,
			TimeoutSeconds: types.Int64(5),
		}

		pageLeft := true
		for pageLeft {
			response, err := clientset.ApiextensionsV1().CustomResourceDefinitions().List(ctx, input)
			if err != nil {
				// At the plugin load time, if the config is not setup properly, return nil
				if strings.Contains(err.Error(), "/apis/apiextensions.k8s.io/v1/customresourcedefinitions?limit=500") {
					return nil, nil
				}
				plugin.Logger(ctx).Error("listK8sDynamicCRDs", "list_err", err)
				return nil, err
			}

			if response.GetContinue() != "" {
				input.Continue = response.Continue
			} else {
				pageLeft = false
			}

			for _, pattern := range filterCrds {
				for _, item := range response.Items {
					if helpers.StringSliceContains(temp_crd_names, item.Name) {
						continue
					} else if ok, _ := path.Match(pattern, item.Name); ok {
						crds = append(crds, item)
						temp_crd_names = append(temp_crd_names, item.Name)
					} else if ok, _ := path.Match(pattern, item.Spec.Names.Singular); ok {
						crds = append(crds, item)
						temp_crd_names = append(temp_crd_names, item.Name)
					}
				}
			}
		}
	}

	// the Kube API doesn't return CRDs in a consistent order, so sort here to guarantee consistent CRD table generation
	sort.SliceStable(crds[:], func(i, j int) bool {
		return crds[i].Name < crds[j].Name
	})

	return crds, nil
}
