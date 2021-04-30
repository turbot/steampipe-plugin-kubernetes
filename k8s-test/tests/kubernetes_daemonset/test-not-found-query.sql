select
  name,
  namespace,
  desired_number_scheduled as desired,
  current_number_scheduled as current,
  number_ready as ready,
  number_available as available,
  selector,
  age(current_timestamp, creation_timestamp)
from
  kubernetes.kubernetes_daemonset
where
  name = '' and namespace ='';

