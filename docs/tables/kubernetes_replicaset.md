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

### Get container and image used in the replicaset

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

### List pods for a replicaset (by name)

```sql
select
  pod.namespace,
  rs.name as replicaset_name,
  pod.name as pod_name,
  pod.phase,
  age(current_timestamp, pod.creation_timestamp),
  pod.pod_ip,
  pod.node_name
from
  kubernetes_pod as pod,
  jsonb_array_elements(pod.owner_references) as pod_owner,
  kubernetes_replicaset as rs
where
  pod_owner ->> 'kind' = 'ReplicaSet'
  and rs.uid = pod_owner ->> 'uid'
  and rs.name = 'frontend-56fc5b6b47'
order by
  pod.namespace,
  rs.name,
  pod.name;
```

### List manifest resources

```sql
select
  name,
  namespace,
  replicas as desired,
  ready_replicas as ready,
  available_replicas as available,
  selector,
  fully_labeled_replicas,
  path
from
  kubernetes_replicaset
where
  path is not null;
```
