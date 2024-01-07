package util

import (
	"context"
	"strings"

	"github.com/nobbs/kubectl-mapr-ticket/internal/ticket"
	"github.com/spf13/cobra"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func CompleteStringValues(values []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	var suggestions []string
	for _, v := range values {
		if toComplete == "" || strings.HasPrefix(v, toComplete) {
			suggestions = append(suggestions, v)
		}
	}

	return suggestions, cobra.ShellCompDirectiveNoFileComp
}

func CompleteNamespaceNames(flags *genericclioptions.ConfigFlags, toComplete string) ([]string, cobra.ShellCompDirective) {
	client, err := ClientFromFlags(flags)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	namespaces, err := client.CoreV1().Namespaces().List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var suggestions []string
	for _, ns := range namespaces.Items {
		if toComplete == "" || strings.HasPrefix(ns.Name, toComplete) {
			suggestions = append(suggestions, ns.Name)
		}
	}

	return suggestions, cobra.ShellCompDirectiveNoFileComp
}

func CompleteTicketNames(flags *genericclioptions.ConfigFlags, allNamespaces bool, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	client, err := ClientFromFlags(flags)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	namespace := GetNamespace(flags, allNamespaces)

	secrets, err := client.CoreV1().Secrets(namespace).List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var suggestions []string
	for _, secret := range secrets.Items {
		if _, ok := secret.Data[ticket.SecretMaprTicketKey]; !ok {
			continue
		}

		// skip already completed tickets
		if contains(args, secret.Name) {
			continue
		}

		if toComplete == "" || strings.HasPrefix(secret.Name, toComplete) {
			suggestions = append(suggestions, secret.Name)
		}
	}

	return suggestions, cobra.ShellCompDirectiveNoFileComp
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
