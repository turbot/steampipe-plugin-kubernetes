# Table: kubernetes_network_policy

Network policy specifiy how pods are allowed to communicate with each other and with other network endpoints.

## Examples

### Basic Info

```sql
select
  name,
  namespace,
  policy_types,
  ingress,
  egress,
  pod_selector,
  labels,
  annotations
from
  k8s_minikube.kubernetes_network_policy;
```

### List conditions for node

```yaml
-- ➜  ingress git:(issue-6) ✗ kubectl describe networkpolicies test-network-policy
Name:         test-network-policy
Namespace:    default
Created on:   2021-04-22 18:35:07 +0530 IST
Labels:       <none>
Annotations:  kubectl.kubernetes.io/last-applied-configuration:
                {"apiVersion":"networking.k8s.io/v1","kind":"NetworkPolicy","metadata":{"annotations":{},"name":"test-network-policy","namespace":"default...
Spec:
  PodSelector:     role=db
  Allowing ingress traffic:
    To Port: 6379/TCP
    From:
      IPBlock:
        CIDR: 172.17.0.0/16
        Except: 172.17.1.0/24
    From:
      NamespaceSelector: project=myproject
    From:
      PodSelector: role=frontend
  Allowing egress traffic:
    To Port: 5978/TCP
    To:
      IPBlock:
        CIDR: 10.0.0.0/24
        Except:
  Policy Types: Ingress, Egress

```

```sql
select
  name,
  namespace,
  concat(ports -> 'port', '/', ports ->> 'protocol') as to_port,
  case when sources ? 'ipBlock' then
    'ip_block'
  when sources ? 'namespaceSelector' then
    'namespace_selector'
  when sources ? 'podSelector' then
    'pod_elector'
  end as source_type,
  coalesce(sources -> 'ipBlock', sources -> 'namespaceSelector', sources -> 'podSelector') as source
from
  k8s_minikube.kubernetes_network_policy,
  jsonb_array_elements(ingress) as ingress_rule,
  jsonb_array_elements(ingress_rule -> 'ports') as ports,
  jsonb_array_elements(ingress_rule -> 'from') as sources
where
  name = 'test-network-policy'
  and namespace = 'default';
```
```
+---------------------+-----------+----------+--------------------+-----------------------------------------------------+
| name                | namespace | to_port  | source_type        | source                                              |
+---------------------+-----------+----------+--------------------+-----------------------------------------------------+
| test-network-policy | default   | 6379/TCP | ip_block           | {"cidr":"172.17.0.0/16","except":["172.17.1.0/24"]} |
| test-network-policy | default   | 6379/TCP | namespace_selector | {"matchLabels":{"project":"myproject"}}             |
| test-network-policy | default   | 6379/TCP | pod_elector        | {"matchLabels":{"role":"frontend"}}                 |
+---------------------+-----------+----------+--------------------+-----------------------------------------------------+
```

### Get system info for nodes

```sql
select
  name,
  node_info ->> 'machineID' as machine_id,
  node_info ->> 'systemUUID' as systemUUID,
  node_info ->> 'bootID' as bootID,
  node_info ->> 'kernelVersion' as kernelVersion,
  node_info ->> 'osImage' as osImage,
  node_info ->> 'operatingSystem' as operatingSystem,
  node_info ->> 'architecture' as architecture,
  node_info ->> 'containerRuntimeVersion' as containerRuntimeVersion,
  node_info ->> 'kubeletVersion' as kubeletVersion,
  node_info ->> 'kubeProxyVersion' as kubeProxyVersion
from
  kubernetes_node
```
