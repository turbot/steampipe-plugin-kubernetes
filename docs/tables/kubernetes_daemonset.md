# Table: kubernetes_daemonset

A DaemonSet ensures that all (or some) Nodes run a copy of a Pod.

Some typical uses of a DaemonSet are:

- running a cluster storage daemon on every node
- running a logs collection daemon on every node
- running a node monitoring daemon on every node

## Examples

### Basic Info

```sql
select
  name,
  namespace,
  desired_number_scheduled as desired,
  current_number_scheduled as current,
  number_ready as ready,
  number_available as available,
  selector,
  age(current_timestamp, creation_timestamp)
from
  kubernetes_daemonset;
```

### Get container and image used in the daemonset

```sql
select
  name,
  namespace,
  c ->> 'name' as container_name,
  c ->> 'image' as image
from
  kubernetes_daemonset,
  jsonb_array_elements(template -> 'spec' -> 'containers') as c
order by
  namespace,
  name;
```

### Get update strategy for the daemonset

```sql
select
  namespace,
  name,
  update_strategy -> 'maxUnavailable' as max_unavailable,
  update_strategy -> 'type' as type
from
  kubernetes_daemonset;
```

### List manifest resources

```sql
select
  name,
  namespace,
  desired_number_scheduled as desired,
  current_number_scheduled as current,
  number_available as available,
  selector,
  path
from
  kubernetes_daemonset
where
  path is not null;
```
