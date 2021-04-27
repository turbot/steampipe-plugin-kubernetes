# Table: kubernetes_persistent_volume_claim

A PersistentVolumeClaim (PVC) is a request for storage by a user.

## Examples

### Basic Info

```sql
select
  name,
  namespace,
  volume_name as volume,
  volume_mode,
  access_modes,
  phase as status,
  capacity ->> 'storage' as capacity,
  creation_timestamp,
  data_source,
  selector,
  resources,
  age(current_timestamp, creation_timestamp)
from
  kubernetes_persistent_volume_claim;
```
