# steampipe-plugin-kuberenetes

A steampipe plugin for kubernetes.   This is a WIP and not yet ready for use....

To build, run `make`.  Copy the `config/kubernetes.spc` file to `~/.steampipe/config/` to create a  connection.  Currently, uses kubectl current context.



Notes...
- Steampipe Standard Fields:
    - `title`:  `name` (from metadata)
    - Omit for now:  `akas`: `uid`  (from metadata)
    - `tags`: merge `labels` and `annotations`  (from metadata).  If there is a name collision, prefer `label`

- K8S standard fields:
    - all the metadata fields?
        - include `namespace` for non-namespaced resources (node, namespace, etc) ??
        - Omit `cluster_name` 
            - Per https://v1-18.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#objectmeta-v1-meta:  This field is not set anywhere right now and apiserver is going to ignore it if set in create or update request.

        - Omit: `self_link`:
            - DEPRECATED: Kubernetes will stop propagating this field in 1.20 release and the field is planned to be removed in 1.21 release.

        - omit `managed_fields`:
           - 'users typically shouldn't need to set or understand this field.' Will add later if requested / required.

    - We need the `cluster_name` or `kubectl_context` to disambiguate connections...

    - `Spec` and `status` should be expanded into individual columns, not included as a giant JSON column
        - do not prefix with with `spec_` and `status_` ?

    - Why are the TypeMeta columns?  They're always missing??


- Auth / creds / Scope: 
    - support these.... https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs

    - 1 connection = 1 context from kubeconfig file

    - use kubeconfig file
        - default `config_path` to `KUBE_CONFIG_PATH`, then "~/.kube/config"
        - default `config_context` to current active (default) config


- session cache is not working????

