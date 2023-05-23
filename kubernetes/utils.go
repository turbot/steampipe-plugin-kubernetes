package kubernetes

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	goYaml "gopkg.in/yaml.v3"

	filehelpers "github.com/turbot/go-kit/files"
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

	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/engine"

	helmClient "github.com/mittwald/go-helm-client"

	"github.com/mitchellh/go-homedir"
	"github.com/turbot/steampipe-plugin-sdk/v5/connection"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

type SourceType string

const (
	Deployed SourceType = "deployed"
	Manifest SourceType = "manifest"
	All      SourceType = "all"
)

// Validate the source type.
func (sourceType SourceType) IsValid() error {
	switch sourceType {
	case Deployed, Manifest, All:
		return nil
	}
	return fmt.Errorf("invalid source type: %s", sourceType)
}

// Convert the source type to its string equivalent
func (sourceType SourceType) String() string {
	return string(sourceType)
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

	// Return nil, if the config is set only to list the manifest resources.
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

	// Return nil, if the config is set to only list the manifest resources.
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

	// Return nil, if the config is set to only list the manifest resources.
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

	// Check for the sourceType argument in the config. Valid values are: "deployed", "manifest" and "all".
	// Default set to "all".
	var source SourceType = "all"
	if kubernetesConfig.SourceType != nil {
		source = SourceType(*kubernetesConfig.SourceType)
		if err := source.IsValid(); err != nil {
			plugin.Logger(ctx).Debug("getK8Config", "invalid_source_type_error", "connection", d.Connection.Name, "error", err)
			return nil, err
		}
	}

	// By default source type is set to "all", which indicates querying the table will return both deployed and manifest resources.
	// If the source type is explicitly set to "manifest", the table will only return the manifest resources.
	// Similarly, setting the value as "deployed" will return all the deployed resources.
	if source.String() == "manifest" {
		plugin.Logger(ctx).Debug("getK8Config", "The source_type set to 'manifest'. Returning nil for API server client.", "connection", d.Connection.Name)
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

	// Check for the sourceType argument in the config. Valid values are: "deployed", "manifest" and "all".
	// Default set to "all".
	var source SourceType = "all"
	if kubernetesConfig.SourceType != nil {
		source = SourceType(*kubernetesConfig.SourceType)
		if err := source.IsValid(); err != nil {
			plugin.Logger(ctx).Debug("getK8ConfigRaw", "invalid_source_type_error", "connection", c.Name, "error", err)
			return nil, err
		}
	}

	// By default source type is set to "all", which indicates querying the table will return both deployed and manifest resources.
	// If the source type is explicitly set to "manifest", the table will only return the manifest resources.
	// Similarly, setting the value as "deployed" will return all the deployed resources.
	if source.String() == "manifest" {
		plugin.Logger(ctx).Debug("getK8ConfigRaw", "The source_type set to 'manifest'. Returning nil for API server client.", "connection", c.Name)
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

//// Utility functions for manifest files

type parsedContent struct {
	Data      any
	Kind      string
	Path      string
	StartLine int
	EndLine   int
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
	plugin.Logger(ctx).Debug("parsedManifestFileContentUncached", "Parsing file content...", "connection", d.Connection.Name)

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
			if !strings.Contains(resource, "kind:") {
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
				Data:      targetObj,
				Kind:      obj.GetKind(),
				Path:      path,
				StartLine: pos + 1, // Since starts from 0
				EndLine:   pos + len(blockLength),
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
	k8sConfig := GetConfig(d.Connection)

	// Return nil, if the source_type is set to "deployed"
	if k8sConfig.SourceType != nil &&
		*k8sConfig.SourceType == "deployed" {
		return nil, nil
	}

	// Return error if source_tpe arg is explicitly set to "manifest" in the config, but
	// manifest_file_paths arg is not set.
	if k8sConfig.SourceType != nil &&
		*k8sConfig.SourceType == "manifest" &&
		k8sConfig.ManifestFilePaths == nil {
		return nil, errors.New("manifest_file_paths must be set in the config while the source_type is 'manifest'")
	}

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

	return resolvedPaths, nil
}

type parsedHelmChart struct {
	Chart *chart.Chart
	Path  string
}

// Get the parsed contents of the given Helm chart.
func getParsedHelmChart(ctx context.Context, d *plugin.QueryData) (*parsedHelmChart, error) {
	conn, err := parsedHelmChartCached(ctx, d, nil)
	if err != nil {
		return nil, err
	}
	return conn.(*parsedHelmChart), nil
}

// Cached form of the parsed Helm chart.
var parsedHelmChartCached = plugin.HydrateFunc(parsedHelmChartUncached).Memoize()

// parsedHelmChartUncached is the actual implementation of getParsedHelmChart, which should
// be run only once per connection. Do not call this directly, use
// getParsedHelmChart instead.
func parsedHelmChartUncached(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	// Read the config
	helmConfig := GetConfig(d.Connection)

	chartDir := helmConfig.HelmChartDir

	// Return empty parsedHelmChart object if no Helm chart directory path provided in the config
	if chartDir == nil {
		plugin.Logger(ctx).Debug("parsedHelmChartUncached", "helm_chart_dir not configured in the config", "connection", d.Connection.Name)
		return &parsedHelmChart{}, nil
	}
	plugin.Logger(ctx).Debug("parsedHelmChartUncached", "Parsing Helm chart", chartDir, "connection", d.Connection.Name)

	// Load the given chart directory
	chart, err := loader.Load(*chartDir)
	if err != nil {
		plugin.Logger(ctx).Error("parsedHelmChartUncached", "load_chart_error", err)
		return nil, err
	}

	return &parsedHelmChart{
		Chart: chart,
		Path:  *chartDir,
	}, nil
}

// Get the rendered template contents.
func getHelmRenderedTemplates(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (map[string]string, error) {
	helmRenderedTemplates, err := parsedHelmRenderedTemplatesCached(ctx, d, nil)
	if err != nil {
		return nil, err
	}

	if helmRenderedTemplates != nil {
		return helmRenderedTemplates.(map[string]string), nil
	}
	return nil, nil
}

// Cached form of the rendered template content.
var parsedHelmRenderedTemplatesCached = plugin.HydrateFunc(parsedHelmRenderedTemplatesUncached).Memoize()

// parsedHelmRenderedTemplatesUncached is the actual implementation of getHelmRenderedTemplates, which should
// be run only once per connection. Do not call this directly, use
// getHelmRenderedTemplates instead.
func parsedHelmRenderedTemplatesUncached(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	chart, err := getParsedHelmChart(ctx, d)
	if err != nil {
		return nil, err
	}

	// Return nil, if the config doesn't have any chart path configured
	if chart.Chart == nil {
		plugin.Logger(ctx).Debug("parsedHelmRenderedTemplatesUncached", "no chart configuration found", "connection", d.Connection.Name)
		return nil, nil
	}

	// Get the values required to render the templates
	values, err := getHelmChartOverrideValues(ctx, d, chart)
	if err != nil {
		return nil, err
	}

	values = map[string]interface{}{
		"Values": values,
		"Release": map[string]interface{}{
			"Service": "Helm",
			"Name":    chart.Chart.Metadata.Name, // Keeping it as same as the chart name for now. In CLI, either the value can be passed in the arg, or can be auto-generated.
		},
		"Chart":        chart.Chart.Metadata,
		"Capabilities": chartutil.Capabilities{},
		"Template": map[string]interface{}{
			"BasePath": "/path/to/base",
		},
	}

	renderedChart, err := engine.Render(chart.Chart, values)
	if err != nil {
		plugin.Logger(ctx).Error("parsedHelmRenderedTemplatesUncached", "connection", d.Connection.Name, "template_render_error", err)
		return nil, err
	}

	return renderedChart, nil
}

// Return the values required to render the templates
func getHelmChartOverrideValues(ctx context.Context, d *plugin.QueryData, chart *parsedHelmChart) (map[string]interface{}, error) {
	helmConfig := GetConfig(d.Connection)

	// Get the default values defined in the values.yaml file
	values := chart.Chart.Values

	// Check for value override files configured in the connection config
	var matches, valueFiles []string
	if helmConfig.HelmValueOverride != nil {
		for _, i := range helmConfig.HelmValueOverride {
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
			valueFiles = append(valueFiles, i)
		}
	}

	// If any value override files provided in the config, use those value to render the templates
	// The priority will be given to the last file specified.
	// For example, if both values.yaml and override.yaml
	// contained a key called 'foo', the value set in override.yaml would take precedence.
	for _, f := range valueFiles {
		var override map[string]interface{}
		// Read files
		bs, err := os.ReadFile(f)
		if err != nil {
			plugin.Logger(ctx).Error("parsedHelmRenderedTemplatesUncached", "read_file_error", "connection_name", d.Connection.Name, "failed to read file %s: %v", f, err)
			return nil, err
		}
		if err := goYaml.Unmarshal(bs, &override); err != nil {
			plugin.Logger(ctx).Debug("getHelmChartOverrideValues", "unmarshal_error", "failed to unmarshal value override file", f, "error", err)
			return nil, err
		}
		values = mergeMaps(values, override)
	}

	return values, nil
}

// Get the parsed contents of the given files.
func getRenderedHelmTemplateContent(ctx context.Context, d *plugin.QueryData) ([]parsedContent, error) {
	conn, err := RenderedHelmTemplateContentCached(ctx, d, nil)
	if err != nil {
		return nil, err
	}
	return conn.([]parsedContent), nil
}

// Cached form of the parsed file content.
var RenderedHelmTemplateContentCached = plugin.HydrateFunc(RenderedHelmTemplateContentUncached).Memoize()

// RenderedHelmTemplateContentUncached is the actual implementation of getRenderedHelmTemplateContent, which should
// be run only once per connection. Do not call this directly, use
// getRenderedHelmTemplateContent instead.
func RenderedHelmTemplateContentUncached(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Debug("RenderedHelmTemplateContentUncached", "Parsing file content...", "connection", d.Connection.Name)

	// Read the config
	renderedTemplates, err := getHelmRenderedTemplates(ctx, d, nil)
	if err != nil {
		return nil, err
	}

	var parsedContents []parsedContent
	for path, template := range renderedTemplates {

		obj := &unstructured.Unstructured{}
		err = yaml.Unmarshal([]byte(template), obj)
		if err != nil {
			plugin.Logger(ctx).Error("RenderedHelmTemplateContentUncached", "unmarshal_error", err)
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
			plugin.Logger(ctx).Error("RenderedHelmTemplateContentUncached", "failed to convert content into a concrete type", err, "path", path)
			return nil, err
		}

		parsedContents = append(parsedContents, parsedContent{
			Data: targetObj,
			Kind: obj.GetKind(),
			Path: path,
		})
	}

	// // Check for the start of the document
	// pos := 0
	// for _, resource := range strings.Split(string(content), "---") {
	// 	// Skip empty documents, `Decode` will fail on them
	// 	// Also, increment the pos to include the separator position (e.g. ---)
	// 	if len(resource) == 0 {
	// 		pos++
	// 		continue
	// 	}

	// 	// Calculate the length of the YAML resource block
	// 	blockLength := strings.Split(strings.ReplaceAll(resource, " ", ""), "\n")

	// 	// Remove the extra lines added during the split operation based on the separator
	// 	blockLength = blockLength[:len(blockLength)-1]
	// 	if blockLength[0] == "" {
	// 		blockLength = blockLength[1:]
	// 	}

	// 	// skip if no kind defined
	// 	if !strings.Contains(resource, "kind:") {
	// 		pos = pos + len(blockLength) + 1
	// 		continue
	// 	}

	// 	obj := &unstructured.Unstructured{}
	// 	err = yaml.Unmarshal([]byte(resource), obj)
	// 	if err != nil {
	// 		plugin.Logger(ctx).Error("parsedHelmChartContentUncached", "failed to unmarshal the content", err, "path", path)
	// 		return nil, err
	// 	}

	// 	obj.SetAPIVersion(obj.GetAPIVersion())
	// 	obj.SetKind(obj.GetKind())
	// 	gvk := obj.GetObjectKind().GroupVersionKind()
	// 	obj.SetGroupVersionKind(schema.GroupVersionKind{
	// 		Group:   gvk.Group,
	// 		Version: gvk.Version,
	// 		Kind:    gvk.Kind,
	// 	})

	// 	// Convert the content to concrete type based on the resource kind
	// 	targetObj, err := convertUnstructuredDataToType(obj)
	// 	if err != nil {
	// 		plugin.Logger(ctx).Error("parsedHelmChartContentUncached", "failed to convert content into a concrete type", err, "path", path)
	// 		return nil, err
	// 	}

	// 	parsedContents = append(parsedContents, parsedContent{
	// 		Data:      targetObj,
	// 		Kind:      obj.GetKind(),
	// 		Path:      path,
	// 		StartLine: pos + 1, // Since starts from 0
	// 		EndLine:   pos + len(blockLength),
	// 	})

	// 	// Increment the position by the length of the block
	// 	// the value is added with 1 to include the separator
	// 	pos = pos + len(blockLength) + 1
	// }

	return parsedContents, nil
}

// getHelmClient creates  the client for Helm
func getHelmClient(ctx context.Context, namespace string) (helmClient.Client, error) {
	// Set the namespace if specified.
	// By default current namespace context is used.
	options := &helmClient.Options{}
	if namespace != "" {
		options.Namespace = namespace
	}

	// Create client
	client, err := helmClient.New(options)
	if err != nil {
		plugin.Logger(ctx).Error("getHelmClient", "client_error", err)
		return nil, err
	}

	return client, nil
}

// Returns combined values of two files.
// The objects got merged, but same attributes gets replaced.
func mergeMaps(a, b map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(a))
	for k, v := range a {
		out[k] = v
	}
	for k, v := range b {
		if v, ok := v.(map[string]interface{}); ok {
			if bv, ok := out[k]; ok {
				if bv, ok := bv.(map[string]interface{}); ok {
					out[k] = mergeMaps(bv, v)
					continue
				}
			}
		}
		out[k] = v
	}
	return out
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
