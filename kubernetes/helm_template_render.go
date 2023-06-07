package kubernetes

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

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
	charts, err := getParsedHelmChart(ctx, d)
	if err != nil {
		return nil, err
	}
	helmConfig := GetConfig(d.Connection)

	var renderedTemplates []HelmRenderedTemplate
	for _, chart := range charts {

		// Return nil, if the config doesn't have any chart path configured
		if chart == nil {
			plugin.Logger(ctx).Debug("getHelmRenderedTemplatesUncached", "no chart configuration found", "connection", d.Connection.Name)
			return nil, nil
		}

		for name, c := range helmConfig.HelmRenderedCharts {
			if c.ChartPath == chart.Path {

				client := newClient()
				client.ReleaseName = name
				client.Namespace = "default" // TODO: Update this to use namespace defined in the current context

				vals := &values.Options{
					ValueFiles: c.ValuesFilePaths,
				}

				manifest, _, err := runInstall([]string{c.ChartPath}, client, vals)
				if err != nil {
					plugin.Logger(ctx).Debug("getHelmRenderedTemplatesUncached", "run_install_error", err, "connection", d.Connection.Name)
					return nil, err
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
	// Read the config
	renderedTemplates, err := getHelmRenderedTemplates(ctx, d, nil)
	if err != nil {
		return nil, err
	}

	var parsedContents []parsedContent
	for _, t := range renderedTemplates {
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
				Data:       targetObj,
				Kind:       obj.GetKind(),
				Path:       t.Path,
				SourceType: fmt.Sprintf("helm_rendered:%s", t.ConfigKey),
			})
		}
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
