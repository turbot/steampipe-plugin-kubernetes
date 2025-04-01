---
title: "Steampipe Table: helm_template - Query Kubernetes Helm Templates using SQL"
description: "Allows users to query Helm Templates in Kubernetes, specifically providing details about the chart, metadata, and template files, offering insights into Kubernetes deployment configurations."
folder: "Helm"
---

# Table: helm_template - Query Kubernetes Helm Templates using SQL

Helm Templates are part of Kubernetes, a platform for managing containerized applications across a cluster of nodes. Helm, a package manager for Kubernetes, uses templates to generate Kubernetes manifest files, which describe the resources needed for applications. These templates offer a way to manage complex applications and their dependencies in a standardized, repeatable, and efficient manner.

## Table Usage Guide

The `helm_template` table provides insights into Helm Templates within Kubernetes. As a DevOps engineer, explore template-specific details through this table, including chart details, metadata, and template files. Utilize it to understand Kubernetes deployment configurations, manage complex applications, and their dependencies more efficiently.

**Important Notes**
- The table will show the raw template as defined in the file. To list the fully rendered templates, use table `helm_template_rendered`.

## Examples

### Basic info
Explore the basic information of your Helm charts, including their names and paths. This can help you gain insights into your Helm configuration, understand its structure, and identify any potential issues.

```sql+postgres
select
  chart_name,
  path,
  raw
from
  helm_template;
```

```sql+sqlite
select
  chart_name,
  path,
  raw
from
  helm_template;
```

### List templates defined for a specific chart
Explore which templates are defined for a specific chart in a Helm-based application deployment. This can be useful in understanding the configuration and setup of a specific application like 'redis'.

```sql+postgres
select
  chart_name,
  path,
  raw
from
  helm_template
where
  chart_name = 'redis';
```

```sql+sqlite
select
  chart_name,
  path,
  raw
from
  helm_template
where
  chart_name = 'redis';
```