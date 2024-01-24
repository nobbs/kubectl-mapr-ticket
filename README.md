# kubectl-mapr-ticket

[![GitHub Release](https://img.shields.io/github/v/release/nobbs/kubectl-mapr-ticket)](https://github.com/nobbs/kubectl-mapr-ticket/releases/latest)
[![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/nobbs/kubectl-mapr-ticket/ci.yaml?branch=release-please--branches--main&label=CI%20Build)](https://github.com/nobbs/kubectl-mapr-ticket/actions/workflows/ci.yaml)
[![Codecov](https://img.shields.io/codecov/c/github/nobbs/kubectl-mapr-ticket)](https://app.codecov.io/gh/nobbs/kubectl-mapr-ticket)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/nobbs/kubectl-mapr-ticket)

`kubectl-mapr-ticket` is a `kubectl` plugin that allows you to list and inspect MapR tickets deployed as Kubernetes secrets in a cluster.

MapR tickets are used by the [MapR CSI driver](https://github.com/mapr/mapr-csi) to authenticate and authorize access to Persistent Volumes backed by MapR storage.

## Installation

### Using `krew`

The easiest way to install the plugin is using the [krew](https://krew.sigs.k8s.io/) plugin manager for `kubectl`. Once you have `krew` installed, you can install the plugin as follows:

```console
$ kubectl krew install mapr-ticket
$ kubectl mapr-ticket --help
```

### Using Release Binaries

You can download the latest release binaries from the [releases page](https://github.com/nobbs/kubectl-mapr-ticket/releases). Binaries are available for Linux and macOS for both AMD64 and ARM64 architectures.

<!-- x-release-please-start-version -->

Example installation of `v0.2.1` for Apple Silicon (ARM64) macOS:

```console
$ curl -LO https://github.com/nobbs/kubectl-mapr-ticket/releases/download/v0.2.1/kubectl-mapr-ticket-arm64-darwin.tar.gz
$ tar -xvf kubectl-mapr-ticket-arm64-darwin.tar.gz
$ mv ./kubectl-mapr-ticket /usr/local/bin
$ kubectl mapr-ticket --help
```

<!-- x-release-please-end -->

### From Source

To install from source, you will need to have [Go](https://golang.org/) installed on your system. Once you have Go installed, you can build the plugin as follows:

```console
$ git clone https://github.com/nobbs/kubectl-mapr-ticket.git
$ cd kubectl-mapr-ticket && CGO_ENABLED=0 go build -buildvcs=true -o ./bin/kubectl-mapr-ticket ./cmd && mv ./bin/kubectl-mapr-ticket /usr/local/bin
$ kubectl mapr-ticket --help
```

## Usage

The plugin can be invoked using the `kubectl mapr-ticket` command. The plugin supports the following subcommands:

- `secret`, alias `s` - List all secrets containing MapR tickets in the current namespace and print some information about them.
- `volume`, alias `pv` - List all persistent volumes that use the specified MapR ticket secret and print some information about them.
- `claim`, alias `pvc` - List all persistent volume claims that use a MapR ticket in the current namespace.

### Secrets

The `secret` subcommand will list all MapR tickets deployed as `Secrets` in the current namespace. The output by default will be a table that can be extended with the `--output wide` flag. Additional flags can be used to customize the output, see `kubectl mapr-ticket secret --help` for more details.

```console
$ kubectl mapr-ticket secret
NAME                      MAPR CLUSTER         USER     STATUS              AGE
mapr-dev-ticket-user-a    demo.dev.mapr.com    user_a   Valid (4y left)     75d
mapr-dev-ticket-user-b    demo.dev.mapr.com    user_b   Valid (4y left)     75d
mapr-dev-ticket-user-c    demo.dev.mapr.com    user_c   Valid (4y left)     75d
mapr-prod-ticket-user-a   demo.prod.mapr.com   user_a   Expired (43d ago)   73d
mapr-prod-ticket-user-b   demo.prod.mapr.com   user_b   Expired (43d ago)   73d
mapr-prod-ticket-user-c   demo.prod.mapr.com   user_c   Expired (43d ago)   73d
```

### Volumes

The `volume` subcommand will list all Persistent Volumes that are using a specific MapR ticket if a secret name is specified, or any ticket in the current namespace if no argument is provided. The output by default will be a table with the following columns. Additional flags can be used to customize the output, see `kubectl mapr-ticket volume --help` for more details.

```console
$ kubectl mapr-ticket mapr-ticket-secret -n test-csi
NAME             SECRET NAMESPACE   SECRET               CLAIM NAMESPACE   CLAIM        TICKET STATUS         AGE
test-static-pv   test-csi           mapr-ticket-secret   default           test-claim   Not found / Invalid   13h
var-lib-mapr     test-csi           mapr-ticket-secret   default           test-var     Valid (4y left)       12d
expired-pv       test-csi           mapr-ticket-secret   default           test-exp     Expired (43d ago)     12d
```

### Claims

The `claim` subcommand will list all Persistent Volume Claims in the current namespace that are using a MapR ticket. The output by default will be a table with the following columns. Additional flags can be used to customize the output, see `kubectl mapr-ticket claim --help` for more details.

```console
$ kubectl mapr-ticket claim -n default
NAME         SECRET NAMESPACE   SECRET               VOLUME NAME      TICKET STATUS         AGE
test-claim   test-csi           mapr-ticket-secret   test-static-pv   Not found / Invalid   13h
test-var     test-csi           mapr-ticket-secret   var-lib-mapr     Valid (4y left)       12d
test-exp     test-csi           mapr-ticket-secret   expired-pv       Expired (43d ago)     12d
```

### Shell Completion

The plugin supports shell completion for various shells. To enable shell completion, you will need to source the completion script for your shell. For example, to enable completion for `zsh`, you can run the following command:

```console
$ source <(kubectl mapr-ticket completion zsh)
```

Note, that this is only local to your current shell session. To enable completion permanently, you either need to add the command to your shell profile or place the completion script in the appropriate location for your shell.

Unfortunately, the above setup will only provide completion for the `kubectl-mapr_ticket` command, not the actual `kubectl mapr-ticket` alias. To enable completion for the alias, you need to create a special `kubectl_complete-mapr_ticket` executable in your `PATH` that will delegate execution to the plugin. You can find an example of it in the [hack](hack) directory. Place the script somewhere in your `PATH` and make sure it is executable, e.g. by running:

```console
$ curl -LO https://github.com/nobbs/kubectl-mapr-ticket/raw/main/hack/kubectl_complete-mapr_ticket
$ chmod +x ./kubectl_complete-mapr_ticket
$ mv ./kubectl_complete-mapr_ticket /usr/local/bin
```

## Does this require a connection to a MapR cluster?

**No, this `kubectl` plugin does not require a connection to a MapR cluster.** The plugin will inspect the secrets in the current namespace, filter them down to those that are MapR tickets, and then decode the ticket contents using [this reverse-engineered ticket parser](https://github.com/nobbs/mapr-ticket-parser) which is based on this [blog post of mine](https://nobbs.dev/posts/reverse-engineering-mapr-ticket-format/).

Based on testing, the plugin is able to parse tickets starting at least from MapR 6.0.0 as the format did not receive any breaking changes since then.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
