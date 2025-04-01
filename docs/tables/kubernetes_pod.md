---
title: "Steampipe Table: kubernetes_pod - Query Kubernetes Pods using SQL"
description: "Allows users to query Kubernetes Pods, providing insights into the status, configuration, and usage of Pods within a Kubernetes cluster."
folder: "Pod"
---

# Table: kubernetes_pod - Query Kubernetes Pods using SQL

Kubernetes Pods are the smallest and simplest unit in the Kubernetes object model that you create or deploy. A Pod represents a running process on your cluster and encapsulates an application's container (or a group of tightly-coupled containers), storage resources, a unique network IP, and options that govern how the container(s) should run. Pods are designed to support co-located (co-scheduled), co-managed helper programs, such as content management systems, file and data loaders, local cache managers, etc.

## Table Usage Guide

The `kubernetes_pod` table provides insights into the Pods within a Kubernetes cluster. As a DevOps engineer, explore Pod-specific details through this table, including status, configuration, and usage. Utilize it to uncover information about Pods, such as their current state, the containers running within them, and the resources they are consuming.

## Examples

### Basic Info
Analyze the settings to understand the distribution and status of your Kubernetes pods. This query helps you identify the number of each type of container within each pod, as well as their age, phase, and the node they're running on, providing a comprehensive view of your Kubernetes environment.

```sql+postgres
select
  namespace,
  name,
  phase,
  age(current_timestamp, creation_timestamp),
  pod_ip,
  node_name,
  jsonb_array_length(containers) as container_count,
  jsonb_array_length(init_containers) as init_container_count,
  jsonb_array_length(ephemeral_containers) as ephemeral_container_count
from
  kubernetes_pod
order by
  namespace,
  name;
```

```sql+sqlite
select
  namespace,
  name,
  phase,
  (julianday('now') - julianday(datetime(creation_timestamp, 'unixepoch'))) as age,
  pod_ip,
  node_name,
  json_array_length(containers) as container_count,
  json_array_length(init_containers) as init_container_count,
  json_array_length(ephemeral_containers) as ephemeral_container_count
from
  kubernetes_pod
order by
  namespace,
  name;
```

### List Unowned (Naked) Pods
Discover the segments that consist of unassigned pods within your Kubernetes system. This query is useful for identifying potential resource inefficiencies or orphaned pods that could impact system performance.

```sql+postgres
select
  name,
  namespace,
  phase,
  pod_ip,
  node_name
from
  kubernetes_pod
where
  owner_references is null;
```

```sql+sqlite
select
  name,
  namespace,
  phase,
  pod_ip,
  node_name
from
  kubernetes_pod
where
  owner_references is null;
```

### List Privileged Containers
Discover the segments that are running privileged containers within your Kubernetes pods. This can help in identifying potential security risks and ensuring that your pods are following best practices for security configurations.

```sql+postgres
select
  name as pod_name,
  namespace,
  phase,
  jsonb_pretty(owner_references) as owners,
  c ->> 'name' as container_name,
  c ->> 'image' as container_image
from
  kubernetes_pod,
  jsonb_array_elements(containers) as c
where
  c -> 'securityContext' ->> 'privileged' = 'true';
```

```sql+sqlite
select
  name as pod_name,
  namespace,
  phase,
  owner_references as owners,
  json_extract(c.value, '$.name') as container_name,
  json_extract(c.value, '$.image') as container_image
from
  kubernetes_pod,
  json_each(containers) as c
where
  json_extract(c.value, '$.securityContext.privileged') = 'true';
```

### List Pods with access to the host process ID, IPC, or network namespace
Explore which Kubernetes pods have access to critical host resources such as the process ID, IPC, or network namespace. This is useful for identifying potential security risks and ensuring proper resource isolation.

```sql+postgres
select
  name,
  namespace,
  phase,
  host_pid,
  host_ipc,
  host_network,
  jsonb_pretty(owner_references) as owners
from
  kubernetes_pod
where
  host_pid or host_ipc or host_network;
```

