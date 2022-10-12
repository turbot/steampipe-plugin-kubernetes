# Table: kubernetes_horizontal_pod_autoscaler

Kubernetes  HorizontalPodAutoscaler is the configuration for a horizontal pod autoscaler, which automatically manages the replica count of any resource implementing the scale subresource based on the metrics specified.

## Examples

### Basic Info

```sql
select
  name,
  namespace,
  min_replicas,
  max_replicas,
  current_replicas,
  desired_replicas 
from
  kubernetes_horizontal_pod_autoscaler;
```

### Get list of HPA metrics configurations

```sql
select
  name,
  namespace,
  min_replicas,
  max_replicas,
  current_replicas,
  desired_replicas,
  jsonb_array_elements(metrics) as metrics,
  jsonb_array_elements(current_metrics) as current_metrics,
  conditions 
from
  kubernetes_horizontal_pod_autoscaler;
```
