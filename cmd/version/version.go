// Copyright (c) 2024 Alexej Disterhoft
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
//
// SPDX-License-Identifier: MIT

// Package version provides the version command for the application.
package version

import (
	"github.com/spf13/cobra"

	"github.com/nobbs/kubectl-mapr-ticket/cmd/common"
	"github.com/nobbs/kubectl-mapr-ticket/pkg/version"
)

const (
	versionUse   = `version`
	versionShort = "Print the version of kubectl-mapr-ticket and exit"
	versionLong  = `
		Print the version of kubectl-mapr-ticket and exit.
		`
)

// options holds the options for 'version' sub command
type options struct {
	// embed common options from RootCmdOptions
	*common.Options
}

func newOptions(opts *common.Options) *options {
	return &options{
		Options: opts,
	}
}

func NewCmd(rootOpts *common.Options) *cobra.Command {
	o := newOptions(rootOpts)

	cmd := &cobra.Command{
		Aliases: []string{"v"},
		Use:     versionUse,
		Short:   versionShort,
		Long:    common.CliLongDesc(versionLong),
		Args:    cobra.NoArgs,
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
		Run: func(cmd *cobra.Command, args []string) {
			PrintVersionInfo(cmd)
		},
	}

	// set IOStreams for the command
	cmd.SetIn(o.IOStreams.In)
	cmd.SetOut(o.IOStreams.Out)
	cmd.SetErr(o.IOStreams.ErrOut)

	return cmd
}

func PrintVersionInfo(cmd *cobra.Command) {
	versionInfo := version.NewVersion()
	cmd.Println(versionInfo)
}
