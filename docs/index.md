---
organization: Turbot
category: ["software development"]
icon_url: "/images/plugins/turbot/kubernetes.svg"
brand_color: "#326CE5"
display_name: "Kubernetes"
short_name: "kubernetes"
description: "Steampipe plugin for Kubernetes components."
og_description: "Query Kubernetes with SQL! Open source CLI. No DB required."
og_image: "/images/plugins/turbot/kubernetes-social-graphic.png"
---

# Kubernetes + Steampipe

[Steampipe](https://steampipe.io) is an open source CLI to instantly query cloud APIs using SQL.

[Kubernetes](https://kubernetes.io) is an open-source system for automating deployment, scaling, and management of containerized applications.

For example:

```sql
select
  name,
  namespace,
  phase,
  creation_timestamp,
  pod_ip
from
  kubernetes_pod;
```

```
+-----------------------------------------+-------------+-----------+---------------------+-----------+
| name                                    | namespace   | phase     | creation_timestamp  | pod_ip    |
+-----------------------------------------+-------------+-----------+---------------------+-----------+
| metrics-server-86cbb8457f-bf8dm         | kube-system | Running   | 2021-06-11 14:21:48 | 10.42.0.5 |
| coredns-7448499f4d-klb8l                | kube-system | Running   | 2021-06-11 14:21:48 | 10.42.0.6 |
| helm-install-traefik-crd-hb87d          | kube-system | Succeeded | 2021-06-11 14:21:48 | 10.42.0.3 |
| local-path-provisioner-5ff76fc89d-c9hnm | kube-system | Running   | 2021-06-11 14:21:48 | 10.42.0.2 |
+-----------------------------------------+-------------+-----------+---------------------+-----------+
```

## Documentation

- **[Table definitions & examples →](/plugins/turbot/kubernetes/tables)**

## Get started

### Install

Download and install the latest Kubernetes plugin:

```bash
steampipe plugin install kubernetes
```

### Configuration

Installing the latest kubernetes plugin will create a config file (`~/.steampipe/config/kubernetes.spc`) with a single connection named `kubernetes`:

```hcl
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

  # If no kubeconfig file can be found, the plugin will attempt to use the service account Kubernetes gives to pods.
  # This authentication method is intended for clients that expect to be running inside a pod running on Kubernetes.
}
```

- `config_context` - (Optional) The kubeconfig context to use. If not set, the current context will be used.
- `config_path` - (Optional) The kubeconfig file path. If not set, the plugin will check `~/.kube/config`. Can also be set with the `KUBE_CONFIG_PATHS` or `KUBERNETES_MASTER` environment variables. 

## Get involved

- Open source: https://github.com/turbot/steampipe-plugin-kubernetes
- Community: [Slack Channel](https://steampipe.io/community/join)

## Configuring Kubernetes Credentials

By default, the plugin will use the kubeconfig in `~/.kube/config` with the current context. If using the default kubectl CLI configurations, the kubeconfig will be in this location and the Kubernetes plugin connections will work by default.

You can also set the kubeconfig file path and context with the `config_path` and `config_context` config arguments respectively.

This plugin supports querying Kubernetes clusters using [OpenID Connect](https://kubernetes.io/docs/reference/access-authn-authz/authentication/#openid-connect-tokens) (OIDC) authentication. No extra configuration is required to query clusters using OIDC.

If no kubeconfig file is found, then the plugin will [attempt to access the API from within a pod](https://kubernetes.io/docs/tasks/run-application/access-api-from-pod/#accessing-the-api-from-within-a-pod) using the service account Kubernetes gives to pods. 
