---
title: "Steampipe Table: kubernetes_role_binding - Query Kubernetes Role Bindings using SQL"
description: "Allows users to query Role Bindings in Kubernetes, providing insights into the permissions granted to specific users, groups, or service accounts within a namespace."
folder: "Role"
---

# Table: kubernetes_role_binding - Query Kubernetes Role Bindings using SQL

A Role Binding in Kubernetes is a link between a Role (or ClusterRole) and the subjects (users, groups, or service accounts) it applies to within a namespace. Role Bindings are used to grant the permissions defined in a Role to a user or set of users. They can reference a Role in the same namespace or a ClusterRole and then bind it to one or more subjects.

## Table Usage Guide

The `kubernetes_role_binding` table provides insights into Role Bindings within Kubernetes. As a Kubernetes administrator, explore role binding-specific details through this table, including the roles they reference, the subjects they apply to, and the namespace they belong to. Utilize it to uncover information about role bindings, such as those granting excessive permissions, the association between roles and subjects, and the verification of role permissions within a namespace.

## Examples

### Basic Info
Explore which roles are bound to specific subjects in your Kubernetes environment. This can help you gain insights into access controls and permissions, ensuring that only authorized entities have access to certain resources.

```sql+postgres
select
  name,
  namespace,
  role_name,
  role_kind,
  jsonb_pretty(subjects) as subjects,
  creation_timestamp
from
  kubernetes_role_binding
order by
  name;
```

```sql+sqlite
select
  name,
  namespace,
  role_name,
  role_kind,
  subjects,
  creation_timestamp
from
  kubernetes_role_binding
order by
  name;
```

### Get details subject and role details for bindings
Uncover the details of role bindings in a Kubernetes environment to understand the relationship between subjects and their associated roles. This can be useful in managing access control and ensuring proper role assignments.

```sql+postgres
select
  name as binding_name,
  namespace,
  role_name,
  subject ->> 'name' as subject_name,
  subject ->> 'namespace' as subject_namespace,
  subject ->> 'apiGroup' as subject_api_group,
  subject ->> 'kind' as subject_kind
from
  kubernetes_role_binding,
  jsonb_array_elements(subjects) as subject
order by
  subject_kind,
  role_name,
  subject_name;
```

```sql+sqlite
select
  name as binding_name,
  namespace,
  role_name,
  json_extract(subject.value, '$.name') as subject_name,
  json_extract(subject.value, '$.namespace') as subject_namespace,
  json_extract(subject.value, '$.apiGroup') as subject_api_group,
  json_extract(subject.value, '$.kind') as subject_kind
from
  kubernetes_role_binding,
  json_each(subjects) as subject
order by
  subject_kind,
  role_name,
  subject_name;
```

### Get role bindings for each role
Explore which roles have been assigned to each role binding in your Kubernetes configuration. This can help in managing access controls and ensuring appropriate permissions are assigned.

```sql+postgres
select
  role_name,
  jsonb_agg(name) as bindings
from
  kubernetes_role_binding
group by
  role_name;
```

```sql+sqlite
select
  role_name,
  json_group_array(name) as bindings
from
  kubernetes_role_binding
group by
  role_name;
```

### List manifest resources
Explore the role bindings within your Kubernetes environment to understand the relationships between different resources. This can be particularly useful in identifying potential security risks or misconfigurations.

```sql+postgres
select
  name,
  namespace,
  role_name,
  role_kind,
  jsonb_pretty(subjects) as subjects,
  path
from
  kubernetes_role_binding
where
  path is not null
order by
  name;
```

```sql+sqlite
select
  name,
  namespace,
  role_name,
  role_kind,
  subjects,
  path
from
  kubernetes_role_binding
where
  path is not null
order by
  name;
```