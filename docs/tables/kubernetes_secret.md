# Table: kubernetes_secret

Secrets are used to store sensitive information either as individual properties or coarse-grained entries like entire files or JSON blobs.

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
  kubernetes_secret,
  jsonb_each(data) as data
order by
  namespace,
  name;
```


### List and base64 decode secret values
```sql
select
  name,
  namespace,
  data.key,
  decode(data.value, 'base64') as decoded_data,
  age(current_timestamp, creation_timestamp)
from
  kubernetes_secret,
  jsonb_each_text(data) as data
order by
  namespace,
  name;
```
