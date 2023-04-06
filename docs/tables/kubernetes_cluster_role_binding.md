# Table: kubernetes_cluster_role_binding

A ClusterRoleBinding grants the permissions defined in a cluster role to a user or set of users. Access granted by ClusterRoleBinding is cluster-wide.

## Examples

### Basic Info

```sql
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

### Get details subject and role details for bindings

```sql
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

### Get cluster role bindings associated for each role

```sql
select
  role_name,
  jsonb_agg(name) as bindings
from
  kubernetes_cluster_role_binding
group by
  role_name;
```

### List manifest resources

```sql
select
  name,
  role_name,
  role_kind,
  jsonb_pretty(subjects) as subjects,
  manifest_file_path
from
  kubernetes_cluster_role_binding
where
  manifest_file_path is not null
order by
  name;
```
