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

  # Specify the custom resource definitions for which the dynamic tables will be created.
  # The custom_resource_tables list may include wildcards (e.g. *, ip*, storagestates.migration.???.io), singular names or the full name.
  # By default plugin will load the dynamic tables for all the available custom resource definitions.
  # The plugin will not create dynamic tables if custom_resource_tables is empty or not set.
  # custom_resource_tables = ["certificate","ip*","storagestates.migration.k8s.io"]
  custom_resource_tables = ["*"]

  # If no kubeconfig file can be found, the plugin will attempt to use the service account Kubernetes gives to pods.
  # This authentication method is intended for clients that expect to be running inside a pod running on Kubernetes.
}
