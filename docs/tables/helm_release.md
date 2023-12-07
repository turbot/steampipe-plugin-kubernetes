---
title: "Steampipe Table: helm_release - Query Kubernetes Helm Releases using SQL"
description: "Allows users to query Helm Releases in Kubernetes, specifically information about the release, the chart used for the release, and the configuration values used."
---

# Table: helm_release - Query Kubernetes Helm Releases using SQL

A Helm Release in Kubernetes is a running instance of a chart, which can be installed in the same Kubernetes cluster multiple times. Each time a chart is installed, a new release is created. It records the version of the chart used and the configuration values set during installation.

## Table Usage Guide

The `helm_release` table provides insights into Helm Releases within Kubernetes. As a DevOps engineer, explore release-specific details through this table, including the chart used for the release, and the configuration values used. Utilize it to uncover information about releases, such as the version of the chart used, and the configuration values set during installation.

## Examples

### Basic info
Analyze the settings to understand the status and version of deployments in your Helm environment. This can be useful for keeping track of the various releases and their respective deployment times, ensuring you have a comprehensive overview of your system.

```sql+postgres
select
  name,
  namespace,
  version,
  status,
  first_deployed,
  chart_name
from
  helm_release;
```

```sql+sqlite
select
  name,
  namespace,
  version,
  status,
  first_deployed,
  chart_name
from
  helm_release;
```

### List kubernetes deployment resources deployed using a specific chart
Discover deployments in Kubernetes that have been implemented using a specific chart. This information can be vital for managing resources and understanding the distribution of deployments across your Kubernetes environment.

```sql+postgres
select
  d.name as deployment_name,
  d.namespace,
  d.uid,
  r.name as release_name,
  r.chart_name,
  r.version as release_version,
  r.first_deployed as deployed_at
from
  kubernetes_deployment as d
  left join helm_release as r
    on (
      d.labels ->> 'app.kubernetes.io/managed-by' = 'Helm'
      and d.labels ->> 'app.kubernetes.io/instance' = r.name
    )
where
  r.chart_name = 'ingress-nginx'
  and d.source_type = 'deployed';
```

```sql+sqlite
select
  d.name as deployment_name,
  d.namespace,
  d.uid,
  r.name as release_name,
  r.chart_name,
  r.version as release_version,
  r.first_deployed as deployed_at
from
  kubernetes_deployment as d
  left join helm_release as r
    on (
      json_extract(d.labels, '$.app.kubernetes.io/managed-by') = 'Helm'
      and json_extract(d.labels, '$.app.kubernetes.io/instance') = r.name
    )
where
  r.chart_name = 'ingress-nginx'
  and d.source_type = 'deployed';
```

### List all deployed releases of a specific chart
Explore the history of a specific deployment chart by listing all its releases. This allows you to assess the evolution and status of a specific deployment, providing valuable insights for ongoing management and future planning.

```sql+postgres
select
  name,
  namespace,
  version,
  status,
  first_deployed,
  chart_name
from
  helm_release
where
  chart_name = 'ingress-nginx';
```

```sql+sqlite
select
  name,
  namespace,
  version,
  status,
  first_deployed,
  chart_name
from
  helm_release
where
  chart_name = 'ingress-nginx';
```

### List releases from a specific namespace
Explore which releases have been deployed from a specific namespace in the Helm chart, allowing you to assess their status and understand the history of deployments. This is useful for managing and tracking the versions and statuses of your deployments within that namespace.

```sql+postgres
select
  name,
  namespace,
  version,
  status,
  last_deployed,
  description
from
  helm_release
where
  namespace = 'steampipe';
```

```sql+sqlite
select
  name,
  namespace,
  version,
  status,
  last_deployed,
  description
from
  helm_release
where
  namespace = 'steampipe';
```

### List all failed releases
Uncover the details of unsuccessful software deployments. This query is particularly useful in identifying problematic releases for further investigation and troubleshooting.

```sql+postgres
select
  name,
  namespace,
  version,
  status,
  last_deployed,
  description
from
  helm_release
where
  status = 'failed';
```

```sql+sqlite
select
  name,
  namespace,
  version,
  status,
  last_deployed,
  description
from
  helm_release
where
  status = 'failed';
```

### List all unfinished releases
Determine the areas in which there are pending releases that require further action or attention, allowing for efficient management and completion of these outstanding tasks.

```sql+postgres
select
  name,
  namespace,
  version,
  status,
  last_deployed,
  description
from
  helm_release
where
  status = 'pending';
```

```sql+sqlite
select
  name,
  namespace,
  version,
  status,
  last_deployed,
  description
from
  helm_release
where
  status = 'pending';
```

### List releases updated in last 3 days
Identify recent updates to your system by pinpointing releases that have been deployed in the last three days. This allows for efficient tracking and management of recent changes, ensuring your system remains up-to-date.

```sql+postgres
select
  name,
  namespace,
  version,
  status,
  last_deployed,
  description
from
  helm_release
where
  last_deployed > (now() - interval '3 days');
```

```sql+sqlite
select
  name,
  namespace,
  version,
  status,
  last_deployed,
  description
from
  helm_release
where
  last_deployed > datetime('now', '-3 days');
```

### Get a specific release
Analyze the settings to understand the status and details of a specific software release deployed using Helm, a package manager for Kubernetes. This query can be useful in tracking software versions and deployment statuses, providing insights into the efficiency of software management processes.

```sql+postgres
select
  name,
  namespace,
  version,
  status,
  last_deployed,
  description
from
  helm_release
where
  name = 'brigade-github-app-1683552635';
```

```sql+sqlite
select
  name,
  namespace,
  version,
  status,
  last_deployed,
  description
from
  helm_release
where
  name = 'brigade-github-app-1683552635';
```