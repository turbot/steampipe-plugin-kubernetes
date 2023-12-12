---
title: "Steampipe Table: kubernetes_service_account - Query Kubernetes Service Accounts using SQL"
description: "Allows users to query Kubernetes Service Accounts, providing detailed information about each service account's metadata, secrets, image pull secrets, and automount service account token."
---

# Table: kubernetes_service_account - Query Kubernetes Service Accounts using SQL

Kubernetes Service Account is a type of Kubernetes resource that provides an identity for processes that run in a Pod. Service accounts are namespaced and can provide identity for applications running within a namespace. They are used to provide specific permissions to applications, allowing more granular control over system access.

## Table Usage Guide

The `kubernetes_service_account` table provides insights into service accounts within Kubernetes. As a DevOps engineer, explore service account-specific details through this table, including their metadata, secrets, image pull secrets, and automount service account token. Utilize it to uncover information about service accounts, such as their respective namespaces, the secrets they hold, and their automount settings.

## Examples

### Basic Info - `kubectl get serviceaccounts --all-namespaces` columns
Determine the areas in which Kubernetes service accounts are deployed and assess the number of secrets associated with each, while also gaining insights into their age. This query is useful for maintaining security and managing resource allocation within your Kubernetes environment.

```sql+postgres
select
  name,
  namespace,
  jsonb_array_length(secrets) as secrets,
  age(current_timestamp, creation_timestamp)
from
  kubernetes_service_account
order by
  namespace,
  name;
```

```sql+sqlite
select
  name,
  namespace,
  json_array_length(secrets) as secrets,
  strftime('%s', 'now') - strftime('%s', creation_timestamp) as age
from
  kubernetes_service_account
order by
  namespace,
  name;
```

### List role bindings
Explore the connections between service accounts and role bindings in a Kubernetes environment. This can help you understand the permissions and access levels of different service accounts, which is crucial for managing security and access control.

```sql+postgres
select
  sub ->> 'name' as service_account_name,
  sub ->> 'namespace' as service_account_namespace,
  name as role_binding,
  role_name,
  role_kind
from
  kubernetes_role_binding,
  jsonb_array_elements(subjects) as sub
where
  sub ->> 'kind' = 'ServiceAccount';
```

```sql+sqlite
select
  json_extract(sub.value, '$.name') as service_account_name,
  json_extract(sub.value, '$.namespace') as service_account_namespace,
  name as role_binding,
  role_name,
  role_kind
from
  kubernetes_role_binding,
  json_each(subjects) as sub
where
  json_extract(sub.value, '$.kind') = 'ServiceAccount';
```

### List cluster role bindings and rules
Explore the associations between cluster role bindings and their rules in a Kubernetes environment. This is useful for understanding the permissions and access levels of different service accounts, aiding in security and access management.

```sql+postgres
select
  crb.name as cluster_role_binding,
  crb.role_name,
  crb_sub ->> 'name' as service_account_name,
  crb_sub ->> 'namespace' as service_account_namespace,
  cr_rule ->> 'apiGroups' as rule_api_groups,
  cr_rule ->> 'resources' as rule_resources,
  cr_rule ->> 'verbs' as rule_verbs,
  cr_rule ->> 'resourceNames' as rule_resource_names
from
  kubernetes_cluster_role_binding as crb,
  jsonb_array_elements(subjects) as crb_sub,
  kubernetes_cluster_role as cr,
  jsonb_array_elements(cr.rules) as cr_rule
where
  cr.name = crb.role_name
  and crb_sub ->> 'kind' = 'ServiceAccount';
```

```sql+sqlite
select
  crb.name as cluster_role_binding,
  crb.role_name,
  json_extract(crb_sub.value, '$.name') as service_account_name,
  json_extract(crb_sub.value, '$.namespace') as service_account_namespace,
  json_extract(cr_rule.value, '$.apiGroups') as rule_api_groups,
  json_extract(cr_rule.value, '$.resources') as rule_resources,
  json_extract(cr_rule.value, '$.verbs') as rule_verbs,
  json_extract(cr_rule.value, '$.resourceNames') as rule_resource_names
from
  kubernetes_cluster_role_binding as crb,
  json_each(crb.subjects) as crb_sub,
  kubernetes_cluster_role as cr,
  json_each(cr.rules) as cr_rule
where
  cr.name = crb.role_name
  and json_extract(crb_sub.value, '$.kind') = 'ServiceAccount';
```

### List manifest resources
Discover the segments that have assigned secrets within the Kubernetes service accounts, allowing for a thorough review and management of these resources. This is useful for maintaining security and ensuring proper access controls are in place.

```sql+postgres
select
  name,
  namespace,
  jsonb_array_length(secrets) as secrets,
  path
from
  kubernetes_service_account
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
  json_array_length(secrets) as secrets,
  path
from
  kubernetes_service_account
where
  path is not null
order by
  namespace,
  name;
```