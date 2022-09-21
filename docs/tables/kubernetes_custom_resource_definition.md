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

### Get CustomResourceDefinitions for a particular group

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

### Get spec detail for CustomResourceDefinitions

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
