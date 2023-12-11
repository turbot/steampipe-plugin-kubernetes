![image](https://hub.steampipe.io/images/plugins/turbot/kubernetes-social-graphic.png)

# Kubernetes Plugin for Steampipe

Use SQL to query Kubernetes components.

Apart from querying the deployed resources, the plugin also supports scanning the [Kubernetes manifest files](https://hub.steampipe.io/plugins/turbot/kubernetes#supported-manifest-file-path-formats) from different sources, parsing the configured [Helm charts](https://hub.steampipe.io/plugins/turbot/kubernetes#helm-configuration) and scanning all the templates to get the list of Kubernetes resources.

- **[Get started →](https://hub.steampipe.io/plugins/turbot/kubernetes)**
- Documentation: [Table definitions & examples](https://hub.steampipe.io/plugins/turbot/kubernetes/tables)
- Community: [Join #steampipe on Slack →](https://turbot.com/community/join)
- Get involved: [Issues](https://github.com/turbot/steampipe-plugin-kubernetes/issues)

## Quick start

### Install

Download and install the latest Kubernetes plugin:

```bash
steampipe plugin install kubernetes
```

Installing the latest Kubernetes plugin will create a config file (`~/.steampipe/config/kubernetes.spc`) with a single connection named `kubernetes`:

```hcl
connection "kubernetes" {
  plugin = "kubernetes"

  # By default, the plugin will use credentials in "~/.kube/config" with the current context.
  # OpenID Connect (OIDC) authentication is supported without any extra configuration.
  # The kubeconfig path and context can also be specified with the following config arguments:

  # Specify the file path to the kubeconfig.
  # Can also be set with the "KUBECONFIG" or "KUBE_CONFIG_PATH" environment variables. Plugin will prioritize KUBECONFIG if both are available.
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

  # Specify the source(s) of the resource(s). Possible values: `deployed`, `helm` and `manifest`.
  # Defaults to all possible values. Set the argument to override the default value.
  # If `deployed` is contained in the value, tables will show all the deployed resources.
  # If `helm` is contained in the value, tables will show resources from the configured helm charts.
  # If `manifest` is contained in the value, tables will show all the resources from the kubernetes manifest. Make sure that the `manifest_file_paths` arg is set.
  # source_types = ["deployed", "helm", "manifest"]

  # Manifest File Configuration

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

  # Helm configuration

  # A map for Helm charts along with the path to the chart directory and the paths of the value override files (if any).
  # Every map should have chart_path defined, and the values_file_paths is optional.
  # You can define multiple charts in the config.
  # helm_rendered_charts = {
  #   "chart_name" = {
  #     chart_path        = "/path/to/chart/dir"
  #     values_file_paths = ["/path/to/value/override/files.yaml"]
  #   }
  # }
}
```

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

## Engines

This plugin is available for the following engines:

| Engine        | Description
|---------------|------------------------------------------
| [Steampipe](https://steampipe.io/docs) | The Steampipe CLI exposes APIs and services as a high-performance relational database, giving you the ability to write SQL-based queries to explore dynamic data. Mods extend Steampipe's capabilities with dashboards, reports, and controls built with simple HCL. The Steampipe CLI is a turnkey solution that includes its own Postgres database, plugin management, and mod support.
| [Postgres FDW](https://steampipe.io/docs/steampipe_postgres/index) | Steampipe Postgres FDWs are native Postgres Foreign Data Wrappers that translate APIs to foreign tables. Unlike Steampipe CLI, which ships with its own Postgres server instance, the Steampipe Postgres FDWs can be installed in any supported Postgres database version.
| [SQLite Extension](https://steampipe.io/docs//steampipe_sqlite/index) | Steampipe SQLite Extensions provide SQLite virtual tables that translate your queries into API calls, transparently fetching information from your API or service as you request it.
| [Export](https://steampipe.io/docs/steampipe_export/index) | Steampipe Plugin Exporters provide a flexible mechanism for exporting information from cloud services and APIs. Each exporter is a stand-alone binary that allows you to extract data using Steampipe plugins without a database.
| [Turbot Pipes](https://turbot.com/pipes/docs) | Turbot Pipes is the only intelligence, automation & security platform built specifically for DevOps. Pipes provide hosted Steampipe database instances, shared dashboards, snapshots, and more.

## Developing

Prerequisites:

- [Steampipe](https://steampipe.io/downloads)
- [Golang](https://golang.org/doc/install)

Clone:

```sh
git clone https://github.com/turbot/steampipe-plugin-kubernetes.git
cd steampipe-plugin-kubernetes
```

Build, which automatically installs the new version to your `~/.steampipe/plugins` directory:

```sh
make
```

Configure the plugin:

```shell
cp config/* ~/.steampipe/config
vi ~/.steampipe/config/kubernetes.spc
```

Try it!

```shell
steampipe query
> .inspect kubernetes
```

Further reading:

- [Writing plugins](https://steampipe.io/docs/develop/writing-plugins)
- [Writing your first table](https://steampipe.io/docs/develop/writing-your-first-table)

## Contributing

Please see the [contribution guidelines](https://github.com/turbot/steampipe/blob/main/CONTRIBUTING.md) and our [code of conduct](https://github.com/turbot/steampipe/blob/main/CODE_OF_CONDUCT.md). Contributions to the plugin are subject to the [Apache 2.0 open source license](https://github.com/turbot/steampipe-plugin-kubernetes/blob/main/LICENSE). Contributions to the plugin documentation are subject to the [CC BY-NC-ND license](https://github.com/turbot/steampipe-plugin-kubernetes/blob/main/docs/LICENSE).
`help wanted` issues:

- [Steampipe](https://github.com/turbot/steampipe/labels/help%20wanted)
- [Kubernetes Plugin](https://github.com/turbot/steampipe-plugin-kubernetes/labels/help%20wanted)
