# Table: kubernetes_stateful_set

In Kubernetes, stateful sets represent a set of pods with unique, persistent identities and stable hostnames that GKE maintains regardless of where they are scheduled.

## Examples

### Basic Info - `kubectl get statefulsets --all-namespaces` columns

```sql
select
  name,
  namespace,
  service_name,
  replicas,
  age(current_timestamp, creation_timestamp)
from
  kubernetes_stateful_set
order by
  namespace,
  name;
```

### List stateful sets that require manual update when the object's configuration is changed

```sql
select
  name,
  namespace,
  service_name,
  update_strategy ->> 'type' as update_strategy_type
from
  kubernetes_stateful_set
where
  update_strategy ->> 'type' = 'OnDelete';
```
