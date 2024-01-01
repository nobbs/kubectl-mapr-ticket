package cli

import (
	"fmt"

	"github.com/nobbs/kubectl-mapr-ticket/internal/util"
	"github.com/nobbs/kubectl-mapr-ticket/internal/volumes"
	"github.com/spf13/cobra"
)

type UsedByOptions struct {
	*rootCmdOptions

	// Args are the arguments passed to the command
	args []string

	// SecretName is the name of the secret to find persistent volumes for
	SecretName string

	// AllSecrets indicates whether to find persistent volumes for all secrets
	// in the current namespace
	AllSecrets bool

	// OutputFormat is the format to use for output
	OutputFormat string
}

func NewUsedByOptions(rootOpts *rootCmdOptions) *UsedByOptions {
	return &UsedByOptions{
		rootCmdOptions: rootOpts,
	}
}

func newUsedByCmd(rootOpts *rootCmdOptions) *cobra.Command {
	o := NewUsedByOptions(rootOpts)

	cmd := &cobra.Command{
		Use:   "used-by {secret-name|--all} [flags]",
		Short: "List all persistent volumes that use the specified MapR ticket secret",
		Long: `List all persistent volumes that use the specified MapR ticket secret and print
some information about them.`,
		Args: cobra.MaximumNArgs(1),
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

	// set IOStreams for the command
	cmd.SetIn(o.IOStreams.In)
	cmd.SetOut(o.IOStreams.Out)
	cmd.SetErr(o.IOStreams.ErrOut)

	// add flags
	cmd.Flags().StringVarP(&o.OutputFormat, "output", "o", "table", "Output format. One of: table|wide")
	cmd.Flags().BoolVarP(&o.AllSecrets, "all", "a", false, "List persistent volumes for all MapR ticket secrets in the current namespace")

	return cmd
}

func (o *UsedByOptions) Complete(cmd *cobra.Command, args []string) error {
	o.args = args

	if len(args) > 0 {
		o.SecretName = args[0]
	}

	return nil
}

func (o *UsedByOptions) Validate() error {
	// ensure that the secret name was provided
	if !o.AllSecrets && o.SecretName == "" {
		return fmt.Errorf("either --all or a secret name must be provided")
	}

	// ensure that the output format is valid
	if o.OutputFormat != "table" && o.OutputFormat != "wide" {
		return fmt.Errorf("output format %s is not valid", o.OutputFormat)
	}

	return nil
}

func (o *UsedByOptions) Run(cmd *cobra.Command, args []string) error {
	client, err := util.ClientFromFlags(o.kubernetesConfigFlags)
	if err != nil {
		return err
	}

	// create list options
	opts := []volumes.ListerOption{}

	// if we are listing volumes for all secrets in the namespace, create an option to do so
	if o.AllSecrets {
		opts = append(opts, volumes.WithAllSecrets())
	}

	// create lister
	lister := volumes.NewLister(client, o.SecretName, *o.kubernetesConfigFlags.Namespace, opts...)

	// run the lister
	pvs, err := lister.Run()
	if err != nil {
		return err
	}

	// print the volumes
	if err := volumes.Print(cmd, pvs); err != nil {
		return err
	}

	return nil
}
