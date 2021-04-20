select
  name,
  namespace,
  phase,
  -- pod_ip,
  node_name
from
  kubernetes_pod
where
  name = 'static-web'
  and owner_references is null;

