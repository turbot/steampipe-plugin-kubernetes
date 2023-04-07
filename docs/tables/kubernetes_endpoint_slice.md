# Table: kubernetes_endpoint_slice

Represents a subset of the endpoints that implement a service.

## Examples

### Basic Info

```sql
select
  name,
  namespace,
  generate_name as endpoint_name,
  address_type,
  endpoints,
  ports
from
  kubernetes_endpoint_slice;
```

### Endpoint Slice IP Information

```sql
select
  name,
  namespace,
  addr,
  port -> 'port' as port,
  port ->> 'protocol' as protocol
from
    kubernetes_endpoint_slice,
    jsonb_array_elements(endpoints) as ep,
    jsonb_array_elements(ep -> 'addresses') as addr,
    jsonb_array_elements(ports) as port;
```

### List manifest resources

```sql
select
  name,
  namespace,
  generate_name as endpoint_name,
  address_type,
  endpoints,
  ports
from
  kubernetes_endpoint_slice
where
  manifest_file_path is not null;
```
