---
title: "Steampipe Table: kubernetes_persistent_volume - Query Kubernetes Persistent Volumes using SQL"
description: "Allows users to query Kubernetes Persistent Volumes, providing insights into the storage resources available in a Kubernetes cluster."
folder: "Persistent Volume"
---

# Table: kubernetes_persistent_volume - Query Kubernetes Persistent Volumes using SQL

A Kubernetes Persistent Volume (PV) is a piece of storage in the cluster that has been provisioned by an administrator or dynamically provisioned using Storage Classes. PVs are volume plugins like Volumes, but have a lifecycle independent of any individual Pod that uses the PV. These resources allow Pods to store data that can survive the lifecycle of a Pod.

## Table Usage Guide

The `kubernetes_persistent_volume` table provides insights into persistent volumes within Kubernetes. As a DevOps engineer, explore volume-specific details through this table, including storage capacity, access modes, and associated metadata. Utilize it to uncover information about volumes, such as those with certain storage classes, the status of volumes, and the reclaim policy set for volumes.

## Examples

### Basic Info
Explore the status and capacity of your persistent storage volumes within your Kubernetes environment. This allows you to manage your storage resources effectively and plan for future capacity needs.

```sql+postgres
select
  name,
  access_modes,
  storage_class,
  capacity ->> 'storage' as storage_capacity,
  creation_timestamp,
  persistent_volume_reclaim_policy,
  phase as status,
  volume_mode,
  age(current_timestamp, creation_timestamp)
from
  kubernetes_persistent_volume;
```

```sql+sqlite
select
  name,
  access_modes,
  storage_class,
  json_extract(capacity, '$.storage') as storage_capacity,
  creation_timestamp,
  persistent_volume_reclaim_policy,
  phase as status,
  volume_mode,
  (julianday('now') - julianday(creation_timestamp)) * 24 * 60 * 60 as age
from
  kubernetes_persistent_volume;
```

### Get hostpath details for the volume
Explore the details of your persistent volume's hostpath in your Kubernetes setup. This can help in understanding the type and path associated with your volume, which is crucial for managing and troubleshooting your storage configuration.

```sql+postgres
select
  name,
  persistent_volume_source -> 'hostPath' ->> 'path' as path,
  persistent_volume_source -> 'hostPath' ->> 'type' as type
from
  kubernetes_persistent_volume;
```

```sql+sqlite
select
  name,
  json_extract(persistent_volume_source, '$.hostPath.path') as path,
  json_extract(persistent_volume_source, '$.hostPath.type') as type
from
  kubernetes_persistent_volume;
```

### List manifest resources
Explore the various resources within your Kubernetes persistent volumes, focusing on those that have a specified path. This allows you to assess storage capacities, access modes, and reclaim policies to better manage your Kubernetes environment.

```sql+postgres
select
  name,
  access_modes,
  storage_class,
  capacity ->> 'storage' as storage_capacity,
  persistent_volume_reclaim_policy,
  phase as status,
  volume_mode,
  path
from
  kubernetes_persistent_volume
where
  path is not null;
```

```sql+sqlite
select
  name,
  access_modes,
  storage_class,
  json_extract(capacity, '$.storage') as storage_capacity,
  persistent_volume_reclaim_policy,
  phase as status,
  volume_mode,
  path
from
  kubernetes_persistent_volume
where
  path is not null;
```