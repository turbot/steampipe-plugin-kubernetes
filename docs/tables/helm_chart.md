# Table: helm_chart

A Helm chart is a collection of files that describe a set of Kubernetes resources and their dependencies. It provides a way to package, version, and deploy these resources in a repeatable way. Charts are designed to be reusable and configurable, allowing you to deploy applications with different settings and configurations.

## Examples

### Basic info

```sql
select
  name,
  api_version,
  version,
  deprecated,
  type,
  description
from
  helm_chart;
```

### List all deployed charts

```sql
select
  hc.name as chart_name,
  hc.type as chart_type,
  hc.version as chart_version,
  hr.name as release_name,
  hr.status as release_status,
  hr.version as deployment_version,
  hr.first_deployed as deployed_at,
  hr.config as deployment_config
from
  kubernetes.helm_chart as hc
  left join kubernetes.helm_release as hr
    on hc.name = hr.chart_name
where
  hr.status = 'deployed';
```

### List all deprecated charts

```sql
select
  name,
  api_version,
  version,
  type,
  description,
  app_version
from
  helm_chart
where
  deprecated;
```

### List application type charts

```sql
select
  name,
  api_version,
  version,
  type,
  description,
  app_version
from
  helm_chart
where
  type = 'application';
```
