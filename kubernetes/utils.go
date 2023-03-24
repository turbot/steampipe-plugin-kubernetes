package kubernetes

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	filehelpers "github.com/turbot/go-kit/files"
	apiextension "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	"github.com/mitchellh/go-homedir"
	"github.com/turbot/steampipe-plugin-sdk/v5/connection"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// GetNewClientset :: gets client for querying k8s apis for the provided context
func GetNewClientset(ctx context.Context, d *plugin.QueryData) (*kubernetes.Clientset, error) {
	logger := plugin.Logger(ctx)

	// have we already created and cached the session?
	serviceCacheKey := "k8sClient"

	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*kubernetes.Clientset), nil
	}

	kubeconfig, err := getK8Config(ctx, d)
	if err != nil {
		logger.Error("GetNewClientset", "getK8Config", err)
		return nil, err
	}

	// Get a rest.Config from the kubeconfig file.
	restconfig, err := kubeconfig.ClientConfig()
	if err != nil {
		// if .kube/config file is not available check for inClusterConfig
		configErr := err
		if strings.Contains(err.Error(), ".kube/config: no such file or directory") {
			clientset, err := inClusterConfig(ctx)
			if err != nil {
				return nil, errors.New(configErr.Error() + ", " + err.Error())
			}

			// save clientset in cache
			d.ConnectionManager.Cache.Set(serviceCacheKey, clientset)

			return clientset, nil
		}

		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(restconfig)
	if err != nil {
		return nil, err
	}

	// save clientset in cache
	d.ConnectionManager.Cache.Set(serviceCacheKey, clientset)

	return clientset, err
}

func inClusterConfig(ctx context.Context) (*kubernetes.Clientset, error) {
	clusterConfig, err := rest.InClusterConfig()
	if err != nil {
		plugin.Logger(ctx).Error("InClusterConfig", "InClusterConfig", err)
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(clusterConfig)
	if err != nil {
		plugin.Logger(ctx).Error("InClusterConfig", "NewForConfig", err)
		return nil, err
	}

	return clientset, nil
}

// GetNewClientCRD :: gets client for querying k8s apis for CustomResourceDefinition
func GetNewClientCRD(ctx context.Context, d *plugin.QueryData) (*apiextension.Clientset, error) {
	// have we already created and cached the session?
	serviceCacheKey := "GetNewClientCRD"

	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*apiextension.Clientset), nil
	}

	kubeconfig, err := getK8Config(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("GetNewClientCRD", "getK8Config", err)
		return nil, err
	}

	// Get a rest.Config from the kubeconfig file.
	restconfig, err := kubeconfig.ClientConfig()
	if err != nil {
		// if .kube/config file is not available check for inClusterConfig
		configErr := err
		if strings.Contains(err.Error(), ".kube/config: no such file or directory") {
			clientset, err := inClusterConfigCRD(ctx)
			if err != nil {
				return nil, errors.New(configErr.Error() + ", " + err.Error())
			}

			// save clientset in cache
			d.ConnectionManager.Cache.Set(serviceCacheKey, clientset)

			return clientset, nil
		}

		return nil, err
	}

	clientset, err := apiextension.NewForConfig(restconfig)
	if err != nil {
		plugin.Logger(ctx).Error("GetNewClientCRD", "NewForConfig", err)
		return nil, err
	}

	// save clientset in cache
	d.ConnectionManager.Cache.Set(serviceCacheKey, clientset)

	return clientset, err
}

