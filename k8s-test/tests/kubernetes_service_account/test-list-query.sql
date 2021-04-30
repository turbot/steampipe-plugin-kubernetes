select
  name,
  namespace,
  jsonb_array_length(secrets) as secrets
from
  kubernetes.kubernetes_service_account
where
  name = 'jenkins'
order by
  namespace,
  name;

