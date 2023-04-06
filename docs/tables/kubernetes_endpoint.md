# Table: kubernetes_endpoint

An Endpoints resource is an abstraction, linked to a Service, which defines the list of endpoints that actually implement the service.

## Examples

### Basic Info

```sql
select
  name,
  namespace,
  subsets
from
  kubernetes_endpoint;
```

### Endpoint IP Info

```sql
select
  name,
  namespace,
  addr ->> 'ip' as address,
  nr_addr ->> 'ip'  as not_ready_address,
  port -> 'port' as port,
  port ->> 'protocol' as protocol
from
  kubernetes_endpoint,
  jsonb_array_elements(subsets) as subset
  left join jsonb_array_elements(subset -> 'addresses') as addr on true
  left join jsonb_array_elements(subset -> 'notReadyAddresses') as nr_addr on true
  left join jsonb_array_elements(subset -> 'ports') as port on true;
```

### List manifest resources

```sql
select
  name,
  namespace,
  subsets
from
  kubernetes_endpoint
where
  manifest_file_path is not null;
```
