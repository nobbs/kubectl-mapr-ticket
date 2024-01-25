package common_test

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	. "github.com/nobbs/kubectl-mapr-ticket/cmd/common"
	"github.com/nobbs/kubectl-mapr-ticket/pkg/ticket"

	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

func TestCompleteStringValues(t *testing.T) {
	t.Parallel()

	type args struct {
		values     []string
		toComplete string
	}

	type expected struct {
		suggestions []string
		directive   cobra.ShellCompDirective
	}

	tests := []struct {
		name string
		args args
		want expected
	}{
		{
			name: "EmptyValues",
			args: args{
				values:     []string{},
				toComplete: "",
			},
			want: expected{
				suggestions: nil,
				directive:   cobra.ShellCompDirectiveNoFileComp,
			},
		},
		{
			name: "EmptyToComplete",
			args: args{
				values:     []string{"apple", "banana", "cherry"},
				toComplete: "",
			},
			want: expected{
				suggestions: []string{"apple", "banana", "cherry"},
				directive:   cobra.ShellCompDirectiveNoFileComp,
			},
		},
		{
			name: "SingleValue",
			args: args{
				values:     []string{"apple"},
				toComplete: "a",
			},
			want: expected{
				suggestions: []string{"apple"},
				directive:   cobra.ShellCompDirectiveNoFileComp,
			},
		},
		{
			name: "MultipleValues",
			args: args{
				values:     []string{"apple", "banana", "cherry"},
				toComplete: "b",
			},
			want: expected{
				suggestions: []string{"banana"},
				directive:   cobra.ShellCompDirectiveNoFileComp,
			},
		},
		{
			name: "MultipleValuesMultipleMatches",
			args: args{
				values:     []string{"apple", "banana", "blueberry", "cherry"},
				toComplete: "b",
			},
			want: expected{
				suggestions: []string{"banana", "blueberry"},
				directive:   cobra.ShellCompDirectiveNoFileComp,
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			suggestions, directive := CompleteStringValues(test.args.values, test.args.toComplete)

			assert.Len(t, suggestions, len(test.want.suggestions))
			assert.ElementsMatch(t, test.want.suggestions, suggestions)
			assert.Equal(t, test.want.directive, directive)
		})
	}
}

func TestCompleteNamespaceNames(t *testing.T) {
	t.Parallel()

	type args struct {
		client     kubernetes.Interface
		toComplete string
	}

	type expected struct {
		suggestions []string
		directive   cobra.ShellCompDirective
	}

	tests := []struct {
		name string
		args args
		want expected
	}{
		{
			name: "EmptyToComplete",
			args: args{
				client:     fake.NewSimpleClientset(),
				toComplete: "",
			},
			want: expected{
				suggestions: nil,
				directive:   cobra.ShellCompDirectiveNoFileComp,
			},
		},
		{
			name: "SingleNamespace",
			args: args{
				client: fake.NewSimpleClientset(
					newNamespace("default"),
				),
				toComplete: "d",
			},
			want: expected{
				suggestions: []string{"default"},
				directive:   cobra.ShellCompDirectiveNoFileComp,
			},
		},
		{
			name: "MultipleNamespaces",
			args: args{
				client: fake.NewSimpleClientset(
					newNamespace("default"),
					newNamespace("kube-system"),
					newNamespace("kube-public"),
				),
				toComplete: "k",
			},
			want: expected{
				suggestions: []string{"kube-public", "kube-system"},
				directive:   cobra.ShellCompDirectiveNoFileComp,
			},
		},
		{
			name: "MultipleNamespacesNoMatch",
			args: args{
				client: fake.NewSimpleClientset(
					newNamespace("default"),
					newNamespace("kube-system"),
					newNamespace("kube-public"),
				),
				toComplete: "z",
			},
			want: expected{
				suggestions: nil,
				directive:   cobra.ShellCompDirectiveNoFileComp,
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			suggestions, directive := CompleteNamespaceNames(test.args.client, test.args.toComplete)

			assert.Len(t, suggestions, len(test.want.suggestions))
			assert.ElementsMatch(t, test.want.suggestions, suggestions)
			assert.Equal(t, test.want.directive, directive)
		})
	}
}

