---
title: "Steampipe Table: kubernetes_replicaset - Query Kubernetes ReplicaSets using SQL"
description: "Allows users to query Kubernetes ReplicaSets, providing details about the current state of each ReplicaSet, including the number of replicas, the desired number of replicas, and the labels and selectors used to identify its pods."
folder: "ReplicaSet"
---

# Table: kubernetes_replicaset - Query Kubernetes ReplicaSets using SQL

A ReplicaSet in Kubernetes is a resource that ensures that a specified number of pod replicas are running at any given time. It is often used to guarantee the availability of a specified number of identical pods. A ReplicaSet creates new pods when needed and removes old pods when too many are running.

## Table Usage Guide

The `kubernetes_replicaset` table provides insights into the ReplicaSets within a Kubernetes cluster. As a DevOps engineer or Kubernetes administrator, you can use this table to monitor the status and health of your ReplicaSets, including the current and desired number of replicas, as well as the labels and selectors used to identify its pods. This can be particularly useful for maintaining high availability and for troubleshooting issues with your applications running on Kubernetes.

## Examples

### Basic Info
Explore the status and configuration of your Kubernetes replica sets to understand their readiness, availability, and age. This can be useful to assess the health and stability of your application deployments.

```sql+postgres
select
  name,
  namespace,
  replicas as desired,
  ready_replicas as ready,
  available_replicas as available,
  selector,
  fully_labeled_replicas,
  age(current_timestamp, creation_timestamp)
from
  kubernetes_replicaset;
```

```sql+sqlite
select
  name,
  namespace,
  replicas as desired,
  ready_replicas as ready,
  available_replicas as available,
  selector,
  fully_labeled_replicas,
  strftime('%s', 'now') - strftime('%s', creation_timestamp) as age
from
  kubernetes_replicaset;
```

### Get container and image used in the replicaset
Gain insights into the relationship between containers and their corresponding images within a replicaset, helping to manage and track the utilization of resources in a Kubernetes environment. This query is particularly useful for administrators looking to optimize their deployments.

```sql+postgres
select
  name,
  namespace,
  c ->> 'name' as container_name,
  c ->> 'image' as image
from
  kubernetes_replicaset,
  jsonb_array_elements(template -> 'spec' -> 'containers') as c
order by
  namespace,
  name;
```

```sql+sqlite
select
  name,
  namespace,
  json_extract(c.value, '$.name') as container_name,
  json_extract(c.value, '$.image') as image
from
  kubernetes_replicaset,
  json_each(json_extract(template, '$.spec.containers')) as c
order by
  namespace,
  name;
```

### List pods for a replicaset (by name)
Discover the details of pods associated with a specific replicaset in a Kubernetes environment. This is useful in monitoring and managing the pods that belong to a particular replicaset, ensuring the replicaset is functioning as expected.

```sql+postgres
select
  pod.namespace,
  rs.name as replicaset_name,
  pod.name as pod_name,
  pod.phase,
  age(current_timestamp, pod.creation_timestamp),
  pod.pod_ip,
  pod.node_name
from
  kubernetes_pod as pod,
  jsonb_array_elements(pod.owner_references) as pod_owner,
  kubernetes_replicaset as rs
where
  pod_owner ->> 'kind' = 'ReplicaSet'
  and rs.uid = pod_owner ->> 'uid'
  and rs.name = 'frontend-56fc5b6b47'
order by
  pod.namespace,
  rs.name,
  pod.name;
```

```sql+sqlite
select
  pod.namespace,
  rs.name as replicaset_name,
  pod.name as pod_name,
  pod.phase,
  strftime('%s', 'now') - strftime('%s', pod.creation_timestamp) as age,
  pod.pod_ip,
  pod.node_name
from
  kubernetes_pod as pod,
  json_each(pod.owner_references) as pod_owner,
  kubernetes_replicaset as rs
where
  json_extract(pod_owner.value, '$.kind') = 'ReplicaSet'
  and rs.uid = json_extract(pod_owner.value, '$.uid')
  and rs.name = 'frontend-56fc5b6b47'
order by
  pod.namespace,
  rs.name,
  pod.name;
```

### List manifest resources
Analyze the status of replica sets in your Kubernetes environment to understand their readiness and availability. This can help in assessing the health and performance of your applications running on Kubernetes.

```sql+postgres
select
  name,
  namespace,
  replicas as desired,
  ready_replicas as ready,
  available_replicas as available,
  selector,
  fully_labeled_replicas,
  path
from
  kubernetes_replicaset
where
  path is not null;
```

```sql+sqlite
select
  name,
  namespace,
  replicas as desired,
  ready_replicas as ready,
  available_replicas as available,
  selector,
  fully_labeled_replicas,
  path
from
  kubernetes_replicaset
where
  path is not null;
```