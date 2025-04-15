---
title: "Steampipe Table: kubernetes_custom_resource_definition - Query Kubernetes Custom Resource Definitions using SQL"
description: "Allows users to query Custom Resource Definitions in Kubernetes, providing details about the structure and configuration of the custom resources."
folder: "CRD"
---

# Table: kubernetes_custom_resource_definition - Query Kubernetes Custom Resource Definitions using SQL

Kubernetes Custom Resource Definitions (CRDs) allow users to create new types of resources that they can later use like the built-in resource types in Kubernetes. These custom resources can be used to store and retrieve structured data. They extend the Kubernetes API, allowing developers to define the kind of resources that they need to work with in their applications.

## Table Usage Guide

The `kubernetes_custom_resource_definition` table provides insights into Custom Resource Definitions within Kubernetes. As a Kubernetes developer or administrator, explore the details of these custom resources through this table, including their structure, configuration, and associated metadata. Utilize it to uncover information about the custom resources, such as their validation schema, versioning details, and the scope of their usage within the Kubernetes cluster.

## Examples

### Basic Info
Explore the fundamental details of custom resources in your Kubernetes environment to gain insights into their identities and creation times. This can be useful in understanding the composition and history of your resources.

```sql+postgres
select
  name,
  namespace,
  uid,
  resource_version,
  creation_timestamp
from
  kubernetes_custom_resource_definition;
```

```sql+sqlite
select
  name,
  namespace,
  uid,
  resource_version,
  creation_timestamp
from
  kubernetes_custom_resource_definition;
```

### List CRDs for a particular group
Discover the segments that contain custom resource definitions (CRDs) for a specific group within your Kubernetes environment. This can be particularly useful for managing and tracking your resources, especially in larger deployments.

```sql+postgres
select
  name,
  namespace,
  uid,
  resource_version,
  creation_timestamp
from
  kubernetes_custom_resource_definition
where
  spec ->> 'group' = 'stable.example.com';
```

```sql+sqlite
select
  name,
  namespace,
  uid,
  resource_version,
  creation_timestamp
from
  kubernetes_custom_resource_definition
where
  json_extract(spec, '$.group') = 'stable.example.com';
```

### List Certificate type CRDs
Explore which custom resources in your Kubernetes cluster are of the 'Certificate' type. This is useful for managing and tracking security certificates within your system.

```sql+postgres
select
  name,
  namespace,
  uid,
  resource_version,
  creation_timestamp
from
  kubernetes_custom_resource_definition
where
  status -> 'acceptedNames' ->> 'kind' = 'Certificate';
```

```sql+sqlite
select
  name,
  namespace,
  uid,
  resource_version,
  creation_timestamp
from
  kubernetes_custom_resource_definition
where
  json_extract(json_extract(status, '$.acceptedNames'), '$.kind') = 'Certificate';
```

### List namespaced CRDs
Explore which custom resource definitions (CRDs) are namespaced in your Kubernetes system. This can be useful for understanding the scope of your CRDs and managing resources more effectively.

```sql+postgres
select
  name,
  namespace,
  uid,
  resource_version,
  creation_timestamp
from
  kubernetes_custom_resource_definition
where
  spec ->> 'scope' = 'Namespaced';
```

```sql+sqlite
select
  name,
  namespace,
  uid,
  resource_version,
  creation_timestamp
from
  kubernetes_custom_resource_definition
where
  json_extract(spec, '$.scope') = 'Namespaced';
```

### Get active version detail of each CRD
Explore the active versions of each Custom Resource Definition (CRD) in your Kubernetes environment. This is useful for identifying which versions are currently being served and can help in maintaining version control.

```sql+postgres
select
  name,
  namespace,
  creation_timestamp,
  jsonb_pretty(v) as active_version
from
  kubernetes_custom_resource_definition,
  jsonb_array_elements(spec -> 'versions') as v
where
  v ->> 'served' = 'true';
```

```sql+sqlite
select
  name,
  namespace,
  creation_timestamp,
  v.value as active_version
from
  kubernetes_custom_resource_definition,
  json_each(spec, '$.versions') as v
where
  json_extract(v.value, '$.served') = 'true';
```

### List CRDs created within the last 90 days
Identify the custom resource definitions (CRDs) that have been created in your Kubernetes environment within the last 90 days. This can be useful for tracking recent changes and monitoring the development of your resources.

```sql+postgres
select
  name,
  namespace,
  uid,
  resource_version,
  creation_timestamp
from
  kubernetes_custom_resource_definition
where
  creation_timestamp >= (now() - interval '90' day)
order by
  creation_timestamp;
```

```sql+sqlite
select
  name,
  namespace,
  uid,
  resource_version,
  creation_timestamp
from
  kubernetes_custom_resource_definition
where
  creation_timestamp >= datetime('now', '-90 day')
order by
  creation_timestamp;
```

### Get spec detail of each CRD
Discover the segments that detail each Custom Resource Definition (CRD) in your Kubernetes cluster. This query is useful for understanding the scope, versions, and conversion details of your CRDs, which can help you manage and optimize your Kubernetes resources.

```sql+postgres
select
  name,
  uid,
  creation_timestamp,
  spec ->> 'group' as "group",
  spec -> 'names' as "names",
  spec ->> 'scope' as "scope",
  spec -> 'versions' as "versions",
  spec -> 'conversion' as "conversion"
from
  kubernetes_custom_resource_definition;
```

```sql+sqlite
select
  name,
  uid,
  creation_timestamp,
  json_extract(spec, '$.group') as "group",
  json_extract(spec, '$.names') as "names",
  json_extract(spec, '$.scope') as "scope",
  json_extract(spec, '$.versions') as "versions",
  json_extract(spec, '$.conversion') as "conversion"
from
  kubernetes_custom_resource_definition;
```

### List manifest resources
Explore which custom resources in your Kubernetes cluster have defined paths. This can be useful in understanding the structure and organization of your resources.

```sql+postgres
select
  name,
  namespace,
  resource_version,
  path
from
  kubernetes_custom_resource_definition
where
  path is not null;
```

```sql+sqlite
select
  name,
  namespace,
  resource_version,
  path
from
  kubernetes_custom_resource_definition
where
  path is not null;
```