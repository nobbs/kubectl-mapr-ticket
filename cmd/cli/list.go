package cli

import (
	"github.com/nobbs/kubectl-mapr-ticket/internal/ticket"
	"github.com/nobbs/kubectl-mapr-ticket/internal/utils"
	"github.com/spf13/cobra"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/printers"
)

type listOptions struct {
	*rootCmdOptions

	// AllNamespaces indicates whether to list secrets in all namespaces
	AllNamespaces bool
}

func NewListOptions(rootOpts *rootCmdOptions) *listOptions {
	return &listOptions{
		rootCmdOptions: rootOpts,
	}
}

func newListCmd(rootOpts *rootCmdOptions) *cobra.Command {
	o := NewListOptions(rootOpts)

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all secrets containing MapR tickets in the current namespace",
		Long: `List all secrets containing MapR tickets in the current namespace and print
some information about them.`,
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
	cmd.Flags().BoolVarP(&o.AllNamespaces, "all-namespaces", "A", false, "If present, list the requested object(s) across all namespaces. Namespace in current context is ignored even if specified with --namespace.")

	return cmd
}

func (o *listOptions) Complete(cmd *cobra.Command, args []string) error {
	if o.kubernetesConfigFlags.Namespace == nil || *o.kubernetesConfigFlags.Namespace == "" {
		namespace, err := getNamespace(o.kubernetesConfigFlags)
		if err != nil {
			return err
		}

		o.kubernetesConfigFlags.Namespace = &namespace
	}

	if o.AllNamespaces {
		namespaceAll := metaV1.NamespaceAll
		o.kubernetesConfigFlags.Namespace = &namespaceAll
	}

	return nil
}

func (o *listOptions) Validate() error {
	return nil
}

func (o *listOptions) Run(cmd *cobra.Command, args []string) error {
	client, err := utils.ClientFromFlags(o.kubernetesConfigFlags)
	if err != nil {
		return err
	}

	ticketSecrets, err := ticket.NewList(
		client.CoreV1().Secrets(*o.kubernetesConfigFlags.Namespace),
	).Run()
	if err != nil {
		return err
	}

	// generate table for output
	table, err := ticket.GenerateTable(cmd, ticketSecrets)
	if err != nil {
		return err
	}

	// print table
	printer := printers.NewTablePrinter(printers.PrintOptions{
		WithNamespace: o.AllNamespaces,
	})
	return printer.PrintObj(table, o.IOStreams.Out)
}

// getNamespace returns the namespace from the kubeconfig or the default flag
func getNamespace(flags *genericclioptions.ConfigFlags) (string, error) {
	namespace, _, err := flags.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		return "", err
	}

	return namespace, nil
}
