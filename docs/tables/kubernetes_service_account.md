# Table: kubernetes_service_account

In Kubernetes, service accounts are used to provide an identity for pods. Pods that want to interact with the API server will authenticate with a particular service account. By default, applications will authenticate as the default service account in the namespace they are running in.

## Examples

### Basic Info - `kubectl get serviceaccounts --all-namespaces` columns

```sql
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

### List role bindings

```sql
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

### List cluster role bindings and rules

```sql
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

### List manifest resources

```sql
select
  name,
  namespace,
  jsonb_array_length(secrets) as secrets
from
  kubernetes_service_account
where
  manifest_file_path is not null
order by
  namespace,
  name;
```
