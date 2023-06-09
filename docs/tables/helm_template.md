# Table: helm_template

A template is a file that defines a Kubernetes manifest in a way that is generic enough to allow customization at the time of installation. It can reference variables and functions that are provided by Helm or defined in the chart.

**Note:** The table `helm_template` will show the raw template as defined in the file. To list the fully rendered templates, use table `helm_template_rendered`.

## Examples

### Basic info

```sql
select
  chart_name,
  path,
  raw
from
  helm_template;
```

### List templates defined for a specific chart

```sql
select
  chart_name,
  path,
  raw
from
  helm_template
where
  chart_name = 'redis';
```
