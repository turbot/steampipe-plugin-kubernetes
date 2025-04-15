---
title: "Steampipe Table: kubernetes_replication_controller - Query Kubernetes Replication Controllers using SQL"
description: "Allows users to query Kubernetes Replication Controllers, providing insights into the status, configuration, and specifications of these controllers within a Kubernetes environment."
folder: "Replication Controller"
---

# Table: kubernetes_replication_controller - Query Kubernetes Replication Controllers using SQL

Kubernetes Replication Controllers are a core component of Kubernetes that ensure a specified number of pod replicas are running at any given time. They are particularly useful for stateless applications where more instances can be easily created or destroyed. Replication Controllers supersede the functionality of Kubernetes Pods by adding life-cycle control, system introspection, and self-healing mechanisms.

## Table Usage Guide

The `kubernetes_replication_controller` table provides insights into Replication Controllers within Kubernetes. As a DevOps engineer, you can explore controller-specific details through this table, including its status, configuration, and specifications. Utilize it to manage and monitor the state of your Kubernetes environment, ensuring the desired number of pod replicas are always running.

## Examples

### Basic Info
Explore the status of your Kubernetes replication controllers to understand the current state of your system. This can help you assess the number of desired, current, and ready replicas, and determine the age and selector details of each controller.

```sql+postgres
select
  name,
  namespace,
  replicas as desired,
  status_replicas as current,
  ready_replicas as ready,
  age(current_timestamp, creation_timestamp),
  selector
from
  kubernetes_replication_controller;
```

```sql+sqlite
select
  name,
  namespace,
  replicas as desired,
  status_replicas as current,
  ready_replicas as ready,
  (julianday('now') - julianday(creation_timestamp)) as age,
  selector
from
  kubernetes_replication_controller;
```

### Get details of containers and image
Explore the intricacies of your Kubernetes replication controllers by identifying the associated containers and images. This enables you to better understand the structure of your deployment, facilitating more effective management and troubleshooting.

```sql+postgres
select
  name,
  namespace,
  jsonb_agg(container.value -> 'name') as containers,
  jsonb_agg(container.value -> 'image') as images
from
  kubernetes_replication_controller,
  jsonb_array_elements(template -> 'spec' -> 'containers') as container
group by
  name,
  namespace;
```

```sql+sqlite
select
  name,
  namespace,
  json_group_array(json_extract(container.value, '$.name')) as containers,
  json_group_array(json_extract(container.value, '$.image')) as images
from
  kubernetes_replication_controller,
  json_each(json_extract(template, '$.spec.containers')) as container
group by
  name,
  namespace;
```

### List manifest resources
Explore the Kubernetes replication controllers with a specified path to understand their names, namespaces, and desired replicas. This can help in managing and monitoring the distribution and replication of workloads in a Kubernetes environment.

```sql+postgres
select
  name,
  namespace,
  replicas as desired,
  selector,
  path
from
  kubernetes_replication_controller
where
  path is not null;
```

```sql+sqlite
select
  name,
  namespace,
  replicas as desired,
  selector,
  path
from
  kubernetes_replication_controller
where
  path is not null;
```