select
  name,
  namespace,
  security_context,
  (security_context -> 'runAsNonRoot') as run_as_non_root
from
  kubernetes.kubernetes_pod
where
  name like '%pod%'
  and security_context::jsonb ? 'runAsNonRoot'
  and (security_context -> 'runAsNonRoot')::bool
order by
  name,
  namespace;

