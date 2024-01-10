package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nobbs/kubectl-mapr-ticket/internal/util"
	"github.com/nobbs/kubectl-mapr-ticket/internal/volume"
	"github.com/spf13/cobra"
)

const (
	usedByUse   = `used-by {secret-name|--all}`
	usedByShort = "List all persistent volumes that use the specified MapR ticket secret"
	usedByLong  = `
		List all persistent volumes that use the specified MapR ticket secret and print
		some information about them.
		`
	usedByExample = `
		# List all persistent volumes that use the specified MapR ticket secret
		%[1]s used-by my-secret
		`
)

var (
	usedByValidOutputFormats = []string{"table", "wide"}
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
		Use:     usedByUse,
		Short:   usedByShort,
		Long:    util.CliLongDesc(usedByLong),
		Example: util.CliExample(usedByExample, filepath.Base(os.Args[0])),
		Args:    cobra.MaximumNArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			// if we are listing volumes for all secrets in the namespace, we don't want to complete
			if o.AllSecrets {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}

			// we only want one argument, the secret name
			if len(args) > 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}

			client, err := util.ClientFromFlags(o.kubernetesConfigFlags)
			if err != nil {
				return nil, cobra.ShellCompDirectiveError
			}

			namespace := util.GetNamespace(o.kubernetesConfigFlags, false)

			return util.CompleteTicketNames(client, namespace, args, toComplete)
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

	// set IOStreams for the command
	cmd.SetIn(o.IOStreams.In)
	cmd.SetOut(o.IOStreams.Out)
	cmd.SetErr(o.IOStreams.ErrOut)

	// add flags
	cmd.Flags().StringVarP(&o.OutputFormat, "output", "o", "table", fmt.Sprintf("Output format. One of (%s)", util.StringSliceToFlagOptions(usedByValidOutputFormats)))
	cmd.Flags().BoolVarP(&o.AllSecrets, "all", "a", false, "List persistent volumes for all MapR ticket secrets in the current namespace")

	// register completions for flags
	if err := o.registerCompletions(cmd); err != nil {
		panic(err)
	}

	return cmd
}

func (o *UsedByOptions) Complete(cmd *cobra.Command, args []string) error {
	// parse the arguments
	o.args = args

	if len(args) > 0 {
		o.SecretName = args[0]
	}

	// set namespace
	if o.kubernetesConfigFlags.Namespace == nil || *o.kubernetesConfigFlags.Namespace == "" {
		namespace := util.GetNamespace(o.kubernetesConfigFlags, false)
		o.kubernetesConfigFlags.Namespace = &namespace
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
	opts := []volume.ListerOption{}

	// if we are listing volumes for all secrets in the namespace, create an option to do so
	if o.AllSecrets {
		opts = append(opts, volume.WithAllSecrets())
	}

	// create lister
	lister := volume.NewLister(client, o.SecretName, *o.kubernetesConfigFlags.Namespace, opts...)

	// run the lister
	pvs, err := lister.Run()
	if err != nil {
		return err
	}

	// print the volumes
	if err := volume.Print(cmd, pvs); err != nil {
		return err
	}

	return nil
}

func (o *UsedByOptions) registerCompletions(cmd *cobra.Command) error {
	err := cmd.RegisterFlagCompletionFunc("output", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return util.CompleteStringValues(usedByValidOutputFormats, toComplete)
	})
	if err != nil {
		return err
	}

	return nil
}
