# Table: kubernetes_replication_controller

A ReplicationController ensures that a specified number of pod replicas are running at any one time.

## Examples

### Basic Info

```sql
select
  name,
  namespace,
  replicas as desired,
  status_replicas as current,
  ready_replicas as ready,
  age(current_timestamp, creation_timestamp),
  selector
from
  kubernetes_replication_controller;
```

### Get details of containers and image

```sql
select
  name,
  namespace,
  jsonb_agg(container.value -> 'name') as containers,
  jsonb_agg(container.value -> 'image') as images
from
  kubernetes_replication_controller,
  jsonb_array_elements(template -> 'spec' -> 'containers') as container
group by 
  name,
  namespace;
```
