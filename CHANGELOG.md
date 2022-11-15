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
