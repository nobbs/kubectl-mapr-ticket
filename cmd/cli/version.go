package cli

import (
	"github.com/nobbs/kubectl-mapr-ticket/internal/version"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericiooptions"
)

type versionOptions struct {
	configFlags *genericclioptions.ConfigFlags
	IOStreams   genericiooptions.IOStreams
}

func newVersionOptions(streams genericiooptions.IOStreams) *versionOptions {
	return &versionOptions{
		configFlags: genericclioptions.NewConfigFlags(true),
		IOStreams:   streams,
	}
}

func newVersionCmd(streams genericiooptions.IOStreams) *cobra.Command {
	o := newVersionOptions(streams)

	cmd := &cobra.Command{
		Use:     "version",
		Aliases: []string{"v"},
		Short:   "Print the version of kubectl-mapr-ticket and exit",
		Run: func(cmd *cobra.Command, args []string) {
			o.PrintVersionInfo(cmd)
		},
	}

	// set IOStreams for the command
	cmd.SetIn(o.IOStreams.In)
	cmd.SetOut(o.IOStreams.Out)
	cmd.SetErr(o.IOStreams.ErrOut)

	return cmd
}

func (o *versionOptions) PrintVersionInfo(cmd *cobra.Command) {
	versionInfo := version.NewVersion()
	cmd.Println(versionInfo)
}
