# Table: kubernetes_{custom_resource_singular_name}

Query data from the custom resource called `kubernetes_{custom_resource_singular_name}`, e.g., `kubernetes_certificate`, `kubernetes_capacityrequest`. A table is automatically created to represent each custom resource.

If the table name is already created in the above format and exists in the table list, then the subsequent ones will have the fully qualified name `kubernetes_{custom_resource_singular_name}_{custom_resource_group_name}`, e.g., `kubernetes_certificate_cert_manager_io`.

For instance, given the CRD `certManager.yaml`:

```yml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: certificates.cert-manager.io
  labels:
    app: cert-manager
    app.kubernetes.io/name: cert-manager-crd
spec:
  group: cert-manager.io
  names:
    kind: Certificate
    listKind: CertificateList
    plural: certificates
    shortNames:
      - cert
      - certs
    # table name - kubernetes_certificate
    singular: certificate
    categories:
      - cert-manager
  scope: Namespaced
  versions:
    - name: v1
      subresources:
        status: {}
      schema:
        openAPIV3Schema:
          description: "A Certificate resource should be created to ensure an up to date and signed x509 certificate is stored in the Kubernetes Secret resource named in `spec.secretName`."
          type: object
          required:
            - spec
          properties:
            apiVersion:
              description: The versioned schema of this representation of an object.
              type: string
            kind:
              description: A string value representing the REST resource this object represents.
              type: string
            metadata:
              type: object
            spec:
              description: Desired state of the Certificate resource.
              type: object
              required:
                - issuerRef
                - secretName
              properties:
                commonName:
                  description: A common name to be used on the Certificate.
                  type: string
                dnsNames:
                  description: A list of DNS subjectAltNames to be set on the Certificate.
                  type: array
                  items:
                    type: string
                duration:
                  description: "The requested 'duration' (i.e. lifetime) of the Certificate. Defaults to 90 days."
                  type: string
                ipAddresses:
                  description: A list of IP address subjectAltNames to be set on the Certificate.
                  type: array
                  items:
                    type: string
                isCA:
                  description: IsCA will mark this Certificate as valid for certificate signing.
                  type: boolean
                issuerRef:
                  description: A reference to the issuer for this certificate.
                  type: object
                  required:
                    - name
                  properties:
                    name:
                      description: Name of the resource being referred to.
                      type: string
                renewBefore:
                  description: "How long before the currently issued certificate's expiry cert-manager should renew the certificate. Default to 2/3 of the duration."
                  type: string
                secretName:
                  description: The name of the secret resource that will be automatically created and managed by this Certificate resource.
                  type: string
      served: true
      storage: true
```

And the custom resource `spCloudCertificate.yaml`:

```yml
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: temporal-w-spcloud123456
  labels:
    workspace-id: w_spcloud123456
    identity-id: o_spcloud123456
    region: apse1
    shard: "0001"
    workspace-pluginset: "202203170111114"
    app.kubernetes.io/component: steampipe-workspace-db
    app.kubernetes.io/managed-by: steampipe-api
    app.kubernetes.io/part-of: steampipe-cloud
    app.kubernetes.io/instance: w-spcloud123456
    app.kubernetes.io/version: 0.13.3-workspace-spcloud.20220317010004
spec:
  secretName: temporal-w-spcloud123456-tls
  duration: 87600h # 10 years
  dnsNames:
    - w-spcloud123456
```

If my connection configuration is:

```hcl
connection "kubernetes" {
  plugin = "kubernetes"
  custom_resource_tables = ["certificates.*"]
}
```

Steampipe will automatically create the `kubernetes_certificate` table, which can then be inspected and queried like other tables:

```
.inspect kubernetes_certificate;
+--------------------+--------------------------+--------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| column             | type                     | description                                                                                                                                                        |
+--------------------+--------------------------+--------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| _ctx               | jsonb                    | Steampipe context in JSON form, e.g. connection_name.                                                                                                              |
| api_version        | text                     | The API version of the resource.                                                                                                                                   |
| common_name        | text                     | A common name to be used on the Certificate.                                                                                                                       |
| context_name       | text                     | Kubectl config context name.                                                                                                                                       |
| creation_timestamp | timestamp with time zone | CreationTimestamp is a timestamp representing the server time when this object was created.                                                                        |
| dns_names          | jsonb                    | A list of DNS subjectAltNames to be set on the Certificate.                                                                                                        |
| duration           | text                     | The requested 'duration' (i.e. lifetime) of the Certificate. Defaults to 90 days.                                                                                  |
| end_line           | bigint                   | The path to the manifest file.                                                                                                                                     |
| ip_addresses       | jsonb                    | A list of IP address subjectAltNames to be set on the Certificate.                                                                                                 |
| is_ca              | boolean                  | IsCA will mark this Certificate as valid for certificate signing.                                                                                                  |
| issuer_ref         | jsonb                    | A reference to the issuer for this certificate.                                                                                                                    |
| kind               | text                     | Type of resource.                                                                                                                                                  |
| labels             | jsonb                    | Map of string keys and values that can be used to organize and categorize (scope and select) objects. May match selectors of replication controllers and services. |
| name               | text                     | Name of resource.                                                                                                                                                  |
| namespace          | text                     | Namespace defines the space within which each name must be unique.                                                                                                 |
| path               | text                     | The path to the manifest file.                                                                                                                                     |
| renew_before       | text                     | How long before the currently issued certificate's expiry cert-manager should renew the certificate. Default to 2/3 of the duration.                               |
| secret_name        | text                     | The name of the secret resource that will be automatically created and managed by this Certificate resource.                                                       |
| source_type        | text                     | The source of the resource. Possible values are: deployed and manifest. If the resource is fetched from the spec file the value will be manifest.                  |
| start_line         | bigint                   | The path to the manifest file.                                                                                                                                     |
| uid                | text                     | UID is the unique in time and space value for this object.                                                                                                         |
+--------------------+--------------------------+--------------------------------------------------------------------------------------------------------------------------------------------------------------------+
```

```bash
> select name, uid, kind, api_version, namespace from kubernetes_certificate;
+------------------------------------+--------------------------------------+-------------+--------------------+-----------+
| name                               | uid                                  | kind        | api_version        | namespace |
+------------------------------------+--------------------------------------+-------------+--------------------+-----------+
| temporal-w-spcloudt6t6sk7toegg-tls | 5ccd69be-6e73-4edc-8c1d-bccd6a1e6e38 | Certificate | cert-manager.io/v1 | default   |
+------------------------------------+--------------------------------------+-------------+--------------------+-----------+
```

## Examples

### List all certificates

```sql
select
  name,
  uid,
  namespace,
  creation_timestamp,
  api_version
from
  kubernetes_certificate;
```

### List certificates added in the last 24 hours

```sql
select
  name,
  uid,
  namespace,
  creation_timestamp,
  api_version
from
  kubernetes_certificate
where
  creation_timestamp = now() - interval '24 hrs';
```

### List ISCA certificates

```sql
select
  name,
  uid,
  namespace,
  creation_timestamp,
  api_version
from
  kubernetes_certificate
where
  is_ca;
```

### List expired certificates

```sql
select
  name,
  uid,
  namespace,
  creation_timestamp,
  api_version
from
  kubernetes_certificate
where
  now() > to_timestamp(not_after,'YYYY-MM-DD');
```
