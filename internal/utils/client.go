package utils

import (
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
)

// ClientFromFlags creates a Kubernetes client from the flags passed to the
// CLI.
func ClientFromFlags(flags *genericclioptions.ConfigFlags) (kubernetes.Interface, error) {
	config, err := flags.ToRESTConfig()
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}
