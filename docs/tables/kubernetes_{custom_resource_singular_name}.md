# Table: kubernetes\_{custom_resource_singular_name}

Query data from the custom resource called `kubernetes_{custom_resource_singular_name}`, e.g., `kubernetes_certificate`, `kubernetes_capacityrequest`. A table is automatically created to represent each custom resource.

If the table name is already created in the above format and exist in the table list, then the subsequent ones will have the fully qualified name `kubernetes_{custom_resource_singular_name}_{custom_resource_group_name}`, e.g., `kubernetes_certificate_cert_manager_io`.

For instance, given the CRD `certManager.yaml`:

```yml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
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
    # table name - kubernetes_certificate
    singular: certificate
    categories:
      - cert-manager
  scope: Namespaced
  versions:
    - name: v1
      subresources:
        status: {}
      additionalPrinterColumns:
        - jsonPath: .status.conditions[?(@.type=="Ready")].status
          name: Ready
          type: string
        - jsonPath: .spec.secretName
          name: Secret
          type: string
        - jsonPath: .spec.issuerRef.name
          name: Issuer
          priority: 1
          type: string
        - jsonPath: .status.conditions[?(@.type=="Ready")].message
          name: Status
          priority: 1
          type: string
        - jsonPath: .metadata.creationTimestamp
          description: CreationTimestamp is a timestamp representing the server time when this object was created. It is not guaranteed to be set in happens-before order across separate operations. Clients may not set this value. It is represented in RFC3339 form and is in UTC.
          name: Age
          type: date
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
                - issuerRef
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
                issuerRef:
                  description: IssuerRef is a reference to the issuer for this certificate. If the `kind` field is not set, or set to `Issuer`, an Issuer resource with the given name in the same namespace as the Certificate will be used. If the `kind` field is set to `ClusterIssuer`, a ClusterIssuer with the provided name will be used. The `name` field in this stanza is required at all times.
                  type: object
                  required:
                    - name
                  properties:
                    group:
                      description: Group of the resource being referred to.
                      type: string
                    kind:
                      description: Kind of the resource being referred to.
                      type: string
                    name:
                      description: Name of the resource being referred to.
                      type: string
                keystores:
                  description: Keystores configures additional keystore output formats stored in the `secretName` Secret resource.
                  type: object
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
                uris:
                  description: URIs is a list of URI subjectAltNames to be set on the Certificate.
                  type: array
                  items:
                    type: string
                usages:
                  description: Usages is the set of x509 usages that are requested for the certificate. Defaults to `digital signature` and `key encipherment` if not specified.
                  type: array
                  items:
                    description: "KeyUsage specifies valid usage contexts for keys. See: https://tools.ietf.org/html/rfc5280#section-4.2.1.3 https://tools.ietf.org/html/rfc5280#section-4.2.1.12 \n Valid KeyUsage values are as follows: \"signing\", \"digital signature\", \"content commitment\", \"key encipherment\", \"key agreement\", \"data encipherment\", \"cert sign\", \"crl sign\", \"encipher only\", \"decipher only\", \"any\", \"server auth\", \"client auth\", \"code signing\", \"email protection\", \"s/mime\", \"ipsec end system\", \"ipsec tunnel\", \"ipsec user\", \"timestamping\", \"ocsp signing\", \"microsoft sgc\", \"netscape sgc\""
                    type: string
            status:
              description: Status of the Certificate. This is set and managed automatically.
              type: object
              properties:
                conditions:
                  description: List of status conditions to indicate the status of certificates. Known condition types are `Ready` and `Issuing`.
                  type: array
                  items:
                    description: CertificateCondition contains condition information for an Certificate.
                    type: object
                    required:
                      - status
                      - type
                    properties:
                      lastTransitionTime:
                        description: LastTransitionTime is the timestamp corresponding to the last status change of this condition.
                        type: string
                        format: date-time
                      message:
                        description: Message is a human readable description of the details of the last transition, complementing reason.
                        type: string
                      observedGeneration:
                        description: If set, this represents the .metadata.generation that the condition was set based upon. For instance, if .metadata.generation is currently 12, but the .status.condition[x].observedGeneration is 9, the condition is out of date with respect to the current state of the Certificate.
                        type: integer
                        format: int64
                      reason:
                        description: Reason is a brief machine readable explanation for the condition's last transition.
                        type: string
                      status:
                        description: Status of the condition, one of (`True`, `False`, `Unknown`).
                        type: string
                        enum:
                          - "True"
                          - "False"
                          - Unknown
                      type:
                        description: Type of the condition, known values are (`Ready`, `Issuing`).
                        type: string
                  x-kubernetes-list-map-keys:
                    - type
                  x-kubernetes-list-type: map
                lastFailureTime:
                  description: LastFailureTime is the time as recorded by the Certificate controller of the most recent failure to complete a CertificateRequest for this Certificate resource. If set, cert-manager will not re-request another Certificate until 1 hour has elapsed from this time.
                  type: string
                  format: date-time
                notAfter:
                  description: The expiration time of the certificate stored in the secret named by this resource in `spec.secretName`.
                  type: string
                  format: date-time
                notBefore:
                  description: The time after which the certificate stored in the secret named by this resource in spec.secretName is valid.
                  type: string
                  format: date-time
                renewalTime:
                  description: RenewalTime is the time at which the certificate will be next renewed. If not set, no upcoming renewal is scheduled.
                  type: string
                  format: date-time
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

This plugin will automatically create a table called `kubernetes_certificate`:

```
> select name, uid, kind, api_version, namespace from kubernetes_certificate;
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
| kubernetes_certificate                | A Certificate resource should be created to ensure an up to date and signed x509 certificate is stored in the Kubernetes Secret resource named in `spec.secretName`. The stored certificate will be renewed before it expires (as configured by `spec.renewBefore`).Custom resource for certificates.cert-manager.io.                                                                                                         |
| kubernetes_cluster_role               | ClusterRole contains rules that represent a set of permissions.                                                                                                  |
| kubernetes_cluster_role_binding       | A ClusterRoleBinding grants the permissions defined in a cluster role to a user or set of users. Access granted by ClusterRoleBinding is cluster-wide.           |
| ...                                   | ...                                                                                                                                                              |
+---------------------------------------+------------------------------------------------------------------------------------------------------------------------------------------------------------------+
```

To get details of a specific custom resource table, inspect it by name:

```sql
.inspect kubernetes_certificate;
+---------------------------+--------------------------+-------------------------------------------------------------------------------------------------------------+
| column                    | type                     | description                                                                                                 |
+---------------------------+--------------------------+-------------------------------------------------------------------------------------------------------------+
| _ctx                      | jsonb                    | Steampipe context in JSON form, e.g. connection_name.                                                       |
| additional_output_formats | jsonb                    | AdditionalOutputFormats defines extra output formats of the private key and signed certificate chain to be  |
|                           |                          | written to this Certificate's target Secret. This is an Alpha Feature and is only enabled with the `--featu |
|                           |                          | re-gates=AdditionalCertificateOutputFormats=true` option on both the controller and webhook components.     |
| api_version               | text                     | The API version of the resource.                                                                            |
| common_name               | text                     | CommonName is a common name to be used on the Certificate. The CommonName should have a length of 64 charac |
|                           |                          | ters or fewer to avoid generating invalid CSRs. This value is ignored by TLS clients when any subject alt n |
|                           |                          | ame is set. This is x509 behaviour: https://tools.ietf.org/html/rfc6125#section-6.4.4                       |
| conditions                | jsonb                    | List of status conditions to indicate the status of certificates. Known condition types are `Ready` and `Is |
|                           |                          | suing`.                                                                                                     |
| creation_timestamp        | timestamp with time zone | CreationTimestamp is a timestamp representing the server time when this object was created.                 |
| dns_names                 | jsonb                    | DNSNames is a list of DNS subjectAltNames to be set on the Certificate.                                     |
| duration                  | text                     | The requested 'duration' (i.e. lifetime) of the Certificate. This option may be ignored/overridden by some  |
|                           |                          | issuer types. If unset this defaults to 90 days. Certificate will be renewed either 2/3 through its duratio |
|                           |                          | n or `renewBefore` period before its expiry, whichever is later. Minimum accepted duration is 1 hour. Value |
|                           |                          |  must be in units accepted by Go time.ParseDuration https://golang.org/pkg/time/#ParseDuration              |
| email_addresses           | jsonb                    | EmailAddresses is a list of email subjectAltNames to be set on the Certificate.                             |
| encode_usages_in_request  | boolean                  | EncodeUsagesInRequest controls whether key usages should be present in the CertificateRequest               |
| ip_addresses              | jsonb                    | IPAddresses is a list of IP address subjectAltNames to be set on the Certificate.                           |
| is_ca                     | boolean                  | IsCA will mark this Certificate as valid for certificate signing. This will automatically add the `cert sig |
|                           |                          | n` usage to the list of `usages`.                                                                           |
| issuer_ref                | jsonb                    | IssuerRef is a reference to the issuer for this certificate. If the `kind` field is not set, or set to `Iss |
|                           |                          | uer`, an Issuer resource with the given name in the same namespace as the Certificate will be used. If the  |
|                           |                          | `kind` field is set to `ClusterIssuer`, a ClusterIssuer with the provided name will be used. The `name` fie |
|                           |                          | ld in this stanza is required at all times.                                                                 |
| keystores                 | jsonb                    | Keystores configures additional keystore output formats stored in the `secretName` Secret resource.         |
| kind                      | text                     | Type of resource.                                                                                           |
| labels                    | jsonb                    | Map of string keys and values that can be used to organize and categorize (scope and select) objects. May m |
|                           |                          | atch selectors of replication controllers and services.                                                     |
| last_failure_time         | text                     | LastFailureTime is the time as recorded by the Certificate controller of the most recent failure to complet |
|                           |                          | e a CertificateRequest for this Certificate resource. If set, cert-manager will not re-request another Cert |
|                           |                          | ificate until 1 hour has elapsed from this time.                                                            |
| literal_subject           | text                     | LiteralSubject is an LDAP formatted string that represents the [X.509 Subject field](https://datatracker.ie |
|                           |                          | tf.org/doc/html/rfc5280#section-4.1.2.6). Use this *instead* of the Subject field if you need to ensure the |
|                           |                          |  correct ordering of the RDN sequence, such as when issuing certs for LDAP authentication. See https://gith |
|                           |                          | ub.com/cert-manager/cert-manager/issues/3203, https://github.com/cert-manager/cert-manager/issues/4424. Thi |
|                           |                          | s field is alpha level and is only supported by cert-manager installations where LiteralCertificateSubject  |
|                           |                          | feature gate is enabled on both cert-manager controller and webhook.                                        |
| name                      | text                     | Name of resource.                                                                                           |
| namespace                 | text                     | Namespace defines the space within which each name must be unique.                                          |
| not_after                 | text                     | The expiration time of the certificate stored in the secret named by this resource in `spec.secretName`.    |
| not_before                | text                     | The time after which the certificate stored in the secret named by this resource in spec.secretName is vali |
|                           |                          | d.                                                                                                          |
| private_key               | jsonb                    | Options to control private keys used for the Certificate.                                                   |
| renew_before              | text                     | How long before the currently issued certificate's expiry cert-manager should renew the certificate. The de |
|                           |                          | fault is 2/3 of the issued certificate's duration. Minimum accepted value is 5 minutes. Value must be in un |
|                           |                          | its accepted by Go time.ParseDuration https://golang.org/pkg/time/#ParseDuration                            |
| renewal_time              | text                     | RenewalTime is the time at which the certificate will be next renewed. If not set, no upcoming renewal is s |
|                           |                          | cheduled.                                                                                                   |
| revision_history_limit    | bigint                   | revisionHistoryLimit is the maximum number of CertificateRequest revisions that are maintained in the Certi |
|                           |                          | ficate's history. Each revision represents a single `CertificateRequest` created by this Certificate, eithe |
|                           |                          | r when it was created, renewed, or Spec was changed. Revisions will be removed by oldest first if the numbe |
|                           |                          | r of revisions exceeds this number. If set, revisionHistoryLimit must be a value of `1` or greater. If unse |
|                           |                          | t (`nil`), revisions will not be garbage collected. Default value is `nil`.                                 |
| secret_name               | text                     | SecretName is the name of the secret resource that will be automatically created and managed by this Certif |
|                           |                          | icate resource. It will be populated with a private key and certificate, signed by the denoted issuer.      |
| secret_template           | jsonb                    | SecretTemplate defines annotations and labels to be copied to the Certificate's Secret. Labels and annotati |
|                           |                          | ons on the Secret will be changed as they appear on the SecretTemplate when added or removed. SecretTemplat |
|                           |                          | e annotations are added in conjunction with, and cannot overwrite, the base set of annotations cert-manager |
|                           |                          |  sets on the Certificate's Secret.                                                                          |
| uid                       | text                     | UID is the unique in time and space value for this object.                                                  |
| uris                      | jsonb                    | URIs is a list of URI subjectAltNames to be set on the Certificate.                                         |
| usages                    | jsonb                    | Usages is the set of x509 usages that are requested for the certificate. Defaults to `digital signature` an |
|                           |                          | d `key encipherment` if not specified.                                                                      |
+---------------------------+--------------------------+-------------------------------------------------------------------------------------------------------------+
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
