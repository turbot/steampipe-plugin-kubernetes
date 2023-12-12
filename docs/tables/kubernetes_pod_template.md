---
title: "Steampipe Table: kubernetes_pod_template - Query Kubernetes Pod Templates using SQL"
description: "Allows users to query Kubernetes Pod Templates, specifically providing insights into the template's metadata, specification, and status."
---

# Table: kubernetes_pod_template - Query Kubernetes Pod Templates using SQL

A Kubernetes Pod Template is a pod specification which produces the same pod each time it is instantiated. It is used to create a pod directly, or it is nested inside replication controllers, jobs, replicasets, etc. A Pod Template in a workload object must have a Labels field and it must match the selector of its controlling workload object.

## Table Usage Guide

The `kubernetes_pod_template` table provides insights into pod templates within Kubernetes. As a Kubernetes administrator, you can explore pod template-specific details through this table, including metadata, specifications, and status. Utilize it to uncover information about pod templates, such as those with specific labels, the replication controllers they are nested in, and the status of each pod template.

## Examples

### Basic info
Discover the segments that show the age of pod templates in your Kubernetes environment, along with the count of various container types. This can help in managing resources and understanding the capacity usage within your system.

```sql+postgres
select
  name,
  namespace,
  age(current_timestamp, creation_timestamp),
  template -> 'spec' ->> 'node_name' as pod_node_name,
  jsonb_array_length(template -> 'spec' -> 'containers') as container_count,
  jsonb_array_length(template -> 'spec' -> 'pod_init_containers') as init_container_count,
  jsonb_array_length(template -> 'spec' -> 'pod_ephemeral_containers') as ephemeral_container_count
from
  kubernetes_pod_template
order by
  namespace,
  name;
```

```sql+sqlite
select
  name,
  namespace,
  julianday('now') - julianday(creation_timestamp) as age,
  json_extract(template, '$.spec.node_name') as pod_node_name,
  json_array_length(json_extract(template, '$.spec.containers')) as container_count,
  json_array_length(json_extract(template, '$.spec.pod_init_containers')) as init_container_count,
  json_array_length(json_extract(template, '$.spec.pod_ephemeral_containers')) as ephemeral_container_count
from
  kubernetes_pod_template
order by
  namespace,
  name;
```

### List pod templates with privileged pod containers
Uncover the details of your system's pod templates which contain privileged pod containers. This allows you to assess the security implications and manage the risk associated with these privileged containers.

```sql+postgres
select
  name,
  namespace,
  template -> 'metadata' ->> 'name' as pod_name,
  c ->> 'name' as container_name,
  c ->> 'image' as container_image
from
  kubernetes_pod_template,
  jsonb_array_elements(template -> 'spec' -> 'containers') as c
where
  c -> 'securityContext' ->> 'privileged' = 'true';
```

```sql+sqlite
select
  name,
  namespace,
  json_extract(template, '$.metadata.name') as pod_name,
  json_extract(c.value, '$.name') as container_name,
  json_extract(c.value, '$.image') as container_image
from
  kubernetes_pod_template,
  json_each(json_extract(template, '$.spec.containers')) as c
where
  json_extract(c.value, '$.securityContext.privileged') = 'true';
```

### List pod templates with pod access to the host process ID, IPC, or network namespace
Explore which pod templates have access to the host process ID, IPC, or network namespace. This is useful for identifying potential security risks and ensuring appropriate access control in a Kubernetes environment.

```sql+postgres
select
  name,
  namespace,
  template -> 'metadata' ->> 'name' as pod_name,
  template -> 'spec' -> 'host_pid' as pod_host_pid,
  template -> 'spec' -> 'host_ipc' as pod_host_ipc,
  template -> 'spec' -> 'host_network' as pod_host_network
from
  kubernetes_pod_template
where
  (template -> 'spec' -> 'host_pid')::boolean
  or (template -> 'spec' -> 'host_ipc')::boolean
  or (template -> 'spec' -> 'host_network')::boolean;
```

```sql+sqlite
select
  name,
  namespace,
  json_extract(template, '$.metadata.name') as pod_name,
  json_extract(template, '$.spec.host_pid') as pod_host_pid,
  json_extract(template, '$.spec.host_ipc') as pod_host_ipc,
  json_extract(template, '$.spec.host_network') as pod_host_network
from
  kubernetes_pod_template
where
  json_extract(template, '$.spec.host_pid') = 1
  or json_extract(template, '$.spec.host_ipc') = 1
  or json_extract(template, '$.spec.host_network') = 1;
```

### `kubectl get podtemplates` columns
Determine the areas in which Kubernetes pod templates are deployed. This query helps in identifying the containers and images used, along with the associated pod labels, providing a comprehensive summary of your Kubernetes deployment.

```sql+postgres
select
  name,
  coalesce(uid, concat(path, ':', start_line)) as uid,
  array_agg(c ->> 'name') as containers,
  array_agg(c ->> 'image') as images,
  template -> 'metadata' -> 'labels' as pod_labels 
from
  kubernetes_pod_template,
  jsonb_array_elements(template -> 'spec' -> 'containers') as c 
where
  source_type = 'deployed' 
group by
  name,
  uid,
  path,
  start_line,
  template;
```

```sql+sqlite
select
  name,
  coalesce(uid, path || ':' || start_line) as uid,
  group_concat(json_extract(c.value, '$.name')) as containers,
  group_concat(json_extract(c.value, '$.image')) as images,
  json_extract(template, '$.metadata.labels') as pod_labels 
from
  kubernetes_pod_template,
  json_each(json_extract(template, '$.spec.containers')) as c 
where
  source_type = 'deployed' 
group by
  name,
  uid,
  path,
  start_line,
  template;
```

### List pod templates that have a container with profiling argument set to false
Determine the areas in which pod templates contain a container with a disabled profiling argument. This is useful for ensuring optimal performance and security within your Kubernetes environment.

```sql+postgres
select
  name as pod_template_name,
  namespace,
  template -> 'metadata' ->> 'name' as pod_name,
  c ->> 'name' as pod_container_name,
  c ->> 'image' as pod_container_image 
from
  kubernetes_pod_template,
  jsonb_array_elements(template -> 'spec' -> 'containers') as c 
where
  (
    c -> 'command'
  )
  @ > '["kube-scheduler"]' 
  and 
  (
    c -> 'command'
  )
  @ > '["--profiling=false"]';
```

```sql+sqlite
Error: The corresponding SQLite query is unavailable.
```

### List manifest pod template resources
This query allows you to analyze the resources of manifest pod templates in a Kubernetes cluster. It's particularly useful for gaining insights into the containers, images, and pod labels associated with each pod template, helping to enhance management and optimization of the cluster.

```sql+postgres
select
  name,
  coalesce(uid, concat(path, ':', start_line)) as uid,
  array_agg(c ->> 'name') as containers,
  array_agg(c ->> 'image') as images,
  template -> 'metadata' -> 'labels' as pod_labels 
from
  kubernetes_pod_template,
  jsonb_array_elements(template -> 'spec' -> 'containers') as c 
where
  path is not null 
group by
  name,
  uid,
  path,
  start_line,
  template;
```

```sql+sqlite
select
  name,
  coalesce(uid, path || ':' || start_line) as uid,
  group_concat(json_extract(c.value, '$.name')) as containers,
  group_concat(json_extract(c.value, '$.image')) as images,
  json_extract(template, '$.metadata.labels') as pod_labels 
from
  kubernetes_pod_template,
  json_each(json_extract(template, '$.spec.containers')) as c 
where
  path is not null 
group by
  name,
  uid,
  path,
  start_line,
  template;
```