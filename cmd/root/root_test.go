// Copyright (c) 2024 Alexej Disterhoft
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
//
// SPDX-License-Identifier: MIT

package root_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nobbs/kubectl-mapr-ticket/cmd/claim"
	"github.com/nobbs/kubectl-mapr-ticket/cmd/common"
	"github.com/nobbs/kubectl-mapr-ticket/cmd/inspect"
	. "github.com/nobbs/kubectl-mapr-ticket/cmd/root"
	"github.com/nobbs/kubectl-mapr-ticket/cmd/secret"
	"github.com/nobbs/kubectl-mapr-ticket/cmd/version"
	"github.com/nobbs/kubectl-mapr-ticket/cmd/volume"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericiooptions"
)

func TestNewCmd(t *testing.T) {
	t.Parallel()

	t.Run("print usage information if no command is given", func(t *testing.T) {
		t.Parallel()

		// set up
		ioStreams, _, out, errOut := genericiooptions.NewTestIOStreams()
		flags := genericclioptions.NewConfigFlags(false)

		// run the command
		cmd := NewCmd(flags, ioStreams)
		cmd.SetArgs([]string{})
		err := cmd.Execute()

		// check the error
		assert.NoError(t, err)

		// check the output
		assert.NotEmpty(t, cmd.Use)
		assert.NotEmpty(t, cmd.Long)

		assert.Contains(t, out.String(), cmd.Use)
		assert.Contains(t, out.String(), cmd.Long)

		assert.Empty(t, errOut.String())
	})

	t.Run("check if all subcommands are registered", func(t *testing.T) {
		t.Parallel()

		// set up
		ioStreams, _, _, _ := genericiooptions.NewTestIOStreams()
		flags := genericclioptions.NewConfigFlags(false)

		opts := common.NewOptions(
			flags,
			ioStreams,
		)

		// run the command
		cmd := NewCmd(flags, ioStreams)

		// check the subcommands
		assert.True(t, cmd.HasSubCommands())

		// check if all subcommands are registered
		assert.ElementsMatch(t,
			[]string{
				claim.NewCmd(opts).Use,
				inspect.NewCmd(opts).Use,
				secret.NewCmd(opts).Use,
				version.NewCmd(opts).Use,
				volume.NewCmd(opts).Use,
			},
			func() (cmdCommandsUse []string) {
				for _, c := range cmd.Commands() {
					cmdCommandsUse = append(cmdCommandsUse, c.Use)
				}
				return
			}(),
		)
	})
}
