---
title: "Steampipe Table: kubernetes_stateful_set - Query Kubernetes Stateful Sets using SQL"
description: "Allows users to query Kubernetes Stateful Sets, specifically providing details about the stateful applications running in a Kubernetes environment."
---

# Table: kubernetes_stateful_set - Query Kubernetes Stateful Sets using SQL

A Kubernetes Stateful Set is a workload API object that manages stateful applications. It is used to manage applications which require one or more of the following: stable, unique network identifiers, stable, persistent storage, and ordered, graceful deployment and scaling. Stateful Sets are valuable for applications that require stable network identity or stable storage, like databases.

## Table Usage Guide

The `kubernetes_stateful_set` table provides insights into the stateful applications running in a Kubernetes environment. As a DevOps engineer, explore details of these applications through this table, including network identifiers, persistent storage, and deployment details. Utilize it to manage and monitor stateful applications, such as databases, that require stable network identity or persistent storage.

## Examples

### Basic Info - `kubectl get statefulsets --all-namespaces` columns
Explore the organization and status of your Kubernetes stateful sets by identifying their names, associated services, and the number of replicas. This query also allows you to assess the age of these sets, helping you manage system resources and plan for updates or decommissioning.

```sql+postgres
select
  name,
  namespace,
  service_name,
  replicas,
  age(current_timestamp, creation_timestamp)
from
  kubernetes_stateful_set
order by
  namespace,
  name;
```

```sql+sqlite
select
  name,
  namespace,
  service_name,
  replicas,
  strftime('%s', 'now') - strftime('%s', creation_timestamp) as age
from
  kubernetes_stateful_set
order by
  namespace,
  name;
```

### List stateful sets that require manual update when the object's configuration is changed
Explore which stateful sets in your Kubernetes environment require manual updates whenever there are changes in the object's configuration. This is useful for ensuring optimal management and timely updates of stateful sets, particularly those with an 'OnDelete' update strategy.

```sql+postgres
select
  name,
  namespace,
  service_name,
  update_strategy ->> 'type' as update_strategy_type
from
  kubernetes_stateful_set
where
  update_strategy ->> 'type' = 'OnDelete';
```

```sql+sqlite
select
  name,
  namespace,
  service_name,
  json_extract(update_strategy, '$.type') as update_strategy_type
from
  kubernetes_stateful_set
where
  json_extract(update_strategy, '$.type') = 'OnDelete';
```

### List manifest resources
Explore which stateful applications in your Kubernetes cluster have specified storage configurations. This can help you understand how your persistent data is managed and identify any potential issues with data persistence.

```sql+postgres
select
  name,
  namespace,
  service_name,
  replicas,
  path
from
  kubernetes_stateful_set
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
  service_name,
  replicas,
  path
from
  kubernetes_stateful_set
where
  path is not null
order by
  namespace,
  name;
```