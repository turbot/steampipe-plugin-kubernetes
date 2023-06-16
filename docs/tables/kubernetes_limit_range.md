# Table: kubernetes_limit_range

A LimitRange provides constraints that can:

- Enforce minimum and maximum compute resources usage per Pod or Container in a namespace.
- Enforce minimum and maximum storage request per PersistentVolumeClaim in a namespace.
- Enforce a ratio between request and limit for a resource in a namespace.
- Set default request/limit for compute resources in a namespace and automatically inject them to Containers at runtime.

## Examples

### Basic Info

```sql
select
  name,
  namespace,
  resource_version,
  creation_timestamp,
  jsonb_pretty(spec_limits) as spec_limits
from
  kubernetes_limit_range
order by
  namespace;
```

### Get spec limits details of limit range

```sql
select
  name,
  namespace,
  limits ->> 'type' as type,
  limits ->> 'default' as default,
  limits ->> 'defaultRequest' as default_request
from
  kubernetes_limit_range,
  jsonb_array_elements(spec_limits) as limits;
```

### List manifest resources

```sql
select
  name,
  namespace,
  resource_version,
  jsonb_pretty(spec_limits) as spec_limits,
  path
from
  kubernetes_limit_range
where
  path is not null
order by
  namespace;
```
