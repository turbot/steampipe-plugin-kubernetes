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
  plugin      = "kubernetes"
}
```

This will create a `kubernetes` connection that uses the default kubeconfig context.

## Get involved

- Open source: https://github.com/turbot/steampipe-plugin-kubernetes
- Community: [Slack Channel](https://steampipe.io/community/join)

## Advanced configuration options

If you have a kube config setup using the kubectl CLI Steampipe just works with that connection.

The Kubernetes plugin allows you set the name of kube kubectl context with the `config_context` argument in any connection profile. You may also specify the path to kube config file with `config_path` argument.

This plugin also supports querying Kubernetes clusters using [OpenID Connect](https://kubernetes.io/docs/reference/access-authn-authz/authentication/#openid-connect-tokens) (OIDC) authentication. No extra configuration is required in a connection profile to query clusters using OIDC.

This plugin also supports querying Kubernetes clusters using [InClusterConfig](https://kubernetes.io/docs/tasks/run-application/access-api-from-pod/#accessing-the-api-from-within-a-pod) configuration. No extra configuration is required in a connection profile to query clusters using InClusterConfig. If the `~/.kube/config file is not available, plugin will automatically look for InClusterConfig configuration.

### Credentials via kube config

```hcl
connection "k8s_minikube" {
  plugin         = "kubernetes"
  config_context = "minikube"
  # config_path    = "~/.kube/config"
}
```
