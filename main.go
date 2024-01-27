// Copyright (c) 2024 Alexej Disterhoft
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
//
// SPDX-License-Identifier: MIT

// Package main implements the main entry point for the kubectl-mapr-ticket CLI. It is responsible for
// creating the root command and executing it.
//
// The root command is responsible for creating the subcommands and executing them. The subcommands
// are responsible for the actual work.
//
// If you are interested in the actual plugin documentation, please refer to the README.md file.
package main

import (
	"os"

	"github.com/spf13/pflag"

	"github.com/nobbs/kubectl-mapr-ticket/cmd/root"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericiooptions"
)

const (
	// Name of the CLI
	cliName = "kubectl mapr-ticket"
)

func main() {
	// Create a set of flags to pass to the CLI
	flags := pflag.NewFlagSet(cliName, pflag.ExitOnError)
	pflag.CommandLine = flags

	// Create a set of default Kubernetes flags and IOStreams
	kubernetesConfigFlags := genericclioptions.NewConfigFlags(true)
	streams := genericiooptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}

	// Create the root command and execute it
	root := root.NewCmd(kubernetesConfigFlags, streams)
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
