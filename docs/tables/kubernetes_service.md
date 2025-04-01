---
title: "Steampipe Table: kubernetes_service - Query Kubernetes Services using SQL"
description: "Allows users to query Kubernetes Services, particularly the configuration and status of services within a Kubernetes cluster."
folder: "Service"
---

# Table: kubernetes_service - Query Kubernetes Services using SQL

Kubernetes Service is a resource within Kubernetes that is used to expose an application running on a set of Pods. The set of Pods targeted by a Service is determined by a Label Selector. It provides the abstraction of a logical set of Pods and a policy by which to access them, often referred to as micro-services.

## Table Usage Guide

The `kubernetes_service` table offers insights into the services within a Kubernetes cluster. As a DevOps engineer, you can probe service-specific details through this table, including service configurations, status, and associated metadata. Use it to discover information about services, such as those with specific selectors, the type of service, and the ports exposed by the service.

## Examples

### Basic Info - `kubectl describe service --all-namespaces` columns
Analyze the settings of your Kubernetes services to understand their organization and longevity. This query is useful for gaining insights into how your services are distributed across namespaces, their types, and how long they have been active.

```sql+postgres
select
  name,
  namespace,
  type,
  cluster_ip,
  age(current_timestamp, creation_timestamp)
from
  kubernetes_service
order by
  namespace,
  name;
```

```sql+sqlite
select
  name,
  namespace,
  type,
  cluster_ip,
  strftime('%s', 'now') - strftime('%s', creation_timestamp) as age
from
  kubernetes_service
order by
  namespace,
  name;
```

### List manifest resources
Analyze the settings to understand the distribution of resources within a Kubernetes cluster. This can help to identify instances where resources are not properly allocated, improving the efficiency of the cluster.

```sql+postgres
select
  name,
  namespace,
  type,
  cluster_ip,
  path
from
  kubernetes_service
where
  path is not null
order by
  namespace,
  name;
```

```sql+sqlite
select
  name,
  namespace,
  type,
  cluster_ip,
  path
from
  kubernetes_service
where
  path is not null
order by
  namespace,
  name;
```