select
  name,
  namespace,
  suspend,
  time_zone
from
  kubernetes_cronjob
where
  name = 'hello'
  and namespace = 'default';

