---
title: "Steampipe Table: kubernetes_node - Query Kubernetes Nodes using SQL"
description: "Allows users to query Nodes in Kubernetes, providing details about the node's status, capacity, and allocatable resources."
folder: "Node"
---

# Table: kubernetes_node - Query Kubernetes Nodes using SQL

Kubernetes Nodes are the worker machines that run containerized applications and other workloads. Nodes include the necessary services to run Pods (the smallest and simplest unit in the Kubernetes object model that you create or deploy) and are managed by the master components. The services on a node include the container runtime, kubelet and kube-proxy.

## Table Usage Guide

The `kubernetes_node` table provides insights into the Nodes within Kubernetes. As a DevOps engineer, explore node-specific details through this table, including its status, capacity, and allocatable resources. Utilize it to uncover information about nodes, such as the number of pods that can be scheduled for execution on the Node, the amount of CPU and memory resources available, and the overall status of the Node.

## Examples

### Basic Info
Assess the elements within your Kubernetes nodes to gain insights into their creation time, capacity, and associated addresses. This query is beneficial for understanding your nodes' infrastructure and helps in efficient resource allocation and management.

```sql+postgres
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

```sql+sqlite
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
This example demonstrates how to analyze the status and conditions of nodes within a Kubernetes system. It's useful for maintaining system health and troubleshooting, as it enables users to track node conditions over time and identify potential issues based on changes in status or the occurrence of specific conditions.

```sql+postgres
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
order by
  name,
  status desc;
```

```sql+sqlite
select
  name,
  json_extract(cond.value, '$.type') as type,
  lower(json_extract(cond.value, '$.status')) = 'true' as status,
  datetime(json_extract(cond.value, '$.lastHeartbeatTime')) as last_heartbeat_time,
  datetime(json_extract(cond.value, '$.lastTransitionTime')) as last_transition_time,
  json_extract(cond.value, '$.reason') as reason,
  json_extract(cond.value, '$.message') as message
from
  kubernetes_node,
  json_each(conditions) as cond
order by
  name,
  status desc;
```

### Get system info for nodes
Explore the system information for specific nodes to gain insights into machine specifications, operating system details, and various versions of installed software. This could be beneficial for troubleshooting, system audits, or understanding the overall system configuration.

```sql+postgres
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
  kubernetes_node;
```

```sql+sqlite
select
  name,
  json_extract(node_info, '$.machineID') as machine_id,
  json_extract(node_info, '$.systemUUID') as systemUUID,
  json_extract(node_info, '$.bootID') as bootID,
  json_extract(node_info, '$.kernelVersion') as kernelVersion,
  json_extract(node_info, '$.osImage') as osImage,
  json_extract(node_info, '$.operatingSystem') as operatingSystem,
  json_extract(node_info, '$.architecture') as architecture,
  json_extract(node_info, '$.containerRuntimeVersion') as containerRuntimeVersion,
  json_extract(node_info, '$.kubeletVersion') as kubeletVersion,
  json_extract(node_info, '$.kubeProxyVersion') as kubeProxyVersion
from
  kubernetes_node;
```

### Node IP info
Gain insights into the distribution of internal and external IP addresses, as well as hostnames and internal DNS, across your Kubernetes nodes. This can help in managing network configurations and troubleshooting connectivity issues.

```sql+postgres
select
  name,
  jsonb_path_query_array(
    addresses,
    '$[*] ? (@.type == "ExternalIP").address'
  ) as public_ips,
  jsonb_path_query_array(
    addresses,
    '$[*] ? (@.type == "InternalIP").address'
  ) as internal_ips,
    jsonb_path_query_array(
    addresses,
    '$[*] ? (@.type == "InternalDNS").address'
  ) as internal_dns,
  jsonb_path_query_array(
    addresses,
    '$[*] ? (@.type == "Hostname").address'
  ) as hostnames
from
  kubernetes_node;
```

```sql+sqlite
Error: The corresponding SQLite query is unavailable.
```

### List manifest resources
Explore which Kubernetes nodes have a specified path, providing insights into the distribution of resources across your infrastructure. This can help optimize resource allocation and ensure balanced workload distribution.

```sql+postgres
select
  name,
  pod_cidr,
  pod_cidrs,
  provider_id,
  addresses,
  capacity,
  path
from
  kubernetes_node
where
  path is not null;
```

```sql+sqlite
select
  name,
  pod_cidr,
  pod_cidrs,
  provider_id,
  addresses,
  capacity,
  path
from
  kubernetes_node
where
  path is not null;
```