select
  name,
  namespace
from
  kubernetes.kubernetes_resource_quota
where
  name = 'pods-medium'
order by
  namespace,
  name;

