package util

import (
	"context"
	"strings"

	"github.com/nobbs/kubectl-mapr-ticket/internal/ticket"
	"github.com/spf13/cobra"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// CompleteStringValues returns a list of suggestions for the given available
// values and the toComplete string. If toComplete is empty, all values are
// returned. Otherwise, only values that start with toComplete are returned.
func CompleteStringValues(values []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	var suggestions []string
	for _, v := range values {
		if toComplete == "" || strings.HasPrefix(v, toComplete) {
			suggestions = append(suggestions, v)
		}
	}

	return suggestions, cobra.ShellCompDirectiveNoFileComp
}

// CompleteNamespaceNames returns a list of suggestions for the given available
// namespaces and the toComplete string. If toComplete is empty, all namespaces
// are returned. Otherwise, only namespaces that start with toComplete are
// returned.
func CompleteNamespaceNames(client kubernetes.Interface, toComplete string) ([]string, cobra.ShellCompDirective) {
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

// CompleteTicketNames returns a list of suggestions for the given available
// tickets and the toComplete string. If toComplete is empty, all tickets are
// returned. Otherwise, only tickets that start with toComplete are returned.
// Tickets that have already been completed as part of the command are not
// returned.
func CompleteTicketNames(client kubernetes.Interface, namespace string, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
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
