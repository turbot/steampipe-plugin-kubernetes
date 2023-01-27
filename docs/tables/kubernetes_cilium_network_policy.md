# Table: cilium_kubernetes_network_policy

Cilium Network policy specifiy how pods are allowed to communicate with each other and with other network endpoints.

## Examples

### Basic Info

```sql
select
  name,
  namespace,
  ingress_allow,
  egress_allow,
  endpoint_selector,
  labels,
  annotations
from
  kubernetes_cilium_network_policy;
```

### List policies that allow all egress
```sql
select
  name,
  namespace,
  endpoint_selector,
  egress_allow
from
  kubernetes_cilium_network_policy
where
  endpoint_selector = '{}'
  and egress_allow @> '[{}]';
```


### List deny egress policies
```sql
select
  name,
  namespace,
  endpoint_selector,
  egress_deny
from
  kubernetes_cilium_network_policy
where
  endpoint_selector = '{}'
  and egress_deny is not null;

```
### List policies that allow all ingress
```sql
select
  name,
  namespace,
  endpoint_selector,
  ingress_allow
from
  kubernetes_cilium_network_policy
where
  endpoint_selector = '{}'
  and ingress_allow @> '[{}]';
```

### List deny ingress policies
```sql
select
  name,
  namespace,
  endpoint_selector,
  ingress_deny
from
  kubernetes_cilium_network_policy
where
  endpoint_selector = '{}'
  and ingress_deny is not null;
```


### View rules for a specific network policy

```sql
select
  name,
  namespace,
  policy_types,
  jsonb_pretty(ingress_allow),
  jsonb_pretty(egress_allow)
from
  kubernetes_cilium_network_policy
where
  name = 'test-network-policy'
  and namespace = 'default';
```
