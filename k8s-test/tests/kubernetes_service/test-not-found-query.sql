select
  name,
  namespace,
  ports,
  age(current_timestamp, creation_timestamp)
from
  kubernetes.kubernetes_service
where
  name = 'jenkins_123_123'
  and namespace = 'default'
order by
  namespace,
  name;
