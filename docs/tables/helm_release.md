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
  first_deployed,
  chart_name
from
  helm_release;
```

### List kubernetes deployment resources deployed using a specific chart

```sql
select
  d.name as deployment_name,
  d.namespace,
  d.uid,
  r.name as release_name,
  r.chart_name,
  r.version as release_version,
  r.first_deployed as deployed_at
from
  kubernetes_deployment as d
  left join helm_release as r
    on (
      d.labels ->> 'app.kubernetes.io/managed-by' = 'Helm'
      and d.labels ->> 'app.kubernetes.io/instance' = r.name
    )
where
  r.chart_name = 'ingress-nginx'
  and d.source_type = 'deployed';
```

### List all deployed releases of a specific chart

```sql
select
  name,
  namespace,
  version,
  status,
  first_deployed,
  chart_name
from
  helm_release
where
  chart_name = 'ingress-nginx';
```

### List releases from a specific namespace

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
