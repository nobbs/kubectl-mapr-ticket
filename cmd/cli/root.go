package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nobbs/kubectl-mapr-ticket/internal/util"
	"github.com/spf13/cobra"
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
}

func NewCmdOptions(kubernetesConfigFlags *genericclioptions.ConfigFlags, streams genericiooptions.IOStreams) *rootCmdOptions {
	return &rootCmdOptions{
		kubernetesConfigFlags: kubernetesConfigFlags,
		IOStreams:             streams,
	}
}

func NewRootCmd(flags *genericclioptions.ConfigFlags, streams genericiooptions.IOStreams) *cobra.Command {
	rootOpts := NewCmdOptions(
		flags,
		streams,
	)

	rootCmd := &cobra.Command{
		Use:   fmt.Sprintf(rootUse, filepath.Base(os.Args[0])),
		Short: rootShort,
		Long:  util.CliLongDesc(rootLong),
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
		newUsedByCmd(rootOpts),
	)

	return rootCmd
}