func TestCompleteTicketNames(t *testing.T) {
	t.Parallel()

	type args struct {
		client     kubernetes.Interface
		namespace  string
		args       []string
		toComplete string
	}

	type expected struct {
		suggestions []string
		directive   cobra.ShellCompDirective
	}

	tests := []struct {
		name string
		args args
		want expected
	}{
		{
			name: "EmptyToComplete",
			args: args{
				client:     fake.NewSimpleClientset(),
				namespace:  "default",
				args:       []string{},
				toComplete: "",
			},
			want: expected{
				suggestions: nil,
				directive:   cobra.ShellCompDirectiveNoFileComp,
			},
		},
		{
			name: "SingleTicket",
			args: args{
				client: fake.NewSimpleClientset(
					newTicket("default", "ticket-1"),
				),
				namespace:  "default",
				args:       []string{},
				toComplete: "t",
			},
			want: expected{
				suggestions: []string{"ticket-1"},
				directive:   cobra.ShellCompDirectiveNoFileComp,
			},
		},
		{
			name: "MultipleTicketsOneAlreadyCompleted",
			args: args{
				client: fake.NewSimpleClientset(
					newTicket("default", "ticket-1"),
					newTicket("default", "ticket-2"),
					newTicket("default", "ticket-3"),
				),
				namespace:  "default",
				args:       []string{"ticket-2"},
				toComplete: "t",
			},
			want: expected{
				suggestions: []string{"ticket-1", "ticket-3"},
				directive:   cobra.ShellCompDirectiveNoFileComp,
			},
		},
		{
			name: "MultipleTicketsDifferentNamespace",
			args: args{
				client: fake.NewSimpleClientset(
					newTicket("default", "ticket-1"),
					newTicket("kube-system", "ticket-2"),
					newTicket("kube-public", "ticket-3"),
				),
				namespace:  "default",
				args:       []string{},
				toComplete: "t",
			},
			want: expected{
				suggestions: []string{"ticket-1"},
				directive:   cobra.ShellCompDirectiveNoFileComp,
			},
		},
		{
			name: "MultipleTicketsNotInNamespace",
			args: args{
				client: fake.NewSimpleClientset(
					newTicket("default", "ticket-1"),
					newTicket("kube-system", "ticket-2"),
					newTicket("kube-public", "ticket-3"),
				),
				namespace:  "test",
				args:       []string{},
				toComplete: "t",
			},
			want: expected{
				suggestions: nil,
				directive:   cobra.ShellCompDirectiveNoFileComp,
			},
		},
		{
			name: "MultipleTicketsAllNamespaces",
			args: args{
				client: fake.NewSimpleClientset(
					newTicket("default", "ticket-1"),
					newTicket("kube-system", "ticket-2"),
					newTicket("kube-public", "ticket-3"),
				),
				namespace:  metaV1.NamespaceAll,
				args:       []string{},
				toComplete: "t",
			},
			want: expected{
				suggestions: []string{"ticket-1", "ticket-2", "ticket-3"},
				directive:   cobra.ShellCompDirectiveNoFileComp,
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			suggestions, directive := CompleteTicketNames(test.args.client, test.args.namespace, test.args.args, test.args.toComplete)

			assert.Len(t, suggestions, len(test.want.suggestions))
			assert.ElementsMatch(t, test.want.suggestions, suggestions)
			assert.Equal(t, test.want.directive, directive)
		})
	}
}

func newNamespace(name string) *coreV1.Namespace {
	return &coreV1.Namespace{
		ObjectMeta: metaV1.ObjectMeta{
			Name: name,
		},
	}
}

func newTicket(namespace, name string) *coreV1.Secret {
	return &coreV1.Secret{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: map[string][]byte{
			ticket.SecretMaprTicketKey: {},
		},
	}
}
