select
  name,
  namespace,
  status_replicas,
  ready_replicas,
  updated_replicas,
  available_replicas,
  unavailable_replicas
from
  kubernetes.kubernetes_deployment
where
  name = 'nginx-deployment-test'
  and namespace = 'default';

