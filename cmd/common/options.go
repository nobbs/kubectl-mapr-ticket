// Copyright (c) 2024 Alexej Disterhoft
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
//
// SPX-License-Identifier: MIT

package common

import (
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericiooptions"
)

// Options is a struct to hold common options for all commands
type Options struct {
	KubernetesConfigFlags *genericclioptions.ConfigFlags
	IOStreams             genericiooptions.IOStreams

	// Debug flag to enable Debug logging
	Debug bool
}

// NewOptions returns a new common options struct
func NewOptions(flags *genericclioptions.ConfigFlags, streams genericiooptions.IOStreams) *Options {
	return &Options{
		KubernetesConfigFlags: flags,
		IOStreams:             streams,
	}
}
