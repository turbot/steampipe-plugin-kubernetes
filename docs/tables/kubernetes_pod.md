# Table: kubernetes_pod

Pods are the smallest deployable units of computing that you can create and manage in Kubernetes.

A Pod is a group of one or more containers, with shared storage and network resources, and a specification for how to run the containers.

## Examples

### Basic Info

```sql
select
  namespace,
  name,
  phase,
  age(current_timestamp, creation_timestamp),
  pod_ip,
  node_name,
  jsonb_array_length(containers) as container_count,
  jsonb_array_length(init_containers) as init_container_count,
  jsonb_array_length(ephemeral_containers) as ephemeral_container_count
from
  kubernetes_pod
order by
  namespace,
  name;
```

### List Unowned (Naked) Pods

```sql
select
  name,
  namespace,
  phase,
  pod_ip,
  node_name
from
  kubernetes_pod
where
  owner_references is null;
```

### List Privileged Containers

```sql
select
  name as pod_name,
  namespace,
  phase,
  jsonb_pretty(owner_references) as owners,
  c ->> 'name' as container_name,
  c ->> 'image' as container_image
from
  kubernetes_pod,
  jsonb_array_elements(containers) as c
where
  c -> 'securityContext' ->> 'privileged' = 'true';
```

### List Pods with access to the host process ID, IPC, or network namespace

```sql
select
  name,
  namespace,
  phase,
  host_pid,
  host_ipc,
  host_network,
  jsonb_pretty(owner_references) as owners
from
  kubernetes_pod
where
  host_pid or host_ipc or host_network;
```

### Container Statuses

```sql
select
  namespace,
  name as pod_name,
  phase,
  cs ->> 'name' as container_name,
  cs ->> 'image' as image,
  cs ->> 'ready' as ready,
  cs_state as state,
  cs ->> 'started' as started,
  cs ->> 'restartCount' as restarts
from
  kubernetes_pod,
  jsonb_array_elements(container_statuses) as cs,
  jsonb_object_keys(cs -> 'state') as cs_state
order by
  namespace,
  name,
  container_name;
```

### `kubectl get pods` columns

```sql
select
  namespace,
  name,
  phase,
  count(cs) filter (
    where
      cs -> 'state' -> 'running' is not null
  ) as running_container_count,
  jsonb_array_length(containers) as container_count,
  age(current_timestamp, creation_timestamp),
  COALESCE(sum((cs ->> 'restartCount') :: int), 0) as restarts,
  pod_ip,
  node_name
from
  kubernetes_pod
  left join jsonb_array_elements(container_statuses) as cs on true
group by
  namespace,
  name,
  phase,
  containers,
  creation_timestamp,
  pod_ip,
  node_name
 order by
  namespace,
  name;
```

### List manifest resources

```sql
select
  namespace,
  name,
  phase,
  pod_ip,
  node_name,
  jsonb_array_length(containers) as container_count,
  jsonb_array_length(init_containers) as init_container_count,
  jsonb_array_length(ephemeral_containers) as ephemeral_container_count,
  path
from
  kubernetes_pod
where
  path is not null
order by
  namespace,
  name;
```
