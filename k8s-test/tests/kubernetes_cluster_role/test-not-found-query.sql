select
  name,
  rules,
  age(current_timestamp, creation_timestamp)
from
  kubernetes.kubernetes_cluster_role
where
  name = ''
order by
  name;

