---
title: "Steampipe Table: kubernetes_horizontal_pod_autoscaler - Query Kubernetes Horizontal Pod Autoscalers using SQL"
description: "Allows users to query Kubernetes Horizontal Pod Autoscalers, providing insights into the configuration and current status of autoscalers in the Kubernetes cluster."
---

# Table: kubernetes_horizontal_pod_autoscaler - Query Kubernetes Horizontal Pod Autoscalers using SQL

A Kubernetes Horizontal Pod Autoscaler automatically scales the number of pods in a replication controller, deployment, replica set, or stateful set based on observed CPU utilization. It is designed to maintain a specified amount of CPU utilization across the pods, irrespective of the load. The Horizontal Pod Autoscaler is implemented as a Kubernetes API resource and a controller.

## Table Usage Guide

The `kubernetes_horizontal_pod_autoscaler` table provides insights into Horizontal Pod Autoscalers within a Kubernetes cluster. As a Kubernetes administrator or DevOps engineer, explore autoscaler-specific details through this table, including the current and desired number of replicas, target CPU utilization, and associated metadata. Utilize it to monitor the performance and efficiency of the autoscalers, ensuring that your applications are scaling correctly and efficiently under varying load conditions.

## Examples

### Basic Info
Explore the configuration of your Kubernetes horizontal pod autoscaler to determine its current and desired replica settings. This will help you understand how your system is scaling and whether it is operating within your set parameters.

```sql+postgres
select
  name,
  namespace,
  min_replicas,
  max_replicas,
  current_replicas,
  desired_replicas
from
  kubernetes_horizontal_pod_autoscaler;
```

```sql+sqlite
The PostgreSQL query provided does not use any PostgreSQL-specific functions, data types, or JSON functions. Therefore, the query can be used in SQLite without any changes.

Here is the SQLite equivalent:

```sql
select
  name,
  namespace,
  min_replicas,
  max_replicas,
  current_replicas,
  desired_replicas
from
  kubernetes_horizontal_pod_autoscaler;
```
```

### Get list of HPA metrics configurations
Explore the configurations of your Horizontal Pod Autoscalers (HPA) to understand their current and desired replica settings. This can help you assess the efficiency of your current setup and identify areas for potential optimization.

```sql+postgres
select
  name,
  namespace,
  min_replicas,
  max_replicas,
  current_replicas,
  desired_replicas,
  jsonb_array_elements(metrics) as metrics,
  jsonb_array_elements(current_metrics) as current_metrics,
  conditions
from
  kubernetes_horizontal_pod_autoscaler;
```

```sql+sqlite
select
  name,
  namespace,
  min_replicas,
  max_replicas,
  current_replicas,
  desired_replicas,
  metrics,
  current_metrics,
  conditions
from
  kubernetes_horizontal_pod_autoscaler,
  json_each(metrics),
  json_each(current_metrics);
```

### List manifest resources
Explore which Kubernetes horizontal pod autoscalers have a defined path. This helps in understanding the autoscaling configuration for the pods and aids in optimizing resource usage within your Kubernetes environment.

```sql+postgres
select
  name,
  namespace,
  min_replicas,
  max_replicas,
  current_replicas,
  desired_replicas,
  path
from
  kubernetes_horizontal_pod_autoscaler
where
  path is not null;
```

```sql+sqlite
select
  name,
  namespace,
  min_replicas,
  max_replicas,
  current_replicas,
  desired_replicas,
  path
from
  kubernetes_horizontal_pod_autoscaler
where
  path is not null;
```