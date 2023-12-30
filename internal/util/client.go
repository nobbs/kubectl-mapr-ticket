package util

import (
	apiV1 "k8s.io/api/core/v1"
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

// GetNamespace returns the namespace from the flags passed to the CLI.
func GetNamespace(flags *genericclioptions.ConfigFlags) string {
	namespace, _, err := flags.ToRawKubeConfigLoader().Namespace()

	if err != nil || namespace == "" {
		namespace = apiV1.NamespaceDefault
	}

	return namespace
}
