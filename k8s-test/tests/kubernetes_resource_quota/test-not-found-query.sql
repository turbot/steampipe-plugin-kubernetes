select
  name,
  namespace
from
  kubernetes.kubernetes_resource_quota
where
  name = ''
  and namespace = ''
order by
  namespace,
  name;

