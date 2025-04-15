## v1.4.0 [2025-04-15]

_Enhancements_

- Added `folder` metadata to the documentation of all the Kubernetes tables for improved organization on the Steampipe Hub. ([#294](https://github.com/turbot/steampipe-plugin-kubernetes/pull/294))
- Added `selector` and `selector_query` columns to `kubernetes_stateful_set` tables. ([#293](https://github.com/turbot/steampipe-plugin-kubernetes/pull/293)) (Thanks [@dongho-jung](https://github.com/dongho-jung) for the contribution!!) 

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.11.5](https://github.com/turbot/steampipe-plugin-sdk/releases/tag/v5.11.5) ([#294](https://github.com/turbot/steampipe-plugin-kubernetes/pull/294))

## v1.3.0 [2025-03-07]

_Bug fixes_

- Fixed the `containers_resources_requests_std` column of the `kubernetes_pod` table to correctly return data instead of an error when the memory field is defined in the unit of `millibyte`. ([#283](https://github.com/turbot/steampipe-plugin-kubernetes/pull/283)) (Thanks [@omer-do](https://github.com/omer-do) for the contribution!)

_Dependencies_

- Recompiled plugin with Go version `1.23.1`.
- Recompiled plugin with [steampipe-plugin-sdk v5.11.3](https://github.com/turbot/steampipe-plugin-sdk/blob/develop/CHANGELOG.md#v5113-2025-02-11) that addresses critical and high vulnerabilities in dependent packages.

## v1.2.1 [2025-02-12]

_Bug fixes_

- Fixed the `containers_resources_requests_std` column of the `kubernetes_pod` table to correctly return data instead of an error when the memory field is defined in the unit of `bytes`. ([#273](https://github.com/turbot/steampipe-plugin-kubernetes/pull/273)) (Thanks [@omer-do](https://github.com/omer-do) for the contribution!) 

## v1.2.0 [2024-11-29]

_Dependencies_

- Updated the plugin's Osquery extension build to use the latest versions of `Go`, `helm.sh/helm/v3` and `docker` that fixes several security vulnerabilities.

## v1.1.0 [2024-11-29]

_Deprecations_

- Deprecated `kubernetes_pod_security_policy` table due to the lack of API support in [Kubernetes v1.25](https://kubernetes.io/blog/2022/08/23/kubernetes-v1-25-release/#what-s-new-major-themes). ([#262](https://github.com/turbot/steampipe-plugin-kubernetes/pull/262))

_Dependencies_

- Recompiled plugin with [helm.sh/helm/v3 v3.14.2](https://github.com/helm/helm/releases/tag/v3.14.2) to fix security vulnerabilities. ([#262](https://github.com/turbot/steampipe-plugin-kubernetes/pull/262))

## v1.0.0 [2024-10-22]

There are no significant changes in this plugin version; it has been released to align with [Steampipe's v1.0.0](https://steampipe.io/changelog/steampipe-cli-v1-0-0) release. This plugin adheres to [semantic versioning](https://semver.org/#semantic-versioning-specification-semver), ensuring backward compatibility within each major version.

_Dependencies_

- Recompiled plugin with Go version `1.22`. ([#246](https://github.com/turbot/steampipe-plugin-kubernetes/pull/246))
- Recompiled plugin with [steampipe-plugin-sdk v5.10.4](https://github.com/turbot/steampipe-plugin-sdk/blob/develop/CHANGELOG.md#v5104-2024-08-29) that fixes logging in the plugin export tool. ([#246](https://github.com/turbot/steampipe-plugin-kubernetes/pull/246))

## v0.29.0 [2024-07-26]

_Enhancements_

- Updated the plugin to use fully qualified names while creating CRD tables to make sure that you can query the exact custom resource if there are multiple resources with the same singular name. ([#228](https://github.com/turbot/steampipe-plugin-kubernetes/pull/228)) (Thanks [@afarid](https://github.com/afarid) for the contribution!!)

_Bug fixes_

- Fixed the `kubernetes_custom_resource_definition` table to correctly list out all the CRDs on a cluuster instead of returning a truncated set. ([#235](https://github.com/turbot/steampipe-plugin-kubernetes/pull/235)) (Thanks [@afarid](https://github.com/afarid) for the contribution!!)

## v0.28.1 [2024-06-17]

_Bug fixes_

- Fixed the issue of missing and inconsistent columns in Kubernetes CRD tables. ([#229](https://github.com/turbot/steampipe-plugin-kubernetes/pull/229)) (Thanks [@dongho-jung](https://github.com/dongho-jung) for the contribution!!)

## v0.28.0 [2024-05-09]

_Enhancements_

- The `context_name` column has now been assigned as a connection key column across all the tables which facilitates more precise and efficient querying across multiple Kubernetes connections. ([#217](https://github.com/turbot/steampipe-plugin-kubernetes/pull/217))
- The Plugin and the Steampipe Anywhere binaries are now built with the `netgo` package. ([#219](https://github.com/turbot/steampipe-plugin-kubernetes/pull/219))
- Added the `version` flag to the plugin's Export tool. ([#65](https://github.com/turbot/steampipe-export/pull/65))

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.10.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v5100-2024-04-10) that adds support for connection key columns. ([#217](https://github.com/turbot/steampipe-plugin-kubernetes/pull/217))
- Recompiled plugin with [github.com/hashicorp/go-getter v1.7.4](https://github.com/hashicorp/go-getter/releases/tag/v1.7.4). ([#218](https://github.com/turbot/steampipe-plugin-kubernetes/pull/218))

## v0.27.0 [2024-01-22]

_Enhancements_

- Added the `annotations` columns on all CRD resources. ([#202](https://github.com/turbot/steampipe-plugin-kubernetes/pull/202))
- Updated the API version for table `kubernetes_horizontal_pod_autoscaler`. ([#190](https://github.com/turbot/steampipe-plugin-kubernetes/pull/190))

## v0.26.0 [2023-12-12]

_What's new?_

- The plugin can now be downloaded and used with the [Steampipe CLI](https://steampipe.io/docs), as a [Postgres FDW](https://steampipe.io/docs/steampipe_postgres/overview), as a [SQLite extension](https://steampipe.io/docs//steampipe_sqlite/overview) and as a standalone [exporter](https://steampipe.io/docs/steampipe_export/overview). ([#197](https://github.com/turbot/steampipe-plugin-kubernetes/pull/197))
- The table docs have been updated to provide corresponding example queries for Postgres FDW and SQLite extension. ([#197](https://github.com/turbot/steampipe-plugin-kubernetes/pull/197))
- Docs license updated to match Steampipe [CC BY-NC-ND license](https://github.com/turbot/steampipe-plugin-kubernetes/blob/main/docs/LICENSE). ([#197](https://github.com/turbot/steampipe-plugin-kubernetes/pull/197))

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.8.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v580-2023-12-11) that includes plugin server encapsulation for in-process and GRPC usage, adding Steampipe Plugin SDK version to `_ctx` column, and fixing connection and potential divide-by-zero bugs. ([#196](https://github.com/turbot/steampipe-plugin-kubernetes/pull/196))

## v0.25.2 [2023-11-21]

_Bug fixes_

- Fixed the plugin to pass the namespace qualifier to the Kubernetes API client when querying namespace scoped resources. ([#181](https://github.com/turbot/steampipe-plugin-kubernetes/pull/181)) (Thanks [@pdecat](https://github.com/pdecat) for the contribution!!)

## v0.25.1 [2023-10-03]

_Bug fixes_

- Fixed the plugin to prevent crashes when `source_types` config argument contains `manifest` but `manifest_file_paths` is not defined. ([#177](https://github.com/turbot/steampipe-plugin-kubernetes/pull/177))

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.6.2](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v562-2023-10-03) which prevents nil pointer reference errors for implicit hydrate configs. ([#178](https://github.com/turbot/steampipe-plugin-kubernetes/pull/178))

## v0.25.0 [2023-10-02]

_Dependencies_

- Upgraded to [steampipe-plugin-sdk v5.6.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v561-2023-09-29) with support for rate limiters. ([#173](https://github.com/turbot/steampipe-plugin-kubernetes/pull/173))
- Recompiled plugin with Go version `1.21`. ([#173](https://github.com/turbot/steampipe-plugin-kubernetes/pull/173))

## v0.24.0 [2023-09-29]

_Deprecated_

- The `source_type` config argument has been deprecated and will be removed in the next major version. Please use the `source_types` config argument instead. If both config arguments are set, `source_types` will take precedence. For backward compatibility, please see below for old and new value equivalents: ([#167](https://github.com/turbot/steampipe-plugin-kubernetes/pull/167))
  - `source_type = 'all'` : `source_types = ["deployed", "helm", "manifest"]`
  - `source_type = 'deployed'` : `source_types = ["deployed"]`
  - `source_type = 'helm'` : `source_types = ["helm"]`
  - `source_type = 'manifest'` : `source_types = ["manifest"]`

_What's new?_

- Added the `source_types` config argument, which allows specifying a combination of source types to load per connection. ([#167](https://github.com/turbot/steampipe-plugin-kubernetes/pull/167))

## v0.23.0 [2023-09-26]

_What's new?_

- New tables added
  - [kubernetes_pod_template](https://hub.steampipe.io/plugins/turbot/kubernetes/tables/kubernetes_pod_template) ([#170](https://github.com/turbot/steampipe-plugin-kubernetes/pull/170))

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.5.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v551-2023-07-26). ([#152](https://github.com/turbot/steampipe-plugin-kubernetes/pull/152))

## v0.22.1 [2023-08-24]

_Bug fixes_

- Fixed the `data` column of `kubernetes_config_map` and `kubernetes_secret` tables to correctly return data instead of `nil`. ([#150](https://github.com/turbot/steampipe-plugin-kubernetes/pull/150)) (Thanks [@hileef](https://github.com/hileef) for the contribution!!)

## v0.22.0 [2023-07-05]

_Enhancements_

- The plugin has been updated to support the new environment variables `KUBECONFIG` and `KUBE_CONFIG_PATH`, replacing the previous variables `KUBE_CONFIG_PATHS` and `KUBERNETES_MASTER`. The usage of `KUBE_CONFIG_PATHS` and `KUBERNETES_MASTER` has been marked as deprecated and will be removed in a future release. To ensure seamless functionality, it is strongly recommended to update any existing scripts and workflows to utilize the new environment variables instead. ([#142](https://github.com/turbot/steampipe-plugin-kubernetes/pull/142)) (Thanks [@mrkwtz](https://github.com/mrkwtz) for the contribution!!)

## v0.21.0 [2023-06-23]

_What's new?_

- Added support to query Kubernetes Helm charts. This can be set using the `helm_rendered_charts` config argument in the `kubernetes.spc` file. Please check the [Helm Charts](https://hub.steampipe.io/plugins/turbot/kubernetes#helm-charts) section for more information. ([#139](https://github.com/turbot/steampipe-plugin-kubernetes/pull/139))

_Bug fixes_

- Fixed typo in the description of `service_name` column of `kubernetes_stateful_set` table. ([#141](https://github.com/turbot/steampipe-plugin-kubernetes/pull/141)) (Thanks [@jacksgt](https://github.com/jacksgt) for the contribution!!)

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.5.0](https://github.com/turbot/steampipe-plugin-sdk/blob/v5.5.0/CHANGELOG.md#v550-2023-06-16) which significantly reduces API calls and boosts query performance, resulting in faster data retrieval. ([#139](https://github.com/turbot/steampipe-plugin-kubernetes/pull/139))

## v0.20.0 [2023-05-24]

_What's new?_

- Added support to query Kubernetes manifest files. This can be set using the `manifest_file_paths` config argument in the `kubernetes.spc` file. Please check the [Supported Manifest File Path Formats](https://hub.steampipe.io/plugins/turbot/kubernetes#manifest-files) section for more information. ([#129](https://github.com/turbot/steampipe-plugin-kubernetes/pull/129))

## v0.19.0 [2023-05-11]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.4.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v541-2023-05-05) which fixes increased plugin initialization time due to multiple connections causing the schema to be loaded repeatedly. ([#134](https://github.com/turbot/steampipe-plugin-kubernetes/pull/134))

## v0.18.1 [2023-04-14]

_Bug fixes_

- The plugin will no longer fail to initialize when attempting to create dynamic custom resource tables and the client cannot be created. ([#130](https://github.com/turbot/steampipe-plugin-kubernetes/pull/130))
- Fixed the `selector_query` column of `kubernetes_replication_controller` table to correctly return data instead of an error. ([#127](https://github.com/turbot/steampipe-plugin-kubernetes/pull/127))
- Fixed `kubernetes_endpoint_slice`, `kubernetes_horizontal_pod_autoscaler`, `kubernetes_pod_disruption_budget` and `kubernetes_pod_security_policy` tables to correctly return data instead of an error by removing incompatible API dependencies. ([#126](https://github.com/turbot/steampipe-plugin-kubernetes/pull/126))

## v0.18.0 [2023-04-04]

_What's new?_

- Added a new config argument `custom_resource_tables` to allow filtering of which CRDs to create table for. For more information please check the [Configuration](https://hub.steampipe.io/plugins/turbot/kubernetes#configuration) section. ([#121](https://github.com/turbot/steampipe-plugin-kubernetes/pull/121))

_Enhancements_

- Updated `docs/index.md` to include better and more detailed multiple context connection examples. ([#122](https://github.com/turbot/steampipe-plugin-kubernetes/pull/122))

## v0.17.0 [2023-03-09]

_What's new?_

- New tables added
  - [kubernetes_{custom_resource_singular_name}](https://hub.steampipe.io/plugins/turbot/kubernetes/tables/kubernetes_{custom_resource_singular_name}) ([#110](https://github.com/turbot/steampipe-plugin-kubernetes/pull/110))
- Added support for creating dynamic tables for custom resources. A table is automatically created for each custom resource in a cluster. To learn more, please see [kubernetes_{custom_resource_singular_name}](https://hub.steampipe.io/plugins/turbot/kubernetes/tables/kubernetes_{custom_resource_singular_name}).

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.2.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v520-2023-03-02) which includes fixes for query cache pending item mechanism and aggregator connections not working for dynamic tables. ([#110](https://github.com/turbot/steampipe-plugin-kubernetes/pull/110))

## v0.16.0 [2023-02-10]

_Enhancements_

- Added column `title` to `kubernetes_config_map` table. ([#107](https://github.com/turbot/steampipe-plugin-kubernetes/pull/107))

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.1.3](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v513-2023-02-09) which fixes the query caching functionality. ([#111](https://github.com/turbot/steampipe-plugin-kubernetes/pull/111))

## v0.15.0 [2023-01-05]

_Bug fixes_

- Renamed column `backend` to `default_backend` in `kubernetes_ingress` table to correctly follow the naming convention used in the API response. ([#98](https://github.com/turbot/steampipe-plugin-kubernetes/pull/98))
- Fixed the `default_backend` column (earlier named as `backend`) in `kubernetes_ingress` table to correctly return data instead of an error. ([#98](https://github.com/turbot/steampipe-plugin-kubernetes/pull/98))

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.0.2](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v502-2023-01-04) which fixes optional key column quals not working correctly for list hydrate call for plugins using `TableMapFunc`. ([#103](https://github.com/turbot/steampipe-plugin-kubernetes/pull/103))

## v0.14.0 [2022-12-26]

_What's new?_

- New tables added
  - [kubernetes_storage_class](https://hub.steampipe.io/plugins/turbot/kubernetes/tables/kubernetes_storage_class) ([#99](https://github.com/turbot/steampipe-plugin-kubernetes/pull/99))

## v0.13.1 [2022-11-18]

_Bug fixes_

- Temporarily disabled dynamic custom resource table creation due to aggregator connection incompatibility. ([#95](https://github.com/turbot/steampipe-plugin-kubernetes/pull/95))

## v0.13.0 [2022-11-16]

_What's new?_

- New tables added
  - [kubernetes_event](https://hub.steampipe.io/plugins/turbot/kubernetes/tables/kubernetes_event) ([#93](https://github.com/turbot/steampipe-plugin-kubernetes/pull/93)) (Thanks to [@svend](https://github.com/svend) for the new table!)
  - [{custom_resource_name.group_name}](https://hub.steampipe.io/plugins/turbot/kubernetes/tables/{custom_resource_name.group_name}) ([#85](https://github.com/turbot/steampipe-plugin-kubernetes/pull/85))
- Added support for creating dynamic tables for custom resources. A table is automatically created for each custom resource in a cluster. To learn more, please see [{custom_resource_name.group_name}](https://hub.steampipe.io/plugins/turbot/kubernetes/tables/{custom_resource_name.group_name}).

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.0.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v500-2022-11-16) which includes support for fetching remote files with go-getter and file watching. ([#85](https://github.com/turbot/steampipe-plugin-kubernetes/pull/85))

## v0.12.0 [2022-10-19]

_What's new?_

- New tables added
  - [kubernetes_custom_resource_definition](https://hub.steampipe.io/plugins/turbot/kubernetes/tables/kubernetes_custom_resource_definition) ([#84](https://github.com/turbot/steampipe-plugin-kubernetes/pull/84))
  - [kubernetes_horizontal_pod_autoscaler](https://hub.steampipe.io/plugins/turbot/kubernetes/tables/kubernetes_horizontal_pod_autoscaler) ([#86](https://github.com/turbot/steampipe-plugin-kubernetes/pull/86)) (Thanks [@aminvielledebatAtBedrock](https://github.com/aminvielledebatAtBedrock) for the contribution!)
- Added support for accessing the Kubernetes APIs from within a pod using [InClusterConfig](https://kubernetes.io/docs/tasks/run-application/access-api-from-pod/#accessing-the-api-from-within-a-pod). This is an alternative method of configuring Kubernetes credentials for the plugin when no kubeconfig file is found. ([#82](https://github.com/turbot/steampipe-plugin-kubernetes/pull/82))

## v0.11.0 [2022-09-26]

_What's new?_

- New tables added
  - [kubernetes_pod_disruption_budget](https://hub.steampipe.io/plugins/turbot/kubernetes/tables/kubernetes_pod_disruption_budget) ([#76](https://github.com/turbot/steampipe-plugin-kubernetes/pull/76)) (Thanks to [@mafrosis](https://github.com/mafrosis) for the contribution!)

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v4.1.7](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v417-2022-09-08) which includes several caching and memory management improvements. ([#80](https://github.com/turbot/steampipe-plugin-kubernetes/pull/80))
- Recompiled plugin with Go version `1.19`. ([#80](https://github.com/turbot/steampipe-plugin-kubernetes/pull/80))

## v0.10.0 [2022-07-07]

_Enhancements_

- Recompiled plugin with [steampipe-plugin-sdk v3.3.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v331--2022-06-30) which includes several caching fixes. ([#74](https://github.com/turbot/steampipe-plugin-kubernetes/pull/74))

## v0.9.0 [2022-06-27]

_Enhancements_

- Recompiled plugin with [steampipe-plugin-sdk v3.3.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v330--2022-6-22). ([#71](https://github.com/turbot/steampipe-plugin-kubernetes/pull/71))

## v0.8.0 [2022-06-01]

_Enhancements_

- Added additional optional key quals, page limits, filter support, and context cancellation handling across all the tables. ([#53](https://github.com/turbot/steampipe-plugin-kubernetes/pull/53))

## v0.7.1 [2022-05-23]

_Bug fixes_

- Fixed the Slack community links in README and docs/index.md files. ([#68](https://github.com/turbot/steampipe-plugin-kubernetes/pull/68))

## v0.7.0 [2022-05-16]

_Enhancements_

- Added column `selector_query` to the following tables: ([#65](https://github.com/turbot/steampipe-plugin-kubernetes/pull/65))
  - `kubernetes_daemonset`
  - `kubernetes_deployment`
  - `kubernetes_job`
  - `kubernetes_replicaset`
  - `kubernetes_replication_controller`
- Added column `label_selector` to `kubernetes_pod` table. ([#64](https://github.com/turbot/steampipe-plugin-kubernetes/pull/64))

## v0.6.0 [2022-04-28]

_Enhancements_

- Added support for native Linux ARM and Mac M1 builds. ([#58](https://github.com/turbot/steampipe-plugin-kubernetes/pull/58))
- Recompiled plugin with [steampipe-plugin-sdk v3.1.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v310--2022-03-30) and Go version `1.18`. ([#57](https://github.com/turbot/steampipe-plugin-kubernetes/pull/57))
- Added column `available_replicas` to `kubernetes_stateful_set` table ([#60](https://github.com/turbot/steampipe-plugin-kubernetes/pull/60))

## v0.5.0 [2022-03-23]

_Enhancements_

- Recompiled plugin with [steampipe-plugin-sdk v2.1.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v211--2022-03-10) ([#54](https://github.com/turbot/steampipe-plugin-kubernetes/pull/54))

## v0.4.0 [2022-01-19]

_What's new?_

- New tables added
  - [kubernetes_cronjob](https://hub.steampipe.io/plugins/turbot/kubernetes/tables/kubernetes_cronjob) ([#44](https://github.com/turbot/steampipe-plugin-kubernetes/pull/44))

_Enhancements_

- Imported the azure package to get the authentication that works with AzureAD OIDC ([#48](https://github.com/turbot/steampipe-plugin-kubernetes/pull/48))
- Recompiled plugin with [steampipe-plugin-sdk v1.8.3](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v183--2021-12-23) ([#49](https://github.com/turbot/steampipe-plugin-kubernetes/pull/49))
- Added the `template` column to the `kubernetes_stateful_set` table ([#46](https://github.com/turbot/steampipe-plugin-kubernetes/pull/46))

## v0.3.0 [2021-12-10]

_What's new?_

- Added support for querying Kubernetes clusters that use OIDC authentication mechanism ([#34](https://github.com/turbot/steampipe-plugin-kubernetes/pull/34))

## v0.2.0 [2021-12-08]

_Enhancements_

- Recompiled plugin with Go version 1.17 ([#36](https://github.com/turbot/steampipe-plugin-kubernetes/pull/36))
- Recompiled plugin with [steampipe-plugin-sdk v1.8.2](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v182--2021-11-22) ([#35](https://github.com/turbot/steampipe-plugin-kubernetes/pull/35))

## v0.1.0 [2021-09-01]

_What's new?_

- New tables added
  - [kubernetes_limit_range](https://hub.steampipe.io/plugins/turbot/kubernetes/tables/kubernetes_limit_range) ([#30](https://github.com/turbot/steampipe-plugin-kubernetes/pull/30))
  - [kubernetes_resource_quota](https://hub.steampipe.io/plugins/turbot/kubernetes/tables/kubernetes_resource_quota) ([#29](https://github.com/turbot/steampipe-plugin-kubernetes/pull/29))

_Enhancements_

- Recompiled plugin with [steampipe-plugin-sdk v1.5.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v150--2021-08-06) ([#25](https://github.com/turbot/steampipe-plugin-kubernetes/pull/25))
- Updated plugin license to Apache 2.0 per [turbot/steampipe#22](https://github.com/turbot/steampipe-plugin-kubernetes/pull/22)

## v0.0.2 [2021-06-03]

_What's new?_

- New tables added
  - [kubernetes_service](https://hub.steampipe.io/plugins/turbot/kubernetes/tables/kubernetes_service) ([#13](https://github.com/turbot/steampipe-plugin-kubernetes/pull/13))
  - [kubernetes_stateful_set](https://hub.steampipe.io/plugins/turbot/kubernetes/tables/kubernetes_stateful_set) ([#18](https://github.com/turbot/steampipe-plugin-kubernetes/pull/18))

## v0.0.1 [2021-04-01]

_What's new?_

- New tables added
  - [kubernetes_cluster_role](https://hub.steampipe.io/plugins/turbot/kubernetes/tables/kubernetes_cluster_role)
  - [kubernetes_cluster_role_binding](https://hub.steampipe.io/plugins/turbot/kubernetes/tables/kubernetes_cluster_role_binding)
  - [kubernetes_config_map](https://hub.steampipe.io/plugins/turbot/kubernetes/tables/kubernetes_config_map)
  - [kubernetes_daemonset](https://hub.steampipe.io/plugins/turbot/kubernetes/tables/kubernetes_daemonset)
  - [kubernetes_deployment](https://hub.steampipe.io/plugins/turbot/kubernetes/tables/kubernetes_deployment)
  - [kubernetes_endpoint](https://hub.steampipe.io/plugins/turbot/kubernetes/tables/kubernetes_endpoint)
  - [kubernetes_endpoint_slice](https://hub.steampipe.io/plugins/turbot/kubernetes/tables/kubernetes_endpoint_slice)
  - [kubernetes_ingress](https://hub.steampipe.io/plugins/turbot/kubernetes/tables/kubernetes_ingress)
  - [kubernetes_job](https://hub.steampipe.io/plugins/turbot/kubernetes/tables/kubernetes_job)
  - [kubernetes_namespace](https://hub.steampipe.io/plugins/turbot/kubernetes/tables/kubernetes_namespace)
  - [kubernetes_network_policy](https://hub.steampipe.io/plugins/turbot/kubernetes/tables/kubernetes_network_policy)
  - [kubernetes_node](https://hub.steampipe.io/plugins/turbot/kubernetes/tables/kubernetes_node)
  - [kubernetes_persistent_volume](https://hub.steampipe.io/plugins/turbot/kubernetes/tables/kubernetes_persistent_volume)
  - [kubernetes_persistent_volume_claim](https://hub.steampipe.io/plugins/turbot/kubernetes/tables/kubernetes_persistent_volume_claim)
  - [kubernetes_pod](https://hub.steampipe.io/plugins/turbot/kubernetes/tables/kubernetes_pod)
  - [kubernetes_pod_security_policy](https://hub.steampipe.io/plugins/turbot/kubernetes/tables/kubernetes_pod_security_policy)
  - [kubernetes_replicaset](https://hub.steampipe.io/plugins/turbot/kubernetes/tables/kubernetes_replicaset)
  - [kubernetes_replication_controller](https://hub.steampipe.io/plugins/turbot/kubernetes/tables/kubernetes_replication_controller)
  - [kubernetes_role](https://hub.steampipe.io/plugins/turbot/kubernetes/tables/kubernetes_role)
  - [kubernetes_role_binding](https://hub.steampipe.io/plugins/turbot/kubernetes/tables/kubernetes_role_binding)
  - [kubernetes_secret](https://hub.steampipe.io/plugins/turbot/kubernetes/tables/kubernetes_secret)
  - [kubernetes_service_account](https://hub.steampipe.io/plugins/turbot/kubernetes/tables/kubernetes_service_account)
