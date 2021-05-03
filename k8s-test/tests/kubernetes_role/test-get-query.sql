select
  name,
  namespace,
  rules
from
  kubernetes.kubernetes_role
where
  name = 'jenkins'
  and namespace = 'default'
order by
  namespace,
  name;

