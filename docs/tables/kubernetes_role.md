---
title: "Steampipe Table: kubernetes_role - Query Kubernetes Roles using SQL"
description: "Allows users to query Roles in Kubernetes, specifically the permissions and privileges assigned to a Role, providing insights into access control and security configurations."
folder: "Role"
---

# Table: kubernetes_role - Query Kubernetes Roles using SQL

A Role in Kubernetes is a set of permissions that can be assigned to resources within a namespace. Roles dictate what actions are permitted and which resources those actions can be performed on. It is an integral part of Kubernetes' Role-Based Access Control (RBAC) system used to manage permissions and access to the Kubernetes API.

## Table Usage Guide

The `kubernetes_role` table provides insights into Roles within Kubernetes RBAC. As a DevOps engineer, explore role-specific details through this table, including permissions, associated resources, and the namespaces they belong to. Utilize it to uncover information about roles, such as their access privileges, the resources they can interact with, and the namespaces they are active in.

## Examples

### Basic Info
Explore the roles within your Kubernetes environment, including their creation time and associated rules, to gain insights into your system's configuration and organization. This can help you understand how roles are distributed and manage them effectively.

```sql+postgres
select
  name,
  namespace,
  creation_timestamp,
  rules
from
  kubernetes_role
order by
  name;
```

```sql+sqlite
select
  name,
  namespace,
  creation_timestamp,
  rules
from
  kubernetes_role
order by
  name;
```

### List rules and verbs for roles
Explore which roles have specific permissions in your Kubernetes environment. This query helps in understanding the distribution of access rights, assisting in access control and security management.

```sql+postgres
select
  name as role_name,
  rule ->> 'apiGroups' as api_groups,
  rule ->> 'resources' as resources,
  rule ->> 'nonResourceURLs' as non_resource_urls,
  rule ->> 'verbs' as verbs,
  rule ->> 'resourceNames' as resource_names
from
  kubernetes_role,
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
  kubernetes_role,
  json_each(rules) as rule
order by
  role_name,
  api_groups;
```

### List manifest resources
Explore which Kubernetes roles have a defined path to better organize and manage your resources. This can help in identifying instances where roles may be improperly configured or misplaced.

```sql+postgres
select
  name,
  namespace,
  rules,
  path
from
  kubernetes_role
where
  path is not null
order by
  name;
```

```sql+sqlite
select
  name,
  namespace,
  rules,
  path
from
  kubernetes_role
where
  path is not null
order by
  name;
```