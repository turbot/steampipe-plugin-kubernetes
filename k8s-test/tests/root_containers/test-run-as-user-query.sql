select
  name,
  namespace,
  security_context,
  (security_context -> 'runAsUser') as run_as_user
from
  kubernetes.kubernetes_pod
where
  name like '%pod%'
  and security_context::jsonb ? 'runAsUser'
  and (security_context -> 'runAsUser')::int > 0
order by
  name,
  namespace;

