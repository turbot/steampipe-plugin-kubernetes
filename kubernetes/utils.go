package kubernetes

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	filehelpers "github.com/turbot/go-kit/files"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextension "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
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

	if kubeconfig == nil {
		return nil, nil
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

	if kubeconfig == nil {
		return nil, nil
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

	if kubeconfig == nil {
		return nil, nil
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

	if kubeconfig == nil {
		return nil, nil
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

	if err := validateKubernetesConfig(d.Connection); err != nil {
		plugin.Logger(ctx).Error("getK8Config", "connection_config_error", "connection", d.Connection.Name, "error", err)
		return nil, err
	}
	if isManifestFilePathsConfigured := isManifestFilePathsConfigured(d.Connection); isManifestFilePathsConfigured {
		return nil, nil
	}

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

	if err := validateKubernetesConfig(c); err != nil {
		plugin.Logger(ctx).Error("getK8ConfigRaw", "connection_config_error", "connection", c.Name, "error", err)
		return nil, err
	}
	if isManifestFilePathsConfigured := isManifestFilePathsConfigured(c); isManifestFilePathsConfigured {
		return nil, nil
	}

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

	if kubeconfig == nil {
		return nil, nil
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

// Check whether the config file is configured correctly.
// The plugin can list both live resources and the resources from the manifest file.
// To avoid the conflict, a connection can only contains either config_path / config_paths to define the kube config files, or,
// the manifest_file_paths to define the location to the manifest files.
func validateKubernetesConfig(connection *plugin.Connection) error {
	kubernetesConfig := GetConfig(connection)

	if (kubernetesConfig.ConfigPath != nil || kubernetesConfig.ConfigPaths != nil) &&
		kubernetesConfig.ManifestFilePaths != nil {
		return fmt.Errorf("both config_path, config_paths and manifest_file_paths can not be passed in a single connection")
	}
	return nil
}

func isManifestFilePathsConfigured(connection *plugin.Connection) bool {
	kubernetesConfig := GetConfig(connection)

	return kubernetesConfig.ManifestFilePaths != nil
}

func fetchResourceFromManifestFileByKind(ctx context.Context, d *plugin.QueryData, kind string) ([]parsedContent, error) {

	if kind == "" {
		return nil, fmt.Errorf("missing required property: kind")
	}

	var data []parsedContent
	parsedContents, err := getParsedManifestFileContent(ctx, d)
	if err != nil {
		return nil, err
	}

	for _, content := range parsedContents {
		if content.GroupVersionKind.Kind == kind {
			data = append(data, content)
		}
	}

	return data, nil
}

type parsedContent struct {
	Data             runtime.Object
	GroupVersionKind *schema.GroupVersionKind
	Path             string
	Line             int
}

func getParsedManifestFileContent(ctx context.Context, d *plugin.QueryData) ([]parsedContent, error) {
	conn, err := parsedManifestFileContentCached(ctx, d, nil)
	if err != nil {
		return nil, err
	}
	return conn.([]parsedContent), nil
}

var parsedManifestFileContentCached = plugin.HydrateFunc(parsedManifestFileContentUncached).Memoize()

func parsedManifestFileContentUncached(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("parsedManifestFileContentUncached", "Parsing file content...", "connection", d.Connection.Name)

	// Read the config
	k8sConfig := GetConfig(d.Connection)

	// Gather file path matches for the glob
	var matches, resolvedPaths []string
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
		resolvedPaths = append(resolvedPaths, i)
	}

	var parsedContents []parsedContent
	for _, path := range resolvedPaths {
		// Load the file into a buffer
		content, err := os.ReadFile(path)
		if err != nil {
			plugin.Logger(ctx).Error("parsedManifestFileContentUncached", "failed to read file", err, "path", path)
			return nil, err
		}
		decoder := clientgoscheme.Codecs.UniversalDeserializer()

		// Get the start lines for all resource specified in the file
		temp := strings.Split(string(content), "\n")
		resourceLocationMap := map[string][]int{}

		// Initialize the line with 0
		line := 0
		for _, item := range temp {
			item = strings.ReplaceAll(item, " ", "")

			// Extract the starting position of the resource defined in a file.
			// If the file has multiple resource configuration defined, store the position based on the occurrence.
			if strings.Contains(item, "kind") {
				splitStr := strings.Split(item, "kind:")
				if len(splitStr) == 2 {
					resourceLocationMap[splitStr[1]] = append(resourceLocationMap[splitStr[1]], line)
				}
			}
			line++
		}
		plugin.Logger(ctx).Debug("parsedManifestFileContentUncached", "file", path, "resource locations", resourceLocationMap)

		// Check for the start of the document
		for _, resource := range strings.Split(string(content), "---") {
			// skip empty documents, `Decode` will fail on them
			if len(resource) == 0 {
				continue
			}

			// Decode the file content
			obj, groupVersionKind, err := decodeFileContent(decoder, resource)
			if err != nil {
				plugin.Logger(ctx).Error("parsedManifestFileContentUncached", "failed to decode the file", err, "path", path)
				return nil, err
			}

			// Extract the start line number based on the resource kind
			// If the file has multiple resource configuration defined,
			// location of the resource will be fetched as per FIFO basis
			//
			// For example if a file content is
			// ---
			// kind: Deployment
			// ...
			// ---
			// kind: Pod
			// ...
			// ---
			// Kind: Deployment
			//
			// The value of resourceLocationMap will be map[Deployment:[2 8] Pod:[5]]
			line := 0
			if resourceLocationMap[groupVersionKind.Kind] != nil {
				arr := resourceLocationMap[groupVersionKind.Kind]
				if len(arr) > 0 {
					line = arr[0]
					resourceLocationMap[groupVersionKind.Kind] = arr[1:]
				}
			}
			parsedContents = append(parsedContents, parsedContent{
				Data:             obj,
				GroupVersionKind: groupVersionKind,
				Path:             path,
				Line:             line + 1, // Since starts from 0
			})
		}
	}

	return parsedContents, nil
}

func decodeFileContent(decoder runtime.Decoder, resource string) (runtime.Object, *schema.GroupVersionKind, error) {
	// Decode the file content
	obj, groupVersionKind, err := decoder.Decode([]byte(resource), nil, nil)
	if err != nil {
		// The resource kind CustomResourceDefinition is not valid.
		// For the CustomResourceDefinition it needed to add the scheme.
		if strings.HasPrefix(err.Error(), "no kind \"CustomResourceDefinition\"") {
			sch := runtime.NewScheme()
			_ = clientgoscheme.AddToScheme(sch)
			_ = apiextensionsv1.AddToScheme(sch)

			decoder := serializer.NewCodecFactory(sch).UniversalDeserializer()

			return decodeFileContent(decoder, resource)
		}
		return nil, nil, err
	}

	return obj, groupVersionKind, nil
}
