# Table: helm_template_rendered

A template is a file that defines a Kubernetes manifest in a way that is generic enough to allow customization at the time of installation. It can reference variables and functions that are provided by Helm or defined in the chart.

During the installation process, Helm takes the template files in the chart and renders them using the values provided by the user or the defaults defined in the chart's values.yaml file.

## Examples

### Basic info

```sql
select
  name,
  chart_name,
  source_type,
  rendered
from
  helm_template_rendered;
```

### List templates defined for a specific chart

```sql
select
  name,
  chart_name,
  source_type,
  rendered
from
  helm_template_rendered
where
  chart_name = 'steampipe';
```
