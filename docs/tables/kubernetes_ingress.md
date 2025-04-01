---
title: "Steampipe Table: kubernetes_ingress - Query Kubernetes Ingresses using SQL"
description: "Allows users to query Kubernetes Ingresses, specifically to obtain details about the network traffic routing rules, providing insights into application or service access patterns."
folder: "Ingress"
---

# Table: kubernetes_ingress - Query Kubernetes Ingresses using SQL

Kubernetes Ingress is a collection of routing rules that govern how external users access services running in a Kubernetes cluster. Typically, these rules are used to expose services to external traffic coming from the internet. It provides a way to manage external access to the services in a cluster, typically HTTP.

## Table Usage Guide

The `kubernetes_ingress` table provides insights into Ingresses within Kubernetes. As a DevOps engineer, explore Ingress-specific details through this table, including host information, backend service details, and associated annotations. Utilize it to uncover information about Ingresses, such as those with specific routing rules, the services they expose, and their configurations.

## Examples

### Basic Info
Explore which Kubernetes ingress resources are associated with specific namespaces and classes, and how long they have been created. This can help in tracking resource allocation and usage over time.

```sql+postgres
select
  name,
  namespace,
  ingress_class_name as class,
  age(current_timestamp, creation_timestamp)
from
  kubernetes_ingress
order by
  namespace,
  name;
```

```sql+sqlite
select
  name,
  namespace,
  ingress_class_name as class,
  (julianday('now') - julianday(creation_timestamp)) * 24 * 60 * 60 as age
from
  kubernetes_ingress
order by
  namespace,
  name;
```

### View rules for the ingress
Explore which ingress rules are currently in place within your Kubernetes environment. This can help in understanding and managing traffic routing, ensuring efficient and secure communication between services.

```sql+postgres
select
  name,
  namespace,
  jsonb_pretty(rules) as rules
from
  kubernetes_ingress;
```

```sql+sqlite
select
  name,
  namespace,
  rules
from
  kubernetes_ingress;
```

### List manifest resources
Explore which Kubernetes ingress resources are configured with a specific path. This can help identify areas where traffic routing rules have been established, which is essential for understanding and managing application traffic flow.

```sql+postgres
select
  name,
  namespace,
  ingress_class_name as class,
  path
from
  kubernetes_ingress
where
  path is not null
order by
  namespace,
  name;
```

```sql+sqlite
select
  name,
  namespace,
  ingress_class_name as class,
  path
from
  kubernetes_ingress
where
  path is not null
order by
  namespace,
  name;
```