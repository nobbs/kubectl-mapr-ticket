package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/nobbs/kubectl-mapr-ticket/internal/util"

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

type rootCmdOptions struct {
	kubernetesConfigFlags *genericclioptions.ConfigFlags
	IOStreams             genericiooptions.IOStreams

	// debug flag to enable debug logging
	debug bool
}

func NewCmdOptions(kubernetesConfigFlags *genericclioptions.ConfigFlags, streams genericiooptions.IOStreams) *rootCmdOptions {
	return &rootCmdOptions{
		kubernetesConfigFlags: kubernetesConfigFlags,
		IOStreams:             streams,
	}
}

func NewRootCmd(flags *genericclioptions.ConfigFlags, streams genericiooptions.IOStreams) *cobra.Command {
	o := NewCmdOptions(
		flags,
		streams,
	)

	rootCmd := &cobra.Command{
		Use:   fmt.Sprintf(rootUse, filepath.Base(os.Args[0])),
		Short: rootShort,
		Long:  util.CliLongDesc(rootLong),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return util.SetupLogging(o.IOStreams.ErrOut, o.debug)
		},
	}

	// set IOStreams for the command
	rootCmd.SetIn(o.IOStreams.In)
	rootCmd.SetOut(o.IOStreams.Out)
	rootCmd.SetErr(o.IOStreams.ErrOut)

	// add default kubernetes flags as global flags
	o.kubernetesConfigFlags.AddFlags(rootCmd.PersistentFlags())

	// add own global flags
	rootCmd.PersistentFlags().BoolVar(&o.debug, "debug", false, "Enable debug logging")

	// add subcommands
	rootCmd.AddCommand(
		newClaimCmd(o),
		newSecretCmd(o),
		newVersionCmd(o),
		newVolumeCmd(o),
	)

	// add completions
	err := rootCmd.RegisterFlagCompletionFunc("namespace", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		client, err := util.ClientFromFlags(o.kubernetesConfigFlags)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		return util.CompleteNamespaceNames(client, toComplete)
	})
	if err != nil {
		panic(err)
	}

	return rootCmd
}
