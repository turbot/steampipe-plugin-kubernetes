select
  name,
  namespace,
  service_name,
  replicas
from
  kubernetes.kubernetes_statefulset
where
  name = 'web'
order by
  namespace,
  name;
