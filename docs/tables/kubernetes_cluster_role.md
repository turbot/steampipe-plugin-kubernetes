---
title: "Steampipe Table: kubernetes_cluster_role - Query Kubernetes Cluster Roles using SQL"
description: "Allows users to query Kubernetes Cluster Roles, specifically the permissions and rules associated with each role, providing insights into the access controls in a Kubernetes cluster."
folder: "Cluster"
---

# Table: kubernetes_cluster_role - Query Kubernetes Cluster Roles using SQL

Kubernetes Cluster Roles is a feature within Kubernetes that allows you to define permissions for resources across the entire cluster. It provides a way to set up and manage access controls for various Kubernetes resources, including pods, services, and more. Kubernetes Cluster Roles helps you stay informed about the access controls in your Kubernetes cluster and take appropriate actions when managing permissions.

ClusterRole contains rules that represent a set of permissions. You can use them to grant access to:

- cluster-scoped resources (like nodes)
- non-resource endpoints (like /healthz)
- namespaced resources (like Pods), across all namespaces

## Table Usage Guide

The `kubernetes_cluster_role` table provides insights into Cluster Roles within Kubernetes. As a DevOps engineer, explore role-specific details through this table, including permissions and associated rules. Utilize it to uncover information about roles, such as those with wildcard permissions, and the verification of access controls.

## Examples

### Basic Info
Explore the creation timelines and rules associated with different roles in a Kubernetes cluster to better manage and organize your resources. This allows for effective prioritization and allocation of tasks within the cluster.

```sql+postgres
select
  name,
  creation_timestamp,
  rules,
  aggregation_rule
from
  kubernetes_cluster_role
order by
  name;
```

```sql+sqlite
select
  name,
  creation_timestamp,
  rules,
  aggregation_rule
from
  kubernetes_cluster_role
order by
  name;
```

### List rules and verbs for cluster roles
Explore the permissions associated with different roles in a Kubernetes cluster. This can help in understanding the access levels of various roles, thereby aiding in managing security and access control.

```sql+postgres
select
  name as role_name,
  rule ->> 'apiGroups' as api_groups,
  rule ->> 'resources' as resources,
  rule ->> 'nonResourceURLs' as non_resource_urls,
  rule ->> 'verbs' as verbs,
  rule ->> 'resourceNames' as resource_names
from
  kubernetes_cluster_role,
  jsonb_array_elements(rules) as rule
order by
  role_name,
  api_groups;
```

```sql+sqlite
select
  name as role_name,
  json_extract(rule.value, '$.apiGroups') as api_groups,
  json_extract(rule.value, '$.resources') as resources,
  json_extract(rule.value, '$.nonResourceURLs') as non_resource_urls,
  json_extract(rule.value, '$.verbs') as verbs,
  json_extract(rule.value, '$.resourceNames') as resource_names
from
  kubernetes_cluster_role,
  json_each(rules) as rule
order by
  role_name,
  api_groups;
```

### Group cluster roles by same set of aggregation rules
Discover the segments that share the same set of aggregation rules within cluster roles, which can help in assessing similar roles and their configuration in a Kubernetes environment. This may be useful in optimizing role assignment and ensuring consistent rule enforcement across roles.

```sql+postgres
select
  jsonb_agg(name) as roles,
  aggregation_rule
from
  kubernetes_cluster_role
group by
  aggregation_rule;
```

```sql+sqlite
select
  json_group_array(name) as roles,
  aggregation_rule
from
  kubernetes_cluster_role
group by
  aggregation_rule;
```

### List manifest resources
Explore the rules and aggregation rules of your Kubernetes cluster roles that have a defined path. This can help you understand the permissions structure and assess any potential security risks in your cluster.

```sql+postgres
select
  name,
  rules,
  aggregation_rule,
  path
from
  kubernetes_cluster_role
where
  path is not null
order by
  name;
```

```sql+sqlite
select
  name,
  rules,
  aggregation_rule,
  path
from
  kubernetes_cluster_role
where
  path is not null
order by
  name;
```