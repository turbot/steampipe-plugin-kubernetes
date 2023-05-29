# Table: helm_release

A Helm release is an instance of a chart running in a Kubernetes cluster. When you use the `helm install` command, it creates a release for the chart and generates a set of Kubernetes resources based on the chart's templates and the values provided.

## Examples

### Basic info

```sql
select
  name,
  namespace,
  version,
  status,
  last_deployed
from
  helm_release;
```

### List all failed releases

```sql
select
  name,
  namespace,
  version,
  status,
  last_deployed,
  description
from
  helm_release
where
  status = 'failed';
```

### List all unfinished releases

```sql
select
  name,
  namespace,
  version,
  status,
  last_deployed,
  description
from
  helm_release
where
  status = 'pending';
```

### List releases updated in last 3 days

```sql
select
  name,
  namespace,
  version,
  status,
  last_deployed,
  description
from
  helm_release
where
  last_deployed > (now() - interval '3 days');
```

### Get a specific release

```sql
select
  name,
  namespace,
  version,
  status,
  last_deployed,
  description
from
  helm_release
where
  name = 'brigade-github-app-1683552635';
```

### List releases for a specific namespace

```sql
select
  name,
  namespace,
  version,
  status,
  last_deployed,
  description
from
  helm_release
where
  namespace = 'steampipe';
```

### List releases information of a deployment

```sql
select
  d.name as deployment_name,
  d.namespace,
  d.uid,
  h.name as release_name,
  h.version as release_version,
  h.last_deployed as release_last_updated
from
  kubernetes_deployment as d
  left join helm_release as h on (d.labels ->> 'release' = h.name and d.source_type = 'deployed');
```
