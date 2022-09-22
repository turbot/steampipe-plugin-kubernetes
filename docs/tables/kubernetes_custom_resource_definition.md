# Table: kubernetes_custom_resource_definition

The CustomResourceDefinition API resource allows you to define custom resources. Defining a CRD object creates a new custom resource with a name and schema that you specify. The Kubernetes API serves and handles the storage of your custom resource. The name of a CRD object must be a valid DNS subdomain name.

## Examples

### Basic Info

```sql
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

```sql
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

### List Certificate type CRDs

```sql
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

### List namespaced CRDs

```sql
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

### Get active version detail of each CRD

```sql
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

### List CRDs created within the last 90 days

```sql
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

### Get spec detail of each CRD

```sql
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
