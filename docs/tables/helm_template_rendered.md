# Table: helm_template_rendered

A template is a file that defines a Kubernetes manifest in a way that is generic enough to allow customization at the time of installation. It can reference variables and functions that are provided by Helm or defined in the chart.

During the installation process, Helm takes the template files in the chart and renders them using the values provided by the user or the defaults defined in the chart's values.yaml file.

The table `helm_template_rendered` reads the templates defined in the chart, the value files provided in the config, renders the templates and returns the resulting manifest after communicating with the kubernetes cluster without actually creating any resources on the cluster.

## Examples

### List fully rendered kubernetes resource templates defined in a chart

```sql
select
  path,
  source_type,
  rendered
from
  helm_template_rendered
where
  chart_name = 'redis';
```

### List fully rendered kubernetes resource templates for different environments

Let's say you have two different environments for maintaining your app: dev and prod. And, you have a helm chart with 2 different set of values for your environments. For example:

```hcl
connection "kubernetes" {
  plugin = "kubernetes"

  helm_rendered_charts = {
    "my-app-dev" = {
      chart_path        = "~/charts/my-app"
      values_file_paths = ["~/value/file/for/dev.yaml"]
    }
    "my-app-prod" = {
      chart_path        = "~/charts/my-app"
      values_file_paths = ["~/value/file/for/prod.yaml"]
    }
  }
}
```

In both case, it is using same chart with a different set of values.

To list the kubernetes resource configurations defined for the dev environment, you can simply run the below query:

```sql
select
  chart_name,
  path,
  source_type,
  rendered
from
  helm_template_rendered
where
  source_type = 'helm_rendered:my-app-dev';
```

Similarly, to query the kubernetes resource configurations for prod,

```sql
select
  chart_name,
  path,
  source_type,
  rendered
from
  helm_template_rendered
where
  source_type = 'helm_rendered:my-app-prod';
```
