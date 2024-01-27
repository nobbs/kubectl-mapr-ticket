// Copyright (c) 2024 Alexej Disterhoft
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
//
// SPX-License-Identifier: MIT

// Package secret provides the secret command for the application.
package secret

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nobbs/kubectl-mapr-ticket/cmd/common"
	"github.com/nobbs/kubectl-mapr-ticket/pkg/secret"
	"github.com/nobbs/kubectl-mapr-ticket/pkg/util"
	"github.com/nobbs/kubectl-mapr-ticket/pkg/volume"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// command string constants for use in help and usage text
const (
	secretUse   = `secret`
	secretShort = "List all secrets containing MapR tickets in the current namespace"
	secretLong  = `
		List all secrets containing MapR tickets in the current namespace and print
		some information about them.
		`
	secretExample = `
		# List all MapR tickets in the current namespace
		%[1]s secret

		# List all MapR tickets in all namespaces
		%[1]s secret --all-namespaces

		# List only expired MapR tickets
		%[1]s secret --only-expired

		# List only MapR tickets that expire in the next 7 days
		%[1]s secret --expires-before 7d

		# List MapR tickets for a specific MapR user in all namespaces
		%[1]s secret --mapr-user mapr --all-namespaces

		# List MapR tickets with number of persistent volumes that use them
		%[1]s secret --show-in-use
		`
)

var (
	// valid output formats for the command
	secretValidOutputFormats = []string{"table", "wide"}
)

type options struct {
	*common.Options

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
	FilterExpiresBefore common.DurationValue

	// ShowInUse indicates whether to show only secrets that are in use by a
	// persistent volume
	ShowInUse bool
}

func newOptions(opts *common.Options) *options {
	return &options{
		Options: opts,
	}
}

// NewCmd creates a new secret command for the application.
func NewCmd(opts *common.Options) *cobra.Command {
	o := newOptions(opts)

	cmd := &cobra.Command{
		Aliases: []string{"s"},
		Use:     secretUse,
		Short:   secretShort,
		Long:    common.CliLongDesc(secretLong),
		Example: common.CliExample(secretExample, common.CliBinName),
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
	cmd.Flags().StringVarP(&o.OutputFormat, "output", "o", "table", fmt.Sprintf("Output format. One of (%s)", common.StringSliceToFlagOptions(secretValidOutputFormats)))
	cmd.Flags().BoolVarP(&o.AllNamespaces, "all-namespaces", "A", false, "If true, list the requested object(s) across all namespaces. Namespace in current context is ignored even if specified with --namespace.")
	cmd.Flags().StringSliceVar(&o.SortBy, "sort-by", nil, fmt.Sprintf("Sort list of secrets by the specified fields. One of (%s)", common.StringSliceToFlagOptions(secret.SortOptionsList)))
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

// Complete sets any default values for the command flags not handled automatically
func (o *options) Complete(cmd *cobra.Command, args []string) error {
	// set namespace based on flags
	ns := util.GetNamespace(o.KubernetesConfigFlags, o.AllNamespaces)
	o.KubernetesConfigFlags.Namespace = &ns

	return nil
}

// Validate ensures that all required arguments and flag values are provided
func (o *options) Validate() error {
	// validate output format
	if o.OutputFormat != "table" && o.OutputFormat != "wide" {
		return fmt.Errorf("invalid output format: %s. Must be one of: table|wide", o.OutputFormat)
	}

	// validate sort options
	if err := util.ValidateSortOptions(secret.SortOptionsList, o.SortBy); err != nil {
		return err
	}

	return nil
}

// Run executes the command
func (o *options) Run(cmd *cobra.Command, args []string) error {
	client, err := util.ClientFromFlags(o.KubernetesConfigFlags)
	if err != nil {
		return err
	}

	// create list options and pass them to the lister
	opts := []secret.ListerOption{}

	if cmd.Flags().Changed("sort-by") && o.SortBy != nil {
		// convert sort options to SortOptions
		sortOptions := make([]secret.SortOption, 0, len(o.SortBy))
		for _, sortBy := range o.SortBy {
			sortOptions = append(sortOptions, secret.SortOption(sortBy))
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

		// add volume lister, since we need to know which secrets are in use
		volumeLister := volume.NewLister(client, util.SecretAll, metaV1.NamespaceAll)
		opts = append(opts, secret.WithVolumeLister(volumeLister))
	}

	if cmd.Flags().Changed("expires-before") {
		opts = append(opts, secret.WithFilterExpiresBefore(o.FilterExpiresBefore.Duration()))
	}

	if cmd.Flags().Changed("show-in-use") && o.ShowInUse {
		opts = append(opts, secret.WithShowInUse())

		// add volume lister, since we need to know which secrets are in use
		volumeLister := volume.NewLister(client, util.SecretAll, metaV1.NamespaceAll)
		opts = append(opts, secret.WithVolumeLister(volumeLister))
	}

	// create lister
	lister := secret.NewLister(client, *o.KubernetesConfigFlags.Namespace, opts...)

	// run lister
	tickets, err := lister.List()
	if err != nil {
		return err
	}

	// print output
	if err := secret.Print(cmd, tickets); err != nil {
		return err
	}

	return nil
}

// registerCompletions registers completions for the command flags
func (o *options) registerCompletions(cmd *cobra.Command) error {
	err := cmd.RegisterFlagCompletionFunc("output", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return common.CompleteStringValues(secretValidOutputFormats, toComplete)
	})
	if err != nil {
		return err
	}

	err = cmd.RegisterFlagCompletionFunc("sort-by", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return common.CompleteStringSliceValues(secret.SortOptionsList, toComplete)
	})
	if err != nil {
		return err
	}

	return nil
}
