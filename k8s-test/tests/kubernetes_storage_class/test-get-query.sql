select
  name,
  allow_volume_expansion,
  reclaim_policy,
  volume_binding_mode
from
  kubernetes.kubernetes_storage_class
where
  name = 'mystorage'
order by
  name;
