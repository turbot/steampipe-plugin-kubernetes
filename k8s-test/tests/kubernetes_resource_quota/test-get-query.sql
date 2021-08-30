select
  name,
  namespace
from
  kubernetes.kubernetes_resource_quota
where
  name = 'pods-medium'
  and namespace = 'default'
order by
  namespace,
  name;

