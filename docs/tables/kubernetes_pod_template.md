# Table: kubernetes_pod_template

A PodTemplate is a Kubernetes resource that defines the desired specification for Pods created or managed by various controllers, such as Deployments, StatefulSets, and CronJobs. It allows you to define a reusable configuration for Pods, reducing duplication and simplifying updates across multiple resources.

A Pod is a group of one or more containers, with shared storage and network resources, and a specification for how to run the containers.

## Examples

### Basic info

```sql
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

### List pod templates with privileged pod containers

```sql
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

### List pod templates with pod has access to the host process ID, IPC, or network namespace

```sql
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

### `kubectl get podtemplates` columns

```sql
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

### List pod templates with pod having a container with --profiling argument is set to false

```sql
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
  (c -> 'command') @> '["kube-scheduler"]'
  and (c -> 'command') @> '["--profiling=false"]';
```

### List manifest resources

```sql
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
