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

### Get subsets info for a specific endpoint

```sql
select
  name,
  namespace,
  (addresse ->> 'ip')::inet as address,
  port -> 'port' as port,
  port ->> 'protocol' as protocol
from
  kubernetes_endpoints,
  jsonb_array_elements(subsets) as subset,
  jsonb_array_elements(subset -> 'addresses') as addresse,
  jsonb_array_elements(subset -> 'ports') as port
where
  name = 'frontend'
  and namespace = 'default';
```
