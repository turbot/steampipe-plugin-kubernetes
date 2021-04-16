# Table: kubernetes_node

A Node is a worker machine in Kubernetes and may be either a virtual or a physical machine, depending on the cluster.

## Examples

### Basic Info

```sql
select
  name,
  pod_cidr,
  pod_cidrs,
  provider_id,
  creation_timestamp,
  addresses,
  capacity
from
  kubernetes_node;
```

### List conditions for node

```sql
select
  name,
  cond ->> 'type' as type,
  lower(cond ->> 'status')::bool as status,
  (cond ->> 'lastHeartbeatTime')::timestamp as last_heartbeat_time,
  (cond ->> 'lastTransitionTime')::timestamp as last_transition_time,
  cond ->> 'reason' as reason,
  cond ->> 'message' as message
from
  kubernetes_node,
  jsonb_array_elements(conditions) as cond
```

### Get system info for nodes

```sql
select
  name,
  node_info ->> 'machineID' as machine_id,
  node_info ->> 'systemUUID' as systemUUID,
  node_info ->> 'bootID' as bootID,
  node_info ->> 'kernelVersion' as kernelVersion,
  node_info ->> 'osImage' as osImage,
  node_info ->> 'operatingSystem' as operatingSystem,
  node_info ->> 'architecture' as architecture,
  node_info ->> 'containerRuntimeVersion' as containerRuntimeVersion,
  node_info ->> 'kubeletVersion' as kubeletVersion,
  node_info ->> 'kubeProxyVersion' as kubeProxyVersion
from
  kubernetes_node
```
