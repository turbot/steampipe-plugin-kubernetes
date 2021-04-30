select
  name,
  namespace,
  rules,
  ingress_class_name as class
from
  kubernetes.kubernetes_ingress
where
  name = 'minimal-ingress';

