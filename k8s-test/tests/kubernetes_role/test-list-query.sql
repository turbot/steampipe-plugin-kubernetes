select
  name,
  namespace,
  rules
from
  kubernetes.kubernetes_role
where
  name = 'jenkins'
order by
  namespace,
  name;

