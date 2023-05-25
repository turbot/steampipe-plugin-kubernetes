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
---

# Kubernetes + Steampipe

[Steampipe](https://steampipe.io) is an open source CLI to instantly query cloud APIs using SQL.

[Kubernetes](https://kubernetes.io) is an open-source system for automating deployment, scaling, and management of containerized applications.

For example:

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

## Documentation

- **[Table definitions & examples →](/plugins/turbot/kubernetes/tables)**

## Get started

### Install

Download and install the latest Kubernetes plugin:

```bash
steampipe plugin install kubernetes
```

### Configuration

Installing the latest kubernetes plugin will create a config file (`~/.steampipe/config/kubernetes.spc`) with a single connection named `kubernetes`:

```hcl
connection "kubernetes" {
  plugin = "kubernetes"

  # By default, the plugin will use credentials in "~/.kube/config" with the current context.
  # OpenID Connect (OIDC) authentication is supported without any extra configuration.
  # The kubeconfig path and context can also be specified with the following config arguments:

  # Specify the file path to the kubeconfig.
  # Can also be set with the "KUBE_CONFIG_PATHS" or "KUBERNETES_MASTER" environment variables.
  # config_path = "~/.kube/config"

  # Specify a context other than the current one.
  # config_context = "minikube"

  # List of custom resources that will be created as dynamic tables.
  # No dynamic tables will be created if this arg is empty or not set.
  # Wildcard based searches are supported.

  # For example:
  #  - "*" matches all custom resources available
  #  - "*.storage.k8s.io" matches all custom resources in the storage.k8s.io group
  #  - "certificates.cert-manager.io" matches a specific custom resource "certificates.cert-manager.io"
  #  - "backendconfig" matches the singular name "backendconfig" in any group

  # Defaults to all custom resources
  custom_resource_tables = ["*"]

  # If no kubeconfig file can be found, the plugin will attempt to use the service account Kubernetes gives to pods.
  # This authentication method is intended for clients that expect to be running inside a pod running on Kubernetes.

  # Manifest file paths is a list of locations to search for Kubernetes manifest files
  # Manifest file paths can be configured with a local directory, a remote Git repository URL, or an S3 bucket URL
  # Refer https://hub.steampipe.io/plugins/turbot/kubernetes#supported-path-formats for more information
  # Wildcard based searches are supported, including recursive searches
  # Local paths are resolved relative to the current working directory (CWD)

  # For example:
  #  - "*.yml" or "*.yaml" or "*.json" matches all Kubernetes manifest files in the CWD
  #  - "**/*.yml" or "**/*.yaml" or "**/*.json" matches all Kubernetes manifest files in the CWD and all sub-directories
  #  - "../*.yml" or "../*.yaml" or "../*.json" matches all Kubernetes manifest files in the CWD's parent directory
  #  - "steampipe*.yml" or "steampipe*.yaml" or "steampipe*.json" matches all Kubernetes manifest files starting with "steampipe" in the CWD
  #  - "/path/to/dir/*.yml" or "/path/to/dir/*.yaml" or "/path/to/dir/*.json" matches all Kubernetes manifest files in a specific directory
  #  - "/path/to/dir/main.yml" or "/path/to/dir/main.yaml" or "/path/to/dir/main.json" matches a specific file

  # If the given paths includes "*", all files (including non-kubernetes manifest files) in
  # the CWD will be matched, which may cause errors if incompatible file types exist

  # Defaults to CWD
  # manifest_file_paths = [ "*.yml", "*.yaml", "*.json" ]

  # Specify the source of the resource. Possible values: `deployed`, `manifest`, and `all`.
  # Default set to `all`. Set the argument to override the default value.
  # If the value is set to `deployed`, tables will show all the deployed resources.
  # If set to `manifest`, tables will show all the resources from the kubernetes manifest. Make sure that the `manifest_file_paths` arg is set.
  # If `all`, tables will show all the deployed and manifest resources.
  # source_type = "all"
}
```

- `config_context` - (Optional) The kubeconfig context to use. If not set, the current context will be used.
- `config_path` - (Optional) The kubeconfig file path. If not set, the plugin will check `~/.kube/config`. Can also be set with the `KUBE_CONFIG_PATHS` or `KUBERNETES_MASTER` environment variables.
- `custom_resource_tables` - (Optional) The custom resources to create as dynamic tables. If set to empty or not set, the plugin will not create any dynamic tables.
- `manifest_file_paths` - (Optional) A list of locations to search for Kubernetes manifest files. If set, the plugin will read the resource configurations from the configured paths and list the resources in the respective tables.
- `source_type` - (Optional) Specify the source of the resource. Default set to `all`. The possible values are: `deployed`, `manifest`, and `all`.

  - If the value is set to `deployed`, tables will show all the deployed resources.
  - If set to `manifest`, tables will show all the resources from the kubernetes manifest. Make sure that the `manifest_file_paths` arg is set.

## Configuring Kubernetes Credentials

By default, the plugin will use the kubeconfig in `~/.kube/config` with the current context. If using the default kubectl CLI configurations, the kubeconfig will be in this location and the Kubernetes plugin connections will work by default.

You can also set the kubeconfig file path and context with the `config_path` and `config_context` config arguments respectively.

This plugin supports querying Kubernetes clusters using [OpenID Connect](https://kubernetes.io/docs/reference/access-authn-authz/authentication/#openid-connect-tokens) (OIDC) authentication. No extra configuration is required to query clusters using OIDC.

If no kubeconfig file is found, then the plugin will [attempt to access the API from within a pod](https://kubernetes.io/docs/tasks/run-application/access-api-from-pod/#accessing-the-api-from-within-a-pod) using the service account Kubernetes gives to pods.

## Multiple Context Connections

You may create multiple Kubernetes connections. Example of creating multiple connections per the same `kubeconfig` file and different contexts:

```hcl
connection "kubernetes_all" {
  type        = "aggregator"
  plugin      = "kubernetes"
  connections = ["kubernetes_*"]
}

connection "kubernetes_cluster_aks" {
  plugin          = "kubernetes"
  config_path = "~/.kube/config"
  config_context = "myAKSCluster"
}

connection "kubernetes_cluster_eks" {
  plugin = "kubernetes"
  config_path = "~/.kube/config"
  config_context = "arn:aws:eks:us-east-1:123456789012:cluster/myEKSCluster"
}
```

Each connection is implemented as a distinct [Postgres schema](https://www.postgresql.org/docs/current/ddl-schemas.html). As such, you can use qualified table names to query a specific connection:

```sql
select * from kubernetes_cluster_aks.kubernetes_namespace
```

Alternatively, you can use an unqualified name and it will be resolved according to the [Search Path](https://steampipe.io/docs/using-steampipe/managing-connections#setting-the-search-path):

```sql
select * from kubernetes_namespace
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
select * from kubernetes_all.kubernetes_namespace
```

Steampipe supports the `*` wildcard in the connection names. For example, to aggregate all the Kubernetes plugin connections whose names begin with `kubernetes_`:

```hcl
connection "kubernetes_all" {
  type        = "aggregator"
  plugin      = "kubernetes"
  connections = ["kubernetes_*"]
}
```

## Custom Resource Definitions

Kubernetes also supports creating [Custom Resource Definitions](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/#customresourcedefinitions) with a name and schema that you specify in the `custom_resource_tables` configuration argument which allows you to extend Kubernetes capabilities by adding any kind of API object useful for your application.

Refer [Custom Resource](https://hub.steampipe.io/plugins/turbot/kubernetes/tables/kubernetes_%7Bcustom_resource_singular_name%7D#table-kubernetes_custom_resource_singular_name) table to get more information about how the plugin handles CRD and custom resources.

## Supported Manifest File Path Formats

The `manifest_file_paths` config argument is flexible and can search for Kubernetes manifest files from several different sources, e.g., local directory paths, Git, S3.

The following sources are supported:

- [Local files](#configuring-local-file-paths)
- [Remote Git repositories](#configuring-remote-git-repository-urls)
- [S3](#configuring-s3-urls)

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
}
```

**Note**: If any path matches on `*` without `.yml` or `.yaml` or `.json`, all files (including non-Kubernetes manifest files) in the directory will be matched, which may cause errors if incompatible file types exist.

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
}
```

## Get involved

- Open source: https://github.com/turbot/steampipe-plugin-kubernetes
- Community: [Slack Channel](https://steampipe.io/community/join)
