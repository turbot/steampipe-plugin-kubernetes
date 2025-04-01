---
title: "Steampipe Table: kubernetes_event - Query Kubernetes Events using SQL"
description: "Allows users to query Kubernetes Events, specifically the details of events occurring within a Kubernetes system, providing insights into system behaviors and potential issues."
folder: "Event"
---

# Table: kubernetes_event - Query Kubernetes Events using SQL

Kubernetes Events are objects that provide insight into what is happening inside a cluster, such as what decisions were made by scheduler or why some pods were evicted from the node. Events are a resource type in Kubernetes that are automatically created when certain situations occur, and they give developers a tool to understand the activity of the system.

## Table Usage Guide

The `kubernetes_event` table provides insights into events within a Kubernetes system. As a DevOps engineer or system administrator, explore event-specific details through this table, including the involved object, source, message, and related metadata. Utilize it to monitor system behaviors, troubleshoot issues, and understand the state changes in the workloads running on the cluster.

## Examples

### Basic Info
Explore the recent events in your Kubernetes environment to understand the status and health of your objects. This query can help you identify any issues or anomalies, providing valuable insights for troubleshooting and maintenance.

```sql+postgres
select
  namespace,
  last_timestamp,
  type,
  reason,
  concat(involved_object ->> 'kind', '/', involved_object ->> 'name') as object,
  message
from
  kubernetes_event;
```

```sql+sqlite
select
  namespace,
  last_timestamp,
  type,
  reason,
  involved_object || '/' || involved_object as object,
  message
from
  kubernetes_event;
```

### List warning events by last timestamp
Identify instances where warning events have occurred in your Kubernetes environment. This query is useful for tracking and understanding the chronology of these events to manage and troubleshoot issues effectively.

```sql+postgres
select
  namespace,
  last_timestamp,
  type,
  reason,
  concat(involved_object ->> 'kind', '/', involved_object ->> 'name') as object,
  message
from
  kubernetes_event
where
  type = 'Warning'
order by
  namespace,
  last_timestamp;
```

```sql+sqlite
select
  namespace,
  last_timestamp,
  type,
  reason,
  involved_object || '/' || involved_object as object,
  message
from
  kubernetes_event
where
  type = 'Warning'
order by
  namespace,
  last_timestamp;
```

### List manifest resources
Explore which Kubernetes events have a defined path to gain insights into the health and status of your Kubernetes resources. This can help identify any potential issues or anomalies within your system.

```sql+postgres
select
  namespace,
  type,
  reason,
  concat(involved_object ->> 'kind', '/', involved_object ->> 'name') as object,
  message,
  path
from
  kubernetes_event
where
  path is not null;
```

```sql+sqlite
select
  namespace,
  type,
  reason,
  involved_object || '/' || involved_object as object,
  message,
  path
from
  kubernetes_event
where
  path is not null;
```