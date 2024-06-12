
# Resource Version Management and Querying in Kubernetes

## Why We Avoid Creating Tables for Inactive Resource Versions and All Active Versions

### Inactive Versions:
- **Querying an Inactive Resource:** Always returns empty rows.
- **Meaninglessness of Empty Tables:** Adding tables that always return empty rows is pointless.

### Active Versions:
- **Plugin Behavior:** Should be consistent with `kubectl`.

### How `kubectl` Works:

1. **CRD with Two Active Versions:**
   - Both versions are set to active (`served: true`).
   - Only one version (`v1`) is marked as the storage version.
   - Each version has a different schema.
   - **Version 1:** Fields (`field1`, `field2`, `field6`).
   - **Version 2:** Fields (`field1`, `field2`, `field3`, `field4`).
   - Both versions have common fields (`field1` and `field2`).

```yaml
# myresource_crd.yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: myresources.example.com
spec:
  group: example.com
  names:
    plural: myresources
    singular: myresource
    kind: MyResource
    shortNames:
    - mr
  scope: Namespaced
  versions:
  - name: v1
    served: true
    storage: true
    schema:
      openAPIV3Schema:
        type: object
        properties:
          spec:
            type: object
            properties:
              field1:
                type: string
              field2:
                type: integer
              field6:
                type: integer
          status:
            type: object
            properties:
              state:
                type: string
  - name: v2
    served: true
    storage: false
    schema:
      openAPIV3Schema:
        type: object
        properties:
          spec:
            type: object
            properties:
              field1:
                type: string
              field2:
                type: integer
              field3:
                type: string
              field4:
                type: integer
          status:
            type: object
            properties:
              state:
                type: string
```

### Apply the CRD:
```sh
kubectl apply -f myresource_crd.yaml
```

### Create Custom Resources for Both Versions:
1. **Create Resource for `v2`**
2. **Create Resource for `v1`**

```yaml
# v2_myresource_example.yaml
apiVersion: example.com/v2
kind: MyResource
metadata:
  name: example-myresource
spec:
  field1: "value11"
  field2: 57
  field3: "value333"
  field4: 97
```

```yaml
# v1_myresource_example.yaml
apiVersion: example.com/v1
kind: MyResource
metadata:
  name: v1-example-myresource
spec:
  field1: "V1-value11"
  field2: 57
  field6: 66666
```

### Apply the Resources:
```sh
kubectl apply -f v2_myresource_example.yaml
kubectl apply -f v1_myresource_example.yaml
```

### List All CRDs:
```sh
kubectl get crds
```

**Result:**
```
NAME                             CREATED AT
myresources.example.com          2024-06-11T04:55:07Z
```

### List All Resources:
```sh
kubectl get myresource
```

**Result:**
```
NAME                    AGE
example-myresource      25h
v1-example-myresource   14h
```

### Get Resource by Name (`example-myresource`):
```sh
kubectl get myresource example-myresource -o yaml
```

**Result:**
```yaml
apiVersion: example.com/v2
kind: MyResource
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"example.com/v2","kind":"MyResource","metadata":{"annotations":{},"name":"example-myresource","namespace":"default"},"spec":{"field1":"value11","field2":57,"field3":"value333","field4":97}}
  creationTimestamp: "2024-06-11T04:55:23Z"
  generation: 4
  name: example-myresource
  namespace: default
  resourceVersion: "3553"
  uid: 9322225a-8ab8-4513-b2fc-560ab7b7162f
spec:
  field1: value11
  field2: 57
  field3: value333
  field4: 97
```

### Get Resource by Name (`v1-example-myresource`):
```sh
kubectl get myresource v1-example-myresource -o yaml
```

**Result:**
```yaml
apiVersion: example.com/v2
kind: MyResource
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"example.com/v1","kind":"MyResource","metadata":{"annotations":{},"name":"v1-example-myresource","namespace":"default"},"spec":{"field1":"V1-value11","field2":57,"field6":66666}}
  creationTimestamp: "2024-06-11T16:00:20Z"
  generation: 1
  name: v1-example-myresource
  namespace: default
  resourceVersion: "32877"
  uid: 4adbb2c0-8308-403f-82d1-9c79dbf5827c
spec:
  field1: V1-value11
  field2: 57
```

### Observations:
- The command `kubectl get myresource v1-example-myresource -o yaml` does not return the value for `field6` even though it was applied for version 1.
- Steampipe creates the table schema based on how `kubectl` works.