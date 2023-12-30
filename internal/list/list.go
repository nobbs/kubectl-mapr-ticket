package list

import (
	"context"

	"github.com/nobbs/kubectl-mapr-ticket/internal/ticket"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	typedV1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type Item struct {
	secret *coreV1.Secret
	ticket *ticket.MaprTicket
}

type List struct {
	client    typedV1.SecretInterface
	namespace string

	filterOnlyExpired   bool
	filterOnlyUnexpired bool
	filterByMaprCluster *string
	filterByMaprUser    *string
}

type ListOption func(*List)

func WithFilterByMaprCluster(cluster string) ListOption {
	return func(l *List) {
		l.filterByMaprCluster = &cluster
	}
}

func WithFilterByMaprUser(user string) ListOption {
	return func(l *List) {
		l.filterByMaprUser = &user
	}
}

func WithFilterOnlyExpired() ListOption {
	return func(l *List) {
		l.filterOnlyExpired = true
	}
}

func WithFilterOnlyUnexpired() ListOption {
	return func(l *List) {
		l.filterOnlyUnexpired = true
	}
}

func NewList(client kubernetes.Interface, namespace string, opts ...ListOption) *List {
	const (
		defaultFilterOnlyExpired = false
	)

	l := &List{
		client:            client.CoreV1().Secrets(namespace),
		namespace:         namespace,
		filterOnlyExpired: defaultFilterOnlyExpired,
	}

	for _, opt := range opts {
		opt(l)
	}

	return l
}

func (l *List) Run() ([]Item, error) {
	secrets, err := l.client.List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		return nil, err
	}

	// filter secrets that don't contain a ticket
	filtered := l.filterSecretsWithMaprTicketKey(secrets.Items)

	// convert secrets to items, parse all tickets
	var items []Item
	for i := range filtered {
		ticket, err := ticket.NewTicketFromSecret(&filtered[i])
		if err != nil {
			continue
		}

		items = append(items, Item{
			secret: &filtered[i],
			ticket: ticket,
		})
	}

	// filter items to only expired tickets, if requested
	if l.filterOnlyExpired {
		items = l.filterItemsOnlyExpired(items)
	}

	// filter items to only unexpired tickets, if requested
	if l.filterOnlyUnexpired {
		items = l.filterItemsOnlyUnexpired(items)
	}

	// filter items to only tickets for the specified MapR cluster, if requested
	if l.filterByMaprCluster != nil && *l.filterByMaprCluster != "" {
		items = l.filterItemsByMaprCluster(items)
	}

	// filter items to only tickets for the specified MapR user, if requested
	if l.filterByMaprUser != nil && *l.filterByMaprUser != "" {
		items = l.filterItemsByMaprUser(items)
	}

	return items, nil
}

func (l *List) filterSecretsWithMaprTicketKey(secrets []coreV1.Secret) []coreV1.Secret {
	var filtered []coreV1.Secret

	for _, secret := range secrets {
		if ticket.SecretContainsMaprTicket(&secret) {
			filtered = append(filtered, secret)
		}
	}

	return filtered
}

func (l *List) filterItemsOnlyExpired(items []Item) []Item {
	var filtered []Item

	for _, item := range items {
		if item.ticket.IsExpired() {
			filtered = append(filtered, item)
		}
	}

	return filtered
}

func (l *List) filterItemsOnlyUnexpired(items []Item) []Item {
	var filtered []Item

	for _, item := range items {
		if !item.ticket.IsExpired() {
			filtered = append(filtered, item)
		}
	}

	return filtered
}

func (l *List) filterItemsByMaprCluster(items []Item) []Item {
	var filtered []Item

	for _, item := range items {
		if item.ticket.Cluster == *l.filterByMaprCluster {
			filtered = append(filtered, item)
		}
	}

	return filtered
}

func (l *List) filterItemsByMaprUser(items []Item) []Item {
	var filtered []Item

	for _, item := range items {
		if item.ticket.UserCreds.GetUserName() == *l.filterByMaprUser {
			filtered = append(filtered, item)
		}
	}

	return filtered
}
