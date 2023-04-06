# Table: kubernetes_deployment

A Deployment provides declarative updates for Pods and ReplicaSets.

You describe a desired state in a Deployment, and the Deployment Controller changes the actual state to the desired state at a controlled rate. You can define Deployments to create new ReplicaSets, or to remove existing Deployments and adopt all their resources with new Deployments.

## Examples

### Basic Info

```sql
select
  name,
  namespace,
  status_replicas,
  ready_replicas,
  updated_replicas,
  available_replicas,
  unavailable_replicas,
  age(current_timestamp, creation_timestamp)
from
  kubernetes_deployment
order by
  namespace,
  name;
```

### Configuration Info

```sql
select
  name,
  paused,
  generate_name,
  generation,
  revision_history_limit,
  replicas,
  strategy,
  selector
from
  kubernetes_deployment;
```

### Container Images used in Deployments

```sql
select 
  name as deployment_name,
  namespace,
  c ->> 'name' as container_name,
  c ->> 'image' as image
from 
  kubernetes_deployment,
  jsonb_array_elements(template -> 'spec' -> 'containers') as c
order by
  namespace,
  name;
```

### List pods for a deployment

```sql
select
  pod.namespace,
  d.name as deployment_name,
  rs.name as replicaset_name,
  pod.name as pod_name,
  pod.phase,
  age(current_timestamp, pod.creation_timestamp),
  pod.pod_ip,
  pod.node_name
from 
  kubernetes_pod as pod,
  jsonb_array_elements(pod.owner_references) as pod_owner,
  kubernetes_replicaset as rs,
  jsonb_array_elements(rs.owner_references) as rs_owner,
  kubernetes_deployment as d
where 
  pod_owner ->> 'kind' = 'ReplicaSet'
  and rs.uid = pod_owner ->> 'uid'
  and rs_owner ->> 'uid' = d.uid 
  and d.name = 'frontend'
order by
  pod.namespace,
  d.name,
  rs.name,
  pod.name;
```

### List Pods with access to the to the host process ID, IPC, or network namespace

```sql
select 
  namespace,
  name,
  template -> 'spec' ->> 'hostPID' as hostPID,
  template -> 'spec' ->> 'hostIPC' as hostIPC,
  template -> 'spec' ->> 'hostNetwork' as hostNetwork
from 
  kubernetes_deployment
where
  template -> 'spec' ->> 'hostPID' = 'true' or
  template -> 'spec' ->> 'hostIPC' = 'true' or
  template -> 'spec' ->> 'hostNetwork' = 'true';
```

### List manifest resources

```sql
select
  name,
  namespace,
  replicas
from
  kubernetes_deployment
where
  manifest_file_path is not null
order by
  namespace,
  name;
```
