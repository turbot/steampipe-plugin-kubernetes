---
title: "Steampipe Table: kubernetes_daemonset - Query Kubernetes DaemonSets using SQL"
description: "Allows users to query Kubernetes DaemonSets, specifically to retrieve data about each DaemonSet's status, spec, and metadata."
folder: "DaemonSet"
---

# Table: kubernetes_daemonset - Query Kubernetes DaemonSets using SQL

A Kubernetes DaemonSet ensures that all (or some) nodes run a copy of a pod. This is used to run system-level applications, such as log collectors, monitoring agents, and more. DaemonSets are crucial for maintaining the desired state and ensuring the smooth operation of Kubernetes clusters.

Some typical uses of a DaemonSet are:

- running a cluster storage daemon on every node
- running a logs collection daemon on every node
- running a node monitoring daemon on every node

## Table Usage Guide

The `kubernetes_daemonset` table provides insights into DaemonSets within Kubernetes. As a DevOps engineer, explore DaemonSet-specific details through this table, including the current status, spec details, and associated metadata. Utilize it to uncover information about DaemonSets, such as the number of desired and current scheduled pods, the DaemonSet's labels, and the node selector terms.

## Examples

### Basic Info
Explore which Kubernetes daemonsets are currently scheduled and ready, and determine how long they have been running. This information can be used to assess the status and performance of your Kubernetes environment.

```sql+postgres
select
  name,
  namespace,
  desired_number_scheduled as desired,
  current_number_scheduled as current,
  number_ready as ready,
  number_available as available,
  selector,
  age(current_timestamp, creation_timestamp)
from
  kubernetes_daemonset;
```

```sql+sqlite
select
  name,
  namespace,
  desired_number_scheduled as desired,
  current_number_scheduled as current,
  number_ready as ready,
  number_available as available,
  selector,
  strftime('%s', 'now') - strftime('%s', creation_timestamp) as age
from
  kubernetes_daemonset;
```

### Get container and image used in the daemonset
Explore the relationship between container names and images used within a daemonset. This can be helpful in understanding how resources are being utilized and managed across different namespaces.

```sql+postgres
select
  name,
  namespace,
  c ->> 'name' as container_name,
  c ->> 'image' as image
from
  kubernetes_daemonset,
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
  kubernetes_daemonset,
  json_each(json_extract(template, '$.spec.containers')) as c
order by
  namespace,
  name;
```

### Get update strategy for the daemonset
Analyze the update strategy settings for daemonsets to understand the maximum number of unavailable updates and their types. This is beneficial in managing and planning updates without disrupting the functioning of the system.

```sql+postgres
select
  namespace,
  name,
  update_strategy -> 'maxUnavailable' as max_unavailable,
  update_strategy -> 'type' as type
from
  kubernetes_daemonset;
```

```sql+sqlite
select
  namespace,
  name,
  json_extract(update_strategy, '$.maxUnavailable') as max_unavailable,
  json_extract(update_strategy, '$.type') as type
from
  kubernetes_daemonset;
```

### List manifest resources
Explore the status of various resources in your Kubernetes Daemonset to understand if resource allocation aligns with your current needs. This can help assess if resources are being efficiently utilized or if adjustments are needed.

```sql+postgres
select
  name,
  namespace,
  desired_number_scheduled as desired,
  current_number_scheduled as current,
  number_available as available,
  selector,
  path
from
  kubernetes_daemonset
where
  path is not null;
```

```sql+sqlite
select
  name,
  namespace,
  desired_number_scheduled as desired,
  current_number_scheduled as current,
  number_available as available,
  selector,
  path
from
  kubernetes_daemonset
where
  path is not null;
```