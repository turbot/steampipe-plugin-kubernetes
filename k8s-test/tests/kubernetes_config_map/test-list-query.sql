select
  name,
  namespace,
  data.key,
  data.value
from
  kubernetes.kubernetes_config_map,
  jsonb_each(data) as data
where
  name = 'game-demo'
