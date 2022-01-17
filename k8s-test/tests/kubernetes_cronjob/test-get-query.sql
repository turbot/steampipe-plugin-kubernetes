select
  name,
  namespace,
  suspend
from
  kubernetes_cronjob
where
  name = 'hello'
  and namespace = 'default';

