// Copyright (c) 2024 Alexej Disterhoft
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
//
// SPDX-License-Identifier: MIT

// Package claim provides the claim command for the application.
package claim

import (
	"fmt"

	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"

	"github.com/nobbs/kubectl-mapr-ticket/cmd/common"
	"github.com/nobbs/kubectl-mapr-ticket/pkg/claim"
	"github.com/nobbs/kubectl-mapr-ticket/pkg/secret"
	"github.com/nobbs/kubectl-mapr-ticket/pkg/util"
)

// command string constants for use in help and usage text
const (
	claimUse   = `claim`
	claimShort = "List all persistent volumes claims that use a MapR ticket in the current namespace"
	claimLong  = `
		List all persistent volumes claims that use a MapR ticket in the current namespace.

		By default, this command lists all persistent volume claims that use a MapR ticket in the current namespace.
		`
	claimExample = `
		# List all persistent volumes claims in the current namespace that use a MapR ticket
		%[1]s claim

		# List all persistent volumes claims in all namespaces that use a MapR ticket
		%[1]s claim --all-namespaces

		# List all persistent volumes claims in all namespaces that use a MapR ticket, sorted by expiration date
		%[1]s claim --all-namespaces --sort-by expiryTime
		`
)

var (
	// valid output formats for the command
	claimValidOutputFormats = []string{"table", "wide"}
)

type options struct {
	*common.Options

	// OutputFormat is the format to use for output
	OutputFormat string

	// AllNamespaces indicates whether to list secrets in all namespaces
	AllNamespaces bool

	// SortBy is the list of fields to sort by
	SortBy []string
}

func newOptions(opts *common.Options) *options {
	return &options{
		Options: opts,
	}
}

// NewCmd creates a new claim command for the application.
func NewCmd(opts *common.Options) *cobra.Command {
	o := newOptions(opts)

	cmd := &cobra.Command{
		Aliases: []string{"pvc"},
		Use:     claimUse,
		Short:   claimShort,
		Long:    common.CliLongDesc(claimLong),
		Example: common.CliExample(claimExample, common.CliBinName),
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

	// set IOStreams for this command
	cmd.SetIn(o.IOStreams.In)
	cmd.SetOut(o.IOStreams.Out)
	cmd.SetErr(o.IOStreams.ErrOut)

	// add flags
	cmd.Flags().StringVarP(&o.OutputFormat, "output", "o", "table", fmt.Sprintf("Output format. One of (%s)", common.StringSliceToFlagOptions(claimValidOutputFormats)))
	cmd.Flags().BoolVarP(&o.AllNamespaces, "all-namespaces", "A", false, "List persistent volumes claims that use a MapR ticket in all namespaces")
	cmd.Flags().StringSliceVar(&o.SortBy, "sort-by", []string{}, fmt.Sprintf("Sort list of persistent volumes claims by the specified fields. One or more of (%s)", common.StringSliceToFlagOptions(claim.SortOptionsList)))

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
	if !slices.Contains(claimValidOutputFormats, o.OutputFormat) {
		return fmt.Errorf("invalid output format %q. Must be one of (%s)", o.OutputFormat, common.StringSliceToFlagOptions(claimValidOutputFormats))
	}

	// ensure that the sort options are valid
	if err := util.ValidateSortOptions(claim.SortOptionsList, o.SortBy); err != nil {
		return err
	}

	return nil
}

// Run executes the command logic
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
	opts := []claim.ListerOption{
		claim.WithSecretLister(secretLister),
	}

	// set sort options
	if cmd.Flags().Changed("sort-by") && o.SortBy != nil {
		sortOptions := make([]claim.SortOption, 0, len(o.SortBy))
		for _, sortBy := range o.SortBy {
			sortOptions = append(sortOptions, claim.SortOption(sortBy))
		}

		opts = append(opts, claim.WithSortBy(sortOptions))
	}

	// create lister
	lister := claim.NewLister(
		client,
		*o.KubernetesConfigFlags.Namespace,
		opts...,
	)

	// run lister
	volumeClaims, err := lister.List()
	if err != nil {
		return err
	}

	// print output
	if err := claim.Print(cmd, volumeClaims); err != nil {
		return err
	}

	return nil
}

// registerCompletions registers completions for the command flags
func (o *options) registerCompletions(cmd *cobra.Command) error {
	err := cmd.RegisterFlagCompletionFunc("output", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return common.CompleteStringValues(claimValidOutputFormats, toComplete)
	})
	if err != nil {
		return err
	}

	err = cmd.RegisterFlagCompletionFunc("sort-by", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return common.CompleteStringSliceValues(claim.SortOptionsList, toComplete)
	})
	if err != nil {
		return err
	}

	return nil
}
