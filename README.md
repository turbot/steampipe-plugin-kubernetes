![image](https://hub.steampipe.io/images/plugins/turbot/kubernetes-social-graphic.png)

# Kubernetes Plugin for Steampipe

Use SQL to query Kubernetes components.

- **[Get started →](https://hub.steampipe.io/plugins/turbot/kubernetes)**
- Documentation: [Table definitions & examples](https://hub.steampipe.io/plugins/turbot/kubernetes/tables)
- Community: [Slack Channel](https://steampipe.io/community/join)
- Get involved: [Issues](https://github.com/turbot/steampipe-plugin-kubernetes/issues)

## Quick start

### Install

Download and install the latest Kubernetes plugin:

```bash
steampipe plugin install kubernetes
```

Configure your [config file](https://hub.steampipe.io/plugins/turbot/kubernetes#configuration).

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

  # Paths is a list of locations to search for Kubernetes manifest files
  # Wildcard based searches are supported, including recursive searches
  # Local paths are resolved relative to the current working directory (CWD)

  # For example:
  #  - "*.yml" matches all Kubernetes manifest files in the CWD
  #  - "**/*.yml" matches all Kubernetes manifest files in the CWD and all sub-directories
  #  - "../*.yml" matches all Kubernetes manifest files in the CWD's parent directory
  #  - "steampipe*.yml" matches all Kubernetes manifest files starting with "steampipe" in the CWD
  #  - "/path/to/dir/*.yml" matches all Kubernetes manifest files in a specific directory
  #  - "/path/to/dir/main.yml" matches a specific file

  # If paths includes "*", all files (including non-kubernetes manifest files) in
  # the CWD will be matched, which may cause errors if incompatible file types exist

  # Defaults to CWD
  manifest_file_paths = [ "*.yml", "*.yaml" ]
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

Please see the [contribution guidelines](https://github.com/turbot/steampipe/blob/main/CONTRIBUTING.md) and our [code of conduct](https://github.com/turbot/steampipe/blob/main/CODE_OF_CONDUCT.md). All contributions are subject to the [Apache 2.0 open source license](https://github.com/turbot/steampipe-plugin-kubernetes/blob/main/LICENSE).

`help wanted` issues:

- [Steampipe](https://github.com/turbot/steampipe/labels/help%20wanted)
- [Kubernetes Plugin](https://github.com/turbot/steampipe-plugin-kubernetes/labels/help%20wanted)
