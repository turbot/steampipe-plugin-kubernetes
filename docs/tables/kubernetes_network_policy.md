---
title: "Steampipe Table: kubernetes_network_policy - Query Kubernetes Network Policies using SQL"
description: "Allows users to query Kubernetes Network Policies, specifically to retrieve information about the network traffic rules applied to pods in a Kubernetes cluster."
folder: "Network Policy"
---

# Table: kubernetes_network_policy - Query Kubernetes Network Policies using SQL

A Kubernetes Network Policy is a specification of how groups of pods are allowed to communicate with each other and other network endpoints. It provides a way to enforce rules on network traffic within a Kubernetes cluster, thereby enhancing the security of the cluster. Network Policies use labels to select pods and define rules which specify what traffic is allowed to the selected pods.

## Table Usage Guide

The `kubernetes_network_policy` table provides insights into the network policies within a Kubernetes cluster. As a security analyst or a DevOps engineer, explore policy-specific details through this table, including pod selectors, policy types, and ingress and egress rules. Utilize it to uncover information about network traffic rules, such as those allowing or blocking specific types of traffic, thereby helping in assessing and enhancing the security posture of your Kubernetes clusters.

## Examples

### Basic Info
Explore which network policies are in place within your Kubernetes environment. This allows you to gain insights into the security settings and manage access controls more effectively.

```sql+postgres
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

```sql+sqlite
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
Explore which network policies permit all outgoing traffic in a Kubernetes environment. This can be useful for identifying potential security risks and ensuring that your network configurations adhere to best practices.

```sql+postgres
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

```sql+sqlite
Error: The corresponding SQLite query is unavailable.
```

### List default deny egress policies
Analyze the settings to understand the network policies that default to denying egress. This is particularly useful for enhancing security by identifying policies that prevent outbound network traffic.

```sql+postgres
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

```sql+sqlite
Error: The corresponding SQLite query is unavailable.
```

### List policies that allow all ingress
Analyze the settings to understand which policies permit all incoming traffic, useful for enhancing security by identifying potential vulnerabilities in your network.

```sql+postgres
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

```sql+sqlite
Error: The corresponding SQLite query is unavailable.
```

### List default deny ingress policies
Discover the segments that have default deny ingress policies in place. This is useful in identifying potential security risks, as it highlights the policies that block all incoming traffic by default.

```sql+postgres
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

```sql+sqlite
Error: The corresponding SQLite query is unavailable.
```

### View rules for a specific network policy
Explore rules associated with a certain network policy to understand its ingress and egress configurations. This is useful for assessing security measures and traffic flow within a specific network environment.

```sql+postgres
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

```sql+sqlite
select
  name,
  namespace,
  policy_types,
  ingress,
  egress
from
  kubernetes_network_policy
where
  name = 'test-network-policy'
  and namespace = 'default';
```

### List manifest resources
Explore the network policies in your Kubernetes environment to understand their configurations and rules. This can be useful to ensure security standards are met and to identify any potential vulnerabilities or misconfigurations.

```sql+postgres
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

```sql+sqlite
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
