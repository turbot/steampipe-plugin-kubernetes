# Table: kubernetes_cluster_role

ClusterRole contains rules that represent a set of permissions. You can use them to grant access to:

- cluster-scoped resources (like nodes)
- non-resource endpoints (like /healthz)
- namespaced resources (like Pods), across all namespaces

For example: you can use a ClusterRole to allow a particular user to run `kubectl get pods --all-namespaces`

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
  cr.name as role_name,
  cr_rule ->> 'apiGroups' as rule_api_groups,
  cr_rule ->> 'resources' as rule_resources,
  cr_rule ->> 'verbs' as rule_verbs,
  cr_rule ->> 'resourceNames' as rule_resource_names
from
  kubernetes_cluster_role as cr,
  jsonb_array_elements(cr.rules) as cr_rule
order by
  role_name,
  rule_api_groups
```
