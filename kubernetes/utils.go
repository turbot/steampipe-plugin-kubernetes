package kubernetes

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	"github.com/mitchellh/go-homedir"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func GetNewClientset(ctx context.Context, d *plugin.QueryData) (*kubernetes.Clientset, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("GetNewClientset")

	// have we already created and cached the session?
	serviceCacheKey := "k8sClient" //should probably per connection/context keys...

	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		logger.Warn("!!!! Clientset Found in Cache !!!!")
		return cachedData.(*kubernetes.Clientset), nil
	}

	kubeconfig, err := getK8Config(ctx, d)
	if err != nil {
		return nil, err
	}

	// Get a rest.Config from the kubeconfig file.
	restconfig, err := kubeconfig.ClientConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(restconfig)
	if err != nil {
		return nil, err
	}

	// save client in cache
	d.ConnectionManager.Cache.Set(serviceCacheKey, clientset)

	// logger.Warn("@@@@@@@  GetNewClientset SET cache status ", "success", success)
	// time.Sleep(5000 * time.Millisecond)
	if value, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		logger.Warn("!!!! Clientset added to cache !!!!")
	} else {
		logger.Warn("!!!! Clientset NOT Found in Cache after adding !!!!", "serviceCacheKey", serviceCacheKey, "Value", value)
	}

	return clientset, err
}

func v1TimeToRFC3339(_ context.Context, d *transform.TransformData) (interface{}, error) {
	if d.Value == nil {
		return nil, nil
	}

	switch v := d.Value.(type) {
	case v1.Time:
		return v.ToUnstructured(), nil
	case *v1.Time:
		if v == nil {
			return nil, nil
		}
		return v.ToUnstructured(), nil
	default:
		return nil, fmt.Errorf("Invalid time format %T!\n", v)
	}
}

func isNotFoundError(err error) bool {
	return strings.HasSuffix(err.Error(), "not found")
}

func mergeTags(labels map[string]string, annotations map[string]string) map[string]string {
	tags := make(map[string]string)
	for k, v := range annotations {
		tags[k] = v
	}
	for k, v := range labels {
		tags[k] = v
	}
	return tags
}

//// HYDRATE FUNCTIONS

func getKubectlContext(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	cacheKey := "getKubectlContext"
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		plugin.Logger(ctx).Warn("getKubectlContext", "######## CACHED CURRENT CONTEXT", cachedData.(string))
		return cachedData.(string), nil
	}

	kubeconfig, err := getK8Config(ctx, d)
	if err != nil {
		return nil, err
	}

	rawConfig, _ := kubeconfig.RawConfig()
	currentContext := rawConfig.CurrentContext

	// get kubernetes config info
	kubernetesConfig := GetConfig(d.Connection)

	// If set in plugin's (~/.steampipe/config/kubernetes.spc) connection profile
	if kubernetesConfig.ConfigContext != nil {
		currentContext = *kubernetesConfig.ConfigContext
	}

	plugin.Logger(ctx).Warn("getKubectlContext", "######## CURRENT CONTEXT", currentContext)

	// save current context in cache
	d.ConnectionManager.Cache.Set(cacheKey, currentContext)

	return currentContext, nil
}

//// TRANSFORM FUNCTIONS

func ensureStringArray(_ context.Context, d *transform.TransformData) (interface{}, error) {
	if d.Value == nil {
		// Should we return empty list instead???
		// return []string{}, nil

		return nil, fmt.Errorf("ensureStringArray - Cannot transform nil value")
	}

	switch v := d.Value.(type) {
	case []string:
		return v, nil
	case string:
		return []string{v}, nil
	default:
		str := fmt.Sprintf("%v", d.Value)
		return []string{string(str)}, nil
	}
}

// Get kubernetes config based on environment variable and plugin config
func getK8Config(ctx context.Context, d *plugin.QueryData) (clientcmd.ClientConfig, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getK8Config")

	// have we already created and cached the session?
	cacheKey := "getK8Config" //should probably per connection/context keys...

	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		logger.Warn("!!!! Clientset Found in Cache !!!!")
		return cachedData.(clientcmd.ClientConfig), nil
	}

	// get kubernetes config info
	kubernetesConfig := GetConfig(d.Connection)

	// Set default loader and overriding rules
	loader := &clientcmd.ClientConfigLoadingRules{}
	overrides := &clientcmd.ConfigOverrides{}

	// variable to store paths for kubernetes config
	configPaths := []string{}
	if kubernetesConfig.ConfigPath != nil {
		configPaths = []string{*kubernetesConfig.ConfigPath}
	} else if kubernetesConfig.ConfigPaths != nil && len(kubernetesConfig.ConfigPaths) > 0 {
		configPaths = kubernetesConfig.ConfigPaths
	} else if v := os.Getenv("KUBE_CONFIG_PATHS"); v != "" {
		// NOTE we have to do this here because the schema
		// does not yet allow you to set a default for a TypeList
		configPaths = filepath.SplitList(v)
	} else {
		// default path
		configPaths = []string{"~/.kube/config"}
		// Error: invalid configuration: no configuration has been provided, try setting KUBERNETES_MASTER environment variable

	}

	if len(configPaths) > 0 {
		expandedPaths := []string{}
		for _, p := range configPaths {
			path, err := homedir.Expand(p)
			if err != nil {
				return nil, err
			}

			logger.Debug("GetNewClientset", "Using kubeconfig: %s", path)
			expandedPaths = append(expandedPaths, path)
		}

		if len(expandedPaths) == 1 {
			loader.ExplicitPath = expandedPaths[0]
		} else {
			loader.Precedence = expandedPaths
		}

		ctxSuffix := "; default context"

		// if kubernetesConfig.ConfigContext != nil {
		// 	kubectx = *kubernetesConfig.ConfigContext
		// }
		// kubectx, ctxOk := d.GetOk("config_context")
		// authInfo, authInfoOk := d.GetOk("config_context_auth_info")
		// cluster, clusterOk := d.GetOk("config_context_cluster")
		// if ctxOk || authInfoOk || clusterOk {
		if kubernetesConfig.ConfigContext != nil {
			ctxSuffix = "; overriden context"
			// if ctxOk {
			overrides.CurrentContext = *kubernetesConfig.ConfigContext
			ctxSuffix += fmt.Sprintf("; config ctx: %s", overrides.CurrentContext)
			logger.Debug("GetNewClientset", "Using custom current context: %q", overrides.CurrentContext)
			// }

			overrides.Context = clientcmdapi.Context{}
			// if authInfoOk {
			// 	overrides.Context.AuthInfo = authInfo.(string)
			// 	ctxSuffix += fmt.Sprintf("; auth_info: %s", overrides.Context.AuthInfo)
			// }
			// if clusterOk {
			// 	overrides.Context.Cluster = cluster.(string)
			// 	ctxSuffix += fmt.Sprintf("; cluster: %s", overrides.Context.Cluster)
			// }
			logger.Debug("GetNewClientset", "Using overidden context: %#v", overrides.Context)
		}
	}

	// Instantiate loader for kubeconfig file.
	// kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
	// 	clientcmd.NewDefaultClientConfigLoadingRules(),
	// 	&clientcmd.ConfigOverrides{},
	// )

	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loader, overrides)

	// save the config in cache
	d.ConnectionManager.Cache.Set(cacheKey, kubeconfig)

	return kubeconfig, nil
}
