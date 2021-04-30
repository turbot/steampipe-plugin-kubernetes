select
  name,
  rules
from
  kubernetes.kubernetes_cluster_role
where
  name = 'jenkins'
order by
  name;

