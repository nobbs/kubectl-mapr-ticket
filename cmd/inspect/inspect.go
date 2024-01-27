// Copyright (c) 2024 Alexej Disterhoft
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
//
// SPX-License-Identifier: MIT

package inspect

import (
	"context"
	"fmt"
	"slices"

	"github.com/spf13/cobra"

	"github.com/nobbs/kubectl-mapr-ticket/cmd/common"
	"github.com/nobbs/kubectl-mapr-ticket/pkg/ticket"
	"github.com/nobbs/kubectl-mapr-ticket/pkg/util"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

const (
	inspectUse   = `inspect`
	inspectShort = "Inspect a MapR ticket either from a secret or locally"
	inspectLong  = `
		Inspect a MapR ticket either from a secret or locally.

		This command will print all the information present in a MapR ticket in a human
		readable format. For local files, both secret manifest as well as MapR ticket
		files are supported.
		`
	inspectExample = `
		# Inspect a MapR ticket from a secret
		%[1]s inspect mapr-ticket-secret --namespace kube-system

		# Inspect a MapR ticket from a file and output in JSON format (default)
		%[1]s inspect -f ./mapr-ticket

		# Inspect a MapR ticket from a secret and output in YAML format
		%[1]s inspect mapr-ticket-secret --namespace kube-system -o yaml
		`
)

var (
	inspectValidOutputFormats = []string{"json", "yaml"}
)

type options struct {
	*common.Options

	// Args are the arguments passed to the command
	args []string

	// SecretName is the name of the secret to inspect
	SecretName string

	// OutputFormat is the format to use for output
	OutputFormat string

	// File is the path to the MapR ticket file
	File string
}

func newOptions(opts *common.Options) *options {
	return &options{
		Options: opts,
	}
}

func NewCmd(opts *common.Options) *cobra.Command {
	o := newOptions(opts)

	cmd := &cobra.Command{
		Aliases:      []string{"describe", "i"},
		Use:          inspectUse,
		Short:        inspectShort,
		Long:         common.CliLongDesc(inspectLong),
		Example:      common.CliExample(inspectExample, common.CliBinName),
		Args:         cobra.MaximumNArgs(1),
		SilenceUsage: true,
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			// we only want one argument, so don't complete once we have one
			if len(args) > 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}

			// set namespace based on flags
			namespace := util.GetNamespace(o.KubernetesConfigFlags, false)
			o.KubernetesConfigFlags.Namespace = &namespace

			// get client
			client, err := util.ClientFromFlags(o.KubernetesConfigFlags)
			if err != nil {
				return nil, cobra.ShellCompDirectiveError
			}

			return common.CompleteTicketNames(client, namespace, args, toComplete)
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
	cmd.Flags().StringVarP(&o.OutputFormat, "output", "o", "json", fmt.Sprintf("Output format. One of (%s)", common.StringSliceToFlagOptions(inspectValidOutputFormats)))
	cmd.Flags().StringVarP(&o.File, "file", "f", "", "Path to the MapR ticket file")

	// register completions for flags
	if err := o.registerCompletions(cmd); err != nil {
		panic(err)
	}

	return cmd
}

func (o *options) Complete(cmd *cobra.Command, args []string) error {
	// parse the arguments
	o.args = args

	// set secret name based on args
	switch len(args) {
	case 0:
		if o.File == "" {
			return fmt.Errorf("either provide a secret name or a file via --file")
		}
	case 1:
		o.SecretName = args[0]
	default:
		return fmt.Errorf("too many arguments provided, either provide a secret name or a file via --file")
	}

	// set namespace based on flags
	ns := util.GetNamespace(o.KubernetesConfigFlags, false)
	o.KubernetesConfigFlags.Namespace = &ns

	return nil
}

func (o *options) Validate() error {
	// validate output format
	if !slices.Contains(inspectValidOutputFormats, o.OutputFormat) {
		return fmt.Errorf("invalid output format %q. Must be one of (%s)", o.OutputFormat, common.StringSliceToFlagOptions(inspectValidOutputFormats))
	}

	return nil
}

func (o *options) Run(cmd *cobra.Command, args []string) error {
	// if we have a secret name, inspect the secret
	if o.SecretName != "" {
		return o.inspectSecret()
	}

	// if we have a file, inspect the file
	if o.File != "" {
		return o.inspectFile()
	}

	return nil
}

// inspectSecret inspects a MapR ticket read from a secret in a Kubernetes cluster
func (o *options) inspectSecret() error {
	client, err := util.ClientFromFlags(o.KubernetesConfigFlags)
	if err != nil {
		return err
	}

	// get single secret
	secret, err := client.CoreV1().Secrets(*o.KubernetesConfigFlags.Namespace).Get(context.TODO(), o.SecretName, metaV1.GetOptions{})
	if err != nil {
		return err
	}

	// get ticket from secret
	ticket, err := ticket.NewMaprTicketFromSecret(secret)
	if err != nil {
		return err
	}

	// print ticket
	if err := o.printTicket(ticket); err != nil {
		return err
	}

	return nil
}

// inspectFile inspects a MapR ticket read from a file
func (o *options) inspectFile() error {
	bytes, err := util.ReadFile(o.File)
	if err != nil {
		return err
	}

	// get ticket from file
	ticket, err := ticket.NewMaprTicketFromBytes(bytes)
	if err != nil {
		return err
	}

	// print ticket
	if err := o.printTicket(ticket); err != nil {
		return err
	}

	return nil
}

// printTicket prints the ticket in the configured output format
func (o *options) printTicket(ticket *ticket.Ticket) error {
	switch o.OutputFormat {
	case "json":
		return printJSON(ticket)
	case "yaml":
		return printYAML(ticket)
	default:
		return fmt.Errorf("invalid output format %q. Must be one of (%s)", o.OutputFormat, common.StringSliceToFlagOptions(inspectValidOutputFormats))
	}
}

// printJSON prints the ticket in JSON format
func printJSON(ticket *ticket.Ticket) error {
	fmt.Println(ticket.AsMaprTicket().String())

	return nil
}

// printYAML prints the ticket in YAML format
func printYAML(ticket *ticket.Ticket) error {
	jsonString := ticket.AsMaprTicket().String()
	yamlBytes, err := yaml.JSONToYAML([]byte(jsonString))
	if err != nil {
		return err
	}

	fmt.Println(string(yamlBytes))

	return nil
}

func (o *options) registerCompletions(cmd *cobra.Command) error {
	err := cmd.RegisterFlagCompletionFunc("output", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return common.CompleteStringValues(inspectValidOutputFormats, toComplete)
	})
	if err != nil {
		return err
	}

	if err := cmd.MarkFlagFilename("file", "yaml", "yml", "json"); err != nil {
		return err
	}

	return nil
}