// GetNewClientCRDRaw :: gets client for querying k8s apis for CustomResourceDefinition
func GetNewClientCRDRaw(ctx context.Context, cc *connection.ConnectionCache, c *plugin.Connection) (*apiextension.Clientset, error) {
	logger := plugin.Logger(ctx)

	// have we already created and cached the session?
	serviceCacheKey := "GetNewClientCRDRaw"

	if cachedData, ok := cc.Get(ctx, serviceCacheKey); ok {
		return cachedData.(*apiextension.Clientset), nil
	}

	kubeconfig, err := getK8ConfigRaw(ctx, cc, c)
	if err != nil {
		logger.Error("GetNewClientCRDRaw", "getK8ConfigRaw", err)
		return nil, err
	}

	// Get a rest.Config from the kubeconfig file.
	restconfig, err := kubeconfig.ClientConfig()
	if err != nil {
		// if .kube/config file is not available check for inClusterConfig
		configErr := err
		if strings.Contains(err.Error(), ".kube/config: no such file or directory") {
			clientset, err := inClusterConfigCRD(ctx)
			if err != nil {
				return nil, errors.New(configErr.Error() + ", " + err.Error())
			}

			// save clientset in cache
			cacheErr := cc.Set(ctx, serviceCacheKey, clientset)
			if cacheErr != nil {
				plugin.Logger(ctx).Error("inClusterConfigCRD", "cache-set", cacheErr)
				return nil, err
			}

			return clientset, nil
		}

		return nil, err
	}

	clientset, err := apiextension.NewForConfig(restconfig)
	if err != nil {
		return nil, err
	}

	// save clientset in cache
	cacheErr := cc.Set(ctx, serviceCacheKey, clientset)
	if cacheErr != nil {
		plugin.Logger(ctx).Error("GetNewClientCRDRaw", "cache-set", cacheErr)
		return nil, err
	}

	return clientset, nil
}

func inClusterConfigCRD(ctx context.Context) (*apiextension.Clientset, error) {
	clusterConfig, err := rest.InClusterConfig()
	if err != nil {
		plugin.Logger(ctx).Error("inClusterConfigCRD", "InClusterConfig", err)
		return nil, err
	}

	clientset, err := apiextension.NewForConfig(clusterConfig)
	if err != nil {
		plugin.Logger(ctx).Error("inClusterConfigCRD", "NewForConfig", err)
		return nil, err
	}

	return clientset, nil
}

// GetNewClientDynamic :: gets client for querying k8s apis for Dynamic Interface
func GetNewClientDynamic(ctx context.Context, d *plugin.QueryData) (dynamic.Interface, error) {
	// have we already created and cached the session?
	serviceCacheKey := "GetNewClientDynamic"

	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(dynamic.Interface), nil
	}

	kubeconfig, err := getK8Config(ctx, d)
	if err != nil {
		return nil, err
	}

	// Get a rest.Config from the kubeconfig file.
	restconfig, err := kubeconfig.ClientConfig()
	if err != nil {
		// if .kube/config file is not available check for inClusterConfig
		configErr := err
		if strings.Contains(err.Error(), ".kube/config: no such file or directory") {
			clientset, err := inClusterConfigCRDDynamic(ctx)
			if err != nil {
				return nil, errors.New(configErr.Error() + ", " + err.Error())
			}

			// save clientset in cache
			d.ConnectionManager.Cache.Set(serviceCacheKey, clientset)

			return clientset, nil
		}

		return nil, err
	}

	clientset, err := dynamic.NewForConfig(restconfig)
	if err != nil {
		return nil, err
	}

	// save clientset in cache
	d.ConnectionManager.Cache.Set(serviceCacheKey, clientset)

	return clientset, err

}

func inClusterConfigCRDDynamic(ctx context.Context) (dynamic.Interface, error) {
	clusterConfig, err := rest.InClusterConfig()
	if err != nil {
		plugin.Logger(ctx).Error("inClusterConfigCRDDynamic", "InClusterConfig", err)
		return nil, err
	}

	clientset, err := dynamic.NewForConfig(clusterConfig)
	if err != nil {
		plugin.Logger(ctx).Error("inClusterConfigCRDDynamic", "NewForConfig", err)
		return nil, err
	}

	return clientset, nil
}

