---
organization: Turbot
category: ["software development"]
icon_url: "/images/plugins/turbot/kubernetes.svg"
brand_color: "#326CE5"
display_name: "Kubernetes"
short_name: "kubernetes"
description: "Steampipe plugin for Kubernetes components."
og_description: "Query Kubernetes with SQL! Open source CLI. No DB required."
og_image: "/images/plugins/turbot/kubernetes-social-graphic.png"
engines: ["steampipe", "sqlite", "postgres", "export"]
---

# Kubernetes + Steampipe

[Steampipe](https://steampipe.io) is an open-source zero-ETL engine to instantly query cloud APIs using SQL.

[Kubernetes](https://kubernetes.io) is an open-source system for automating deployment, scaling, and management of containerized applications.

The Kubernetes plugin makes it simpler to query the variety of Kubernetes resources deployed in a Kubernetes cluster using [Steampipe](https://steampipe.io).

Apart from querying the deployed resources, the plugin also supports scanning [Kubernetes manifest files](#manifest-files) from different sources, parsing the configured [Helm charts](#helm-charts) and scanning all the templates to get the list of Kubernetes resources.

## Documentation

- **[Table definitions & examples →](/plugins/turbot/kubernetes/tables)**

## Get Started

### Install

Download and install the latest Kubernetes plugin:

```bash
steampipe plugin install kubernetes
```

### Configuration

Installing the latest Kubernetes plugin will create a config file (`~/.steampipe/config/kubernetes.spc`) with a single connection named `kubernetes`:

```hcl
connection "kubernetes" {
  plugin         = "kubernetes"
  config_path    = "~/.kube/config"
  config_context = "myCluster"
  source_types   = ["deployed"]
}
```

For a full list of configuration arguments, please see the [default configuration file](https://github.com/turbot/steampipe-plugin-kubernetes/blob/main/config/kubernetes.spc).

### Run a Query

Run steampipe:

```shell
steampipe query
```

List all pods:

```sql
select
  name,
  namespace,
  phase,
  creation_timestamp,
  pod_ip
from
  kubernetes_pod;
```

```sh
+-----------------------------------------+-------------+-----------+---------------------+-----------+
| name                                    | namespace   | phase     | creation_timestamp  | pod_ip    |
+-----------------------------------------+-------------+-----------+---------------------+-----------+
| metrics-server-86cbb8457f-bf8dm         | kube-system | Running   | 2021-06-11 14:21:48 | 10.42.0.5 |
| coredns-7448499f4d-klb8l                | kube-system | Running   | 2021-06-11 14:21:48 | 10.42.0.6 |
| helm-install-traefik-crd-hb87d          | kube-system | Succeeded | 2021-06-11 14:21:48 | 10.42.0.3 |
| local-path-provisioner-5ff76fc89d-c9hnm | kube-system | Running   | 2021-06-11 14:21:48 | 10.42.0.2 |
+-----------------------------------------+-------------+-----------+---------------------+-----------+
```

## Configuring Kubernetes Cluster Credentials

By default, the plugin will use the kubeconfig in `~/.kube/config` with the current context. If using the default kubectl CLI configurations, the kubeconfig will be in this location and the Kubernetes plugin connections will work by default.

You can also set the kubeconfig file path and context with the `config_path` and `config_context` config arguments respectively.

This plugin supports querying Kubernetes clusters using [OpenID Connect](https://kubernetes.io/docs/reference/access-authn-authz/authentication/#openid-connect-tokens) (OIDC) authentication. No extra configuration is required to query clusters using OIDC.

If no kubeconfig file is found, then the plugin will [attempt to access the API from within a pod](https://kubernetes.io/docs/tasks/run-application/access-api-from-pod/#accessing-the-api-from-within-a-pod) using the service account Kubernetes gives to pods.

### Single Context Connection

Each connection is implemented as a distinct [Postgres schema](https://www.postgresql.org/docs/current/ddl-schemas.html). As such, you can use qualified table names to query a specific connection:

```sql
select * from kubernetes_cluster_aks.kubernetes_namespace;
```

Alternatively, you can use an unqualified name and it will be resolved according to the [Search Path](https://steampipe.io/docs/using-steampipe/managing-connections#setting-the-search-path):

```sql
select * from kubernetes_namespace;
```

### Multiple Context Connections

You may create multiple Kubernetes connections. Example of creating multiple connections per the same `kubeconfig` file and different contexts:

```hcl
connection "kubernetes_cluster_aks" {
  plugin         = "kubernetes"
  config_path    = "~/.kube/config"
  config_context = "myAKSCluster"
  source_types   = ["deployed"]
}

connection "kubernetes_cluster_eks" {
  plugin         = "kubernetes"
  config_path    = "~/.kube/config"
  config_context = "arn:aws:eks:us-east-1:123456789012:cluster/myEKSCluster"
  source_types   = ["deployed"]
}
```

You can create multi-subscription connections by using an [**aggregator** connection](https://steampipe.io/docs/using-steampipe/managing-connections#using-aggregators). Aggregators allow you to query data from multiple connections for a plugin as if they are a single connection:

```hcl
connection "kubernetes_all" {
  plugin      = "kubernetes"
  type        = "aggregator"
  connections = ["kubernetes_cluster_aks", "kubernetes_cluster_eks"]
}
```

Querying tables from this connection will return results from the `kubernetes_cluster_aks` and `kubernetes_cluster_eks` connections:

```sql
select * from kubernetes_all.kubernetes_namespace;
```

Steampipe also supports the `*` wildcard in the connection names. For example, to aggregate all the Kubernetes plugin connections whose names begin with `kubernetes_`:

```hcl
connection "kubernetes_all" {
  type        = "aggregator"
  plugin      = "kubernetes"
  connections = ["kubernetes_*"]
}
```

### Custom Resource Definitions

Kubernetes also supports creating [Custom Resource Definitions](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/#customresourcedefinitions) with a name and schema that you specify in the `custom_resource_tables` configuration argument which allows you to extend Kubernetes capabilities by adding any kind of API object useful for your application.

Refer [Custom Resource](https://hub.steampipe.io/plugins/turbot/kubernetes/tables/kubernetes_%7Bcustom_resource_singular_name%7D#table-kubernetes_custom_resource_singular_name) table to get more information about how the plugin handles CRD and custom resources.

## Manifest Files

It is not necessary to always have a Kubernetes cluster online. The plugin also supports reading the manifest files from various sources (e.g., [Local files](#configuring-local-file-paths), [Git](#configuring-remote-git-repository-urls), [S3](#configuring-s3-urls) etc.) and make it available to query the resources using the respective Kubernetes resource tables.

To query resources from the manifest files, set the `manifest_file_paths` argument to point the sources where the manifest files are located.

Manifest file paths may [include wildcards](https://pkg.go.dev/path/filepath#Match) and support `**` for recursive matching. For example:

```hcl
connection "kubernetes" {
  plugin = "kubernetes"

  manifest_file_paths = [
    "*.yml",
    "**/*.json",
    "~/*.yaml",
    "github.com/GoogleCloudPlatform/microservices-demo//release//kubernetes-*.yaml",
    "github.com/GoogleCloudPlatform/microservices-demo//release//kubernetes-manifests.yaml",
    "s3::https://bucket.s3.us-east-1.amazonaws.com/test_folder//*.yml"
  ]
  
  source_types = ["manifest"]
}
```

**Note**: If any path matches on `*` without `.yml` or `.yaml` or `.json`, all files (including non-Kubernetes manifest files) in the directory will be matched, which may cause errors if incompatible file types exist.

By default the plugin always lists the resources deployed in the current Kubernetes cluster context. If you want to restrict this behavior to read resource configurations from the configured manifest files only, add the `source_types` argument to the config and set the value to `manifest`. For example:

```hcl
connection "kubernetes" {
  plugin = "kubernetes"

  manifest_file_paths = [ ... ]

  source_types = ["manifest"]
}
```

### Configuring Local File Paths

You can define a list of local directory paths to search for Kubernetes manifest files. Paths are resolved relative to the current working directory. For example:

- `*.yml` or `*.yaml` or `*.json` matches all Kubernetes manifest files in the CWD.
- `**/*.yml` or `**/*.yaml` or `**/*.json` matches all Kubernetes manifest files in the CWD and all sub-directories.
- `../*.yml` or `../*.yaml` or `../*.json` matches all Kubernetes manifest files in the CWD's parent directory.
- `steampipe*.yml` or `steampipe*.yaml` or `steampipe*.json` matches all Kubernetes manifest files starting with "steampipe" in the CWD.
- `/path/to/dir/*.yml` or `/path/to/dir/*.yaml` or `/path/to/dir/*.json` matches all Kubernetes manifest files in a specific directory.
- `~/*.yml` or `~/*.yaml` or `~/*.json` matches all Kubernetes manifest files in the home directory.
- `~/**/*.yml` or `~/**/*.yaml` or `~/**/*.json` matches all Kubernetes manifest files recursively in the home directory.
- `/path/to/dir/main.yml` or `/path/to/dir/main.yaml` or `/path/to/dir/main.json` matches a specific file.

```hcl
connection "kubernetes" {
  plugin = "kubernetes"

  manifest_file_paths = [ "*.yml", "*.yaml", "*.json", "/path/to/dir/main.yml" ]
  
  source_types = ["manifest"]
}
```

### Configuring Remote Git Repository URLs

You can also configure `manifest_file_paths` with any Git remote repository URLs, e.g., GitHub, BitBucket, GitLab. The plugin will then attempt to retrieve any Kubernetes manifest files from the remote repositories.

For example:

- `github.com/GoogleCloudPlatform/microservices-demo//release//kubernetes-manifests.yaml` matches the file `kubernetes-manifests.yaml` in the specified repository.
- `github.com/GoogleCloudPlatform/microservices-demo//release//*.yaml` matches all top-level Kubernetes manifest files in the specified repository.
- `github.com/GoogleCloudPlatform/microservices-demo//release//**/*.yaml` matches all Kubernetes manifest files in the specified repository and all subdirectories.

You can specify a subdirectory after a double-slash (`//`) if you want to download only a specific subdirectory from a downloaded directory.
Similarly, you can define a list of GitLab and BitBucket URLs to search for Kubernetes manifest files:

```hcl
connection "kubernetes" {
  plugin = "kubernetes"

  manifest_file_paths = [ "bitbucket.org/atlassian/kubectl-run//test/kustomization//deploy.yml" ]
  
  source_types = ["manifest"]
}
```

### Configuring S3 URLs

You can also query all Kubernetes manifest files stored inside an S3 bucket (public or private) using the bucket URL.

#### Accessing a Private Bucket

In order to access your files in a private S3 bucket, you will need to configure your credentials. You can use your configured AWS profile from local `~/.aws/config`, or pass the credentials using the standard AWS environment variables, e.g., `AWS_PROFILE`, `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `AWS_REGION`.

We recommend using AWS profiles for authentication.

**Note:** Make sure that `region` is configured in the config. If not set in the config, `region` will be fetched from the standard environment variable `AWS_REGION`.

You can also authenticate your request by setting the AWS profile and region in `manifest_file_paths`. For example:

```hcl
connection "kubernetes" {
  plugin = "kubernetes"

  manifest_file_paths = [
    "s3::https://bucket-2.s3.us-east-1.amazonaws.com//*.yml?aws_profile=<AWS_PROFILE>",
    "s3::https://bucket-2.s3.us-east-1.amazonaws.com/test_folder//*.yaml?aws_profile=<AWS_PROFILE>"
  ]
  
  source_types = ["manifest"]
}
```

**Note:**

In order to access the bucket, the IAM user or role will require the following IAM permissions:

- `s3:ListBucket`
- `s3:GetObject`
- `s3:GetObjectVersion`

If the bucket is in another AWS account, the bucket policy will need to grant access to your user or role. For example:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "ReadBucketObject",
      "Effect": "Allow",
      "Principal": {
        "AWS": "arn:aws:iam::123456789012:user/YOUR_USER"
      },
      "Action": ["s3:ListBucket", "s3:GetObject", "s3:GetObjectVersion"],
      "Resource": ["arn:aws:s3:::test-bucket1", "arn:aws:s3:::test-bucket1/*"]
    }
  ]
}
```

#### Accessing a Public Bucket

Public access granted to buckets and objects through ACLs and bucket policies allows any user access to data in the bucket. We do not recommend making S3 buckets public, but if there are specific objects you'd like to make public, please see [How can I grant public read access to some objects in my Amazon S3 bucket?](https://aws.amazon.com/premiumsupport/knowledge-center/read-access-objects-s3-bucket/).

You can query any public S3 bucket directly using the URL without passing credentials. For example:

```hcl
connection "kubernetes" {
  plugin = "kubernetes"

  manifest_file_paths = [
    "s3::https://bucket-1.s3.us-east-1.amazonaws.com/test_folder//*.yml",
    "s3::https://bucket-2.s3.us-east-1.amazonaws.com/test_folder//**/*.yaml"
  ]
  
  source_types = ["manifest"]
}
```

## Helm Charts

The plugin also supports configuring Helm charts and allows the users to query the metadata, templates, and deployed versions of the configured charts using Steampipe. It also renders the templates and returns the resulting manifest after communicating with the kubernetes cluster without actually creating any resources on the cluster.

The plugin supports the following `helm_*` tables to query Helm configurations:

- [helm_chart](https://hub.steampipe.io/plugins/turbot/kubernetes/tables/helm_chart)
- [helm_release](https://hub.steampipe.io/plugins/turbot/kubernetes/tables/helm_release)
- [helm_template](https://hub.steampipe.io/plugins/turbot/kubernetes/tables/helm_template)
- [helm_template_rendered](https://hub.steampipe.io/plugins/turbot/kubernetes/tables/helm_template_rendered)
- [helm_value](https://hub.steampipe.io/plugins/turbot/kubernetes/tables/helm_value)

The plugin can also parse the configured Helm charts, render all its templates to Kubernetes manifests, and allow you to query the resource configurations (i.e. resources the chart will deploy when it is installed) using the respective `kubernetes_*` tables, which is particularly helpful while developing a new chart, making changes to the chart, debugging, and so on.

For example:

```hcl
connection "kubernetes" {
  plugin = "kubernetes"

  helm_rendered_charts = {
    "my-app-1" = {
      chart_path        = "~/charts/my-app-1"
      values_file_paths = ["~/value/file/for/my-app-1.yaml"]
    }
    "my-app-2" = {
      chart_path        = "~/charts/my-app-2"
      values_file_paths = [] # works with values from chart's default values.yaml file
    }
  }
  
  source_types = ["helm"]
}
```

`helm_rendered_charts` takes a map of chart configurations. It can have more than 1 chart based on the requirement:

- The above configuration has 2 charts: `my-app-1` and `my-app-2`. The name `my-app-1` and `my-app-2` are considered as release names.
- Every map should have a `chart_path` indicating the directory where the chart is located.
- The map can have an optional `values_file_paths` argument that overrides value files for rendering the templates. The `values_file_paths` can have more than 1 override value file reference. The plugin reads values from all of those files, and uses the resultant value to render the templates. By default, the plugin uses `values.yaml` if no additional value files are passed.

## Get Involved

- Open source: https://github.com/turbot/steampipe-plugin-kubernetes
- Community: [Join #steampipe on Slack →](https://turbot.com/community/join)
