select
  name,
  namespace,
  service_name,
  age(current_timestamp, creation_timestamp)
from
  kubernetes.kubernetes_stateful_set
where
  name = 'jenkins_123_123'
  and namespace = 'default'
order by
  namespace,
  name;