// Get kubernetes config based on environment variable and plugin config
func getK8Config(ctx context.Context, d *plugin.QueryData) (clientcmd.ClientConfig, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getK8Config")

	// have we already created and cached the session?
	cacheKey := "getK8Config"

	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.(clientcmd.ClientConfig), nil
	}

	// get kubernetes config info
	kubernetesConfig := GetConfig(d.Connection)

	// Set default loader and overriding rules
	loader := &clientcmd.ClientConfigLoadingRules{}
	overrides := &clientcmd.ConfigOverrides{}

	// variable to store paths for kubernetes config
	// default kube config path
	var configPaths = []string{"~/.kube/config"}

	if kubernetesConfig.ConfigPath != nil {
		configPaths = []string{*kubernetesConfig.ConfigPath}
	} else if kubernetesConfig.ConfigPaths != nil && len(kubernetesConfig.ConfigPaths) > 0 {
		configPaths = kubernetesConfig.ConfigPaths
	} else if v := os.Getenv("KUBE_CONFIG_PATHS"); v != "" {
		configPaths = filepath.SplitList(v)
	} else if v := os.Getenv("KUBERNETES_MASTER"); v != "" {
		configPaths = []string{v}
	}

	if len(configPaths) > 0 {
		expandedPaths := []string{}
		for _, p := range configPaths {
			path, err := homedir.Expand(p)
			if err != nil {
				return nil, err
			}

			expandedPaths = append(expandedPaths, path)
		}

		if len(expandedPaths) == 1 {
			loader.ExplicitPath = expandedPaths[0]
		} else {
			loader.Precedence = expandedPaths
		}

		if kubernetesConfig.ConfigContext != nil {
			overrides.CurrentContext = *kubernetesConfig.ConfigContext
			overrides.Context = clientcmdapi.Context{}
		}
	}

	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loader, overrides)

	// save the config in cache
	d.ConnectionManager.Cache.Set(cacheKey, kubeconfig)

	return kubeconfig, nil
}

// Get kubernetes config based on environment variable and plugin config
func getK8ConfigRaw(ctx context.Context, cc *connection.ConnectionCache, c *plugin.Connection) (clientcmd.ClientConfig, error) {
	logger := plugin.Logger(ctx)

	// have we already created and cached the session?
	cacheKey := "getK8ConfigRaw"

	if cachedData, ok := cc.Get(ctx, cacheKey); ok {
		return cachedData.(clientcmd.ClientConfig), nil
	}

	// get kubernetes config info
	kubernetesConfig := GetConfig(c)

	// Set default loader and overriding rules
	loader := &clientcmd.ClientConfigLoadingRules{}
	overrides := &clientcmd.ConfigOverrides{}

	// variable to store paths for kubernetes config
	// default kube config path
	var configPaths = []string{"~/.kube/config"}

	if kubernetesConfig.ConfigPath != nil {
		configPaths = []string{*kubernetesConfig.ConfigPath}
	} else if kubernetesConfig.ConfigPaths != nil && len(kubernetesConfig.ConfigPaths) > 0 {
		configPaths = kubernetesConfig.ConfigPaths
	} else if v := os.Getenv("KUBE_CONFIG_PATHS"); v != "" {
		configPaths = filepath.SplitList(v)
	} else if v := os.Getenv("KUBERNETES_MASTER"); v != "" {
		configPaths = []string{v}
	}

	if len(configPaths) > 0 {
		expandedPaths := []string{}
		for _, p := range configPaths {
			path, err := homedir.Expand(p)
			if err != nil {
				return nil, err
			}

			expandedPaths = append(expandedPaths, path)
		}

		if len(expandedPaths) == 1 {
			loader.ExplicitPath = expandedPaths[0]
		} else {
			loader.Precedence = expandedPaths
		}

		if kubernetesConfig.ConfigContext != nil {
			overrides.CurrentContext = *kubernetesConfig.ConfigContext
			overrides.Context = clientcmdapi.Context{}
		}
	}

	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loader, overrides)

	// save the config in cache
	err := cc.Set(ctx, cacheKey, kubeconfig)
	if err != nil {
		logger.Error("getK8ConfigRaw", "cache-set", err)
	}

	return kubeconfig, nil
}

