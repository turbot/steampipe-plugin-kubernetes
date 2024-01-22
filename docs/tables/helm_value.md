---
title: "Steampipe Table: helm_value - Query Kubernetes Helm Values using SQL"
description: "Allows users to query Helm Values in Kubernetes, specifically the configuration values for Helm Charts, providing insights into the configurations of different Kubernetes applications."
---

# Table: helm_value - Query Kubernetes Helm Values using SQL

Kubernetes Helm is a package manager for Kubernetes that allows developers and operators to more easily package, configure, and deploy applications and services onto Kubernetes clusters. Helm uses a packaging format called charts, and a chart is a collection of files that describe a related set of Kubernetes resources. Helm Values are the specific configurations for a Helm Chart.

## Table Usage Guide

The `helm_value` table provides insights into Helm Values within Kubernetes. As a DevOps engineer, explore Helm Value-specific details through this table, including the configurations of different Kubernetes applications and services. Utilize it to uncover information about Helm Values, such as those relating to specific Helm Charts, the configurations of different services, and the verification of configurations.

By design, applications can ship with default values.yaml file tuned for production deployments. Also, considering the multiple environments, it may have different configurations. To override the default value, it is not necessary to change the default values.yaml, but you can refer to the override value files from which it takes the configuration. For example:

Let's say you have two different environments for maintaining your app: dev and prod. And, you have a helm chart with 2 different set of values for your environments. For example:

```hcl
connection "kubernetes" {
  plugin = "kubernetes"

  helm_rendered_charts = {
    "my-app-dev" = {
      chart_path        = "~/charts/my-app"
      values_file_paths = "~/value/file/for/dev.yaml"
    }
    "my-app-prod" = {
      chart_path        = "~/charts/my-app"
      values_file_paths = "~/value/file/for/prod.yaml"
    }
  }
}
```

The table `helm_value` lists the values from the chart's default values.yaml file, as well as it lists the values from the files that are provided to override the default configuration.

**Important Notes**

- You must specify the `path` column in the `where` clause to query this table.

## Examples

### List values configured in the default values.yaml file of a specific chart
Analyze the settings to understand the default configurations set in a specific chart's values.yaml file in Helm, which is beneficial for auditing or modifying these configurations. This allows you to pinpoint the specific locations where changes have been made, enhancing your control over the chart's behavior.

```sql+postgres
select
  path,
  key_path,
  value,
  start_line,
  start_column
from
  helm_value
where
  path = '~/charts/my-app/values.yaml'
order by
  start_line;
```

```sql+sqlite
select
  path,
  key_path,
  value,
  start_line,
  start_column
from
  helm_value
where
  path = '~/charts/my-app/values.yaml'
order by
  start_line;
```

### List values from a specific override file
Explore which values are being used from a specific file in your Helm configuration. This can be particularly useful to understand and manage your development environment settings.

```sql+postgres
select
  path,
  key_path,
  value,
  start_line,
  start_column
from
  helm_value
where
  path = '~/value/file/for/dev.yaml'
order by
  start_line;
```

```sql+sqlite
select
  path,
  key_path,
  value,
  start_line,
  start_column
from
  helm_value
where
  path = '~/value/file/for/dev.yaml'
order by
  start_line;
```