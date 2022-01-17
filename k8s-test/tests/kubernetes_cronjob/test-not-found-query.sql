select
  name,
  namespace
from
  kubernetes_cronjob
where
  name = 'hello_123_123';

