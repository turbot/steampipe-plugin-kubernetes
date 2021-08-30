select
  name,
  namespace
from
  kubernetes.kubernetes_limit_range
where
  name = ''
  and namespace = ''
order by
  namespace,
  name;
