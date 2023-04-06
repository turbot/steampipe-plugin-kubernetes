# Table: kubernetes_cronjob

Cron jobs are useful for creating periodic and recurring tasks, like running backups or sending emails. Cron jobs can also schedule individual tasks for a specific time, such as if you want to schedule a job for a low activity period.

## Examples

### Basic Info

```sql
select
  name,
  namespace,
  uid,
  failed_jobs_history_limit,
  schedule,
  suspend
from
  kubernetes_cronjob;
```

### Get list of container and images for cronJobs

```sql
select
  name,
  namespace,
  jsonb_agg(elems.value -> 'name') as containers,
  jsonb_agg(elems.value -> 'image') as images
from
  kubernetes_cronjob,
  jsonb_array_elements(job_template -> 'spec' -> 'template' -> 'spec' -> 'containers') as elems
group by 
  name, 
  namespace;
```

### List manifest resources

```sql
select
  name,
  namespace,
  uid,
  failed_jobs_history_limit,
  schedule,
  suspend
from
  kubernetes_cronjob
where
  manifest_file_path is not null;
```