//// HYDRATE FUNCTIONS

func getKubectlContext(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	cacheKey := "getKubectlContext"
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
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

	// save current context in cache
	d.ConnectionManager.Cache.Set(cacheKey, currentContext)

	return currentContext, nil
}

//// COMMON TRANSFORM FUNCTIONS

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
		return nil, fmt.Errorf("invalid time format %T! ", v)
	}
}

func v1MicroTimeToRFC3339(_ context.Context, d *transform.TransformData) (interface{}, error) {
	if d.Value == nil {
		return nil, nil
	}

	switch v := d.Value.(type) {
	case v1.MicroTime:
		return v1.NewTime(v.Time).ToUnstructured(), nil
	case *v1.MicroTime:
		if v == nil {
			return nil, nil
		}
		return v1.NewTime(v.Time).ToUnstructured(), nil
	default:
		return nil, fmt.Errorf("invalid time format %T! ", v)
	}
}

func labelSelectorToString(_ context.Context, d *transform.TransformData) (interface{}, error) {
	if d.Value == nil {
		return nil, nil
	}

	selector := d.Value.(*v1.LabelSelector)

	ss, err := v1.LabelSelectorAsSelector(selector)
	if err != nil {
		return nil, err
	}

	return ss.String(), nil
}

func selectorMapToString(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("selectorMapToString")

	selector_map := d.Value.(map[string]string)

	if len(selector_map) == 0 {
		return nil, nil
	}

	selector_string := labels.SelectorFromSet(selector_map).String()

	return selector_string, nil
}

//// Other Utility functions

func isNotFoundError(err error) bool {
	return strings.HasSuffix(err.Error(), "not found")
}

func getCommonOptionalKeyQuals() []*plugin.KeyColumn {
	return []*plugin.KeyColumn{
		{Name: "name", Require: plugin.Optional},
		{Name: "namespace", Require: plugin.Optional},
	}
}

func getOptionalKeyQualWithCommonKeyQuals(otherOptionalQuals []*plugin.KeyColumn) []*plugin.KeyColumn {
	return append(otherOptionalQuals, getCommonOptionalKeyQuals()...)
}

func getCommonOptionalKeyQualsValueForFieldSelector(d *plugin.QueryData) []string {
	fieldSelectors := []string{}

	if d.EqualsQualString("name") != "" {
		fieldSelectors = append(fieldSelectors, fmt.Sprintf("metadata.name=%v", d.EqualsQualString("name")))
	}

	if d.EqualsQualString("namespace") != "" {
		fieldSelectors = append(fieldSelectors, fmt.Sprintf("metadata.namespace=%v", d.EqualsQualString("namespace")))
	}

	return fieldSelectors
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

type filePath struct {
	Path string
}

func listKubernetesManifestFiles(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// #1 - Path via qual

	// If the path was requested through qualifier then match it exactly. Globs
	// are not supported in this context since the output value for the column
	// will never match the requested value.
	quals := d.EqualsQuals
	if quals["manifest_file_path"] != nil {
		d.StreamListItem(ctx, filePath{Path: quals["manifest_file_path"].GetStringValue()})
		return nil, nil
	}

	// #2 - paths in config

	// Glob paths in config
	// Fail if no paths are specified
	k8sConfig := GetConfig(d.Connection)
	if k8sConfig.ManifestFilePaths == nil {
		return nil, errors.New("manifest_file_path must be configured")
	}

	// Gather file path matches for the glob
	var matches []string
	paths := k8sConfig.ManifestFilePaths
	for _, i := range paths {

		// List the files in the given source directory
		files, err := d.GetSourceFiles(i)
		if err != nil {
			return nil, err
		}
		matches = append(matches, files...)
	}

	// Sanitize the matches to ignore the directories
	for _, i := range matches {

		// Ignore directories
		if filehelpers.DirectoryExists(i) {
			continue
		}
		d.StreamListItem(ctx, filePath{Path: i})
	}

	return nil, nil
}
