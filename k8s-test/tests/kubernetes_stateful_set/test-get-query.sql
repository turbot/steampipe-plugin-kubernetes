select
  name,
  namespace,
  service_name,
  replicas,
  selector,
  selector_query,
  collision_count,
  current_replicas,
  observed_generation,
  pod_management_policy,
  ready_replicas,
  revision_history_limit,
  updated_replicas,
  update_strategy
from
  kubernetes.kubernetes_stateful_set
where
  name = 'web'
  and namespace = 'default'
order by
  namespace,
  name;
