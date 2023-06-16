# Table: kubernetes_role

Role contains rules that represent a set of permissions. A Role always sets permissions within a particular namespace.

## Examples

### Basic Info

```sql
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

```sql
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

### List manifest resources

```sql
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
