---
title: "Steampipe Table: kubernetes_limit_range - Query Kubernetes Limit Ranges using SQL"
description: "Allows users to query Kubernetes Limit Ranges, specifically the range of constraints for resources such as CPU and memory that can be consumed by containers in a namespace."
---

# Table: kubernetes_limit_range - Query Kubernetes Limit Ranges using SQL

Kubernetes Limit Range is a policy to constrain resource allocation (CPU, memory, etc.) in a namespace. It configures the minimum and maximum compute resources that are allowed for different types of Kubernetes objects (Pod, Container, PersistentVolumeClaim, etc.). It helps to control the resource consumption and ensure the efficient use of resources across all Pods and Containers in a namespace.

A LimitRange provides constraints that can:

- Enforce minimum and maximum compute resources usage per Pod or Container in a namespace.
- Enforce minimum and maximum storage request per PersistentVolumeClaim in a namespace.
- Enforce a ratio between request and limit for a resource in a namespace.
- Set default request/limit for compute resources in a namespace and automatically inject them to Containers at runtime.

## Table Usage Guide

The `kubernetes_limit_range` table provides insights into the limit ranges within Kubernetes. As a DevOps engineer, explore limit range-specific details through this table, including the types of resources being constrained, their minimum and maximum values, and the namespace in which they are applied. Utilize it to manage and optimize resource allocation across all Pods and Containers in a namespace.

## Examples

### Basic Info
Explore which resources have specific limits within your Kubernetes environment. This can help you manage resources effectively by understanding their configurations and creation times.

```sql+postgres
select
  name,
  namespace,
  resource_version,
  creation_timestamp,
  jsonb_pretty(spec_limits) as spec_limits
from
  kubernetes_limit_range
order by
  namespace;
```

```sql+sqlite
select
  name,
  namespace,
  resource_version,
  creation_timestamp,
  spec_limits
from
  kubernetes_limit_range
order by
  namespace;
```

### Get spec limits details of limit range
Assess the elements within your Kubernetes limit range to understand the specifics of each limit type, including their default values and requests. This allows you to manage resource consumption effectively by identifying the parameters that define the minimum and maximum resource usage.

```sql+postgres
select
  name,
  namespace,
  limits ->> 'type' as type,
  limits ->> 'default' as default,
  limits ->> 'defaultRequest' as default_request
from
  kubernetes_limit_range,
  jsonb_array_elements(spec_limits) as limits;
```

```sql+sqlite
select
  name,
  namespace,
  json_extract(limits.value, '$.type') as type,
  json_extract(limits.value, '$.default') as default,
  json_extract(limits.value, '$.defaultRequest') as default_request
from
  kubernetes_limit_range,
  json_each(spec_limits) as limits;
```

### List manifest resources
Explore the specific limits set for resources in different namespaces of a Kubernetes cluster. This can help in managing resource allocation and ensuring optimal performance.

```sql+postgres
select
  name,
  namespace,
  resource_version,
  jsonb_pretty(spec_limits) as spec_limits,
  path
from
  kubernetes_limit_range
where
  path is not null
order by
  namespace;
```

```sql+sqlite
select
  name,
  namespace,
  resource_version,
  spec_limits,
  path
from
  kubernetes_limit_range
where
  path is not null
order by
  namespace;
```