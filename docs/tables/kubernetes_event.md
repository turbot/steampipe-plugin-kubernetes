# Table: kubernetes_event

Event is a report of an event somewhere in the cluster. Events have a limited retention time and triggers and messages may evolve with time. Events should be treated as informative, best-effort, supplemental data.

## Examples

### Basic Info

```sql
select
  namespace,
  last_timestamp,
  type,
  reason,
  concat(involved_object ->> 'kind', '/', involved_object ->> 'name') as object,
  message
from
  kubernetes_event;
```

### List warning events by last timestamp

```sql
select
  namespace,
  last_timestamp,
  type,
  reason,
  concat(involved_object ->> 'kind', '/', involved_object ->> 'name') as object,
  message
from
  kubernetes_event
where
  type = 'Warning'
order by
  namespace,
  last_timestamp;
```

### List manifest resources

```sql
select
  namespace,
  type,
  reason,
  concat(involved_object ->> 'kind', '/', involved_object ->> 'name') as object,
  message,
  path
from
  kubernetes_event
where
  path is not null;
```
