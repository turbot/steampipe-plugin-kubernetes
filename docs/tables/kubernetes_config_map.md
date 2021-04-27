# Table: kubernetes_config_map

Config map can be used to store fine-grained information like individual properties or coarse-grained information like entire config files or JSON blobs.

## Examples

### Basic Info

```sql
select
  name,
  namespace,
  data.key,
  data.value,
  age(current_timestamp, creation_timestamp)
from
  kubernetes_config_map,
  jsonb_each(data) as data
order by
  namespace,
  name;
```
