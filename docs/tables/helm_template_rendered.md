---
title: "Steampipe Table: helm_template_rendered - Query Kubernetes Helm Templates using SQL"
description: "Allows users to query Helm Templates in Kubernetes, specifically the rendered templates, providing insights into the configuration and deployment of applications within Kubernetes clusters."
---

# Table: helm_template_rendered - Query Kubernetes Helm Templates using SQL

A Helm Template in Kubernetes is a powerful tool that generates Kubernetes manifest files. It is a part of Helm, the package manager for Kubernetes, and is used to streamline the installation and management of applications within Kubernetes clusters. Helm Templates allow users to define, install, and upgrade complex Kubernetes applications, effectively serving as a deployment blueprint.

## Table Usage Guide

The `helm_template_rendered` table provides insights into Helm Templates within Kubernetes. As a DevOps engineer or a Kubernetes administrator, explore the details of rendered templates through this table, including the configuration and deployment of applications within Kubernetes clusters. Utilize it to verify the deployment specifications, understand the configuration of applications, and manage the lifecycle of Kubernetes applications.

## Examples

### List fully rendered kubernetes resource templates defined in a chart
Explore the fully processed resource templates within a specific Kubernetes chart to understand its configuration. This is useful for assessing the elements within a given chart, such as 'redis', for effective management and troubleshooting.

```sql+postgres
select
  path,
  source_type,
  rendered
from
  helm_template_rendered
where
  chart_name = 'redis';
```

```sql+sqlite
select
  path,
  source_type,
  rendered
from
  helm_template_rendered
where
  chart_name = 'redis';
```

### List fully rendered kubernetes resource templates for different environments
Explore the fully rendered resource templates for different environments in Kubernetes. This is useful to understand the configuration for specific applications in development and production environments.
Let's say you have two different environments for maintaining your app: dev and prod. And, you have a helm chart with 2 different set of values for your environments. For example:

```hcl
connection "kubernetes" {
  plugin = "kubernetes"

  helm_rendered_charts = {
    "my-app-dev" = {
      chart_path        = "~/charts/my-app"
      values_file_paths = ["~/value/file/for/dev.yaml"]
    }
    "my-app-prod" = {
      chart_path        = "~/charts/my-app"
      values_file_paths = ["~/value/file/for/prod.yaml"]
    }
  }
}
```

In both case, it is using same chart with a different set of values.

To list the kubernetes resource configurations defined for the dev environment, you can simply run the below query:


```sql+postgres
select
  chart_name,
  path,
  source_type,
  rendered
from
  helm_template_rendered
where
  source_type = 'helm_rendered:my-app-dev';
```

```sql+sqlite
select
  chart_name,
  path,
  source_type,
  rendered
from
  helm_template_rendered
where
  source_type = 'helm_rendered:my-app-dev';
```

Similarly, to query the kubernetes resource configurations for prod,

```sql+postgres
select
  chart_name,
  path,
  source_type,
  rendered
from
  helm_template_rendered
where
  source_type = 'helm_rendered:my-app-prod';
```

```sql+sqlite
select
  chart_name,
  path,
  source_type,
  rendered
from
  helm_template_rendered
where
  source_type = 'helm_rendered:my-app-prod';
```