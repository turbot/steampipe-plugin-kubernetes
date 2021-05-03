# Table: kubernetes_statefulset

In Kubernetes, statefulSets represent a set of Pods with unique, persistent identities and stable hostnames that GKE maintains regardless of where they are scheduled.

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
  kubernetes_statefulset
order by
  namespace,
  name;
```

### List statefulSets requires manual update when the object's configuration is changed

```sql
select
  name,
  namespace,
  service_name,
  update_strategy ->> 'type' as update_strategy_type
from
  kubernetes_statefulset
where
  update_strategy ->> 'type' = 'OnDelete';
```
