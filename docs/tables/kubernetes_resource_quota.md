# Table: kubernetes_resource_quota

A resource quota, defined by a ResourceQuota object, provides constraints that limit aggregate resource consumption per namespace. It can limit the quantity of objects that can be created in a namespace by type, as well as the total amount of compute resources that may be consumed by resources in that namespace.

## Examples

### Basic Info

```sql
select
  name,
  namespace,
  resource_version,
  creation_timestamp,
  jsonb_pretty(spec_hard) as spec_hard
from
  kubernetes_resource_quota
order by
  name;
```

### Get used pod details of namespaces

```sql
select
  name,
  namespace,
  status_used -> 'pods' as used_pods,
  status_used -> 'services' as used_services
from
  kubernetes_resource_quota;
```

### List manifest resources

```sql
select
  name,
  namespace,
  resource_version,
  jsonb_pretty(spec_hard) as spec_hard,
  path
from
  kubernetes_resource_quota
where
  path is not null
order by
  name;
```
