package kubernetes

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/release"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/yaml"
)

var (
	settings = cli.New()
)

// newClient will create a new instance on helm client used to render the chart
func newClient() *action.Install {
	cfg := new(action.Configuration)
	client := action.NewInstall(cfg)
	client.DryRun = true
	client.Replace = true // Skip the name check
	client.ClientOnly = true
	client.APIVersions = chartutil.VersionSet([]string{})
	client.IncludeCRDs = false
	client.Namespace = "default"
	return client
}

// Utils functions from template render

type HelmRenderedTemplate struct {
	Data      string
	Chart     *chart.Chart
	Path      string
	ConfigKey string
}

// getHelmRenderedTemplates returns the resulting manifest after rendering all the templates defined in the configured charts
func getHelmRenderedTemplates(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) ([]HelmRenderedTemplate, error) {
	helmRenderedTemplates, err := getHelmRenderedTemplatesCached(ctx, d, nil)
	if err != nil {
		plugin.Logger(ctx).Error("getHelmRenderedTemplates", "template_render_error", err)
		return nil, err
	}

	if helmRenderedTemplates != nil {
		return helmRenderedTemplates.([]HelmRenderedTemplate), nil
	}

	return nil, nil
}

// Cached form of the rendered templates.
var getHelmRenderedTemplatesCached = plugin.HydrateFunc(getHelmRenderedTemplatesUncached).Memoize()

// getHelmRenderedTemplatesUncached is the actual implementation of getHelmRenderedTemplates, which should
// be run only once per connection. Do not call this directly, use
// getHelmRenderedTemplates instead.
func getHelmRenderedTemplatesUncached(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	charts, err := getUniqueHelmCharts(ctx, d)
	if err != nil {
		return nil, err
	}
	kubernetesConfig := GetConfig(d.Connection)

	var renderedTemplates []HelmRenderedTemplate
	for _, chart := range charts {

		// Return nil, if the config doesn't have any chart path configured
		if chart == nil {
			plugin.Logger(ctx).Debug("getHelmRenderedTemplatesUncached", "no chart configuration found", "connection", d.Connection.Name)
			return nil, nil
		}

		var processedHelmConfigs []string
		for name, c := range kubernetesConfig.HelmRenderedCharts {
			if c.ChartPath == chart.Path && !helpers.StringSliceContains(processedHelmConfigs, name) {

				// Add the processed Helm render configs into processedHelmConfigs to avoid duplicate entries
				processedHelmConfigs = append(processedHelmConfigs, name)

				client := newClient()
				client.ReleaseName = name
				client.Namespace = "default" // TODO: Update this to use namespace defined in the current context

				vals := &values.Options{}
				if len(c.ValuesFilePaths) > 0 {
					vals.ValueFiles = c.ValuesFilePaths
				}

				manifest, _, err := runInstall([]string{c.ChartPath}, client, vals)
				if err != nil {
					plugin.Logger(ctx).Debug("getHelmRenderedTemplatesUncached", "run_install_error", err, "connection", d.Connection.Name)
					// return nil, err
					continue
				}

				splitManifest := strings.Split(manifest.Manifest, "---")
				for _, content := range splitManifest {
					if len(content) == 0 {
						continue
					}

					renderedTemplates = append(renderedTemplates, HelmRenderedTemplate{
						Data:      content,
						Chart:     chart.Chart,
						Path:      path.Join(c.ChartPath, extractTemplatePathFromContent(content)),
						ConfigKey: name,
					})
				}
			}
		}
	}
	return renderedTemplates, nil
}

// Utils functions from formatting the rendered template contents

// Get the parsed contents of the given files.
func getRenderedHelmTemplateContent(ctx context.Context, d *plugin.QueryData) ([]parsedContent, error) {
	conn, err := renderedHelmTemplateContentCached(ctx, d, nil)
	if err != nil {
		return nil, err
	}
	return conn.([]parsedContent), nil
}

// Cached form of the parsed file content.
var renderedHelmTemplateContentCached = plugin.HydrateFunc(renderedHelmTemplateContentUncached).Memoize()

// renderedHelmTemplateContentUncached is the actual implementation of getRenderedHelmTemplateContent, which should
// be run only once per connection. Do not call this directly, use
// getRenderedHelmTemplateContent instead.
func renderedHelmTemplateContentUncached(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	// List the fully rendered templates
	renderedTemplates, err := getHelmRenderedTemplates(ctx, d, nil)
	if err != nil {
		return nil, err
	}

	// Get the start and end line information for the templates
	templateWithLineInfo, err := getRawTemplateLineInfo(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("renderedHelmTemplateContentUncached", "failed to get line information from raw template", err)
		return nil, err
	}
	temp := map[string][]LineInfo{}

	var processedConfigs []string
	var parsedContents []parsedContent
	for _, t := range renderedTemplates {
		// Get the line numbers of each configuration block for the current template
		if len(templateWithLineInfo[t.Path]) == 0 {
			continue
		}

		// If the same chart is configured more than once, use the key used to identify the chart config in the config file
		// to avoid the conflicts when calculating the line numbers.
		// If the current config is not yet processed
		test := fmt.Sprintf("%s:%s", t.Path, t.ConfigKey)
		if !helpers.StringSliceContains(processedConfigs, test) {
			// Set the config as processed
			processedConfigs = append(processedConfigs, test)

			// Reinitialize the temp with the actual templateWithLineInfo data
			for k, v := range templateWithLineInfo {
				temp[k] = v
			}
		}
		lineInfo := temp[t.Path]

		for _, resource := range strings.Split(t.Data, "---") {
			// Skip empty documents, `Decode` will fail on them
			// Also, increment the pos to include the separator position (e.g. ---)
			if len(resource) == 0 {
				continue
			}

			// Skip if no kind defined
			if !(strings.Contains(resource, "kind:") || strings.Contains(resource, "\"kind\":")) {
				continue
			}

			obj := &unstructured.Unstructured{}
			err = yaml.Unmarshal([]byte(resource), obj)
			if err != nil {
				plugin.Logger(ctx).Error("renderedHelmTemplateContentUncached", "unmarshal_error", err)
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
				plugin.Logger(ctx).Error("RenderedHelmTemplateContentUncached", "failed to convert content into a concrete type", err, "path", t.Path)
				return nil, err
			}

			parsedContents = append(parsedContents, parsedContent{
				ParsedData: targetObj,
				Kind:       obj.GetKind(),
				Path:       t.Path,
				SourceType: fmt.Sprintf("helm_rendered:%s", t.ConfigKey),
				StartLine:  lineInfo[0].StartLine,
				EndLine:    lineInfo[0].EndLine,
			})

			// Remove the line information for the processed block
			if len(lineInfo) > 1 {
				lineInfo = lineInfo[1:]
				temp[t.Path] = lineInfo
			}
		}
	}

	return parsedContents, nil
}

