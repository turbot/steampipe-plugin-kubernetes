# Table: {custom_resource_name..group_name}

Query data from the custom resource called `{custom_resource_name.group_name}`, e.g., `certificates.cert-manager.io`, `storeconfigs.crossplane.io`. A table is automatically created to represent each object in the `objects` argument.

## Examples

### Inspect the table structure

List all tables:

```sql
.inspect kubernetes;
+---------------------------------------+------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| table                                 | description                                                                                                                                                      |
+---------------------------------------+------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| certificates.cert-manager.io          | Represents Custom resource certificates.cert-manager.io.                                                                                                         |
| crontabs.stable.example.com           | Represents Custom resource crontabs.stable.example.com.                                                                                                          |
| kubernetes_cluster_role               | ClusterRole contains rules that represent a set of permissions.                                                                                                  |
| kubernetes_cluster_role_binding       | A ClusterRoleBinding grants the permissions defined in a cluster role to a user or set of users. Access granted by ClusterRoleBinding is cluster-wide.           |
| kubernetes_config_map                 | Config Map can be used to store fine-grained information like individual properties or coarse-grained information like entire config files or JSON blobs.        |
| kubernetes_cronjob                    | Cron jobs are useful for creating periodic and recurring tasks, like running backups or sending emails.                                                          |
| kubernetes_custom_resource_definition | Kubernetes Custom Resource Definition.                                                                                                                           |
| kubernetes_daemonset                  | A DaemonSet ensures that all (or some) Nodes run a copy of a Pod.                                                                                                |
| kubernetes_deployment                 | Kubernetes Deployment enables declarative updates for Pods and ReplicaSets.                                                                                      |
| kubernetes_endpoint                   | Set of addresses and ports that comprise a service. More info: https://kubernetes.io/docs/concepts/services-networking/service/#services-without-selectors.      |
| kubernetes_endpoint_slice             | EndpointSlice represents a subset of the endpoints that implement a service.                                                                                     |
| kubernetes_ingress                    | Ingress exposes HTTP and HTTPS routes from outside the cluster to services within the cluster. Traffic routing is controlled by rules defined on the Ingress res |
|                                       | ource.                                                                                                                                                           |
| kubernetes_job                        | A Job creates one or more Pods and will continue to retry execution of the Pods until a specified number of them successfully terminate.                         |
| kubernetes_limit_range                | Kubernetes Limit Range                                                                                                                                           |
| kubernetes_namespace                  | Kubernetes Namespace provides a scope for Names.                                                                                                                 |
| kubernetes_network_policy             | Network policy specifiy how pods are allowed to communicate with each other and with other network endpoints.                                                    |
| kubernetes_node                       | Kubernetes Node is a worker node in Kubernetes.                                                                                                                  |
| kubernetes_persistent_volume          | A PersistentVolume (PV) is a piece of storage in the cluster that has been provisioned by an administrator or dynamically provisioned using Storage Classes. PVs |
|                                       |  are volume plugins like Volumes, but have a lifecycle independent of any individual Pod that uses the PV.                                                       |
| kubernetes_persistent_volume_claim    | A PersistentVolumeClaim (PVC) is a request for storage by a user.                                                                                                |
| kubernetes_pod                        | Kubernetes Pod is a collection of containers that can run on a host. This resource is created by clients and scheduled onto hosts.                               |
| kubernetes_pod_disruption_budget      | A Pod Disruption Budget limits the number of Pods of a replicated application that are down simultaneously from voluntary disruptions.                           |
| kubernetes_pod_security_policy        | A Pod Security Policy is a cluster-level resource that controls security sensitive aspects of the pod specification. The PodSecurityPolicy objects define a set  |
|                                       | of conditions that a pod must run with in order to be accepted into the system, as well as defaults for the related fields.                                      |
| kubernetes_replicaset                 | Kubernetes replica set ensures that a specified number of pod replicas are running at any given time.                                                            |
| kubernetes_replication_controller     | A Replication Controller makes sure that a pod or homogeneous set of pods are always up and available. If there are too many pods, it will kill some. If there a |
|                                       | re too few, the Replication Controller will start more.                                                                                                          |
| kubernetes_resource_quota             | Kubernetes Resource Quota                                                                                                                                        |
| kubernetes_role                       | Role contains rules that represent a set of permissions.                                                                                                         |
| kubernetes_role_binding               | A role binding grants the permissions defined in a role to a user or set of users. It holds a list of subjects (users, groups, or service accounts), and a refer |
|                                       | ence to the role being granted.                                                                                                                                  |
| kubernetes_secret                     | Secrets can be used to store sensitive information either as individual properties or coarse-grained entries like entire files or JSON blobs.                    |
| kubernetes_service                    | A service provides an abstract way to expose an application running on a set of Pods as a network service.                                                       |
| kubernetes_service_account            | A service account provides an identity for processes that run in a Pod.                                                                                          |
| kubernetes_stateful_set               | A statefulSet is the workload API object used to manage stateful applications.                                                                                   |
| storeconfigs.secrets.crossplane.io    | Represents Custom resource storeconfigs.secrets.crossplane.io.                                                                                                   |
+---------------------------------------+------------------------------------------------------------------------------------------------------------------------------------------------------------------+
```

To get details of a specific custom resource table, inspect it by name:

```sql
select column_name, data_type from INFORMATION_SCHEMA.COLUMNS where table_name like 'certificates.cert-manager.io';
+------------------------------+--------------------------+
| column_name                  | data_type                |
+------------------------------+--------------------------+
| _ctx                         | jsonb                    |
| spec_ipAddresses             | jsonb                    |
| spec_usages                  | jsonb                    |
| spec_privateKey              | jsonb                    |
| spec_keystores               | jsonb                    |
| spec_secretTemplate          | jsonb                    |
| spec_uris                    | jsonb                    |
| spec_subject                 | jsonb                    |
| spec_revisionHistoryLimit    | bigint                   |
| spec_isCA                    | boolean                  |
| spec_encodeUsagesInRequest   | boolean                  |
| spec_emailAddresses          | jsonb                    |
| spec_additionalOutputFormats | jsonb                    |
| spec_issuerRef               | jsonb                    |
| spec_dnsNames                | jsonb                    |
| creation_timestamp           | timestamp with time zone |
| spec_secretName              | text                     |
| spec_duration                | text                     |
| kind                         | text                     |
| uid                          | text                     |
| spec_renewBefore             | text                     |
| spec_literalSubject          | text                     |
| spec_commonName              | text                     |
| name                         | text                     |
| namespace                    | text                     |
| api_version                  | text                     |
+------------------------------+--------------------------+
```

### Get all custom resources\_\_c

```sql
select
  *
from
  "custom_resource_name.group_name";
```

### List custom resources added in the last 24 hours\_\_c

```sql
select
  *
from
  "custom_resource_name.group_name"
where
  created_date = now() - interval '24 hrs';
```

### Get details for a custom resource\_\_c

```sql
select
  *
from
  custom_resource_name.group_name
where
  name = 'blah';
```

### Count of all the custom resources\_\_c

```sql
select
  *
from
  custom_resource_name.group_name
where
  name = 'blah';
```