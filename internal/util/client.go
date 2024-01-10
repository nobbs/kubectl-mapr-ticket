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

// GetNamespace returns the namespace to use for the command. If allNamespaces
// is true, the namespace is set to metaV1.NamespaceAll. Otherwise, the
// namespace is set from the context or value of the --namespace flag.
func GetNamespace(flags *genericclioptions.ConfigFlags, allNamespaces bool) string {
	// get namespace from kubeconfig context or --namespace flag
	namespace, _, err := flags.ToRawKubeConfigLoader().Namespace()

	// if no namespace is set, use the default namespace
	if err != nil || namespace == "" {
		namespace = apiV1.NamespaceDefault
	}

	// if allNamespaces is set, override the namespace with metaV1.NamespaceAll
	if allNamespaces {
		namespace = apiV1.NamespaceAll
	}

	return namespace
}
