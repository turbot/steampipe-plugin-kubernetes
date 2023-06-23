# Table: helm_value

Values in Helm-packed applications dictate the configuration of an application. Every Helm charts have an associated values.yaml file where the default configuration is defined. It is a source of content for the Values built-in object offered by Helm templates.

By design, applications can ship with default values.yaml file tuned for production deployments. Also, considering the multiple environments, it may have different configurations. To override the default value, it is not necessary to change the default values.yaml, but you can refer to the override value files from which it takes the configuration. For example:

Let's say you have two different environments for maintaining your app: dev and prod. And, you have a helm chart with 2 different set of values for your environments. For example:

```hcl
connection "kubernetes" {
  plugin = "kubernetes"

  helm_rendered_charts = {
    "my-app-dev" = {
      chart_path        = "~/charts/my-app"
      values_file_paths = "~/value/file/for/dev.yaml"
    }
    "my-app-prod" = {
      chart_path        = "~/charts/my-app"
      values_file_paths = "~/value/file/for/prod.yaml"
    }
  }
}
```

The table `helm_value` lists the values from the chart's default values.yaml file, as well as it lists the values from the files that are provided to override the default configuration.

## Examples

### List values configured in the default values.yaml file of a specific chart

```sql
select
  path,
  key_path,
  value,
  start_line,
  start_column
from
  helm_value
where
  path = '~/charts/my-app/values.yaml'
order by
  start_line;
```

### List values from a specific override file

```sql
select
  path,
  key_path,
  value,
  start_line,
  start_column
from
  helm_value
where
  path = '~/value/file/for/dev.yaml'
order by
  start_line;
```
