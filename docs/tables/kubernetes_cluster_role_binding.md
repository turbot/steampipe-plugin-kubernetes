---
title: "Steampipe Table: kubernetes_cluster_role_binding - Query Kubernetes Cluster Role Bindings using SQL"
description: "Allows users to query Kubernetes Cluster Role Bindings, providing insights into binded roles within a Kubernetes cluster."
---

# Table: kubernetes_cluster_role_binding - Query Kubernetes Cluster Role Bindings using SQL

A Kubernetes Cluster Role Binding binds a role to subjects. Subjects can be groups, users, or service accounts. Cluster Role Binding grants permissions to users at the cluster level, which includes all namespaces.

## Table Usage Guide

The `kubernetes_cluster_role_binding` table provides insights into Cluster Role Bindings within Kubernetes. As a DevOps engineer, explore binding-specific details through this table, including subjects, role references, and associated metadata. Utilize it to uncover information about bindings, such as the roles binded to specific subjects, the namespaces of the roles, and the verification of role references.

## Examples

### Basic Info
Explore which roles are bound to different subjects in your Kubernetes cluster. This allows you to gain insights into the permissions and access levels within your system, assisting in managing security and access control.

```sql+postgres
select
  name,
  role_name,
  role_kind,
  jsonb_pretty(subjects) as subjects,
  creation_timestamp
from
  kubernetes_cluster_role_binding
order by
  name;
```

```sql+sqlite
select
  name,
  role_name,
  role_kind,
  subjects,
  creation_timestamp
from
  kubernetes_cluster_role_binding
order by
  name;
```

### Get details subject and role details for bindings
Uncover the details of role bindings within a Kubernetes cluster. This query is particularly useful for administrators who want to keep track of roles and associated subjects in a systematic order, thereby facilitating efficient management of access and permissions within the cluster.

```sql+postgres
select
  name as binding_name,
  role_name,
  subject ->> 'name' as subject_name,
  subject ->> 'namespace' as subject_namespace,
  subject ->> 'apiGroup' as subject_api_group,
  subject ->> 'kind' as subject_kind
from
  kubernetes_cluster_role_binding,
  jsonb_array_elements(subjects) as subject
order by
  role_name,
  subject_name;
```

```sql+sqlite
select
  name as binding_name,
  role_name,
  json_extract(subject.value, '$.name') as subject_name,
  json_extract(subject.value, '$.namespace') as subject_namespace,
  json_extract(subject.value, '$.apiGroup') as subject_api_group,
  json_extract(subject.value, '$.kind') as subject_kind
from
  kubernetes_cluster_role_binding,
  json_each(subjects) as subject
order by
  role_name,
  subject_name;
```

### Get cluster role bindings associated for each role
Discover the segments that have specific role bindings in a Kubernetes cluster, which can aid in understanding the distribution and assignment of roles across the cluster. This can be particularly useful for managing permissions and access controls within the system.

```sql+postgres
select
  role_name,
  jsonb_agg(name) as bindings
from
  kubernetes_cluster_role_binding
group by
  role_name;
```

```sql+sqlite
select
  role_name,
  json_group_array(name) as bindings
from
  kubernetes_cluster_role_binding
group by
  role_name;
```

### List manifest resources
Assess the elements within your Kubernetes cluster to understand the relationship and permissions between different roles and resources. This allows you to maintain a secure and well-organized system by identifying any irregularities or potential vulnerabilities in role assignments.

```sql+postgres
select
  name,
  role_name,
  role_kind,
  jsonb_pretty(subjects) as subjects,
  path
from
  kubernetes_cluster_role_binding
where
  path is not null
order by
  name;
```

```sql+sqlite
select
  name,
  role_name,
  role_kind,
  subjects,
  path
from
  kubernetes_cluster_role_binding
where
  path is not null
order by
  name;
```