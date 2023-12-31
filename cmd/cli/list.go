package cli

import (
	"fmt"

	"github.com/nobbs/kubectl-mapr-ticket/internal/list"
	"github.com/nobbs/kubectl-mapr-ticket/internal/util"
	"github.com/spf13/cobra"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ListOptions struct {
	*rootCmdOptions

	// OutputFormat is the format to use for output
	OutputFormat string

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

	// FilterByMaprUID indicates whether to filter secrets to only those that have
	// a ticket for the specified UID
	FilterByMaprUID uint32

	// FilterByMaprGID indicates whether to filter secrets to only those that have
	// a ticket for the specified GID
	FilterByMaprGID uint32

	// FilterByInUse indicates whether to filter secrets to only those that are
	// in use by a persistent volume
	FilterByInUse bool

	// ShowInUse indicates whether to show only secrets that are in use by a
	// persistent volume
	ShowInUse bool
}

func NewListOptions(rootOpts *rootCmdOptions) *ListOptions {
	return &ListOptions{
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
	cmd.Flags().StringVarP(&o.OutputFormat, "output", "o", "table", "Output format. One of: table|wide|json|yaml")
	cmd.Flags().BoolVarP(&o.AllNamespaces, "all-namespaces", "A", false, "If true, list the requested object(s) across all namespaces. Namespace in current context is ignored even if specified with --namespace.")
	cmd.Flags().BoolVarP(&o.FilterOnlyExpired, "only-expired", "E", false, "If true, only show secrets with tickets that have expired")
	cmd.Flags().BoolVarP(&o.FilterOnlyUnexpired, "only-unexpired", "U", false, "If true, only show secrets with tickets that have not expired")
	cmd.Flags().StringVarP(&o.FilterByMaprCluster, "mapr-cluster", "c", "", "Only show secrets with tickets for the specified MapR cluster")
	cmd.Flags().StringVarP(&o.FilterByMaprUser, "mapr-user", "u", "", "Only show secrets with tickets for the specified MapR user")
	cmd.Flags().Uint32Var(&o.FilterByMaprUID, "mapr-uid", 0, "Only show secrets with tickets for the specified UID")
	cmd.Flags().Uint32Var(&o.FilterByMaprGID, "mapr-gid", 0, "Only show secrets with tickets for the specified GID")
	cmd.Flags().BoolVarP(&o.FilterByInUse, "in-use", "I", false, "If true, only show secrets that are in use by a persistent volume")
	cmd.Flags().BoolVarP(&o.ShowInUse, "show-in-use", "i", false, "If true, add a column to the output indicating whether the secret is in use by a persistent volume")
	cmd.MarkFlagsMutuallyExclusive("only-expired", "only-unexpired")

	return cmd
}

func (o *ListOptions) Complete(cmd *cobra.Command, args []string) error {
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

func (o *ListOptions) Validate() error {
	// validate output format
	if o.OutputFormat != "table" && o.OutputFormat != "wide" && o.OutputFormat != "json" && o.OutputFormat != "yaml" {
		return fmt.Errorf("invalid output format: %s. Must be one of: table|wide|json|yaml", o.OutputFormat)
	}

	return nil
}

//gocyclo:ignore
func (o *ListOptions) Run(cmd *cobra.Command, args []string) error {
	client, err := util.ClientFromFlags(o.kubernetesConfigFlags)
	if err != nil {
		return err
	}

	// create list options and pass them to the lister
	opts := []list.ListerOption{}

	if cmd.Flags().Changed("only-expired") && o.FilterOnlyExpired {
		opts = append(opts, list.WithFilterOnlyExpired())
	}

	if cmd.Flags().Changed("only-unexpired") && o.FilterOnlyUnexpired {
		opts = append(opts, list.WithFilterOnlyUnexpired())
	}

	if cmd.Flags().Changed("mapr-cluster") {
		opts = append(opts, list.WithFilterByMaprCluster(o.FilterByMaprCluster))
	}

	if cmd.Flags().Changed("mapr-user") {
		opts = append(opts, list.WithFilterByMaprUser(o.FilterByMaprUser))
	}

	if cmd.Flags().Changed("mapr-uid") {
		opts = append(opts, list.WithFilterByUID(o.FilterByMaprUID))
	}

	if cmd.Flags().Changed("mapr-gid") {
		opts = append(opts, list.WithFilterByGID(o.FilterByMaprGID))
	}

	if cmd.Flags().Changed("in-use") && o.FilterByInUse {
		opts = append(opts, list.WithFilterByInUse())
	}

	if cmd.Flags().Changed("show-in-use") && o.ShowInUse {
		opts = append(opts, list.WithShowInUse())
	}

	// create lister
	lister := list.NewLister(client, *o.kubernetesConfigFlags.Namespace, opts...)

	// run lister
	items, err := lister.Run()
	if err != nil {
		return err
	}

	// print output
	if err := list.Print(cmd, items); err != nil {
		return err
	}

	return nil
}
