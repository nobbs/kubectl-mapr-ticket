package cli

import (
	"github.com/nobbs/kubectl-mapr-ticket/internal/version"
	"github.com/spf13/cobra"
)

// versionOptions holds the options for 'version' sub command
type versionOptions struct {
	// embed common options from rootCmdOptions
	*rootCmdOptions
}

func newVersionOptions(rootOpts *rootCmdOptions) *versionOptions {
	return &versionOptions{
		rootCmdOptions: rootOpts,
	}
}

func newVersionCmd(rootOpts *rootCmdOptions) *cobra.Command {
	o := newVersionOptions(rootOpts)

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
