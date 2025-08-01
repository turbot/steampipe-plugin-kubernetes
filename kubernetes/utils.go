package kubernetes

import (
	"context"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"

	goYaml "gopkg.in/yaml.v3"

	apiextension "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	"github.com/mitchellh/go-homedir"
	filehelpers "github.com/turbot/go-kit/files"
	"github.com/turbot/steampipe-plugin-sdk/v5/connection"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

type SourceType string

const (
	Deployed SourceType = "deployed"
	Helm     SourceType = "helm"
	Manifest SourceType = "manifest"
	All      SourceType = "all"
)

// Validate the source type.
func (sourceType SourceType) IsValid() error {
	switch sourceType {
	case Deployed, Helm, Manifest, All:
		return nil
	}
	return fmt.Errorf("invalid source type: %s", sourceType)
}

// Convert the source type to its string equivalent
func (sourceType SourceType) String() string {
	return string(sourceType)
}

// ToSourceTypes is used to convert SourceType to []string
func (sourceType SourceType) ToSourceTypes() []string {
	if sourceType == All {
		return []string{Deployed.String(), Helm.String(), Manifest.String()}
	} else {
		return []string{sourceType.String()}
	}
}

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

	// Return nil if deployed resources should not be included
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

	// Return nil if deployed resources should not be included
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

	// Return nil if deployed resources should not be included
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

	// Check for the sourceTypes argument in the config
	// Default set to include values
	var sources = All.ToSourceTypes()
	if kubernetesConfig.SourceTypes != nil {
		sources = kubernetesConfig.SourceTypes
	}
	// TODO: Remove once `SourceType` is obsolete
	if kubernetesConfig.SourceTypes == nil && kubernetesConfig.SourceType != nil {
		if *kubernetesConfig.SourceType != "all" { // if is all, sources is already set by default
			sources = []string{*kubernetesConfig.SourceType}
		}
	}

	if !slices.Contains(sources, "deployed") {
		plugin.Logger(ctx).Debug("getK8Config", "Returning nil for API server client.", "source_types", sources, "connection", d.Connection.Name)
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
	} else if len(kubernetesConfig.ConfigPaths) > 0 {
		plugin.Logger(ctx).Warn("config_paths parameter is deprecated and will be removed after 31st July 2023, please use config_path instead.")
		configPaths = kubernetesConfig.ConfigPaths
	} else if v := os.Getenv("KUBECONFIG"); v != "" {
		configPaths = []string{v}
	} else if v := os.Getenv("KUBE_CONFIG_PATH"); v != "" {
		configPaths = []string{v}
	} else if v := os.Getenv("KUBE_CONFIG_PATHS"); v != "" {
		plugin.Logger(ctx).Warn("KUBE_CONFIG_PATHS parameter is deprecated and will be removed after 31st July 2023, please use KUBECONFIG or KUBE_CONFIG_PATH instead.")
		configPaths = filepath.SplitList(v)
	} else if v := os.Getenv("KUBERNETES_MASTER"); v != "" {
		plugin.Logger(ctx).Warn("KUBERNETES_MASTER parameter is deprecated and will be removed after 31st July 2023, please use KUBECONFIG or KUBE_CONFIG_PATH instead.")
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

	// Check for the sourceTypes argument in the config.
	// Default set to include values.
	var sources = All.ToSourceTypes()
	if kubernetesConfig.SourceTypes != nil {
		sources = kubernetesConfig.SourceTypes
	}
	// TODO: Remove once `SourceType` is obsolete
	if kubernetesConfig.SourceTypes == nil && kubernetesConfig.SourceType != nil {
		if *kubernetesConfig.SourceType != "all" { // if is all, sources is already set by default
			sources = []string{*kubernetesConfig.SourceType}
		}
	}

	if !slices.Contains(sources, "deployed") {
		plugin.Logger(ctx).Debug("getK8ConfigRaw", "Returning nil for API server client.", "source_types", sources, "connection", c.Name)
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
	} else if len(kubernetesConfig.ConfigPaths) > 0 {
		plugin.Logger(ctx).Warn("config_paths parameter is deprecated and will be removed after 31st July 2023, please use config_path instead.")
		configPaths = kubernetesConfig.ConfigPaths
	} else if v := os.Getenv("KUBECONFIG"); v != "" {
		configPaths = []string{v}
	} else if v := os.Getenv("KUBE_CONFIG_PATH"); v != "" {
		configPaths = []string{v}
	} else if v := os.Getenv("KUBE_CONFIG_PATHS"); v != "" {
		plugin.Logger(ctx).Warn("KUBE_CONFIG_PATHS parameter is deprecated and will be removed after 31st July 2023, please use KUBECONFIG or KUBE_CONFIG_PATH instead.")
		configPaths = filepath.SplitList(v)
	} else if v := os.Getenv("KUBERNETES_MASTER"); v != "" {
		plugin.Logger(ctx).Warn("KUBERNETES_MASTER parameter is deprecated and will be removed after 31st July 2023, please use KUBECONFIG or KUBE_CONFIG_PATH instead.")
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

//// Utility functions for manifest files

type parsedContent struct {
	ParsedData any
	Kind       string
	Path       string
	SourceType string
	StartLine  int
	EndLine    int
}

// Returns the manifest file content based on the kind provided
func fetchResourceFromManifestFileByKind(ctx context.Context, d *plugin.QueryData, kind string) ([]parsedContent, error) {

	if kind == "" {
		return nil, fmt.Errorf("missing required property: kind")
	}
	var data []parsedContent

	// Get parsed content from manifest files
	parsedContents, err := getParsedManifestFileContent(ctx, d)
	if err != nil {
		return nil, err
	}

	// Get parsed content from rendered Helm templates
	renderedTemplateContents, err := getRenderedHelmTemplateContent(ctx, d)
	if err != nil {
		return nil, err
	}

	parsedContents = append(parsedContents, renderedTemplateContents...)
	for _, content := range parsedContents {
		if content.Kind == kind {
			data = append(data, content)
		}
	}

	return data, nil
}

// Get the parsed contents of the given files.
func getParsedManifestFileContent(ctx context.Context, d *plugin.QueryData) ([]parsedContent, error) {
	conn, err := parsedManifestFileContentCached(ctx, d, nil)
	if err != nil {
		return nil, err
	}
	return conn.([]parsedContent), nil
}

// Cached form of the parsed file content.
var parsedManifestFileContentCached = plugin.HydrateFunc(parsedManifestFileContentUncached).Memoize()

// parsedManifestFileContentUncached is the actual implementation of getParsedManifestFileContent, which should
// be run only once per connection. Do not call this directly, use
// getParsedManifestFileContent instead.
func parsedManifestFileContentUncached(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("parsedManifestFileContentUncached", "Parsing file content", "connection", d.Connection.Name)

	// Read the config
	resolvedPaths, err := resolveManifestFilePaths(ctx, d)
	if err != nil {
		return nil, err
	}

	var parsedContents []parsedContent
	for _, path := range resolvedPaths {
		// Load the file into a buffer
		content, err := os.ReadFile(path)
		if err != nil {
			plugin.Logger(ctx).Error("parsedManifestFileContentUncached", "failed to read file", err, "path", path)
			return nil, err
		}

		// Check for the start of the document
		pos := 0
		for _, resource := range strings.Split(string(content), "---") {
			// Skip empty documents, `Decode` will fail on them
			// Also, increment the pos to include the separator position (e.g. ---)
			if len(resource) == 0 {
				pos++
				continue
			}

			// Calculate the length of the YAML resource block
			blockLength := strings.Split(strings.ReplaceAll(resource, " ", ""), "\n")

			// Remove the extra lines added during the split operation based on the separator
			blockLength = blockLength[:len(blockLength)-1]
			if blockLength[0] == "" {
				blockLength = blockLength[1:]
			}

			// skip if no kind defined
			if !(strings.Contains(resource, "kind:") || strings.Contains(resource, "\"kind\":")) {
				pos = pos + len(blockLength) + 1
				continue
			}

			obj := &unstructured.Unstructured{}
			err = yaml.Unmarshal([]byte(resource), obj)
			if err != nil {
				plugin.Logger(ctx).Error("parsedManifestFileContentUncached", "failed to unmarshal the content", err, "path", path)
				return nil, err
			}

			obj.SetAPIVersion(obj.GetAPIVersion())
			obj.SetKind(obj.GetKind())
			gvk := obj.GetObjectKind().GroupVersionKind()
			obj.SetGroupVersionKind(schema.GroupVersionKind{
				Group:   gvk.Group,
				Version: gvk.Version,
				Kind:    gvk.Kind,
			})

			// Convert the content to concrete type based on the resource kind
			targetObj, err := convertUnstructuredDataToType(obj)
			if err != nil {
				plugin.Logger(ctx).Error("parsedManifestFileContentUncached", "failed to convert content into a concrete type", err, "path", path)
				return nil, err
			}

			parsedContents = append(parsedContents, parsedContent{
				ParsedData: targetObj,
				Kind:       obj.GetKind(),
				Path:       path,
				SourceType: "manifest",
				StartLine:  pos + 1, // Since starts from 0
				EndLine:    pos + len(blockLength),
			})

			// Increment the position by the length of the block
			// the value is added with 1 to include the separator
			pos = pos + len(blockLength) + 1
		}
	}

	return parsedContents, nil
}

// Returns the list of file paths/glob patterns after resolving all the given manifest file paths.
func resolveManifestFilePaths(ctx context.Context, d *plugin.QueryData) ([]string, error) {
	// Read the config
	kubernetesConfig := GetConfig(d.Connection)

	// Check for the sourceTypes argument in the config. Valid values are: "deployed", "manifest" and "helm".
	// Default set to include values.
	var sources = All.ToSourceTypes()
	if kubernetesConfig.SourceTypes != nil {
		sources = kubernetesConfig.SourceTypes
	}
	// TODO: Remove once `SourceType` is obsolete
	if kubernetesConfig.SourceTypes == nil && kubernetesConfig.SourceType != nil {
		if *kubernetesConfig.SourceType != "all" { // if is all, sources is already set by default
			sources = []string{*kubernetesConfig.SourceType}
		}
	}

	// Return no files if manifest not set in source_types or if we omit setting any file paths
	if !slices.Contains(sources, "manifest") {
		return nil, nil
	}
	if kubernetesConfig.ManifestFilePaths == nil {
		return nil, nil
	}

	// Gather file path matches for the glob
	var matches, resolvedPaths []string
	paths := kubernetesConfig.ManifestFilePaths
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

	return resolvedPaths, nil
}

//// HELM values

type Rows []Row
type Row struct {
	Path        string
	Key         []string
	Value       interface{}
	Tag         *string
	PreComments []string
	HeadComment string
	LineComment string
	FootComment string
	StartLine   int
	StartColumn int
}

func treeToList(tree *goYaml.Node, prefix []string, rows *Rows, preComments []string, headComments []string, footComments []string) {
	switch tree.Kind {
	case goYaml.DocumentNode:
		for i, v := range tree.Content {
			localComments := []string{}
			headComments = []string{}
			footComments = []string{}
			if i == 0 {
				localComments = append(localComments, preComments...)
				if tree.HeadComment != "" {
					localComments = append(localComments, tree.HeadComment)
					headComments = append(headComments, tree.HeadComment)
				}
				if tree.FootComment != "" {
					footComments = append(footComments, tree.FootComment)
				}
				if tree.LineComment != "" {
					localComments = append(localComments, tree.LineComment)
				}
			}
			treeToList(v, prefix, rows, localComments, headComments, footComments)
		}
	case goYaml.SequenceNode:
		if len(tree.Content) == 0 {
			row := Row{
				Key:         prefix,
				Value:       []string{},
				Tag:         &tree.Tag,
				StartLine:   tree.Line,
				StartColumn: tree.Column,
				PreComments: preComments,
				HeadComment: strings.Join(headComments, ","),
				LineComment: tree.LineComment,
				FootComment: strings.Join(footComments, ","),
			}
			*rows = append(*rows, row)
		}

		for i, v := range tree.Content {
			localComments := []string{}
			headComments = []string{}
			footComments = []string{}
			if i == 0 {
				localComments = append(localComments, preComments...)
				if tree.HeadComment != "" {
					localComments = append(localComments, tree.HeadComment)
					headComments = append(headComments, tree.HeadComment)
				}
				if tree.LineComment != "" {
					localComments = append(localComments, tree.LineComment)
				}
			}
			newKey := make([]string, len(prefix))
			copy(newKey, prefix)
			newKey = append(newKey, strconv.Itoa(i))
			treeToList(v, newKey, rows, localComments, headComments, footComments)
		}
	case goYaml.MappingNode:
		localComments := []string{}
		headComments = []string{}
		footComments = []string{}
		localComments = append(localComments, preComments...)
		if tree.HeadComment != "" {
			localComments = append(localComments, tree.HeadComment)
			headComments = append(headComments, tree.HeadComment)
		}
		if tree.FootComment != "" {
			footComments = append(footComments, tree.FootComment)
		}
		if tree.LineComment != "" {
			localComments = append(localComments, tree.LineComment)
		}
		if len(tree.Content) == 0 {
			row := Row{
				Key:         prefix,
				Value:       map[string]interface{}{},
				Tag:         &tree.Tag,
				StartLine:   tree.Line,
				StartColumn: tree.Column,
				PreComments: preComments,
				HeadComment: strings.Join(headComments, ","),
				LineComment: tree.LineComment,
				FootComment: strings.Join(footComments, ","),
			}
			*rows = append(*rows, row)
		}
		i := 0
		for i < len(tree.Content)-1 {
			key := tree.Content[i]
			val := tree.Content[i+1]
			i = i + 2
			if key.HeadComment != "" {
				localComments = append(localComments, key.HeadComment)
				headComments = append(headComments, key.HeadComment)
			}
			if key.FootComment != "" {
				footComments = append(footComments, key.FootComment)
			}
			if key.LineComment != "" {
				localComments = append(localComments, key.LineComment)
			}
			newKey := make([]string, len(prefix))
			copy(newKey, prefix)
			newKey = append(newKey, key.Value)
			treeToList(val, newKey, rows, localComments, headComments, footComments)
			localComments = make([]string, 0)
			headComments = make([]string, 0)
			footComments = make([]string, 0)
		}
	case goYaml.ScalarNode:
		row := Row{
			Key:         prefix,
			Value:       tree.Value,
			Tag:         &tree.Tag,
			StartLine:   tree.Line,
			StartColumn: tree.Column,
			PreComments: preComments,
			HeadComment: strings.Join(headComments, ","),
			LineComment: tree.LineComment,
			FootComment: strings.Join(footComments, ","),
		}
		if tree.Tag == "!!null" {
			row.Value = nil
		}
		*rows = append(*rows, row)
	}
}

func keysToSnakeCase(_ context.Context, d *transform.TransformData) (interface{}, error) {
	keys := d.Value.([]string)
	snakes := []string{}
	re := regexp.MustCompile(`[^A-Za-z0-9_]`)
	for _, k := range keys {
		snakes = append(snakes, re.ReplaceAllString(k, "_"))
	}
	return strings.Join(snakes, "."), nil
}

// normalizeCPUToMilliCores converts CPU quantities to millicores (m), rounding up if necessary.
func normalizeCPUToMilliCores(cpu string) (int64, error) {
	if strings.HasSuffix(cpu, "m") {
		// Already in millicores
		value, err := strconv.ParseFloat(strings.TrimSuffix(cpu, "m"), 64)
		if err != nil {
			return 0, fmt.Errorf("invalid CPU value: %s", cpu)
		}
		return int64(math.Ceil(value)), nil
	}

	// Convert cores to millicores (handles scientific notation like "500e-3")
	value, err := strconv.ParseFloat(cpu, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid CPU value: %s", cpu)
	}

	milliCores := value * 1000
	return int64(math.Ceil(milliCores)), nil
}

// normalizeMemoryToBytes converts memory quantities to bytes, rounding up if necessary.
func normalizeMemoryToBytes(memory string) (int64, error) {
	// Handle scientific notation by trying to parse the full string as a float first
	if value, err := strconv.ParseFloat(memory, 64); err == nil {
		// If it parses successfully as a pure number (including scientific notation), treat as bytes
		return int64(math.Ceil(value)), nil
	}

	// Find the boundary between number and unit more carefully
	// Look for the start of alphabetic characters that don't include 'e' or 'E' (for scientific notation)
	valuePart := memory
	unitPart := "B"

	for i, r := range memory {
		// If we find an alphabetic character that's not 'e' or 'E', or if we find 'e'/'E'
		// followed by an alphabetic character (not a digit or +/-), then we've found the unit
		if (r >= 'A' && r <= 'Z' && r != 'E') || (r >= 'a' && r <= 'z' && r != 'e') {
			valuePart = memory[:i]
			unitPart = memory[i:]
			break
		}
		// Handle the case where 'e' or 'E' is followed by a letter (not +/- or digit)
		if (r == 'e' || r == 'E') && i+1 < len(memory) {
			next := memory[i+1]
			if (next >= 'A' && next <= 'Z') || (next >= 'a' && next <= 'z') {
				valuePart = memory[:i]
				unitPart = memory[i:]
				break
			}
		}
	}

	if valuePart == "" {
		return 0, fmt.Errorf("invalid memory value: %s", memory)
	}

	value, err := strconv.ParseFloat(valuePart, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid memory value: %s", memory)
	}

	unitPart = strings.TrimSpace(unitPart)
	multiplier, exists := memoryUnits[unitPart]
	if !exists {
		return 0, fmt.Errorf("unknown unit: %s", unitPart)
	}

	bytes := value * multiplier
	return int64(math.Ceil(bytes)), nil
}

// Unit multipliers for memory
var memoryUnits = map[string]float64{
	"m":  1e-3,
	"B":  1,
	"Ki": math.Pow(2, 10),
	"Mi": math.Pow(2, 20),
	"Gi": math.Pow(2, 30),
	"Ti": math.Pow(2, 40),
	"Pi": math.Pow(2, 50),
	"Ei": math.Pow(2, 60),
	"k":  1e3,
	"M":  1e6,
	"G":  1e9,
	"T":  1e12,
	"P":  1e15,
	"E":  1e18,
}
