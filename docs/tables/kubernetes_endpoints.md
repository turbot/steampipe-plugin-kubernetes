# Table: kubernetes_endpoints

An Endpoints resource is an abstraction, linked to a Service, which defines the list of endpoints that actually implement the service.

## Examples

### Basic Info

```sql
select
  name,
  namespace,
  subsets
from
  kubernetes_endpoints;
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
