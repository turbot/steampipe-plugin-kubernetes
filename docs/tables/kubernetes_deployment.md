---
title: "Steampipe Table: kubernetes_deployment - Query Kubernetes Deployments using SQL"
description: "Allows users to query Deployments in Kubernetes, specifically information about deployed applications and their replicas, providing insights into application management and scaling."
---

# Table: kubernetes_deployment - Query Kubernetes Deployments using SQL

A Kubernetes Deployment is a resource object in Kubernetes that provides declarative updates for applications. It allows you to describe an application's life-cycle, such as which images to use for the app, the number of pod replicas, and the way to update them. In addition, it offers advanced features such as rollback and scaling.

## Table Usage Guide

The `kubernetes_deployment` table provides insights into Deployments within Kubernetes. As a DevOps engineer, explore Deployment-specific details through this table, including the current and desired states, strategy used for updates, and associated metadata. Utilize it to manage the life-cycle of your applications, such as scaling up/down the number of replicas, rolling updates, and rollback to earlier versions.

## Examples

### Basic Info
Explore the status and availability of various components within a deployment, helping you understand the overall health and readiness of your system. This is useful for ongoing monitoring and troubleshooting of your deployment.

```sql+postgres
select
  name,
  namespace,
  status_replicas,
  ready_replicas,
  updated_replicas,
  available_replicas,
  unavailable_replicas,
  age(current_timestamp, creation_timestamp)
from
  kubernetes_deployment
order by
  namespace,
  name;
```

```sql+sqlite
select
  name,
  namespace,
  status_replicas,
  ready_replicas,
  updated_replicas,
  available_replicas,
  unavailable_replicas,
  (julianday('now') - julianday(creation_timestamp)) * 86400.0
from
  kubernetes_deployment
order by
  namespace,
  name;
```

### Configuration Info
Explore the status of your deployments in the Kubernetes environment. This query is beneficial to understand the settings of your deployments, such as whether they are paused or active, their generation status, and their selection strategy.

```sql+postgres
select
  name,
  paused,
  generate_name,
  generation,
  revision_history_limit,
  replicas,
  strategy,
  selector
from
  kubernetes_deployment;
```

```sql+sqlite
select
  name,
  paused,
  generate_name,
  generation,
  revision_history_limit,
  replicas,
  strategy,
  selector
from
  kubernetes_deployment;
```

### Container Images used in Deployments
Explore which container images are currently in use across different deployments. This query is useful in managing and tracking the versions of images used, aiding in troubleshooting and ensuring consistency across deployments.

```sql+postgres
select
  name as deployment_name,
  namespace,
  c ->> 'name' as container_name,
  c ->> 'image' as image
from
  kubernetes_deployment,
  jsonb_array_elements(template -> 'spec' -> 'containers') as c
order by
  namespace,
  name;
```

```sql+sqlite
select
  name as deployment_name,
  namespace,
  json_extract(c.value, '$.name') as container_name,
  json_extract(c.value, '$.image') as image
from
  kubernetes_deployment,
  json_each(json_extract(template, '$.spec.containers')) as c
order by
  namespace,
  name;
```

### List pods for a deployment
Determine the areas in which pods are being used for a specific deployment in a Kubernetes environment. This can be useful for understanding the deployment's resource utilization and identifying potential areas for optimization or troubleshooting.

```sql+postgres
select
  pod.namespace,
  d.name as deployment_name,
  rs.name as replicaset_name,
  pod.name as pod_name,
  pod.phase,
  age(current_timestamp, pod.creation_timestamp),
  pod.pod_ip,
  pod.node_name
from
  kubernetes_pod as pod,
  jsonb_array_elements(pod.owner_references) as pod_owner,
  kubernetes_replicaset as rs,
  jsonb_array_elements(rs.owner_references) as rs_owner,
  kubernetes_deployment as d
where
  pod_owner ->> 'kind' = 'ReplicaSet'
  and rs.uid = pod_owner ->> 'uid'
  and rs_owner ->> 'uid' = d.uid
  and d.name = 'frontend'
order by
  pod.namespace,
  d.name,
  rs.name,
  pod.name;
```

```sql+sqlite
select
  pod.namespace,
  d.name as deployment_name,
  rs.name as replicaset_name,
  pod.name as pod_name,
  pod.phase,
  (julianday('now') - julianday(pod.creation_timestamp)) as age,
  pod.pod_ip,
  pod.node_name
from
  kubernetes_pod as pod,
  json_each(pod.owner_references) as pod_owner,
  kubernetes_replicaset as rs,
  json_each(rs.owner_references) as rs_owner,
  kubernetes_deployment as d
where
  json_extract(pod_owner.value, '$.kind') = 'ReplicaSet'
  and rs.uid = json_extract(pod_owner.value, '$.uid')
  and json_extract(rs_owner.value, '$.uid') = d.uid
  and d.name = 'frontend'
order by
  pod.namespace,
  d.name,
  rs.name,
  pod.name;
```

### List Pods with access to the to the host process ID, IPC, or network namespace
Identify the Kubernetes deployments that have access to the host's process ID, inter-process communication, or network namespace. This helps in pinpointing potential security vulnerabilities by highlighting the areas where a pod could potentially gain unauthorized access to sensitive host resources.

```sql+postgres
select
  namespace,
  name,
  template -> 'spec' ->> 'hostPID' as hostPID,
  template -> 'spec' ->> 'hostIPC' as hostIPC,
  template -> 'spec' ->> 'hostNetwork' as hostNetwork
from
  kubernetes_deployment
where
  template -> 'spec' ->> 'hostPID' = 'true' or
  template -> 'spec' ->> 'hostIPC' = 'true' or
  template -> 'spec' ->> 'hostNetwork' = 'true';
```

```sql+sqlite
select
  namespace,
  name,
  json_extract(template, '$.spec.hostPID') as hostPID,
  json_extract(template, '$.spec.hostIPC') as hostIPC,
  json_extract(template, '$.spec.hostNetwork') as hostNetwork
from
  kubernetes_deployment
where
  json_extract(template, '$.spec.hostPID') = 'true' or
  json_extract(template, '$.spec.hostIPC') = 'true' or
  json_extract(template, '$.spec.hostNetwork') = 'true';
```

### List manifest resources
Discover the segments that have allocated resources within a specific namespace in a Kubernetes deployment, allowing you to better manage and allocate resources efficiently. This is particularly useful in larger deployments where resource management is crucial.

```sql+postgres
select
  name,
  namespace,
  replicas,
  path
from
  kubernetes_deployment
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
  replicas,
  path
from
  kubernetes_deployment
where
  path is not null
order by
  namespace,
  name;
```