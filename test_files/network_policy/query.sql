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

select
  name,
  namespace,
  policy_types,
  ingress_rule,
  ingress_rule -> 'ports' -> 'port' as port_number,
  ingress_rule -> 'ports' -> 'protocol' as port_protocol,
  ingress_rule -> 'from' -> 'ipBlock' as ipBlock,
  ingress_rule -> 'from' -> 'namespaceSelector' as namespace_selector,
  ingress_rule -> 'from' -> 'podSelector' as pod_selector
from
  k8s_minikube.kubernetes_network_policy,
  jsonb_array_elements(ingress) as ingress_rule
where
  name = 'test-network-policy'
  and namespace = 'default';

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

select
  name,
  namespace,
  policy_types,
  ingress_rule
from
  k8s_minikube.kubernetes_network_policy,
  jsonb_array_elements(ingress) as ingress_rule
where
  name = 'test-network-policy'
  and namespace = 'default';

select
  name,
  namespace,
  concat(ports -> 'port', '/', ports ->> 'protocol') as to_port,
  if sources ? 'ipBlock' then
    'ip_block'
  else if sources ? 'namespaceSelector' then
      'namespace_selector'
  else
if sources ? 'podSelector' then
        'pod_elector'
      end if as source_type,
        coalesce(sources -> 'ipBlock', sources -> 'namespaceSelector', sources -> 'podSelector') as source
      from
        k8s_minikube.kubernetes_network_policy,
        jsonb_array_elements(ingress) as ingress_rule,
  jsonb_array_elements(ingress_rule -> 'ports') as ports,
  jsonb_array_elements(ingress_rule -> 'from') as sources
where
  name = 'test-network-policy' and namespace = 'default';

