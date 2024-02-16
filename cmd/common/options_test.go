// Copyright (c) 2024 Alexej Disterhoft
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
//
// SPDX-License-Identifier: MIT

package common_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/nobbs/kubectl-mapr-ticket/cmd/common"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericiooptions"
)

// nolint:paralleltest
func TestNewOptions(t *testing.T) {
	// Create a new ConfigFlags and IOStreams
	flags := &genericclioptions.ConfigFlags{}
	streams := genericiooptions.IOStreams{}

	// Call the NewOptions function
	options := NewOptions(flags, streams)

	// Verify that the returned options match the input values
	assert.Equal(t, flags, options.KubernetesConfigFlags)
	assert.Equal(t, streams, options.IOStreams)
}
