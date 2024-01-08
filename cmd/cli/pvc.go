package cli

import (
	"fmt"

	"github.com/nobbs/kubectl-mapr-ticket/internal/pvc"
	"github.com/nobbs/kubectl-mapr-ticket/internal/util"
	"github.com/spf13/cobra"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	pvcUse   = `pvc`
	pvcShort = "List all persistent volumes claims that use a MapR ticket in the current namespace"
	pvcLong  = `
		List all persistent volumes claims that use a MapR ticket in the current namespace.

		By default, this command lists all persistent volumes claims that use a MapR ticket in the current namespace.
		`
)

type PVCOptions struct {
	*rootCmdOptions

	// OutputFormat is the format to use for output
	OutputFormat string

	// AllNamespaces indicates whether to list secrets in all namespaces
	AllNamespaces bool
}

func NewPVCOptions(rootOpts *rootCmdOptions) *PVCOptions {
	return &PVCOptions{
		rootCmdOptions: rootOpts,
	}
}

func newPVCCmd(rootOpts *rootCmdOptions) *cobra.Command {
	o := NewPVCOptions(rootOpts)

	cmd := &cobra.Command{
		Use:   pvcUse,
		Short: pvcShort,
		Long:  util.CliLongDesc(pvcLong),
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

	return cmd
}

func (o *PVCOptions) Complete(cmd *cobra.Command, args []string) error {
	// set namespace
	if o.kubernetesConfigFlags.Namespace == nil || *o.kubernetesConfigFlags.Namespace == "" {
		namespace := util.GetNamespace(o.kubernetesConfigFlags, o.AllNamespaces)
		o.kubernetesConfigFlags.Namespace = &namespace
	}

	// reset namespace if --all-namespaces is set
	if o.AllNamespaces {
		namespaceAll := metaV1.NamespaceAll
		o.kubernetesConfigFlags.Namespace = &namespaceAll
	}

	return nil
}

func (o *PVCOptions) Validate() error {
	return nil
}

func (o *PVCOptions) Run(cmd *cobra.Command, args []string) error {
	client, err := util.ClientFromFlags(o.kubernetesConfigFlags)
	if err != nil {
		return err
	}

	opts := []pvc.ListerOption{}

	lister := pvc.NewLister(client, *o.kubernetesConfigFlags.Namespace, opts...)

	pvcs, err := lister.Run()
	if err != nil {
		return err
	}

	// print the list of pvcs
	fmt.Printf("Found %d PVCs\n", len(pvcs))

	return nil
}
