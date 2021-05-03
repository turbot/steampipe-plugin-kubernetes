---
organization: Turbot
category: ["software development"]
icon_url: "/images/plugins/turbot/kubernetes.svg"
brand_color: "#326CE5"
display_name: "Kubernetes"
short_name: "kubernetes"
description: "Steampipe plugin for Kubernetes components."
og_description: Query Kubernetes with SQL! Open source CLI. No DB required.
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
+------------------------------------+-------------+---------+---------------------+
| name                               | namespace   | phase   | pod_ip              |
+------------------------------------+-------------+---------+---------------------+
| antrea-node-init-v7dd7             | kube-system | Running | 2021-04-22 13:35:16 |
| event-exporter-gke-564fb97f9-pvtvg | kube-system | Running | 2021-04-12 06:00:19 |
| fluentbit-gke-dkqsj                | kube-system | Running | 2021-04-12 06:00:08 |
+------------------------------------+-------------+---------+---------------------+
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
- Community: [Discussion forums](https://github.com/turbot/steampipe/discussions)

## Advanced configuration options

If you have a kube config setup using the kubectl CLI Steampipe just works with that connection.

The Kubernetes plugin allows you set the name of kube kubectl context with the `config_context` argument in any connection profile. You may also specify the path to kube config file with `config_path` argument.

### Credentials via kube config

```hcl
connection "k8s_minikube" {
  plugin         = "kubernetes"
  config_context = "minikube"
  # config_path    = "~/.kube/config"
}
```
