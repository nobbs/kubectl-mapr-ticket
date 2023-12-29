package cli

import (
	"github.com/nobbs/kubectl-mapr-ticket/internal/ticket"
	"github.com/spf13/cobra"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericiooptions"
	"k8s.io/cli-runtime/pkg/printers"
	"k8s.io/client-go/kubernetes"
)

type listOptions struct {
	configFlags *genericclioptions.ConfigFlags
	IOStreams   genericiooptions.IOStreams

	AllNamespaces bool

	client kubernetes.Interface
}

func NewListOptions(streams genericiooptions.IOStreams) *listOptions {
	return &listOptions{
		configFlags: genericclioptions.NewConfigFlags(true),
		IOStreams:   streams,
	}
}

func newListCmd(streams genericiooptions.IOStreams) *cobra.Command {
	o := NewListOptions(streams)

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
	o.configFlags.AddFlags(cmd.Flags())
	cmd.Flags().BoolVarP(&o.AllNamespaces, "all-namespaces", "A", false, "If present, list the requested object(s) across all namespaces. Namespace in current context is ignored even if specified with --namespace.")

	return cmd
}

func (o *listOptions) Complete(cmd *cobra.Command, args []string) error {
	config, err := o.configFlags.ToRESTConfig()
	if err != nil {
		return err
	}

	o.client, err = kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	if o.configFlags.Namespace == nil || *o.configFlags.Namespace == "" {
		namespace, err := getNamespace(o.configFlags)
		if err != nil {
			return err
		}

		o.configFlags.Namespace = &namespace
	}

	if o.AllNamespaces {
		namespaceAll := metaV1.NamespaceAll
		o.configFlags.Namespace = &namespaceAll
	}

	return nil
}

func (o *listOptions) Validate() error {
	return nil
}

func (o *listOptions) Run(cmd *cobra.Command, args []string) error {
	secretGetter := o.client.CoreV1().Secrets(*o.configFlags.Namespace)

	ticketSecrets, err := ticket.NewList(secretGetter).Run()
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
