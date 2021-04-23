# Table: kubernetes_persistent_volume

A PersistentVolume (PV) is a piece of storage in the cluster that has been provisioned by an administrator or dynamically provisioned using Storage Classes. PVs are volume plugins like Volumes, but have a lifecycle independent of any individual Pod that uses the PV.

## Examples

```json
// select * from k8s_minikube.kubernetes_persistent_volume

[
  {
    "access_modes": ["ReadWriteOnce"],
    "annotations": {
      "kubectl.kubernetes.io/last-applied-configuration": "{\"apiVersion\":\"v1\",\"kind\":\"PersistentVolume\",\"metadata\":{\"annotations\":{},\"labels\":{\"type\":\"local\"},\"name\":\"task-pv-volume\"},\"spec\":{\"accessModes\":[\"ReadWriteOnce\"],\"capacity\":{\"storage\":\"10Gi\"},\"hostPath\":{\"path\":\"/mnt/data\"},\"storageClassName\":\"manual\"}}\n"
    },
    "capacity": {
      "storage": "10Gi"
    },
    "claim_ref": null,
    "context_name": "minikube",
    "creation_timestamp": "2021-04-23 11:37:43",
    "deletion_grace_period_seconds": null,
    "deletion_timestamp": null,
    "finalizers": ["kubernetes.io/pv-protection"],
    "generate_name": "",
    "generation": 0,
    "labels": {
      "type": "local"
    },
    "message": "",
    "mount_options": null,
    "name": "task-pv-volume",
    "node_affinity": null,
    "owner_references": null,
    "persistent_volume_reclaim_policy": "Retain",
    "persistent_volume_source": {
      "hostPath": {
        "path": "/mnt/data",
        "type": ""
      }
    },
    "phase": "Available",
    "reason": "",
    "resource_version": "354906",
    "storage_class_name": "manual",
    "tags": {
      "kubectl.kubernetes.io/last-applied-configuration": "{\"apiVersion\":\"v1\",\"kind\":\"PersistentVolume\",\"metadata\":{\"annotations\":{},\"labels\":{\"type\":\"local\"},\"name\":\"task-pv-volume\"},\"spec\":{\"accessModes\":[\"ReadWriteOnce\"],\"capacity\":{\"storage\":\"10Gi\"},\"hostPath\":{\"path\":\"/mnt/data\"},\"storageClassName\":\"manual\"}}\n",
      "type": "local"
    },
    "title": "task-pv-volume",
    "uid": "50ba6a64-bf6b-42ad-a170-1f178a64ff11",
    "volume_mode": "Filesystem"
  }
]
```

```yaml
# ➜  steampipe-plugin-kubernetes git:(issue-6) ✗ kubectl get pv task-pv-volume -o wide
NAME             CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS      CLAIM   STORAGECLASS   REASON   AGE   VOLUMEMODE
task-pv-volume   10Gi       RWO            Retain           Available           manual                  21m   Filesystem

# ➜  steampipe-plugin-kubernetes git:(issue-6) ✗ kubectl describe pv task-pv-volume
Name:            task-pv-volume
Labels:          type=local
Annotations:     kubectl.kubernetes.io/last-applied-configuration:
                   {"apiVersion":"v1","kind":"PersistentVolume","metadata":{"annotations":{},"labels":{"type":"local"},"name":"task-pv-volume"},"spec":{"acce...
Finalizers:      [kubernetes.io/pv-protection]
StorageClass:    manual
Status:          Available
Claim:
Reclaim Policy:  Retain
Access Modes:    RWO
VolumeMode:      Filesystem
Capacity:        10Gi
Node Affinity:   <none>
Message:
Source:
    Type:          HostPath (bare host directory volume)
    Path:          /mnt/data
    HostPathType:
Events:            <none>
```

### Basic Info

```sql
select
  name,
  access_modes,
  storage_class,
  capacity ->> 'storage' as storage_capacity,
  creation_timestamp,
  persistent_volume_reclaim_policy,
  phase as status,
  volume_mode,
  age(current_timestamp, )
from
  kubernetes_persistent_volume
```

### List active jobs

```sql
select
  name,
  namespace,
  start_time,
  age(coalesce(completion_time, current_timestamp), start_time) as duration,
  active,
  succeeded,
  failed
from
  k8s_minikube.kubernetes_job
where active > 0
```

### List failed jobs

```sql
select
  name,
  namespace,
  start_time,
  age(coalesce(completion_time, current_timestamp), start_time) as duration,
  active,
  succeeded,
  failed
from
  k8s_minikube.kubernetes_job
where failed > 0
```

### Get list of container and images for jobs

```sql
select
  name,
  namespace,
  jsonb_agg(elems.value -> 'name') as containers,
  jsonb_agg(elems.value -> 'image') as images
from
  k8s_minikube.kubernetes_job,
  jsonb_array_elements(template -> 'spec' -> 'containers') as elems
group by name, namespace
```
