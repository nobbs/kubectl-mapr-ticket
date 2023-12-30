package cli

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

type cmdOptions struct {
	configFlags *genericclioptions.ConfigFlags
	IOStreams   genericclioptions.IOStreams
}

func NewCmdOptions(streams genericclioptions.IOStreams) *cmdOptions {
	return &cmdOptions{
		configFlags: genericclioptions.NewConfigFlags(true),
		IOStreams:   streams,
	}
}

func NewRootCmd(streams genericclioptions.IOStreams) *cobra.Command {
	o := NewCmdOptions(streams)

	cmd := &cobra.Command{
		Use:   "kubectl-mapr-ticket",
		Short: "A kubectl plugin to list and inspect MapR tickets",
		Long: `A kubectl plugin that allows you to list and inspect MapR tickets from a
Kubernetes cluster, including details stored in the ticket itself without
requiring access to the MapR cluster.`,
	}

	// set IOStreams for the command
	cmd.SetIn(o.IOStreams.In)
	cmd.SetOut(o.IOStreams.Out)
	cmd.SetErr(o.IOStreams.ErrOut)

	// add flags
	o.configFlags.AddFlags(cmd.PersistentFlags())

	// add subcommands
	cmd.AddCommand(
		newListCmd(streams),
		newVersionCmd(streams),
	)

	return cmd
}
