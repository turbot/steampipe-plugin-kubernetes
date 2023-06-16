# Table: kubernetes_role_binding

A role binding grants the permissions defined in a role to a user or set of users. It holds a list of subjects (users, groups, or service accounts), and a reference to the role being granted. A RoleBinding grants permissions within a specific namespace.

## Examples

### Basic Info

```sql
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

### Get details subject and role details for bindings

```sql
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

### Get role bindings for each role

```sql
select
  role_name,
  jsonb_agg(name) as bindings
from
  kubernetes_role_binding
group by
  role_name;
```

### List manifest resources

```sql
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
