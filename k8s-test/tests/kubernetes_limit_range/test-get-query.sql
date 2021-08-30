select
  name,
  namespace
from
  kubernetes.kubernetes_limit_range
where
  name = 'cpu-limit-range'
  and namespace = 'default'
order by
  namespace,
  name;

