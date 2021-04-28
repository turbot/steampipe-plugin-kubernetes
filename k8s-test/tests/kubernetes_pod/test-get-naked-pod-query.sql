select
  namespace,
  name,
  phase,
  -- age(current_timestamp, creation_timestamp),
  -- pod_ip,
  node_name,
  jsonb_array_length(containers) as container_count,
  jsonb_array_length(init_containers) as init_container_count,
  jsonb_array_length(ephemeral_containers) as ephemeral_container_count
from
  kubernetes_pod
where
  name = 'static-web'
order by
  namespace,
  name;

