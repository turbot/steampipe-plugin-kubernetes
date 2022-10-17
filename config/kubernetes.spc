connection "kubernetes" {
  plugin = "kubernetes"

  # By default Plugin will pick the current context, present inside ~/.kube/config.

  # If ~/.kube/config file is unavailable, the plugin will check for InClusterConfig configuration. No extra configuration is required.

  # If you want to choose a specific context, you can set the name of context with the `config_context` argument.
  # config_context = "minikube"

  # If the config file path is located in other location, you can specify the path of kube config file with `config_path` argument.
  # config_path    = "~/.kube/config"

  # If you have a kube config setup using the kubectl CLI(https://kubernetes.io/docs/reference/kubectl/), the plugin just works with that connection.

  # This plugin also supports OpenID Connect (OIDC) authentication. No extra configuration is required.
}



