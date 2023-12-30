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
Example installation of `v0.1.0` for Apple Silicon (ARM64) macOS:

```console
$ curl -LO https://github.com/nobbs/kubectl-mapr-ticket/releases/download/v0.1.0/kubectl-mapr-ticket-arm64-darwin.tar.gz
$ tar -xvf kubectl-mapr-ticket-arm64-darwin.tar.gz
$ mv ./kubectl-mapr-ticket /usr/local/bin
$ kubectl mapr-ticket --help
```
<!-- x-release-please-end -->

## Usage

Currently, `kubectl-mapr-ticket` supports only the `list` command. This command will list all MapR tickets deployed in the current namespace. The output will include the name of the secret, the MapR cluster name, the username, and the expiry date of the ticket.

```console
$ kubectl mapr-ticket list
NAME                      MAPR CLUSTER         USER     EXPIRATION
mapr-dev-ticket-user-a    demo.dev.mapr.com    user_a   2028-11-09T15:50:57+01:00 (Expired)
mapr-dev-ticket-user-b    demo.dev.mapr.com    user_b   2028-11-09T15:50:55+01:00 (Expired)
mapr-dev-ticket-user-c    demo.dev.mapr.com    user_c   2028-11-09T15:50:51+01:00 (Expired)
mapr-prod-ticket-user-a   demo.prod.mapr.com   user_a   2023-11-19T08:47:05+01:00
mapr-prod-ticket-user-b   demo.prod.mapr.com   user_b   2023-11-19T08:47:03+01:00
mapr-prod-ticket-user-c   demo.prod.mapr.com   user_c   2023-11-19T08:47:02+01:00
```

## Does this require a connection to a MapR cluster?

**No, this `kubectl` plugin does not require a connection to a MapR cluster.** The plugin will inspect the secrets in the current namespace, filter them down to those that are MapR tickets, and then decode the ticket contents using [this reverse-engineered ticket parser](https://github.com/nobbs/mapr-ticket-parser) which is based on this [blog post of mine](https://nobbs.dev/posts/reverse-engineering-mapr-ticket-format/).

Based on testing, the plugin is able to parse tickets starting at least from MapR 6.0.0 as the format did not receive any breaking changes since then.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
