package cli

import (
	"github.com/nobbs/kubectl-mapr-ticket/internal/list"
	"github.com/nobbs/kubectl-mapr-ticket/internal/util"
	"github.com/spf13/cobra"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/printers"
)

type listOptions struct {
	*rootCmdOptions

	// AllNamespaces indicates whether to list secrets in all namespaces
	AllNamespaces bool

	// FilterOnlyExpired indicates whether to filter secrets to only those that
	// have expired
	FilterOnlyExpired bool

	// FilterOnlyUnexpired indicates whether to filter secrets to only those
	// that have not expired
	FilterOnlyUnexpired bool

	// FilterByMaprCluster indicates whether to filter secrets to only those
	// that have a ticket for the specified MapR cluster
	FilterByMaprCluster string

	// FilterByMaprUser indicates whether to filter secrets to only those that
	// have a ticket for the specified MapR user
	FilterByMaprUser string
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
	cmd.Flags().BoolVarP(&o.AllNamespaces, "all-namespaces", "A", false, "If true, list the requested object(s) across all namespaces. Namespace in current context is ignored even if specified with --namespace.")
	cmd.Flags().BoolVarP(&o.FilterOnlyExpired, "only-expired", "E", false, "If true, only show secrets with tickets that have expired")
	cmd.Flags().BoolVarP(&o.FilterOnlyUnexpired, "only-unexpired", "U", false, "If true, only show secrets with tickets that have not expired")
	cmd.Flags().StringVarP(&o.FilterByMaprCluster, "mapr-cluster", "c", "", "Only show secrets with tickets for the specified MapR cluster")
	cmd.Flags().StringVarP(&o.FilterByMaprUser, "mapr-user", "u", "", "Only show secrets with tickets for the specified MapR user")
	cmd.MarkFlagsMutuallyExclusive("only-expired", "only-unexpired")

	return cmd
}

func (o *listOptions) Complete(cmd *cobra.Command, args []string) error {
	// set namespace
	if o.kubernetesConfigFlags.Namespace == nil || *o.kubernetesConfigFlags.Namespace == "" {
		namespace := util.GetNamespace(o.kubernetesConfigFlags)
		o.kubernetesConfigFlags.Namespace = &namespace
	}

	// reset namespace if --all-namespaces is set
	if o.AllNamespaces {
		namespaceAll := metaV1.NamespaceAll
		o.kubernetesConfigFlags.Namespace = &namespaceAll
	}

	return nil
}

func (o *listOptions) Run(cmd *cobra.Command, args []string) error {
	client, err := util.ClientFromFlags(o.kubernetesConfigFlags)
	if err != nil {
		return err
	}

	// create list options
	opts := []list.ListerOption{}

	if o.FilterOnlyExpired {
		opts = append(opts, list.WithFilterOnlyExpired())
	}

	if o.FilterOnlyUnexpired {
		opts = append(opts, list.WithFilterOnlyUnexpired())
	}

	if o.FilterByMaprCluster != "" {
		opts = append(opts, list.WithFilterByMaprCluster(o.FilterByMaprCluster))
	}

	if o.FilterByMaprUser != "" {
		opts = append(opts, list.WithFilterByMaprUser(o.FilterByMaprUser))
	}

	// create lister
	lister := list.NewLister(client, *o.kubernetesConfigFlags.Namespace, opts...)

	// run lister
	items, err := lister.Run()
	if err != nil {
		return err
	}

	// generate table for output
	table, err := list.GenerateTable(cmd, items)
	if err != nil {
		return err
	}

	// print table
	printer := printers.NewTablePrinter(printers.PrintOptions{
		WithNamespace: o.AllNamespaces,
	})

	err = printer.PrintObj(table, o.IOStreams.Out)
	if err != nil {
		return err
	}

	return nil
}