```sql+sqlite
select
  name,
  namespace,
  phase,
  host_pid,
  host_ipc,
  host_network,
  owner_references as owners
from
  kubernetes_pod
where
  host_pid or host_ipc or host_network;
```

### Container Statuses
Explore the status of various containers within a Kubernetes pod, including their readiness and restart count. This query can be used to monitor and manage the health and performance of your Kubernetes environment.

```sql+postgres
select
  namespace,
  name as pod_name,
  phase,
  cs ->> 'name' as container_name,
  cs ->> 'image' as image,
  cs ->> 'ready' as ready,
  cs_state as state,
  cs ->> 'started' as started,
  cs ->> 'restartCount' as restarts
from
  kubernetes_pod,
  jsonb_array_elements(container_statuses) as cs,
  jsonb_object_keys(cs -> 'state') as cs_state
order by
  namespace,
  name,
  container_name;
```

```sql+sqlite
select
  namespace,
  name as pod_name,
  phase,
  json_extract(cs.value, '$.name') as container_name,
  json_extract(cs.value, '$.image') as image,
  json_extract(cs.value, '$.ready') as ready,
  json_each.key as state,
  json_extract(cs.value, '$.started') as started,
  json_extract(cs.value, '$.restartCount') as restarts
from
  kubernetes_pod,
  json_each(container_statuses) as cs,
  json_each(json_extract(cs.value, '$.state'))
order by
  namespace,
  name,
  container_name;
```

### `kubectl get pods` columns
This example allows you to monitor the status and performance of your Kubernetes pods. It provides insights into various aspects such as the number of running containers, total restarts, and the age of each pod, helping you to maintain the health and efficiency of your Kubernetes environment.

```sql+postgres
select
  namespace,
  name,
  phase,
  count(cs) filter (
    where
      cs -> 'state' -> 'running' is not null
  ) as running_container_count,
  jsonb_array_length(containers) as container_count,
  age(current_timestamp, creation_timestamp),
  COALESCE(sum((cs ->> 'restartCount') :: int), 0) as restarts,
  pod_ip,
  node_name
from
  kubernetes_pod
  left join jsonb_array_elements(container_statuses) as cs on true
group by
  namespace,
  name,
  phase,
  containers,
  creation_timestamp,
  pod_ip,
  node_name
 order by
  namespace,
  name;
```

```sql+sqlite
select
  namespace,
  name,
  phase,
  (select count(*) from json_each(container_statuses) as cs where json_extract(cs.value, '$.state.running') is not null) as running_container_count,
  (select count(*) from json_each(containers)) as container_count,
  julianday('now') - julianday(creation_timestamp) as age,
  COALESCE((select sum(json_extract(cs.value, '$.restartCount')) from json_each(container_statuses) as cs), 0) as restarts,
  pod_ip,
  node_name
from
  kubernetes_pod
group by
  namespace,
  name,
  phase,
  containers,
  creation_timestamp,
  pod_ip,
  node_name
order by
  namespace,
  name;
```

### List manifest resources
Explore which Kubernetes pods contain manifest resources, including the number of different container types. This can help you understand the distribution and configuration of resources within your Kubernetes environment.

```sql+postgres
select
  namespace,
  name,
  phase,
  pod_ip,
  node_name,
  jsonb_array_length(containers) as container_count,
  jsonb_array_length(init_containers) as init_container_count,
  jsonb_array_length(ephemeral_containers) as ephemeral_container_count,
  path
from
  kubernetes_pod
where
  path is not null
order by
  namespace,
  name;
```

```sql+sqlite
select
  namespace,
  name,
  phase,
  pod_ip,
  node_name,
  json_array_length(containers) as container_count,
  json_array_length(init_containers) as init_container_count,
  json_array_length(ephemeral_containers) as ephemeral_container_count,
  path
from
  kubernetes_pod
where
  path is not null
order by
  namespace,
  name;
```