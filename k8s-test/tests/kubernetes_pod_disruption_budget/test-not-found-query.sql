select
  name,
  max_unavailable
from
  kubernetes_pod_disruption_budget
where
  name = 'zk-pdb-aa';

