---
title: "Steampipe Table: kubernetes_job - Query Kubernetes Jobs using SQL"
description: "Allows users to query Kubernetes Jobs, providing insights into the detailed status of each job, including the number of successful completions and the parallelism limit."
---

# Table: kubernetes_job - Query Kubernetes Jobs using SQL

Kubernetes Jobs are a resource that represent a finite task, i.e., they run until successful completion. They create one or more Pods and ensure that a specified number of them successfully terminate. As pods successfully complete, the job tracks the successful completions.

## Table Usage Guide

The `kubernetes_job` table provides insights into Kubernetes Jobs within a Kubernetes cluster. As a DevOps engineer, explore job-specific details through this table, including the status of each job, the number of successful completions, and the parallelism limit. Utilize it to monitor the progress of jobs, ensure that they are running as expected, and troubleshoot any issues that occur.

## Examples

### Basic Info
Explore the status and performance of jobs within a Kubernetes environment. This allows users to assess job completion status, duration, and overall efficiency, aiding in system monitoring and optimization.

```sql+postgres
select
  name,
  namespace,
  active,
  succeeded,
  failed,
  completions,
  start_time,
  completion_time,
  age(coalesce(completion_time, current_timestamp), start_time) as duration,
  active_deadline_seconds,
  parallelism,
  selector,
  labels,
  annotations
from
  kubernetes_job;
```

```sql+sqlite
select
  name,
  namespace,
  active,
  succeeded,
  failed,
  completions,
  start_time,
  completion_time,
  strftime('%s', coalesce(completion_time, current_timestamp)) - strftime('%s', start_time) as duration,
  active_deadline_seconds,
  parallelism,
  selector,
  labels,
  annotations
from
  kubernetes_job;
```

### List active jobs
Determine the areas in which jobs are currently active within a Kubernetes environment. This can be useful in managing resources and identifying any potential issues or bottlenecks.

```sql+postgres
select
  name,
  namespace,
  start_time,
  age(coalesce(completion_time, current_timestamp), start_time) as duration,
  active,
  succeeded,
  failed
from
  kubernetes_job
where active > 0;
```

```sql+sqlite
select
  name,
  namespace,
  start_time,
  strftime('%s', coalesce(completion_time, current_timestamp)) - strftime('%s', start_time) as duration,
  active,
  succeeded,
  failed
from
  kubernetes_job
where active > 0;
```

### List failed jobs
Identify instances where jobs have failed in a Kubernetes environment. This enables quick detection of issues and facilitates timely troubleshooting.

```sql+postgres
select
  name,
  namespace,
  start_time,
  age(coalesce(completion_time, current_timestamp), start_time) as duration,
  active,
  succeeded,
  failed
from
  kubernetes_job
where failed > 0;
```

```sql+sqlite
select
  name,
  namespace,
  start_time,
  coalesce(completion_time, datetime('now')) - start_time as duration,
  active,
  succeeded,
  failed
from
  kubernetes_job
where failed > 0;
```

### Get list of container and images for jobs
The query provides a way to identify the containers and images associated with specific jobs in a Kubernetes environment. This can be particularly useful for system administrators to track the resources being used by different jobs and ensure optimal resource allocation.

```sql+postgres
select
  name,
  namespace,
  jsonb_agg(elems.value -> 'name') as containers,
  jsonb_agg(elems.value -> 'image') as images
from
  kubernetes_job,
  jsonb_array_elements(template -> 'spec' -> 'containers') as elems
group by
  name,
  namespace;
```

```sql+sqlite
select
  name,
  namespace,
  json_group_array(json_extract(elems.value, '$.name')) as containers,
  json_group_array(json_extract(elems.value, '$.image')) as images
from
  kubernetes_job,
  json_each(template, '$.spec.containers') as elems
group by
  name,
  namespace;
```

### List manifest resources
Explore the status and details of active Kubernetes jobs, including their success and failure rates. This can be useful for identifying any jobs that may require attention or troubleshooting.

```sql+postgres
select
  name,
  namespace,
  active,
  succeeded,
  failed,
  completions,
  parallelism,
  selector,
  labels,
  annotations,
  path
from
  kubernetes_job
where
  path is not null;
```

```sql+sqlite
select
  name,
  namespace,
  active,
  succeeded,
  failed,
  completions,
  parallelism,
  selector,
  labels,
  annotations,
  path
from
  kubernetes_job
where
  path is not null;
```