---
title: "Steampipe Table: kubernetes_secret - Query Kubernetes Secrets using SQL"
description: "Allows users to query Kubernetes Secrets, providing insights into the sensitive information like passwords, OAuth tokens, and ssh keys that are stored."
folder: "Secret"
---

# Table: kubernetes_secret - Query Kubernetes Secrets using SQL

Kubernetes Secrets is a resource that manages sensitive data such as passwords, OAuth tokens, ssh keys, etc. It provides a more secure and flexible solution to manage sensitive data in a Kubernetes cluster, compared to the alternative of putting this information directly into pod specification or in docker images. Kubernetes Secrets offers the ability to decouple sensitive content from the pod specification and isolate the visibility of such sensitive information to just the system components which require access to it.

## Table Usage Guide

The `kubernetes_secret` table provides insights into Kubernetes Secrets within a Kubernetes cluster. As a DevOps engineer, explore secret-specific details through this table, including the type of secret, the namespace it belongs to, and associated metadata. Utilize it to uncover information about secrets, such as those that are not in use, those that are exposed, or those that are stored in a non-compliant manner.

## Examples

### Basic Info
Explore the age and details of various Kubernetes secrets to understand their creation and configuration for better resource management and security. This could be particularly useful in identifying outdated or potentially vulnerable secrets that may need updating or removal.

```sql+postgres
select
  name,
  namespace,
  data.key,
  data.value,
  age(current_timestamp, creation_timestamp)
from
  kubernetes_secret,
  jsonb_each(data) as data
order by
  namespace,
  name;
```

```sql+sqlite
select
  name,
  namespace,
  data.key,
  data.value,
  (julianday('now') - julianday(creation_timestamp)) * 24 * 60 * 60
from
  kubernetes_secret,
  json_each(data) as data
order by
  namespace,
  name;
```

### List and base64 decode secret values
Explore the decoded values of secrets in your Kubernetes environment to better understand the information they hold. This can be particularly useful for troubleshooting or auditing purposes.

```sql+postgres
select
  name,
  namespace,
  data.key,
  decode(data.value, 'base64') as decoded_data,
  age(current_timestamp, creation_timestamp)
from
  kubernetes_secret,
  jsonb_each_text(data) as data
order by
  namespace,
  name;
```

```sql+sqlite
select
  name,
  namespace,
  data.key,
  data.value as decoded_data,
  julianday('now') - julianday(creation_timestamp)
from
  kubernetes_secret,
  json_each(data) as data
order by
  namespace,
  name;
```

### List manifest resources
Explore which encrypted data is associated with each resource in your Kubernetes environment. This can help you assess the elements within your system configuration and identify potential areas of concern.

```sql+postgres
select
  name,
  namespace,
  data.key,
  data.value,
  path
from
  kubernetes_secret,
  jsonb_each(data) as data
where
  path is not null
order by
  namespace,
  name;
```

```sql+sqlite
select
  name,
  namespace,
  data.key,
  data.value,
  path
from
  kubernetes_secret,
  json_each(data) as data
where
  path is not null
order by
  namespace,
  name;
```