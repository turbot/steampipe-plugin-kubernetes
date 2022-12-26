select
  name,
  allow_volume_expansion,
  reclaim_policy,
  volume_binding_mode,
  title
from
  kubernetes.kubernetes_storage_class
where
  title = 'mystorage'
order by
  title;
