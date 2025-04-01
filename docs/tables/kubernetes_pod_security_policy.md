---
title: "Steampipe Table: kubernetes_pod_security_policy - Query Kubernetes Pod Security Policies using SQL"
description: "Allows users to query Kubernetes Pod Security Policies, providing details about the security configurations and restrictions that apply to pods in a Kubernetes cluster."
folder: "Pod"
---

# Table: kubernetes_pod_security_policy - Query Kubernetes Pod Security Policies using SQL

A Pod Security Policy is a cluster-level resource in Kubernetes that controls security-sensitive aspects of the pod specification. It establishes the default security settings for a pod, and can include settings such as the types of volumes that a pod can mount, the use of host networking and ports, and the execution of privileged operations. By defining these policies, administrators can enforce certain security standards across all pods within a cluster.

## Table Usage Guide

The `kubernetes_pod_security_policy` table provides insights into Pod Security Policies within a Kubernetes cluster. As a security engineer or Kubernetes administrator, explore policy-specific details through this table, including allowed and disallowed operations, volume types, and host networking configurations. Utilize it to uncover information about policies, such as those that allow privileged operations, the use of host networking, and the mounting of certain volume types.

## Examples

### Basic Info
Explore the security policies of your Kubernetes pods to understand how they're configured. This can help identify potential vulnerabilities and ensure that your system is robustly protected.

```sql+postgres
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

```sql+sqlite
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
Discover the segments that have policies allowing access to the host process ID, IPC, or network namespace. This is particularly useful in identifying potential security risks within your Kubernetes pod security policies.

```sql+postgres
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

```sql+sqlite
select
  name,
  host_pid,
  host_ipc,
  host_network
from
  kubernetes_pod_security_policy
where
  host_pid = 1 or host_ipc = 1 or host_network = 1;
```

### List policies which allows a pod to be privileged
Explore which pod security policies permit a pod to have privileged status. This can be useful for understanding potential security risks and ensuring compliance with best practice guidelines.

```sql+postgres
select
  name
from
  kubernetes_pod_security_policy
where
  privileged;
```

```sql+sqlite
select
  name
from
  kubernetes_pod_security_policy
where
  privileged = 1;
```

### List manifest resources
Explore which Kubernetes pod security policies have a defined path. This can be useful to identify potential security risks and ensure best practices are being followed.

```sql+postgres
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

```sql+sqlite
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