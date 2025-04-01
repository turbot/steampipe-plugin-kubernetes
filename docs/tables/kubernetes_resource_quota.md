---
title: "Steampipe Table: kubernetes_resource_quota - Query Kubernetes Resource Quotas using SQL"
description: "Allows users to query Resource Quotas in Kubernetes, providing insights into resource usage and restrictions within a namespace."
folder: "Resource Quota"
---

# Table: kubernetes_resource_quota - Query Kubernetes Resource Quotas using SQL

A Resource Quota in Kubernetes is a tool that administrators use to manage resources within a namespace. It sets hard limits on the amount of compute resources that can be used by a namespace in a Kubernetes cluster. This includes CPU and memory resources, the number of pods, services, volumes, and more.

## Table Usage Guide

The `kubernetes_resource_quota` table provides insights into resource quotas within Kubernetes. As a Kubernetes administrator, you can use this table to explore quota-specific details, including resource usage and restrictions within a namespace. Utilize it to uncover information about resource quotas, such as those nearing their limit, and effectively manage resources within your Kubernetes cluster.

## Examples

### Basic Info
Explore the basic information of your Kubernetes resource quotas to understand their allocation. This can help in managing and optimizing resource usage within your Kubernetes environment.

```sql+postgres
select
  name,
  namespace,
  resource_version,
  creation_timestamp,
  jsonb_pretty(spec_hard) as spec_hard
from
  kubernetes_resource_quota
order by
  name;
```

```sql+sqlite
select
  name,
  namespace,
  resource_version,
  creation_timestamp,
  spec_hard
from
  kubernetes_resource_quota
order by
  name;
```

### Get used pod details of namespaces
Discover the segments that are consuming resources in your Kubernetes environment by identifying how many pods and services are currently being used within each namespace. This is beneficial for managing resource allocation and identifying potential areas of overuse or inefficiency.

```sql+postgres
select
  name,
  namespace,
  status_used -> 'pods' as used_pods,
  status_used -> 'services' as used_services
from
  kubernetes_resource_quota;
```

```sql+sqlite
select
  name,
  namespace,
  json_extract(status_used, '$.pods') as used_pods,
  json_extract(status_used, '$.services') as used_services
from
  kubernetes_resource_quota;
```

### List manifest resources
Analyze the configuration of Kubernetes to identify resource quotas with specific paths. This is beneficial in managing resources efficiently by understanding their allocation and usage within your Kubernetes environment.

```sql+postgres
select
  name,
  namespace,
  resource_version,
  jsonb_pretty(spec_hard) as spec_hard,
  path
from
  kubernetes_resource_quota
where
  path is not null
order by
  name;
```

```sql+sqlite
select
  name,
  namespace,
  resource_version,
  spec_hard,
  path
from
  kubernetes_resource_quota
where
  path is not null
order by
  name;
```