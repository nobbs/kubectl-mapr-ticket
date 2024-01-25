package volume

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nobbs/kubectl-mapr-ticket/cmd/common"
	"github.com/nobbs/kubectl-mapr-ticket/pkg/secret"
	"github.com/nobbs/kubectl-mapr-ticket/pkg/util"
	"github.com/nobbs/kubectl-mapr-ticket/pkg/volume"
)

const (
	volumeUse   = `volume [secret-name]`
	volumeShort = "List all persistent volumes that use the specified MapR ticket secret"
	volumeLong  = `
		List all persistent volumes that use the specified MapR ticket secret and print
		some information about them.
		`
	volumeExample = `
		# List all persistent volumes that use the specified MapR ticket secret
		%[1]s volume my-secret

		# List all persistent volumes that use any MapR ticket secret in the current namespace
		%[1]s volume

		# List all persistent volumes that use any MapR ticket secret in all namespaces
		%[1]s volume --all-namespaces

		# List all persistent volumes that use any MapR ticket secret in all namespaces, sorted by expiration date
		%[1]s volume --all-namespaces --sort-by expiryTime
		`
)

var (
	volumeValidOutputFormats = []string{"table", "wide"}
)

type options struct {
	*common.Options

	// Args are the arguments passed to the command
	args []string

	// SecretName is the name of the secret to find persistent volumes for
	SecretName string

	// OutputFormat is the format to use for output
	OutputFormat string

	// AllNamespaces indicates whether to find persistent volumes for all secrets
	// in all namespaces
	AllNamespaces bool

	// SortBy is the list of fields to sort by
	SortBy []string
}

func newOptions(opts *common.Options) *options {
	return &options{
		Options: opts,
	}
}

func NewCmd(opts *common.Options) *cobra.Command {
	o := newOptions(opts)

	cmd := &cobra.Command{
		Aliases: []string{"pv"},
		Use:     volumeUse,
		Short:   volumeShort,
		Long:    common.CliLongDesc(volumeLong),
		Example: common.CliExample(volumeExample, common.CliBinName),
		Args:    cobra.MaximumNArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			// we only want one argument, so don't complete once we have one
			if len(args) > 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}

			// set namespace based on flags
			namespace := util.GetNamespace(o.KubernetesConfigFlags, false)
			o.KubernetesConfigFlags.Namespace = &namespace

			// get client
			client, err := util.ClientFromFlags(o.KubernetesConfigFlags)
			if err != nil {
				return nil, cobra.ShellCompDirectiveError
			}

			return common.CompleteTicketNames(client, namespace, args, toComplete)
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
	cmd.Flags().StringVarP(&o.OutputFormat, "output", "o", "table", fmt.Sprintf("Output format. One of (%s)", common.StringSliceToFlagOptions(volumeValidOutputFormats)))
	cmd.Flags().BoolVarP(&o.AllNamespaces, "all-namespaces", "A", false, "List persistent volumes for all MapR ticket secrets in all namespaces")
	cmd.Flags().StringSliceVar(&o.SortBy, "sort-by", []string{}, fmt.Sprintf("Sort list of persistent volumes by the specified fields. One or more of (%s)", common.StringSliceToFlagOptions(volume.SortOptionsList)))

	// register completions for flags
	if err := o.registerCompletions(cmd); err != nil {
		panic(err)
	}

	return cmd
}

func (o *options) Complete(cmd *cobra.Command, args []string) error {
	// parse the arguments
	o.args = args

	// set secret name based on args
	switch len(args) {
	case 0:
		o.SecretName = util.SecretAll
	case 1:
		o.SecretName = args[0]
	default:
		return fmt.Errorf("too many arguments provided, either provide a secret name or nothing")
	}

	// set namespace based on flags
	ns := util.GetNamespace(o.KubernetesConfigFlags, o.AllNamespaces)
	o.KubernetesConfigFlags.Namespace = &ns

	return nil
}

func (o *options) Validate() error {
	// validate output format
	if o.OutputFormat != "table" && o.OutputFormat != "wide" {
		return fmt.Errorf("output format %s is not valid", o.OutputFormat)
	}

	// ensure that the sort options are valid
	if err := volume.ValidateSortOptions(o.SortBy); err != nil {
		return err
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
	opts := []volume.ListerOption{
		volume.WithSecretLister(secretLister),
	}

	if cmd.Flags().Changed("sort-by") && o.SortBy != nil {
		// convert sort options to SortOptions
		sortOptions := make([]volume.SortOption, 0, len(o.SortBy))
		for _, sortBy := range o.SortBy {
			sortOptions = append(sortOptions, volume.SortOption(sortBy))
		}

		opts = append(opts, volume.WithSortBy(sortOptions))
	}

	// create lister
	lister := volume.NewLister(
		client,
		o.SecretName,
		*o.KubernetesConfigFlags.Namespace,
		opts...,
	)

	// run the lister
	pvs, err := lister.List()
	if err != nil {
		return err
	}

	// print the volumes
	if err := volume.Print(cmd, pvs); err != nil {
		return err
	}

	return nil
}

func (o *options) registerCompletions(cmd *cobra.Command) error {
	err := cmd.RegisterFlagCompletionFunc("output", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return common.CompleteStringValues(volumeValidOutputFormats, toComplete)
	})
	if err != nil {
		return err
	}

	err = cmd.RegisterFlagCompletionFunc("sort-by", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return common.CompleteStringValues(volume.SortOptionsList, toComplete)
	})
	if err != nil {
		return err
	}

	return nil
}