// Utils functions

// runInstall renders the templates and returns the resulting manifest after communicating with the k8s cluster without actually creating any resources on the cluster
func runInstall(args []string, client *action.Install, valueOpts *values.Options) (*release.Release, []string, error) {
	defer log.SetOutput(os.Stderr)
	if client.Version == "" && client.Devel {
		client.Version = ">0.0.0-0"
	}

	_, charts, err := client.NameAndChart(args)
	if err != nil {
		return nil, []string{}, err
	}

	cp, err := client.ChartPathOptions.LocateChart(charts, settings)
	if err != nil {
		return nil, []string{}, err
	}

	p := getter.All(settings)
	vals, err := valueOpts.MergeValues(p)
	if err != nil {
		return nil, []string{}, err
	}

	// Check chart dependencies to make sure all are present in /charts
	chartRequested, err := loader.Load(cp)
	if err != nil {
		return nil, []string{}, err
	}

	excluded := getExcluded(chartRequested, cp)

	if instErr := checkIfInstallable(chartRequested); instErr != nil {
		return nil, []string{}, instErr
	}

	helmRelease, err := client.Run(chartRequested, vals)
	if err != nil {
		return nil, []string{}, err
	}
	return helmRelease, excluded, nil
}

// checkIfInstallable validates if a chart can be installed
//
// Application chart type is only installable
func checkIfInstallable(ch *chart.Chart) error {
	switch ch.Metadata.Type {
	case "", "application":
		return nil
	}
	return fmt.Errorf("%s charts are not installable", ch.Metadata.Type)
}

// getExcluded will return all files rendered to be excluded from scan
func getExcluded(charterino *chart.Chart, chartpath string) []string {
	excluded := make([]string, 0)
	for _, file := range charterino.Raw {
		excluded = append(excluded, filepath.Join(chartpath, file.Name))
	}

	return excluded
}

// extractTemplatePathFromContent extracts the path of the template source file from the rendered template content
func extractTemplatePathFromContent(content string) string {
	splitContent := strings.Split(content, "\n")
	sourceInfoFromManifest := splitContent[1]

	source := strings.Split(sourceInfoFromManifest, "/")

	if len(source) > 1 {
		return strings.Join(source[1:], "/")
	}
	return ""
}

type LineInfo struct {
	StartLine int
	EndLine   int
}

// getRawTemplateLineInfo returns a map containing the line numbers for each resource configuration block defined in a template file
func getRawTemplateLineInfo(ctx context.Context, d *plugin.QueryData) (map[string][]LineInfo, error) {

	templateMetadata := map[string][]LineInfo{}
	charts, err := getUniqueHelmCharts(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Debug("getChartTemplatesInfo", "failed to list helm charts", err)
		return nil, err
	}

	// A template file can have more than 1 configuration defined in it separated by `---` separator.
	// This function will read all the raw templates, and will calculate the start and end line of a resource configuration block.
	// If the file has more than 1 configurations defined, the function will return a map containing the template file path, along with an array of map containing the start and end line information. For example:
	// map["/path/to/the/template1":[{StartLine: ..., EndLine: ...}, {...}], "/path/to/the/template2":[{StartLine: ..., EndLine: ...}, {...}], ...]
	for _, chart := range charts {
		templates := chart.Chart.Templates

		for _, template := range templates {
			templateContent := string(template.Data)

			var lineInfo []LineInfo
			startLine := 0
			count := 0

			for _, content := range strings.Split(templateContent, "---") {
				// Skip empty documents, `Decode` will fail on them
				// Also, increment the pos to include the separator position (e.g. ---)
				if len(content) == 0 {
					startLine++
					continue
				}
				count++

				// Calculate the length of the YAML resource block
				blockLength := strings.Split(strings.ReplaceAll(content, " ", ""), "\n")

				// Remove the extra lines added during the split operation based on the separator
				blockLength = blockLength[:len(blockLength)-1]
				if blockLength[0] == "" {
					blockLength = blockLength[1:]
				}

				// Calculate the end line number
				endLine := startLine + len(blockLength)
				if count > 1 {
					endLine++
				}

				lineInfo = append(lineInfo, LineInfo{
					StartLine: startLine + 1, // Since starts from 0
					EndLine:   endLine,
				})

				// Increment the startLine by the length of the block
				// the value is added with 1 to include the separator
				startLine = startLine + len(blockLength) + 1
			}
			templateMetadata[path.Join(chart.Path, template.Name)] = lineInfo
		}
	}

	return templateMetadata, nil
}
