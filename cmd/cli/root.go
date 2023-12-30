package cli

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

type rootCmdOptions struct {
	kubernetesConfigFlags *genericclioptions.ConfigFlags
	IOStreams             genericclioptions.IOStreams
}

func NewCmdOptions(kubernetesConfigFlags *genericclioptions.ConfigFlags, streams genericclioptions.IOStreams) *rootCmdOptions {
	return &rootCmdOptions{
		kubernetesConfigFlags: kubernetesConfigFlags,
		IOStreams:             streams,
	}
}

func NewRootCmd(flags *genericclioptions.ConfigFlags, streams genericclioptions.IOStreams) *cobra.Command {
	rootOpts := NewCmdOptions(
		flags,
		streams,
	)

	rootCmd := &cobra.Command{
		Use:   "kubectl-mapr-ticket",
		Short: "A kubectl plugin to list and inspect MapR tickets",
		Long: `A kubectl plugin that allows you to list and inspect MapR tickets from a
Kubernetes cluster, including details stored in the ticket itself without
requiring access to the MapR cluster.`,
	}

	// set IOStreams for the command
	rootCmd.SetIn(rootOpts.IOStreams.In)
	rootCmd.SetOut(rootOpts.IOStreams.Out)
	rootCmd.SetErr(rootOpts.IOStreams.ErrOut)

	// add default kubernetes flags as global flags
	rootOpts.kubernetesConfigFlags.AddFlags(rootCmd.PersistentFlags())

	// add subcommands
	rootCmd.AddCommand(
		newListCmd(rootOpts),
		newVersionCmd(rootOpts),
	)

	return rootCmd
}
