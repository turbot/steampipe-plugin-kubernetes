# Table: kubernetes_network_policy

Network policy specifiy how pods are allowed to communicate with each other and with other network endpoints.

## Examples

### Basic Info

```sql
select
  name,
  namespace,
  policy_types,
  ingress,
  egress,
  pod_selector,
  labels,
  annotations
from
  kubernetes_network_policy;
```

### List policies that allow all egress

```sql
select
  name,
  namespace,
  policy_types,
  pod_selector,
  egress
from
  kubernetes_network_policy
where
  policy_types @> '["Egress"]'
  and pod_selector = '{}'
  and egress @> '[{}]';
```

### List default deny egress policies

```sql
select
  name,
  namespace,
  policy_types,
  pod_selector,
  egress
from
  kubernetes_network_policy
where
  policy_types @> '["Egress"]'
  and pod_selector = '{}'
  and egress is null;
```

### List policies that allow all ingress

```sql
select
  name,
  namespace,
  policy_types,
  pod_selector,
  ingress
from
  kubernetes_network_policy
where
  policy_types @> '["Ingress"]'
  and pod_selector = '{}'
  and ingress @> '[{}]';
```

### List default deny ingress policies

```sql
select
  name,
  namespace,
  policy_types,
  pod_selector,
  ingress
from
  kubernetes_network_policy
where
  policy_types @> '["Ingress"]'
  and pod_selector = '{}'
  and ingress is null;
```

### View rules for a specific network policy

```sql
select
  name,
  namespace,
  policy_types,
  jsonb_pretty(ingress),
  jsonb_pretty(egress)
from
  kubernetes_network_policy
where
  name = 'test-network-policy'
  and namespace = 'default';
```

### List manifest resources

```sql
select
  name,
  namespace,
  policy_types,
  ingress,
  egress,
  pod_selector,
  labels,
  annotations,
  path
from
  kubernetes_network_policy
where
  path is not null;
```
