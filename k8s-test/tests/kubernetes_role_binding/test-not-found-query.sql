select
  name,
  namespace,
  role_name,
  role_kind,
  subjects,
  age(current_timestamp, creation_timestamp)
from
  kubernetes.kubernetes_role_binding
where
  name = ''
  and namespace = ''
order by
  namespace,
  name;

