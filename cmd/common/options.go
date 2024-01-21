package common

import (
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericiooptions"
)

type Options struct {
	KubernetesConfigFlags *genericclioptions.ConfigFlags
	IOStreams             genericiooptions.IOStreams

	// Debug flag to enable Debug logging
	Debug bool
}

func NewOptions(flags *genericclioptions.ConfigFlags, streams genericiooptions.IOStreams) *Options {
	return &Options{
		KubernetesConfigFlags: flags,
		IOStreams:             streams,
	}
}
