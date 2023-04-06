# Table: kubernetes_cluster_role

ClusterRole contains rules that represent a set of permissions. You can use them to grant access to:

- cluster-scoped resources (like nodes)
- non-resource endpoints (like /healthz)
- namespaced resources (like Pods), across all namespaces

## Examples

### Basic Info

```sql
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

```sql
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

### Group cluster roles by same set of aggregation rules

```sql
select
  jsonb_agg(name) as roles,
  aggregation_rule
from
  kubernetes_cluster_role
group by
  aggregation_rule;
```

### List manifest resources

```sql
select
  name,
  rules,
  aggregation_rule,
  manifest_file_path
from
  kubernetes_cluster_role
where
  manifest_file_path is not null
order by
  name;
```
