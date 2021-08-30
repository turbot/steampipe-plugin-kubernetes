select
  name,
  namespace
from
  kubernetes.kubernetes_limit_range
where
  name = 'cpu-limit-range'
order by
  namespace,
  name;
