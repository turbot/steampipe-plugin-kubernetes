select
  name,
  role_name,
  role_kind,
  subjects,
  age(current_timestamp, creation_timestamp)
from
  kubernetes.kubernetes_cluster_role_binding
where
  name = '';

