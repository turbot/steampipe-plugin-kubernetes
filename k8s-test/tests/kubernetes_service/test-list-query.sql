select
  name,
  namespace,
  cluster_ip
from
  kubernetes.kubernetes_service
where
  name = 'jenkins'
order by
  namespace,
  name;
