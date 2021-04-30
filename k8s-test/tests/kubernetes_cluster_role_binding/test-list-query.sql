select
  name,
  role_name,
  role_kind,
  subjects
from
  kubernetes.kubernetes_cluster_role_binding
where
  name = 'jenkins';

