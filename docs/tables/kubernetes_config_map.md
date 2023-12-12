---
title: "Steampipe Table: kubernetes_config_map - Query Kubernetes Config Maps using SQL"
description: "Allows users to query Kubernetes Config Maps, providing insights into configuration data and application settings within a Kubernetes cluster."
---

# Table: kubernetes_config_map - Query Kubernetes Config Maps using SQL

Kubernetes Config Maps is a resource that allows you to decouple configuration artifacts from image content to keep containerized applications portable. It is used to store non-confidential data in key-value pairs and consumed by pods or used to store configuration details, such as environment variables for a pod. Kubernetes Config Maps offers a centralized and secure method to manage and deploy configuration data.

## Table Usage Guide

The `kubernetes_config_map` table provides insights into Config Maps within Kubernetes. As a DevOps engineer, explore Config Map-specific details through this table, including data, creation timestamps, and associated metadata. Utilize it to uncover information about Config Maps, such as those used in specific namespaces, the configuration details they hold, and the pods that may be consuming them.

## Examples

### Basic Info
Explore the age and details of Kubernetes configuration maps to understand their longevity and content. This can help you manage and optimize your Kubernetes resources over time.

```sql+postgres
select
  name,
  namespace,
  data.key,
  data.value,
  age(current_timestamp, creation_timestamp)
from
  kubernetes_config_map,
  jsonb_each(data) as data
order by
  namespace,
  name;
```

```sql+sqlite
select
  name,
  namespace,
  data.key,
  data.value,
  strftime('%s', 'now') - strftime('%s', creation_timestamp) as age
from
  kubernetes_config_map,
  json_each(data) as data
order by
  namespace,
  name;
```

### List manifest resources
Analyze the settings to understand the distribution of resources across different namespaces within your Kubernetes environment. This can help in managing resources effectively and preventing any potential conflicts or overlaps.

```sql+postgres
select
  name,
  namespace,
  data.key,
  data.value,
  path
from
  kubernetes_config_map,
  jsonb_each(data) as data
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
  data.key,
  data.value,
  path
from
  kubernetes_config_map,
  json_each(data) as data
where
  path is not null
order by
  namespace,
  name;
```