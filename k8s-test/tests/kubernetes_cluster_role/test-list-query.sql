select
  name,
  rules
from
  kubernetes.kubernetes_cluster_role
where
  name like '%jenkins%';

