# Table: kubernetes_pod_security_policy

Pod Security Policies enable fine-grained authorization of pod creation and updates. A Pod Security Policy is a cluster-level resource that controls security sensitive aspects of the pod specification

## Examples

### Basic Info

```sql
select
  name,
  allow_privilege_escalation,
  default_allow_privilege_escalation,
  host_network,
  host_ports,
  host_pid,
  host_ipc,
  privileged,
  read_only_root_filesystem,
  allowed_csi_drivers,
  allowed_host_paths
from
  kubernetes_pod_security_policy
order by
  name;
```

### List policies which allows access to the host process ID, IPC, or network namespace

```sql
select
  name,
  host_pid,
  host_ipc,
  host_network
from
  kubernetes_pod_security_policy
where
  host_pid or host_ipc or host_network;
```

### List policies which allows a pod to be privileged

```sql
select
  name
from
  kubernetes_pod_security_policy
where
  privileged;
```

### List manifest resources

```sql
select
  name,
  allow_privilege_escalation,
  default_allow_privilege_escalation,
  host_network,
  host_ports,
  host_pid,
  host_ipc,
  privileged,
  read_only_root_filesystem,
  allowed_csi_drivers,
  allowed_host_paths,
  path
from
  kubernetes_pod_security_policy
where
  path is not null
order by
  name;
```
