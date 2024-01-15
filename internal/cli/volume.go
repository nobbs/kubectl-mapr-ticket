package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/nobbs/kubectl-mapr-ticket/internal/secret"
	"github.com/nobbs/kubectl-mapr-ticket/internal/util"
	"github.com/nobbs/kubectl-mapr-ticket/internal/volume"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
		`
)

var (
	volumeValidOutputFormats = []string{"table", "wide"}
)

type VolumeOptions struct {
	*rootCmdOptions

	// Args are the arguments passed to the command
	args []string

	// SecretName is the name of the secret to find persistent volumes for
	SecretName string

	// AllNamespaces indicates whether to find persistent volumes for all secrets
	// in all namespaces
	AllNamespaces bool

	// OutputFormat is the format to use for output
	OutputFormat string
}

func NewVolumeOptions(rootOpts *rootCmdOptions) *VolumeOptions {
	return &VolumeOptions{
		rootCmdOptions: rootOpts,
	}
}

func newVolumeCmd(rootOpts *rootCmdOptions) *cobra.Command {
	o := NewVolumeOptions(rootOpts)

	cmd := &cobra.Command{
		Aliases: []string{"pv"},
		Use:     volumeUse,
		Short:   volumeShort,
		Long:    util.CliLongDesc(volumeLong),
		Example: util.CliExample(volumeExample, filepath.Base(os.Args[0])),
		Args:    cobra.MaximumNArgs(1),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			// we only want one argument, so don't complete once we have one
			if len(args) > 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}

			// set namespace based on flags
			namespace := util.GetNamespace(o.kubernetesConfigFlags, false)
			o.kubernetesConfigFlags.Namespace = &namespace

			// get client
			client, err := util.ClientFromFlags(o.kubernetesConfigFlags)
			if err != nil {
				return nil, cobra.ShellCompDirectiveError
			}

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
	cmd.Flags().StringVarP(&o.OutputFormat, "output", "o", "table", fmt.Sprintf("Output format. One of (%s)", util.StringSliceToFlagOptions(volumeValidOutputFormats)))
	cmd.Flags().BoolVarP(&o.AllNamespaces, "all-namespaces", "A", false, "List persistent volumes for all MapR ticket secrets in all namespaces")

	// register completions for flags
	if err := o.registerCompletions(cmd); err != nil {
		panic(err)
	}

	return cmd
}

func (o *VolumeOptions) Complete(cmd *cobra.Command, args []string) error {
	// parse the arguments
	o.args = args

	// set secret name based on args
	switch len(args) {
	case 0:
		o.SecretName = volume.SecretAll
	case 1:
		o.SecretName = args[0]
	default:
		return fmt.Errorf("too many arguments provided, either provide a secret name or nothing")
	}

	// set namespace based on flags
	namespace := util.GetNamespace(o.kubernetesConfigFlags, o.AllNamespaces)
	o.kubernetesConfigFlags.Namespace = &namespace

	return nil
}

func (o *VolumeOptions) Validate() error {
	// ensure that the output format is valid
	if o.OutputFormat != "table" && o.OutputFormat != "wide" {
		return fmt.Errorf("output format %s is not valid", o.OutputFormat)
	}

	return nil
}

func (o *VolumeOptions) Run(cmd *cobra.Command, args []string) error {
	client, err := util.ClientFromFlags(o.kubernetesConfigFlags)
	if err != nil {
		return err
	}

	// create secret lister
	secretLister := secret.NewLister(
		client,
		metaV1.NamespaceAll,
	)

	// create lister
	lister := volume.NewLister(
		client,
		o.SecretName,
		*o.kubernetesConfigFlags.Namespace,
		volume.WithSecretLister(secretLister),
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

func (o *VolumeOptions) registerCompletions(cmd *cobra.Command) error {
	err := cmd.RegisterFlagCompletionFunc("output", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return util.CompleteStringValues(volumeValidOutputFormats, toComplete)
	})
	if err != nil {
		return err
	}

	return nil
}
