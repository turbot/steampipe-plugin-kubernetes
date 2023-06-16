# Table: kubernetes_ingress

Ingress exposes HTTP and HTTPS routes from outside the cluster to services within the cluster. Traffic routing is controlled by rules defined on the Ingress resource.

## Examples

### Basic Info

```sql
select
  name,
  namespace,
  ingress_class_name as class,
  age(current_timestamp, creation_timestamp)
from
  kubernetes_ingress
order by
  namespace,
  name;
```

### View rules for the ingress

```sql
select
  name,
  namespace,
  jsonb_pretty(rules) as rules
from
  kubernetes_ingress;
```

### List manifest resources

```sql
select
  name,
  namespace,
  ingress_class_name as class,
  path
from
  kubernetes_ingress
where
  path is not null
order by
  namespace,
  name;
```
