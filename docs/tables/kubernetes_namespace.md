# Table: kubernetes_namespace

Kubernetes supports multiple virtual clusters backed by the same physical cluster. These virtual clusters are called namespaces.

Namespaces are intended for use in environments with many users spread across multiple teams, or projects. They provide a scope for names. Names of resources need to be unique within a namespace, but not across namespaces.

## Examples

### Basic Info

```sql
select
  name,
  phase as status,
  annotations,
  labels
from
  kubernetes_namespace;
```

### List manifest resources

```sql
select
  name,
  phase as status,
  annotations,
  labels
from
  kubernetes_namespace
where
  manifest_file_path is not null;
```
