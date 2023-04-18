connection "kubernetes" {
  plugin = "kubernetes"

  # By default, the plugin will use credentials in "~/.kube/config" with the current context.
  # OpenID Connect (OIDC) authentication is supported without any extra configuration.
  # The kubeconfig path and context can also be specified with the following config arguments:

  # Specify the file path to the kubeconfig.
  # Can also be set with the "KUBE_CONFIG_PATHS" or "KUBERNETES_MASTER" environment variables.
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

  # Paths is a list of locations to search for Kubernetes manifest files
  # Paths can be configured with a local directory, a remote Git repository URL, or an S3 bucket URL
  # Refer https://hub.steampipe.io/plugins/turbot/kubernetes#supported-path-formats for more information
  # Wildcard based searches are supported, including recursive searches
  # Local paths are resolved relative to the current working directory (CWD)

  # For example:
  #  - "*.yml" matches all Kubernetes manifest files in the CWD
  #  - "**/*.yml" matches all Kubernetes manifest files in the CWD and all sub-directories
  #  - "../*.yml" matches all Kubernetes manifest files in the CWD's parent directory
  #  - "steampipe*.yml" matches all Kubernetes manifest files starting with "steampipe" in the CWD
  #  - "/path/to/dir/*.yml" matches all Kubernetes manifest files in a specific directory
  #  - "/path/to/dir/main.yml" matches a specific file

  # If paths includes "*", all files (including non-kubernetes manifest files) in
  # the CWD will be matched, which may cause errors if incompatible file types exist

  # Defaults to CWD
  # manifest_file_paths = [ "*.yml", "*.yaml" ]

  # Specify the source of the resource. Possible values: `deployed`, `manifest`, and `all`.
  # Default set to `all`. Set the argument to override the default value.
  # If the value is set to `deployed`, tables will show all the deployed resources.
  # If set to `manifest`, tables will show all the resources from the kubernetes manifest. Make sure that the `manifest_file_paths` arg is set.
  # If `all`, tables will show all the deployed and manifest resources.
  # source_type = "all"
}
