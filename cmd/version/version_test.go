// Copyright (c) 2024 Alexej Disterhoft
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
//
// SPDX-License-Identifier: MIT

package version_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	"github.com/nobbs/kubectl-mapr-ticket/cmd/common"
	. "github.com/nobbs/kubectl-mapr-ticket/cmd/version"
	"github.com/nobbs/kubectl-mapr-ticket/pkg/version"

	"k8s.io/cli-runtime/pkg/genericiooptions"
)

func TestNewCmd(t *testing.T) {
	t.Parallel()

	t.Run("valid command call", func(t *testing.T) {
		t.Parallel()

		// set up
		ioStreams, _, out, errOut := genericiooptions.NewTestIOStreams()
		opts := &common.Options{
			IOStreams: ioStreams,
		}

		// run the command
		cmd := NewCmd(opts)
		cmd.SetArgs([]string{})
		err := cmd.Execute()

		// check the error
		assert.NoError(t, err)

		// check the output
		assert.Empty(t, errOut.String())
		assert.Equal(t, version.NewVersion().String(), strings.TrimSpace(out.String()))
	})

	t.Run("usage information", func(t *testing.T) {
		t.Parallel()

		// set up
		ioStreams, _, out, errOut := genericiooptions.NewTestIOStreams()
		opts := &common.Options{
			IOStreams: ioStreams,
		}

		// run the command
		cmd := NewCmd(opts)
		cmd.SetArgs([]string{"--help"})
		err := cmd.Execute()

		// check the error
		assert.NoError(t, err)

		// check the output
		assert.NotEmpty(t, cmd.Use)
		assert.NotEmpty(t, cmd.Short)
		assert.NotEmpty(t, cmd.Long)

		assert.Contains(t, out.String(), "Usage:")
		assert.Contains(t, out.String(), cmd.Use)
		assert.Contains(t, out.String(), cmd.Short)
		assert.Contains(t, out.String(), cmd.Long)

		assert.Empty(t, errOut.String())
	})

	t.Run("invalid command call with arguments", func(t *testing.T) {
		t.Parallel()

		// set up
		ioStreams, _, out, errOut := genericiooptions.NewTestIOStreams()
		opts := &common.Options{
			IOStreams: ioStreams,
		}

		// run the command
		cmd := NewCmd(opts)
		cmd.SetArgs([]string{"arg1"})
		err := cmd.Execute()

		// check the error
		assert.Error(t, err)

		// check the output
		assert.Contains(t, out.String(), "Usage:")
		assert.Contains(t, errOut.String(), "unknown command")
	})

	t.Run("invalid command call with invalid flag", func(t *testing.T) {
		t.Parallel()

		// set up
		ioStreams, _, out, errOut := genericiooptions.NewTestIOStreams()
		opts := &common.Options{
			IOStreams: ioStreams,
		}

		// run the command
		cmd := NewCmd(opts)
		cmd.SetArgs([]string{"--invalid-flag"})
		err := cmd.Execute()

		// check the error
		assert.Error(t, err)

		// check the output
		assert.Contains(t, out.String(), "Usage:")
		assert.Contains(t, errOut.String(), "unknown flag")
	})
}

func TestPrintVersionInfo(t *testing.T) {
	t.Parallel()

	// run the command
	cmd := &cobra.Command{}
	buf := new(bytes.Buffer)
	cmd.SetOutput(buf)

	PrintVersionInfo(cmd)

	// get the expected output
	expectedOutput := version.NewVersion().String()

	// compare the output
	assert.Equal(t, strings.TrimSpace(expectedOutput), strings.TrimSpace(buf.String()))
}
