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

### Get ingress rules for a specific network policy

```sql
select
  name,
  namespace,
  concat(ports -> 'port', '/', ports ->> 'protocol') as to_port,
  case when from_source ? 'ipBlock' then
    'ip_block'
  when from_source ? 'namespaceSelector' then
    'namespace_selector'
  when from_source ? 'podSelector' then
    'pod_selector'
  end as source_type,
  coalesce(from_source -> 'ipBlock', from_source -> 'namespaceSelector', from_source -> 'podSelector') as source
from
  k8s_minikube.kubernetes_network_policy,
  jsonb_array_elements(ingress) as ingress_rule,
  jsonb_array_elements(ingress_rule -> 'ports') as ports,
  jsonb_array_elements(ingress_rule -> 'from') as from_source
where
  name = 'test-network-policy'
  and namespace = 'default';
```

### Get egress rules for a specific network policy

```sql
select
  name,
  namespace,
  concat(ports -> 'port', '/', ports ->> 'protocol') as to_port,
  case when to_destination ? 'ipBlock' then
    'ip_block'
  when to_destination ? 'namespaceSelector' then
    'namespace_selector'
  when to_destination ? 'podSelector' then
    'pod_selector'
  end as destination_type,
  coalesce(to_destination -> 'ipBlock', to_destination -> 'namespaceSelector', to_destination -> 'podSelector') as destination
from
  k8s_minikube.kubernetes_network_policy,
  jsonb_array_elements(egress) as egress_rule,
  jsonb_array_elements(egress_rule -> 'ports') as ports,
  jsonb_array_elements(egress_rule -> 'to') as to_destination
where
  name = 'test-network-policy'
  and namespace = 'default';
```
