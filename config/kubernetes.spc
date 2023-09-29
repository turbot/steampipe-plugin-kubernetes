connection "kubernetes" {
  plugin = "kubernetes"

  # By default, the plugin will use credentials in "~/.kube/config" with the current context.
  # OpenID Connect (OIDC) authentication is supported without any extra configuration.
  # The kubeconfig path and context can also be specified with the following config arguments:

  # Specify the file path to the kubeconfig.
  # Can also be set with the "KUBECONFIG" or "KUBE_CONFIG_PATH" environment variables. Plugin will prioritize KUBECONFIG if both are available.
  # config_path = "~/.kube/config"

  # Specify a context other than the current one.
  # config_context = "minikube"

  # List of custom resources that will be created as dynamic tables.
  # No dynamic tables will be created if this arg is empty or not set.
  # Wildcard based searches are supported.

  # For example:
  #  - "*" matches all custom resources available
  #  - "*.storage.k8s.io" matches all custom resources in the storage.k8s.io group
  #  - "certificates.cert-manager.io" matches a specific custom resource "certificates.cert-manager.io"
  #  - "backendconfig" matches the singular name "backendconfig" in any group

  # Defaults to all custom resources
  custom_resource_tables = ["*"]

  # If no kubeconfig file can be found, the plugin will attempt to use the service account Kubernetes gives to pods.
  # This authentication method is intended for clients that expect to be running inside a pod running on Kubernetes.

  # Specify the source(s) of the resource(s). Possible values: `deployed`, `helm` and `manifest`.
  # Defaults to all possible values. Set the argument to override the default value.
  # If `deployed` is contained in the value, tables will show all the deployed resources.
  # If `helm` is contained in the value, tables will show resources from the configured helm charts.
  # If `manifest` is contained in the value, tables will show all the resources from the kubernetes manifest. Make sure that the `manifest_file_paths` arg is set.
  # source_types = ["deployed", "helm", "manifest"]

  # Manifest File Configuration

  # Manifest file paths is a list of locations to search for Kubernetes manifest files
  # Manifest file paths can be configured with a local directory, a remote Git repository URL, or an S3 bucket URL
  # Refer https://hub.steampipe.io/plugins/turbot/kubernetes#supported-path-formats for more information
  # Wildcard based searches are supported, including recursive searches
  # Local paths are resolved relative to the current working directory (CWD)

  # For example:
  #  - "*.yml" or "*.yaml" or "*.json" matches all Kubernetes manifest files in the CWD
  #  - "**/*.yml" or "**/*.yaml" or "**/*.json" matches all Kubernetes manifest files in the CWD and all sub-directories
  #  - "../*.yml" or "../*.yaml" or "../*.json" matches all Kubernetes manifest files in the CWD's parent directory
  #  - "steampipe*.yml" or "steampipe*.yaml" or "steampipe*.json" matches all Kubernetes manifest files starting with "steampipe" in the CWD
  #  - "/path/to/dir/*.yml" or "/path/to/dir/*.yaml" or "/path/to/dir/*.json" matches all Kubernetes manifest files in a specific directory
  #  - "/path/to/dir/main.yml" or "/path/to/dir/main.yaml" or "/path/to/dir/main.json" matches a specific file

  # If the given paths includes "*", all files (including non-kubernetes manifest files) in
  # the CWD will be matched, which may cause errors if incompatible file types exist

  # Defaults to CWD
  # manifest_file_paths = [ "*.yml", "*.yaml", "*.json" ]

  # Helm configuration

  # A map for Helm charts along with the path to the chart directory and the paths of the value override files (if any).
  # Every map should have chart_path defined, and the values_file_paths is optional.
  # You can define multiple charts in the config.
  # helm_rendered_charts = {
  #   "chart_name" = {
  #     chart_path        = "/path/to/chart/dir"
  #     values_file_paths = ["/path/to/value/override/files.yaml"]
  #   }
  # }
}
