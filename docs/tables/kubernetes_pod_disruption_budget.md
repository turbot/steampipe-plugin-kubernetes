---
title: "Steampipe Table: kubernetes_pod_disruption_budget - Query Kubernetes Pod Disruption Budgets using SQL"
description: "Allows users to query Kubernetes Pod Disruption Budgets, specifically providing information about the minimum available pods and selector details, offering insights into the disruption allowance of the pods."
---

# Table: kubernetes_pod_disruption_budget - Query Kubernetes Pod Disruption Budgets using SQL

Kubernetes Pod Disruption Budgets (PDB) is a feature that allows a Kubernetes user to specify the number of replicas that an application can tolerate having, relative to how many it is intended to have. It defines the minimum number of pods that an orchestrated app can have, without a voluntary disruption. PDB also provides a way to limit the disruptions of your application while the Kubernetes cluster manager balances the needs of your applications.

## Table Usage Guide

The `kubernetes_pod_disruption_budget` table provides insights into the Pod Disruption Budgets within Kubernetes. As a DevOps engineer, explore details through this table, including the minimum available pods, selector details, and associated metadata. Utilize it to uncover information about the disruption allowance of the pods, such as the minimum number of pods an application can have, and the details of the selectors.

## Examples

### Basic info
Explore the minimum and maximum availability of resources within your Kubernetes environment. This query helps in managing resource allocation and ensuring smooth operation by identifying potential disruption areas.

```sql+postgres
select
   name,
   namespace,
   min_available,
   max_unavailable,
   selector
from
   kubernetes_pod_disruption_budget
order by
   namespace,
   name;
```

```sql+sqlite
select
   name,
   namespace,
   min_available,
   max_unavailable,
   selector
from
   kubernetes_pod_disruption_budget
order by
   namespace,
   name;
```

### List deployments and their matching PDB
Analyze the settings to understand the relationship between different deployments and their corresponding Pod Disruption Budgets (PDB) in a Kubernetes environment. This could be useful to ensure that the deployments are properly configured to handle disruptions, thereby enhancing system resilience.

```sql+postgres
select
  d.namespace,
  d.name,
  min_available,
  replicas
from
  kubernetes_pod_disruption_budget pdb
  inner join
   kubernetes_deployment d
   on d.selector = pdb.selector
   and d.namespace = pdb.namespace
order by
  d.namespace,
  d.name;
```

```sql+sqlite
select
  d.namespace,
  d.name,
  min_available,
  replicas
from
  kubernetes_pod_disruption_budget as pdb
  join
   kubernetes_deployment as d
   on d.selector = pdb.selector
   and d.namespace = pdb.namespace
order by
  d.namespace,
  d.name;
```

### List manifest resources
Explore which Kubernetes pod disruption budgets are available, focusing on those with a specified path. This helps in managing the application availability during voluntary disruptions.

```sql+postgres
select
  name,
  namespace,
  min_available,
  max_unavailable,
  selector,
  path
from
   kubernetes_pod_disruption_budget
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
  min_available,
  max_unavailable,
  selector,
  path
from
   kubernetes_pod_disruption_budget
where
  path is not null
order by
   namespace,
   name;
```