package main

import (
	"os"

	"github.com/nobbs/kubectl-mapr-ticket/cmd/cli"
	"github.com/spf13/pflag"
	"k8s.io/cli-runtime/pkg/genericiooptions"
)

func main() {
	// Create a set of flags to pass to the CLI
	flags := pflag.NewFlagSet("kubectl-mapr-ticket", pflag.ExitOnError)
	pflag.CommandLine = flags

	// Create the root command and execute it
	root := cli.NewRootCmd(genericiooptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr})
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
