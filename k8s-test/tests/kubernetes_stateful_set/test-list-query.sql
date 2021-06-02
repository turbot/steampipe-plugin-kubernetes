select
  name,
  namespace,
  service_name,
  replicas
from
  kubernetes.kubernetes_stateful_set
where
  name = 'web'
order by
  namespace,
  name;
