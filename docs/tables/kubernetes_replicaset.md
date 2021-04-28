# Table: kubernetes_replicaset

A ReplicaSet's purpose is to maintain a stable set of replica Pods running at any given time. As such, it is often used to guarantee the availability of a specified number of identical Pods.

## Examples

### Basic Info

```sql
select
  name,
  namespace,
  replicas as desired,
  ready_replicas as ready,
  available_replicas as available,
  selector,
  fully_labeled_replicas,
  age(current_timestamp, creation_timestamp)
from
  kubernetes_replicaset;
```

### Get conatiner and image used in the replicaset

```sql
select
  name,
  namespace,
  c ->> 'name' as container_name,
  c ->> 'image' as image
from
  kubernetes_replicaset,
  jsonb_array_elements(template -> 'spec' -> 'containers') as c
order by
  namespace,
  name;
```
