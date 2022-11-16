# Table: {custom_resource_name.group_name}

Query data from the custom resource called `{custom_resource_name.group_name}`, e.g., `certificates.cert-manager.io`, `storeconfigs.crossplane.io`. A table is automatically created to represent each object in the `objects` argument.

For instance, given the CRD `certManager.yaml`:

```yml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  # Table name of the custom resource
  name: certificates.cert-manager.io
  labels:
    app: cert-manager
    app.kubernetes.io/name: cert-manager-crd
    # Generated labels {{- include "labels" . | nindent 4 }}
spec:
  group: cert-manager.io
  names:
    kind: Certificate
    listKind: CertificateList
    plural: certificates
    shortNames:
      - cert
      - certs
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
          description: "A Certificate resource should be created to ensure an up to date and signed x509 certificate is stored in the Kubernetes Secret resource named in `spec.secretName`. \n The stored certificate will be renewed before it expires (as configured by `spec.renewBefore`)."
          type: object
          required:
            - spec
          properties:
            apiVersion:
              description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources"
              type: string
            kind:
              description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds"
              type: string
            metadata:
              type: object
            spec:
              description: Desired state of the Certificate resource.
              type: object
              required:
                - secretName
              properties:
                additionalOutputFormats:
                  description: AdditionalOutputFormats defines extra output formats of the private key and signed certificate chain to be written to this Certificate's target Secret. This is an Alpha Feature and is only enabled with the `--feature-gates=AdditionalCertificateOutputFormats=true` option on both the controller and webhook components.
                  type: array
                  items:
                    description: CertificateAdditionalOutputFormat defines an additional output format of a Certificate resource. These contain supplementary data formats of the signed certificate chain and paired private key.
                    type: object
                    required:
                      - type
                    properties:
                      type:
                        description: Type is the name of the format type that should be written to the Certificate's target Secret.
                        type: string
                        enum:
                          - DER
                          - CombinedPEM
                commonName:
                  description: "CommonName is a common name to be used on the Certificate. The CommonName should have a length of 64 characters or fewer to avoid generating invalid CSRs. This value is ignored by TLS clients when any subject alt name is set. This is x509 behaviour: https://tools.ietf.org/html/rfc6125#section-6.4.4"
                  type: string
                dnsNames:
                  description: DNSNames is a list of DNS subjectAltNames to be set on the Certificate.
                  type: array
                  items:
                    type: string
                duration:
                  description: The requested 'duration' (i.e. lifetime) of the Certificate. This option may be ignored/overridden by some issuer types. If unset this defaults to 90 days. Certificate will be renewed either 2/3 through its duration or `renewBefore` period before its expiry, whichever is later. Minimum accepted duration is 1 hour. Value must be in units accepted by Go time.ParseDuration https://golang.org/pkg/time/#ParseDuration
                  type: string
                emailAddresses:
                  description: EmailAddresses is a list of email subjectAltNames to be set on the Certificate.
                  type: array
                  items:
                    type: string
                encodeUsagesInRequest:
                  description: EncodeUsagesInRequest controls whether key usages should be present in the CertificateRequest
                  type: boolean
                ipAddresses:
                  description: IPAddresses is a list of IP address subjectAltNames to be set on the Certificate.
                  type: array
                  items:
                    type: string
                isCA:
                  description: IsCA will mark this Certificate as valid for certificate signing. This will automatically add the `cert sign` usage to the list of `usages`.
                  type: boolean
                literalSubject:
                  description: LiteralSubject is an LDAP formatted string that represents the [X.509 Subject field](https://datatracker.ietf.org/doc/html/rfc5280#section-4.1.2.6). Use this *instead* of the Subject field if you need to ensure the correct ordering of the RDN sequence, such as when issuing certs for LDAP authentication. See https://github.com/cert-manager/cert-manager/issues/3203, https://github.com/cert-manager/cert-manager/issues/4424. This field is alpha level and is only supported by cert-manager installations where LiteralCertificateSubject feature gate is enabled on both cert-manager controller and webhook.
                  type: string
                privateKey:
                  description: Options to control private keys used for the Certificate.
                  type: object
                  properties:
                    algorithm:
                      description: Algorithm is the private key algorithm of the corresponding private key for this certificate. If provided, allowed values are either `RSA`,`Ed25519` or `ECDSA` If `algorithm` is specified and `size` is not provided, key size of 256 will be used for `ECDSA` key algorithm and key size of 2048 will be used for `RSA` key algorithm. key size is ignored when using the `Ed25519` key algorithm.
                      type: string
                      enum:
                        - RSA
                        - ECDSA
                        - Ed25519
                    encoding:
                      description: The private key cryptography standards (PKCS) encoding for this certificate's private key to be encoded in. If provided, allowed values are `PKCS1` and `PKCS8` standing for PKCS#1 and PKCS#8, respectively. Defaults to `PKCS1` if not specified.
                      type: string
                      enum:
                        - PKCS1
                        - PKCS8
                    rotationPolicy:
                      description: RotationPolicy controls how private keys should be regenerated when a re-issuance is being processed. If set to Never, a private key will only be generated if one does not already exist in the target `spec.secretName`. If one does exists but it does not have the correct algorithm or size, a warning will be raised to await user intervention. If set to Always, a private key matching the specified requirements will be generated whenever a re-issuance occurs. Default is 'Never' for backward compatibility.
                      type: string
                      enum:
                        - Never
                        - Always
                    size:
                      description: Size is the key bit size of the corresponding private key for this certificate. If `algorithm` is set to `RSA`, valid values are `2048`, `4096` or `8192`, and will default to `2048` if not specified. If `algorithm` is set to `ECDSA`, valid values are `256`, `384` or `521`, and will default to `256` if not specified. If `algorithm` is set to `Ed25519`, Size is ignored. No other values are allowed.
                      type: integer
                renewBefore:
                  description: How long before the currently issued certificate's expiry cert-manager should renew the certificate. The default is 2/3 of the issued certificate's duration. Minimum accepted value is 5 minutes. Value must be in units accepted by Go time.ParseDuration https://golang.org/pkg/time/#ParseDuration
                  type: string
                revisionHistoryLimit:
                  description: revisionHistoryLimit is the maximum number of CertificateRequest revisions that are maintained in the Certificate's history. Each revision represents a single `CertificateRequest` created by this Certificate, either when it was created, renewed, or Spec was changed. Revisions will be removed by oldest first if the number of revisions exceeds this number. If set, revisionHistoryLimit must be a value of `1` or greater. If unset (`nil`), revisions will not be garbage collected. Default value is `nil`.
                  type: integer
                  format: int32
                secretName:
                  description: SecretName is the name of the secret resource that will be automatically created and managed by this Certificate resource. It will be populated with a private key and certificate, signed by the denoted issuer.
                  type: string
                secretTemplate:
                  description: SecretTemplate defines annotations and labels to be copied to the Certificate's Secret. Labels and annotations on the Secret will be changed as they appear on the SecretTemplate when added or removed. SecretTemplate annotations are added in conjunction with, and cannot overwrite, the base set of annotations cert-manager sets on the Certificate's Secret.
                  type: object
                  properties:
                    annotations:
                      description: Annotations is a key value map to be copied to the target Kubernetes Secret.
                      type: object
                      additionalProperties:
                        type: string
                    labels:
                      description: Labels is a key value map to be copied to the target Kubernetes Secret.
                      type: object
                      additionalProperties:
                        type: string
                usages:
                  description: Usages is the set of x509 usages that are requested for the certificate. Defaults to `digital signature` and `key encipherment` if not specified.
                  type: array
                  items:
                    description: "KeyUsage specifies valid usage contexts for keys. See: https://tools.ietf.org/html/rfc5280#section-4.2.1.3 https://tools.ietf.org/html/rfc5280#section-4.2.1.12 \n Valid KeyUsage values are as follows: \"signing\", \"digital signature\", \"content commitment\", \"key encipherment\", \"key agreement\", \"data encipherment\", \"cert sign\", \"crl sign\", \"encipher only\", \"decipher only\", \"any\", \"server auth\", \"client auth\", \"code signing\", \"email protection\", \"s/mime\", \"ipsec end system\", \"ipsec tunnel\", \"ipsec user\", \"timestamping\", \"ocsp signing\", \"microsoft sgc\", \"netscape sgc\""
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

This plugin will automatically create a table called `certificates.cert-manager.io`:

```
> select name, uid, kind, api_version, namespace from "certificates.cert-manager.io";
+------------------------------------+--------------------------------------+-------------+--------------------+-----------+
| name                               | uid                                  | kind        | api_version        | namespace |
+------------------------------------+--------------------------------------+-------------+--------------------+-----------+
| temporal-w-spcloudt6t6sk7toegg-tls | 5ccd69be-6e73-4edc-8c1d-bccd6a1e6e38 | Certificate | cert-manager.io/v1 | default   |
+------------------------------------+--------------------------------------+-------------+--------------------+-----------+
```

## Examples

### Inspect the table structure

List all tables:

```sql
.inspect kubernetes;
+---------------------------------------+------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| table                                 | description                                                                                                                                                      |
+---------------------------------------+------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| certificates.cert-manager.io          | Represents Custom resource certificates.cert-manager.io.                                                                                                         |
| kubernetes_cluster_role               | ClusterRole contains rules that represent a set of permissions.                                                                                                  |
| kubernetes_cluster_role_binding       | A ClusterRoleBinding grants the permissions defined in a cluster role to a user or set of users. Access granted by ClusterRoleBinding is cluster-wide.           |
| ...                                   | ...
```

To get details of a specific custom resource table, inspect it by name:

```sql
select column_name, data_type from INFORMATION_SCHEMA.COLUMNS where table_name like 'certificates.cert-manager.io';
+-------------------------+--------------------------+
| column_name             | data_type                |
+-------------------------+--------------------------+
| _ctx                    | jsonb                    |
| additionalOutputFormats | jsonb                    |
| emailAddresses          | jsonb                    |
| ipAddresses             | jsonb                    |
| subject                 | jsonb                    |
| isCA                    | boolean                  |
| issuerRef               | jsonb                    |
| revisionHistoryLimit    | bigint                   |
| dnsNames                | jsonb                    |
| encodeUsagesInRequest   | boolean                  |
| secretTemplate          | jsonb                    |
| creation_timestamp      | timestamp with time zone |
| labels                  | jsonb                    |
| keystores               | jsonb                    |
| privateKey              | jsonb                    |
| uris                    | jsonb                    |
| usages                  | jsonb                    |
| uid                     | text                     |
| kind                    | text                     |
| api_version             | text                     |
| namespace               | text                     |
| secretName              | text                     |
| commonName              | text                     |
| duration                | text                     |
| name                    | text                     |
| renewBefore             | text                     |
| literalSubject          | text                     |
+-------------------------+--------------------------+
```

### List all certificates

```sql
select
  name,
  uid,
  namespace,
  creation_timestamp,
  api_version
from
  "certificates.cert-manager.io";
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
  "certificates.cert-manager.io"
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
  "certificates.cert-manager.io"
where
  "isCA";
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
  "certificates.cert - manager.io"
where
  now() > creation_timestamp + (interval '1 hour' * nullif(left("renewbefore", length("renewbefore") - 1), '')::int);
```
