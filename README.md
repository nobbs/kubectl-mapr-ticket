# kubectl-mapr-ticket

`kubectl-mapr-ticket` is a `kubectl` plugin that allows you to list and inspect MapR tickets deployed as Kubernetes secrets in a cluster.

MapR tickets are used by the [MapR CSI driver](https://github.com/mapr/mapr-csi) to authenticate and authorize access to Persistent Volumes backed by MapR storage.

## Installation

### From Source

To install from source, you will need to have [Go](https://golang.org/) installed on your system. Once you have Go installed, you can build the plugin as follows:

```console
$ git clone https://github.com/nobbs/kubectl-mapr-ticket.git
$ cd kubectl-mapr-ticket && CGO_ENABLED=0 go build -buildvcs=true -o ./bin/kubectl-mapr-ticket ./cmd && mv ./bin/kubectl-mapr-ticket /usr/local/bin
$ kubectl mapr-ticket --help
```

### Using Release Binaries

You can download the latest release binaries from the [releases page](https://github.com/nobbs/kubectl-mapr-ticket/releases). Binaries are available for Linux and macOS for both AMD64 and ARM64 architectures.

<!-- x-release-please-start-version -->
Example installation of `v0.1.2` for Apple Silicon (ARM64) macOS:

```console
$ curl -LO https://github.com/nobbs/kubectl-mapr-ticket/releases/download/v0.1.2/kubectl-mapr-ticket-arm64-darwin.tar.gz
$ tar -xvf kubectl-mapr-ticket-arm64-darwin.tar.gz
$ mv ./kubectl-mapr-ticket /usr/local/bin
$ kubectl mapr-ticket --help
```
<!-- x-release-please-end -->

## Usage

The plugin can be invoked using the `kubectl mapr-ticket` command. The plugin supports the following subcommands:

- `list` - List all MapR tickets deployed in the current namespace.
- `used-by` - List all Persistent Volumes that are using a specific MapR ticket.

### List

The `list` subcommand will list all MapR tickets deployed in the current namespace. The output by default will be a table with the following columns. Additional flags can be used to customize the output, see `kubectl mapr-ticket list --help` for more details.

```console
$ kubectl mapr-ticket list
NAME                      MAPR CLUSTER         USER     STATUS              AGE
mapr-dev-ticket-user-a    demo.dev.mapr.com    user_a   Valid (4y left)     75d
mapr-dev-ticket-user-b    demo.dev.mapr.com    user_b   Valid (4y left)     75d
mapr-dev-ticket-user-c    demo.dev.mapr.com    user_c   Valid (4y left)     75d
mapr-prod-ticket-user-a   demo.prod.mapr.com   user_a   Expired (43d ago)   73d
mapr-prod-ticket-user-b   demo.prod.mapr.com   user_b   Expired (43d ago)   73d
mapr-prod-ticket-user-c   demo.prod.mapr.com   user_c   Expired (43d ago)   73d
```

### Used By

The `used-by` subcommand will list all Persistent Volumes that are using a specific MapR ticket or any ticket in the current namespace if `--all` is specified. The output by default will be a table with the following columns. Additional flags can be used to customize the output, see `kubectl mapr-ticket used-by --help` for more details.

```console
$ kubectl mapr-ticket mapr-ticket-secret -n test-csi
NAME             SECRET NAMESPACE   SECRET               CLAIM NAMESPACE   CLAIM   AGE
test-static-pv   test-csi           mapr-ticket-secret                             13h
```

## Does this require a connection to a MapR cluster?

**No, this `kubectl` plugin does not require a connection to a MapR cluster.** The plugin will inspect the secrets in the current namespace, filter them down to those that are MapR tickets, and then decode the ticket contents using [this reverse-engineered ticket parser](https://github.com/nobbs/mapr-ticket-parser) which is based on this [blog post of mine](https://nobbs.dev/posts/reverse-engineering-mapr-ticket-format/).

Based on testing, the plugin is able to parse tickets starting at least from MapR 6.0.0 as the format did not receive any breaking changes since then.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
