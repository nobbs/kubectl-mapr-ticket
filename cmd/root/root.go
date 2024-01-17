package root

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/nobbs/kubectl-mapr-ticket/cmd/claim"
	"github.com/nobbs/kubectl-mapr-ticket/cmd/common"
	"github.com/nobbs/kubectl-mapr-ticket/cmd/secret"
	"github.com/nobbs/kubectl-mapr-ticket/cmd/version"
	"github.com/nobbs/kubectl-mapr-ticket/cmd/volume"
	"github.com/nobbs/kubectl-mapr-ticket/pkg/util"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericiooptions"
)

const (
	rootUse   = `%[1]s`
	rootShort = "A kubectl plugin to list and inspect MapR tickets"
	rootLong  = `
		A kubectl plugin that allows you to list and inspect MapR tickets from a
		Kubernetes cluster, including details stored in the ticket itself without
		requiring access to the MapR cluster.
		`
)

func NewCmd(flags *genericclioptions.ConfigFlags, streams genericiooptions.IOStreams) *cobra.Command {
	o := common.NewOptions(
		flags,
		streams,
	)

	rootCmd := &cobra.Command{
		Use:   fmt.Sprintf(rootUse, filepath.Base(os.Args[0])),
		Short: rootShort,
		Long:  common.CliLongDesc(rootLong),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return util.SetupLogging(o.IOStreams.ErrOut, o.Debug)
		},
	}

	// set IOStreams for the command
	rootCmd.SetIn(o.IOStreams.In)
	rootCmd.SetOut(o.IOStreams.Out)
	rootCmd.SetErr(o.IOStreams.ErrOut)

	// add default kubernetes flags as global flags
	o.KubernetesConfigFlags.AddFlags(rootCmd.PersistentFlags())

	// add own global flags
	rootCmd.PersistentFlags().BoolVar(&o.Debug, "debug", false, "Enable debug logging")

	// add subcommands
	rootCmd.AddCommand(
		claim.NewCmd(o),
		secret.NewCmd(o),
		version.NewCmd(o),
		volume.NewCmd(o),
	)

	// add completions
	err := rootCmd.RegisterFlagCompletionFunc("namespace", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		client, err := util.ClientFromFlags(o.KubernetesConfigFlags)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		return common.CompleteNamespaceNames(client, toComplete)
	})
	if err != nil {
		panic(err)
	}

	return rootCmd
}
