package claim

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nobbs/kubectl-mapr-ticket/cmd/common"
	"github.com/nobbs/kubectl-mapr-ticket/pkg/claim"
	"github.com/nobbs/kubectl-mapr-ticket/pkg/secret"
	"github.com/nobbs/kubectl-mapr-ticket/pkg/util"
)

const (
	claimUse   = `claim`
	claimShort = "List all persistent volumes claims that use a MapR ticket in the current namespace"
	claimLong  = `
		List all persistent volumes claims that use a MapR ticket in the current namespace.

		By default, this command lists all persistent volumes claims that use a MapR ticket in the current namespace.
		`
)

var (
	claimValidOutputFormats = []string{"table", "wide"}
)

type options struct {
	*common.Options

	// OutputFormat is the format to use for output
	OutputFormat string

	// AllNamespaces indicates whether to list secrets in all namespaces
	AllNamespaces bool
}

func newOptions(opts *common.Options) *options {
	return &options{
		Options: opts,
	}
}

func NewCmd(opts *common.Options) *cobra.Command {
	o := newOptions(opts)

	cmd := &cobra.Command{
		Aliases: []string{"pvc"},
		Use:     claimUse,
		Short:   claimShort,
		Long:    common.CliLongDesc(claimLong),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.Complete(cmd, args); err != nil {
				return err
			}

			if err := o.Validate(); err != nil {
				return err
			}

			if err := o.Run(cmd, args); err != nil {
				return err
			}

			return nil
		},
	}

	// set IOStreams for this command
	cmd.SetIn(o.IOStreams.In)
	cmd.SetOut(o.IOStreams.Out)
	cmd.SetErr(o.IOStreams.ErrOut)

	// add flags
	cmd.Flags().StringVarP(&o.OutputFormat, "output", "o", "table", fmt.Sprintf("Output format. One of (%s)", common.StringSliceToFlagOptions(claimValidOutputFormats)))
	cmd.Flags().BoolVarP(&o.AllNamespaces, "all-namespaces", "A", false, "List persistent volumes claims that use a MapR ticket in all namespaces")

	// register completions for flags
	if err := o.registerCompletions(cmd); err != nil {
		panic(err)
	}

	return cmd
}

func (o *options) Complete(cmd *cobra.Command, args []string) error {
	// set namespace based on flags
	ns := util.GetNamespace(o.KubernetesConfigFlags, o.AllNamespaces)
	o.KubernetesConfigFlags.Namespace = &ns

	return nil
}

func (o *options) Validate() error {
	// validate output format
	if o.OutputFormat != "table" && o.OutputFormat != "wide" {
		return fmt.Errorf("invalid output format %q. Must be one of (%s)", o.OutputFormat, common.StringSliceToFlagOptions(claimValidOutputFormats))
	}

	return nil
}

func (o *options) Run(cmd *cobra.Command, args []string) error {
	client, err := util.ClientFromFlags(o.KubernetesConfigFlags)
	if err != nil {
		return err
	}

	// create secret lister
	secretLister := secret.NewLister(
		client,
		util.NamespaceAll,
	)

	// create list options and pass them to the lister
	opts := []claim.ListerOption{
		claim.WithSecretLister(secretLister),
	}

	// create lister
	lister := claim.NewLister(
		client,
		*o.KubernetesConfigFlags.Namespace,
		opts...,
	)

	// run lister
	volumeClaims, err := lister.List()
	if err != nil {
		return err
	}

	// print output
	if err := claim.Print(cmd, volumeClaims); err != nil {
		return err
	}

	return nil
}

func (o *options) registerCompletions(cmd *cobra.Command) error {
	err := cmd.RegisterFlagCompletionFunc("output", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return common.CompleteStringValues(claimValidOutputFormats, toComplete)
	})
	if err != nil {
		return err
	}

	return nil
}
