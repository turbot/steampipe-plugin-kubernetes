![image](https://hub.steampipe.io/images/plugins/turbot/kubernetes-social-graphic.png)

# Kubernetes Plugin for Steampipe

Use SQL to query Kubernetes components.

- **[Get started →](https://hub.steampipe.io/plugins/turbot/kubernetes)**
- Documentation: [Table definitions & examples](https://hub.steampipe.io/plugins/turbot/kubernetes/tables)
- Community: [Slack Channel](https://join.slack.com/t/steampipe/shared_invite/zt-oij778tv-lYyRTWOTMQYBVAbtPSWs3g)
- Get involved: [Issues](https://github.com/turbot/steampipe-plugin-kubernetes/issues)

## Quick start

Install the plugin with [Steampipe](https://steampipe.io):

```shell
steampipe plugin install kubernetes
```

Run a query:

```sql
select name, creation_timestamp, addresses, capacity from kubernetes_node;
```

## Developing

Prerequisites:

- [Steampipe](https://steampipe.io/downloads)
- [Golang](https://golang.org/doc/install)

Clone:

```sh
git clone git@github.com:turbot/steampipe-plugin-kubernetes
cd steampipe-plugin-kubernetes
```

Build, which automatically installs the new version to your `~/.steampipe/plugins` directory:

```
make
```

Configure the plugin:

```
cp config/* ~/.steampipe/config
vi ~/.steampipe/config/kubernetes.spc
```

Try it!

```
steampipe query
> .inspect kubernetes
```

Further reading:

- [Writing plugins](https://steampipe.io/docs/develop/writing-plugins)
- [Writing your first table](https://steampipe.io/docs/develop/writing-your-first-table)

## Contributing

Please see the [contribution guidelines](https://github.com/turbot/steampipe/blob/main/CONTRIBUTING.md) and our [code of conduct](https://github.com/turbot/steampipe/blob/main/CODE_OF_CONDUCT.md). All contributions are subject to the [Apache 2.0 open source license](https://github.com/turbot/steampipe-plugin-kubernetes/blob/main/LICENSE).

`help wanted` issues:

- [Steampipe](https://github.com/turbot/steampipe/labels/help%20wanted)
- [Kubernetes Plugin](https://github.com/turbot/steampipe-plugin-kubernetes/labels/help%20wanted)
