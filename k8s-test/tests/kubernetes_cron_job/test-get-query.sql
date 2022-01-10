select
  name,
  namespace,
  suspend
from
  kubernetes_cron_job
where
  name = 'hello'
  and namespace = 'default';

