---
title: "Steampipe Table: kubernetes_namespace - Query Kubernetes Namespaces using SQL"
description: "Allows users to query Kubernetes Namespaces, specifically the metadata and status of each namespace, providing insights into resource allocation and usage."
---

# Table: kubernetes_namespace - Query Kubernetes Namespaces using SQL

Kubernetes Namespaces are an abstraction used by Kubernetes to support multiple virtual clusters on the same physical cluster. These namespaces provide a scope for names, and they are intended to be used in environments with many users spread across multiple teams, or projects. Namespaces are a way to divide cluster resources between multiple uses.

## Table Usage Guide

The `kubernetes_namespace` table provides insights into Namespaces within Kubernetes. As a DevOps engineer, explore namespace-specific details through this table, including metadata, status, and associated resources. Utilize it to uncover information about namespaces, such as their status, the resources allocated to them, and their overall usage within the Kubernetes cluster.

## Examples

### Basic Info
Explore the status and metadata of different segments within your Kubernetes environment. This allows you to gain insights into the current operational phase and additional details of each namespace, aiding in effective resource management and monitoring.

```sql+postgres
select
  name,
  phase as status,
  annotations,
  labels
from
  kubernetes_namespace;
```

```sql+sqlite
select
  name,
  phase as status,
  annotations,
  labels
from
  kubernetes_namespace;
```

### List manifest resources
Uncover the details of each manifest resource within your Kubernetes namespace, including its status and associated annotations and labels. This is particularly useful for tracking resource utilization and identifying any potential issues or anomalies that may impact system performance.

```sql+postgres
select
  name,
  phase as status,
  annotations,
  labels,
  path
from
  kubernetes_namespace
where
  path is not null;
```

```sql+sqlite
select
  name,
  phase as status,
  annotations,
  labels,
  path
from
  kubernetes_namespace
where
  path is not null;
```