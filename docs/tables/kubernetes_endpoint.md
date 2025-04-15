---
title: "Steampipe Table: kubernetes_endpoint - Query Kubernetes Endpoints using SQL"
description: "Allows users to query Kubernetes Endpoints, providing a detailed overview of each endpoint's configuration and status."
folder: "Endpoint"
---

# Table: kubernetes_endpoint - Query Kubernetes Endpoints using SQL

Kubernetes Endpoints are a part of the Kubernetes Service concept that represents a real or virtual IP address and a port number that Kubernetes pods use to access services. Endpoints can be defined as a network endpoint that can receive traffic. They are a fundamental part of the Kubernetes networking model, allowing pods to communicate with each other and with services outside the Kubernetes cluster.

## Table Usage Guide

The `kubernetes_endpoint` table provides insights into endpoints within Kubernetes. As a DevOps engineer, you can explore details about each endpoint through this table, including its associated services, IP addresses, and ports. Use this table to understand the communication paths within your Kubernetes cluster, track the status of endpoints, and identify any potential networking issues.

## Examples

### Basic Info
Explore which Kubernetes endpoints are currently active in your system. This can help you understand the communication points within your clusters and troubleshoot any networking issues.

```sql+postgres
select
  name,
  namespace,
  subsets
from
  kubernetes_endpoint;
```

```sql+sqlite
select
  name,
  namespace,
  subsets
from
  kubernetes_endpoint;
```

### Endpoint IP Info
Determine the areas in which endpoint IP information, such as address, readiness status and protocol, is used in your Kubernetes environment. This can aid in network troubleshooting and enhancing security measures.

```sql+postgres
select
  name,
  namespace,
  addr ->> 'ip' as address,
  nr_addr ->> 'ip'  as not_ready_address,
  port -> 'port' as port,
  port ->> 'protocol' as protocol
from
  kubernetes_endpoint,
  jsonb_array_elements(subsets) as subset
  left join jsonb_array_elements(subset -> 'addresses') as addr on true
  left join jsonb_array_elements(subset -> 'notReadyAddresses') as nr_addr on true
  left join jsonb_array_elements(subset -> 'ports') as port on true;
```

```sql+sqlite
select
  kubernetes_endpoint.name,
  kubernetes_endpoint.namespace,
  json_extract(addr.value, '$.ip') as address,
  json_extract(nr_addr.value, '$.ip') as not_ready_address,
  json_extract(port.value, '$.port') as port,
  json_extract(port.value, '$.protocol') as protocol
from
  kubernetes_endpoint,
  json_each(kubernetes_endpoint.subsets) as subset,
  json_each(json_extract(subset.value, '$.addresses')) as addr,
  json_each(json_extract(subset.value, '$.notReadyAddresses')) as nr_addr,
  json_each(json_extract(subset.value, '$.ports')) as port;
```

### List manifest resources
Explore which Kubernetes endpoints have a specified path. This is useful to understand the distribution of resources within your Kubernetes environment.

```sql+postgres
select
  name,
  namespace,
  subsets,
  path
from
  kubernetes_endpoint
where
  path is not null;
```

```sql+sqlite
select
  name,
  namespace,
  subsets,
  path
from
  kubernetes_endpoint
where
  path is not null;
```