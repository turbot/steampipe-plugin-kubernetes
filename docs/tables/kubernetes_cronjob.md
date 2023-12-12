---
title: "Steampipe Table: kubernetes_cronjob - Query Kubernetes CronJobs using SQL"
description: "Allows users to query Kubernetes CronJobs, providing insights into scheduled tasks within the Kubernetes environment."
---

# Table: kubernetes_cronjob - Query Kubernetes CronJobs using SQL

A Kubernetes CronJob creates Jobs on a repeating schedule, similar to the job scheduling in Unix-like systems. It is a way to run automated tasks at regular, predetermined times. CronJobs use the Cron format to schedule tasks.

## Table Usage Guide

The `kubernetes_cronjob` table provides insights into CronJobs within Kubernetes. As a DevOps engineer, explore CronJob-specific details through this table, including schedules, job histories, and associated metadata. Utilize it to monitor and manage your automated tasks, and ensure they are running as expected.

## Examples

### Basic Info
Explore which scheduled tasks within your Kubernetes environment have failed. This allows for proactive troubleshooting and understanding of task scheduling and execution patterns.

```sql+postgres
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

```sql+sqlite
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
Explore which cronJobs are running in your Kubernetes environment and identify the containers and images they are using. This is useful to understand the dependencies and configurations of your scheduled tasks, and can help in troubleshooting or optimizing resource usage.

```sql+postgres
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

```sql+sqlite
select
  name,
  namespace,
  json_group_array(json_extract(elems.value, '$.name')) as containers,
  json_group_array(json_extract(elems.value, '$.image')) as images
from
  kubernetes_cronjob,
  json_each(job_template, '$.spec.template.spec.containers') as elems
group by
  name,
  namespace;
```

### List manifest resources
Explore which scheduled tasks within your Kubernetes environment have a specified path. This can be useful to identify tasks that may be associated with certain applications or services, helping you to manage and monitor your resources more effectively.

```sql+postgres
select
  name,
  namespace,
  uid,
  failed_jobs_history_limit,
  schedule,
  suspend,
  path
from
  kubernetes_cronjob
where
  path is not null;
```

```sql+sqlite
select
  name,
  namespace,
  uid,
  failed_jobs_history_limit,
  schedule,
  suspend,
  path
from
  kubernetes_cronjob
where
  path is not null;
```