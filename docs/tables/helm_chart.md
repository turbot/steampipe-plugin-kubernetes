---
title: "Steampipe Table: helm_chart - Query Kubernetes Helm Charts using SQL"
description: "Allows users to query Kubernetes Helm Charts, providing insights into chart details like its version, status, and associated metadata."
folder: "Helm"
---

# Table: helm_chart - Query Kubernetes Helm Charts using SQL

Kubernetes Helm is a package manager for Kubernetes that allows developers and operators to more easily package, configure, and deploy applications and services onto Kubernetes clusters. Helm uses a packaging format called charts, which include all of the Kubernetes resources that a particular application needs to run. A Helm chart can provide information about the version of the application, the Kubernetes resources that will be used, and any other application-specific information.

## Table Usage Guide

The `helm_chart` table provides insights into Helm Charts within Kubernetes. As a DevOps engineer, explore chart-specific details through this table, including version, status, and associated metadata. Utilize it to uncover information about charts, such as their current status, the version of the application they are deploying, and other application-specific information.

## Examples

### Basic info
Explore which Helm charts are deprecated by analyzing their basic information, including name, version, and type. This can help in maintaining up-to-date and secure applications by avoiding the use of outdated or deprecated charts.

```sql+postgres
select
  name,
  api_version,
  version,
  deprecated,
  type,
  description
from
  helm_chart;
```

```sql+sqlite
select
  name,
  api_version,
  version,
  deprecated,
  type,
  description
from
  helm_chart;
```

### List all deployed charts
Discover the segments that have been deployed in your Kubernetes environment. This query can be used to get insights into the status, configuration, and version details of all deployed charts, helping you manage and track your deployments effectively.

```sql+postgres
select
  hc.name as chart_name,
  hc.type as chart_type,
  hc.version as chart_version,
  hr.name as release_name,
  hr.status as release_status,
  hr.version as deployment_version,
  hr.first_deployed as deployed_at,
  hr.config as deployment_config
from
  kubernetes.helm_chart as hc
  left join kubernetes.helm_release as hr
    on hc.name = hr.chart_name
where
  hr.status = 'deployed';
```

```sql+sqlite
select
  hc.name as chart_name,
  hc.type as chart_type,
  hc.version as chart_version,
  hr.name as release_name,
  hr.status as release_status,
  hr.version as deployment_version,
  hr.first_deployed as deployed_at,
  hr.config as deployment_config
from
  kubernetes_helm_chart as hc
  left join kubernetes_helm_release as hr
    on hc.name = hr.chart_name
where
  hr.status = 'deployed';
```

### List all deprecated charts
Discover the segments that contain deprecated charts within your system. This is particularly useful for identifying outdated elements and ensuring your system stays up-to-date.

```sql+postgres
select
  name,
  api_version,
  version,
  type,
  description,
  app_version
from
  helm_chart
where
  deprecated;
```

```sql+sqlite
select
  name,
  api_version,
  version,
  type,
  description,
  app_version
from
  helm_chart
where
  deprecated = 1;
```

### List application type charts
Explore the different charts related to application type in Helm to gain insights into their names, versions, API versions, and descriptions. This can help you understand the variety and specifications of application charts available, assisting in better application management.

```sql+postgres
select
  name,
  api_version,
  version,
  type,
  description,
  app_version
from
  helm_chart
where
  type = 'application';
```

```sql+sqlite
select
  name,
  api_version,
  version,
  type,
  description,
  app_version
from
  helm_chart
where
  type = 'application';
```