select
  name,
  namespace
from
  kubernetes_cron_job
where
  name = 'hello_123_123';

