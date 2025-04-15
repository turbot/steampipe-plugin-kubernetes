---
title: "Steampipe Table: kubernetes_persistent_volume_claim - Query Kubernetes Persistent Volume Claims using SQL"
description: "Allows users to query Kubernetes Persistent Volume Claims, specifically providing information about the status, capacity, and access modes of each claim."
folder: "Persistent Volume"
---

# Table: kubernetes_persistent_volume_claim - Query Kubernetes Persistent Volume Claims using SQL

A Kubernetes Persistent Volume Claim (PVC) is a request for storage by a user. It is similar to a pod in Kubernetes. PVCs can request specific size and access modes like read and write for a Persistent Volume (PV).

## Table Usage Guide

The `kubernetes_persistent_volume_claim` table provides insights into the Persistent Volume Claims within a Kubernetes cluster. As a DevOps engineer, you can use this table to explore details about each claim, including its current status, requested storage capacity, and access modes. This table is beneficial when you need to manage storage resources or troubleshoot storage-related issues in your Kubernetes environment.

## Examples

### Basic Info
Explore the status and capacity of persistent storage volumes in a Kubernetes environment. This can help you manage resources effectively and ensure optimal allocation and usage.

```sql+postgres
select
  name,
  namespace,
  volume_name as volume,
  volume_mode,
  access_modes,
  phase as status,
  capacity ->> 'storage' as capacity,
  creation_timestamp,
  data_source,
  selector,
  resources,
  age(current_timestamp, creation_timestamp)
from
  kubernetes_persistent_volume_claim;
```

```sql+sqlite
select
  name,
  namespace,
  volume_name as volume,
  volume_mode,
  access_modes,
  phase as status,
  json_extract(capacity, '$.storage') as capacity,
  creation_timestamp,
  data_source,
  selector,
  resources,
  (julianday('now') - julianday(creation_timestamp)) * 24 * 60 * 60 as age
from
  kubernetes_persistent_volume_claim;
```

### List manifest resources
Explore the various resources within a manifest by identifying their names, namespaces, and statuses. This is useful for understanding the capacity and configuration of your persistent storage volumes, particularly when you need to assess the availability and allocation of resources.

```sql+postgres
select
  name,
  namespace,
  volume_name as volume,
  volume_mode,
  access_modes,
  phase as status,
  capacity ->> 'storage' as capacity,
  data_source,
  selector,
  resources,
  path
from
  kubernetes_persistent_volume_claim
where
  path is not null;
```

```sql+sqlite
select
  name,
  namespace,
  volume_name as volume,
  volume_mode,
  access_modes,
  phase as status,
  json_extract(capacity, '$.storage') as capacity,
  data_source,
  selector,
  resources,
  path
from
  kubernetes_persistent_volume_claim
where
  path is not null;
```