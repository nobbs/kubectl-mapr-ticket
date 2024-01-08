package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nobbs/kubectl-mapr-ticket/internal/secret"
	"github.com/nobbs/kubectl-mapr-ticket/internal/util"
	"github.com/spf13/cobra"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	listUse   = `list`
	listShort = "List all secrets containing MapR tickets in the current namespace"
	listLong  = `
		List all secrets containing MapR tickets in the current namespace and print
		some information about them.
		`
	listExample = `
		# List all MapR tickets in the current namespace
		%[1]s list

		# List all MapR tickets in all namespaces
		%[1]s list --all-namespaces

		# List only expired MapR tickets
		%[1]s list --only-expired

		# List only MapR tickets that expire in the next 7 days
		%[1]s list --expires-before 7d

		# List MapR tickets for a specific MapR user in all namespaces
		%[1]s list --mapr-user mapr --all-namespaces

		# List MapR tickets with number of persistent volumes that use them
		%[1]s list --show-in-use
		`
)

var (
	listValidOutputFormats = []string{"table", "wide", "json", "yaml"}
	listValidSortByFields  = []string{"name", "namespace", "maprCluster", "maprUser", "creationTimestamp", "expiryTime", "numPVC"}
)

type ListOptions struct {
	*rootCmdOptions

	// OutputFormat is the format to use for output
	OutputFormat string

	// AllNamespaces indicates whether to list secrets in all namespaces
	AllNamespaces bool

	// SortBy is the list of fields to sort by
	SortBy []string

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

	// FilterExpiresBefore indicates whether to filter secrets to only those that
	// expire before the specified duration from now
	FilterExpiresBefore util.DurationValue

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
		Aliases: []string{"ls"},
		Use:     listUse,
		Short:   listShort,
		Long:    util.CliLongDesc(listLong),
		Example: util.CliExample(listExample, filepath.Base(os.Args[0])),
		Args:    cobra.NoArgs,
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

	// set IOStreams for the command
	cmd.SetIn(o.IOStreams.In)
	cmd.SetOut(o.IOStreams.Out)
	cmd.SetErr(o.IOStreams.ErrOut)

	// add flags
	cmd.Flags().StringVarP(&o.OutputFormat, "output", "o", "table", fmt.Sprintf("Output format. One of (%s)", util.StringSliceToFlagOptions(listValidOutputFormats)))
	cmd.Flags().BoolVarP(&o.AllNamespaces, "all-namespaces", "A", false, "If true, list the requested object(s) across all namespaces. Namespace in current context is ignored even if specified with --namespace.")
	cmd.Flags().StringSliceVar(&o.SortBy, "sort-by", nil, fmt.Sprintf("Sort list of secrets by the specified fields. One of (%s)", util.StringSliceToFlagOptions(listValidSortByFields)))
	cmd.Flags().BoolVarP(&o.FilterOnlyExpired, "only-expired", "E", false, "If true, only show secrets with tickets that have expired")
	cmd.Flags().BoolVarP(&o.FilterOnlyUnexpired, "only-unexpired", "U", false, "If true, only show secrets with tickets that have not expired")
	cmd.Flags().StringVarP(&o.FilterByMaprCluster, "mapr-cluster", "c", "", "Only show secrets with tickets for the specified MapR cluster")
	cmd.Flags().StringVarP(&o.FilterByMaprUser, "mapr-user", "u", "", "Only show secrets with tickets for the specified MapR user")
	cmd.Flags().Uint32Var(&o.FilterByMaprUID, "mapr-uid", 0, "Only show secrets with tickets for the specified UID")
	cmd.Flags().Uint32Var(&o.FilterByMaprGID, "mapr-gid", 0, "Only show secrets with tickets for the specified GID")
	cmd.Flags().BoolVarP(&o.FilterByInUse, "in-use", "I", false, "If true, only show secrets that are in use by a persistent volume")
	cmd.Flags().Var(&o.FilterExpiresBefore, "expires-before", "Only show secrets with tickets that expire before the specified duration from now")
	cmd.Flags().BoolVarP(&o.ShowInUse, "show-in-use", "i", false, "If true, add a column to the output indicating whether the secret is in use by a persistent volume")
	cmd.MarkFlagsMutuallyExclusive("only-expired", "only-unexpired")

	// register completions for flags
	if err := o.registerCompletions(cmd); err != nil {
		panic(err)
	}

	return cmd
}

func (o *ListOptions) Complete(cmd *cobra.Command, args []string) error {
	// set namespace
	if o.kubernetesConfigFlags.Namespace == nil || *o.kubernetesConfigFlags.Namespace == "" {
		namespace := util.GetNamespace(o.kubernetesConfigFlags, o.AllNamespaces)
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

	// validate sort options
	if err := secret.ValidateSortOptions(o.SortBy); err != nil {
		return err
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
	opts := []secret.ListerOption{}

	if cmd.Flags().Changed("sort-by") && o.SortBy != nil {
		// convert sort options to SortOptions
		sortOptions := make([]secret.SortOptions, 0, len(o.SortBy))
		for _, sortBy := range o.SortBy {
			sortOptions = append(sortOptions, secret.SortOptions(sortBy))
		}

		opts = append(opts, secret.WithSortBy(sortOptions))
	}

	if cmd.Flags().Changed("only-expired") && o.FilterOnlyExpired {
		opts = append(opts, secret.WithFilterOnlyExpired())
	}

	if cmd.Flags().Changed("only-unexpired") && o.FilterOnlyUnexpired {
		opts = append(opts, secret.WithFilterOnlyUnexpired())
	}

	if cmd.Flags().Changed("mapr-cluster") {
		opts = append(opts, secret.WithFilterByMaprCluster(o.FilterByMaprCluster))
	}

	if cmd.Flags().Changed("mapr-user") {
		opts = append(opts, secret.WithFilterByMaprUser(o.FilterByMaprUser))
	}

	if cmd.Flags().Changed("mapr-uid") {
		opts = append(opts, secret.WithFilterByUID(o.FilterByMaprUID))
	}

	if cmd.Flags().Changed("mapr-gid") {
		opts = append(opts, secret.WithFilterByGID(o.FilterByMaprGID))
	}

	if cmd.Flags().Changed("in-use") && o.FilterByInUse {
		opts = append(opts, secret.WithFilterByInUse())
	}

	if cmd.Flags().Changed("expires-before") {
		opts = append(opts, secret.WithFilterExpiresBefore(o.FilterExpiresBefore.Cast()))
	}

	if cmd.Flags().Changed("show-in-use") && o.ShowInUse {
		opts = append(opts, secret.WithShowInUse())
	}

	// create lister
	lister := secret.NewLister(client, *o.kubernetesConfigFlags.Namespace, opts...)

	// run lister
	items, err := lister.Run()
	if err != nil {
		return err
	}

	// print output
	if err := secret.Print(cmd, items); err != nil {
		return err
	}

	return nil
}

func (o *ListOptions) registerCompletions(cmd *cobra.Command) error {
	err := cmd.RegisterFlagCompletionFunc("output", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return util.CompleteStringValues(listValidOutputFormats, toComplete)
	})
	if err != nil {
		return err
	}

	err = cmd.RegisterFlagCompletionFunc("sort-by", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return util.CompleteStringValues(listValidSortByFields, toComplete)
	})
	if err != nil {
		return err
	}

	return nil
}
