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
