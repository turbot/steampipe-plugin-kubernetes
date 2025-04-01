---
title: "Steampipe Table: kubernetes_storage_class - Query Kubernetes Storage Classes using SQL"
description: "Allows users to query Storage Classes in Kubernetes, providing detailed insights into the different types of storage available in a Kubernetes cluster."
folder: "Storage"
---

# Table: kubernetes_storage_class - Query Kubernetes Storage Classes using SQL

A Storage Class in Kubernetes is a way to describe different types of storage that are available in a Kubernetes cluster. Storage Classes are used to dynamically provision storage, based on the class of storage requested by a Persistent Volume Claim. They are essential for managing storage resources and ensuring that the right type of storage is available for different workloads.

## Table Usage Guide

The `kubernetes_storage_class` table provides insights into the Storage Classes within a Kubernetes cluster. As a Kubernetes administrator, you can explore details about each Storage Class through this table, including the provisioner, reclaim policy, and volume binding mode. Utilize it to manage and optimize storage resources in your cluster, ensuring that the right type of storage is available for different workloads.

## Examples

### Basic Info
Explore which storage classes are available in your Kubernetes environment, including their provisioners, reclaim policies, and mount options. This can be useful to understand how your storage resources are configured and managed.

```sql+postgres
select
  name,
  provisioner,
  reclaim_policy,
  parameters,
  mount_options
from
  kubernetes_storage_class;
```

```sql+sqlite
select
  name,
  provisioner,
  reclaim_policy,
  parameters,
  mount_options
from
  kubernetes_storage_class;
```

### List storage classes that don't allow volume expansion
Explore which storage classes in a Kubernetes environment do not support volume expansion. This is useful for identifying potential storage limitations in your system.

```sql+postgres
select
  name,
  allow_volume_expansion,
  provisioner,
  reclaim_policy,
  parameters,
  mount_options
from
  kubernetes_storage_class
where
  not allow_volume_expansion;
```

```sql+sqlite
select
  name,
  allow_volume_expansion,
  provisioner,
  reclaim_policy,
  parameters,
  mount_options
from
  kubernetes_storage_class
where
  allow_volume_expansion = 0;
```

### List storage classes with immediate volume binding mode enabled
Explore which storage classes have immediate volume binding mode enabled. This is beneficial for understanding the storage configurations that allow immediate access to volumes, which can be crucial for certain applications and workloads.

```sql+postgres
select
  name,
  allow_volume_expansion,
  provisioner,
  reclaim_policy,
  volume_binding_mode
from
  kubernetes_storage_class
where
  volume_binding_mode = 'Immediate';
```

```sql+sqlite
select
  name,
  allow_volume_expansion,
  provisioner,
  reclaim_policy,
  volume_binding_mode
from
  kubernetes_storage_class
where
  volume_binding_mode = 'Immediate';
```

### List manifest resources
Explore the configuration of storage classes in a Kubernetes environment to understand which ones have defined paths for storing data. This is useful in identifying potential storage optimization opportunities or troubleshooting storage-related issues.

```sql+postgres
select
  name,
  provisioner,
  reclaim_policy,
  parameters,
  mount_options,
  path
from
  kubernetes_storage_class
where
  path is not null;
```

```sql+sqlite
select
  name,
  provisioner,
  reclaim_policy,
  parameters,
  mount_options,
  path
from
  kubernetes_storage_class
where
  path is not null;
```