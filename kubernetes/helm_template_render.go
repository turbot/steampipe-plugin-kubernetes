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
)

type HelmTemplates struct {
	Name string
	Data string
}

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

type HelmRenderedTemplate struct {
	Data      string
	Chart     *chart.Chart
	Path      string
	ConfigKey string
}

// Get the rendered templates.
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

func runInstall(args []string, client *action.Install,
	valueOpts *values.Options) (*release.Release, []string, error) {
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

func extractTemplatePathFromContent(content string) string {
	splitContent := strings.Split(content, "\n")
	sourceInfoFromManifest := splitContent[1]

	source := strings.Split(sourceInfoFromManifest, "/")

	if len(source) > 1 {
		return strings.Join(source[1:], "/")
	}
	return ""
}
