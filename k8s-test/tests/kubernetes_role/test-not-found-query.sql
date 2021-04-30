select
  name,
  namespace,
  rules,
  age(current_timestamp, creation_timestamp)
from
  kubernetes.kubernetes_role
where
  name = 'jenkins_123_123'
  and namespace = ''
order by
  namespace,
  name;

