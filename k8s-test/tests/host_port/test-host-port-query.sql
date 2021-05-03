select
  name,
  namespace,
  port -> 'hostPort' as host_port
from
  k8s_minikube.kubernetes_daemonset,
  jsonb_array_elements(template -> 'spec' -> 'containers') as container,
  jsonb_array_elements(container -> 'ports') as port
where
  port::jsonb ? 'hostPort';

