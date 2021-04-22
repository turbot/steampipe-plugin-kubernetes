# Table: kubernetes_job

A Job creates one or more Pods and will continue to retry execution of the Pods until a specified number of them successfully terminate.

## Examples

### Basic Info

```sql
select
  name,
  namespace,
  selector,
  labels,
  annotations,
  parallelism,
  completions,
  start_time,
  completion_time,
  age(coalesce(completion_time, current_timestamp), start_time) as duration,
  active_deadline_seconds,
  active,
  succeeded,
  failed
from
  k8s_minikube.kubernetes_job
```

### List active jobs

```sql
select
  name,
  namespace,
  start_time,
  age(coalesce(completion_time, current_timestamp), start_time) as duration,
  active,
  succeeded,
  failed
from
  k8s_minikube.kubernetes_job
where active > 0
```

### List failed jobs

```sql
select
  name,
  namespace,
  start_time,
  age(coalesce(completion_time, current_timestamp), start_time) as duration,
  active,
  succeeded,
  failed
from
  k8s_minikube.kubernetes_job
where failed > 0
```

### Get list of container and images for jobs

```sql
select
  name,
  namespace,
  jsonb_agg(elems.value -> 'name') as containers,
  jsonb_agg(elems.value -> 'image') as images
from
  k8s_minikube.kubernetes_job,
  jsonb_array_elements(template -> 'spec' -> 'containers') as elems
group by name, namespace
```
