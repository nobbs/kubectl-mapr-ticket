package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nobbs/kubectl-mapr-ticket/internal/claim"
	"github.com/nobbs/kubectl-mapr-ticket/internal/util"
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

type ClaimOptions struct {
	*rootCmdOptions

	// AllNamespaces indicates whether to list secrets in all namespaces
	AllNamespaces bool

	// OutputFormat is the format to use for output
	OutputFormat string
}

func NewClaimOptions(rootOpts *rootCmdOptions) *ClaimOptions {
	return &ClaimOptions{
		rootCmdOptions: rootOpts,
	}
}

func newClaimCmd(rootOpts *rootCmdOptions) *cobra.Command {
	o := NewClaimOptions(rootOpts)

	cmd := &cobra.Command{
		Use:   claimUse,
		Short: claimShort,
		Long:  util.CliLongDesc(claimLong),
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
	cmd.Flags().StringVarP(&o.OutputFormat, "output", "o", "table", fmt.Sprintf("Output format. One of (%s)", util.StringSliceToFlagOptions(claimValidOutputFormats)))
	cmd.Flags().BoolVarP(&o.AllNamespaces, "all-namespaces", "A", false, "List persistent volumes claims that use a MapR ticket in all namespaces")

	// register completions for flags
	if err := o.registerCompletions(cmd); err != nil {
		panic(err)
	}

	return cmd
}

func (o *ClaimOptions) Complete(cmd *cobra.Command, args []string) error {
	// set namespace based on flags
	namespace := util.GetNamespace(o.kubernetesConfigFlags, o.AllNamespaces)
	o.kubernetesConfigFlags.Namespace = &namespace

	return nil
}

func (o *ClaimOptions) Validate() error {
	// ensure output format is valid
	if o.OutputFormat != "table" && o.OutputFormat != "wide" {
		return fmt.Errorf("invalid output format %q. Must be one of (%s)", o.OutputFormat, util.StringSliceToFlagOptions(claimValidOutputFormats))
	}

	return nil
}

func (o *ClaimOptions) Run(cmd *cobra.Command, args []string) error {
	client, err := util.ClientFromFlags(o.kubernetesConfigFlags)
	if err != nil {
		return err
	}

	// create lister
	lister := claim.NewLister(client, *o.kubernetesConfigFlags.Namespace)

	// list volume claims
	volumeClaims, err := lister.List()
	if err != nil {
		return err
	}

	// print volume claims
	if err := claim.Print(cmd, volumeClaims); err != nil {
		return err
	}

	return nil
}

func (o *ClaimOptions) registerCompletions(cmd *cobra.Command) error {
	err := cmd.RegisterFlagCompletionFunc("output", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return util.CompleteStringValues(claimValidOutputFormats, toComplete)
	})
	if err != nil {
		return err
	}

	return nil
}
