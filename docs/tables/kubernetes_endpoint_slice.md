---
title: "Steampipe Table: kubernetes_endpoint_slice - Query Kubernetes Endpoint Slices using SQL"
description: "Allows users to query Kubernetes Endpoint Slices, providing insights into the set of endpoints that a service may route traffic to."
---

# Table: kubernetes_endpoint_slice - Query Kubernetes Endpoint Slices using SQL

Kubernetes Endpoint Slices are a scalable and extensible way to network traffic routing. They provide a simple way to track network endpoints within a Kubernetes cluster. Endpoint Slices group network endpoints together, allowing for efficient and flexible traffic routing.

## Table Usage Guide

The `kubernetes_endpoint_slice` table provides insights into the Endpoint Slices within a Kubernetes cluster. As a network engineer or DevOps professional, explore Endpoint Slice-specific details through this table, including associated services, ports, and addresses. Utilize it to manage and optimize network traffic routing within your Kubernetes environment.

## Examples

### Basic Info
Explore the configuration of your Kubernetes environment by identifying its various endpoints, their corresponding addresses and ports. This can provide valuable insights into the network architecture and communication within your Kubernetes cluster.

```sql+postgres
select
  name,
  namespace,
  generate_name as endpoint_name,
  address_type,
  endpoints,
  ports
from
  kubernetes_endpoint_slice;
```

```sql+sqlite
select
  name,
  namespace,
  generate_name as endpoint_name,
  address_type,
  endpoints,
  ports
from
  kubernetes_endpoint_slice;
```

### Endpoint Slice IP Information
Analyze the settings to understand the IP information for endpoint slices in a Kubernetes environment. This can be beneficial in identifying potential networking issues or inconsistencies within your application's communication paths.

```sql+postgres
select
  name,
  namespace,
  addr,
  port -> 'port' as port,
  port ->> 'protocol' as protocol
from
    kubernetes_endpoint_slice,
    jsonb_array_elements(endpoints) as ep,
    jsonb_array_elements(ep -> 'addresses') as addr,
    jsonb_array_elements(ports) as port;
```

```sql+sqlite
select
  name,
  namespace,
  addr.value as addr,
  json_extract(port.value, '$.port') as port,
  json_extract(port.value, '$.protocol') as protocol
from
  kubernetes_endpoint_slice,
  json_each(endpoints) as ep,
  json_each(json_extract(ep.value, '$.addresses')) as addr,
  json_each(ports) as port;
```

### List manifest resources
Explore the various manifest resources within a Kubernetes cluster, specifically identifying those with a defined path. This can help in understanding the distribution and configuration of resources, which is vital for efficient cluster management and troubleshooting.

```sql+postgres
select
  name,
  namespace,
  generate_name as endpoint_name,
  address_type,
  endpoints,
  ports,
  path
from
  kubernetes_endpoint_slice
where
  path is not null;
```

```sql+sqlite
select
  name,
  namespace,
  generate_name as endpoint_name,
  address_type,
  endpoints,
  ports,
  path
from
  kubernetes_endpoint_slice
where
  path is not null;
```