# Table: kubernetes_storage_class

A StorageClass provides a way for administrators to describe the 'classes' of storage they offer. Different classes might map to quality-of-service levels, or to backup policies, or to arbitrary policies determined by the cluster administrators. Kubernetes itself is unopinionated about what classes represent. This concept is sometimes called 'profiles' in other storage systems.

## Examples

### Basic Info

```sql
select
  name,
  provisioner,
  reclaim_policy,
  parameters,
  mount_options
from
  kubernetes_storage_class;
```

### List storage classes that doesn't allow volume expansion

```sql
select
  name,
  allow_volume_expansion,
  provisioner,
  reclaim_policy,
  parameters,
  mount_options
from
  kubernetes_storage_class
where
  not allow_volume_expansion;
```

### List storage classes with immediate volume binding mode
```sql
select
  name,
  allow_volume_expansion,
  provisioner,
  reclaim_policy,
  volume_binding_mode
from
  kubernetes_storage_class
where
  volume_binding_mode = 'Immediate';
```
