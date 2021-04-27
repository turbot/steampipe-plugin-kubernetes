# Table: kubernetes_persistent_volume

A PersistentVolume (PV) is a piece of storage in the cluster that has been provisioned by an administrator or dynamically provisioned using Storage Classes. PVs are volume plugins like Volumes, but have a lifecycle independent of any individual Pod that uses the PV.

## Examples

### Basic Info

```sql
select
  name,
  access_modes,
  storage_class,
  capacity ->> 'storage' as storage_capacity,
  creation_timestamp,
  persistent_volume_reclaim_policy,
  phase as status,
  volume_mode,
  age(current_timestamp, creation_timestamp)
from
  kubernetes_persistent_volume;
```

### Get hostpath details for the volume

```sql
select
  name,
  persistent_volume_source -> 'hostPath' ->> 'path' as path,
  persistent_volume_source -> 'hostPath' ->> 'type' as type
from
  kubernetes_persistent_volume;
```
