package list

import (
	"context"

	"github.com/nobbs/kubectl-mapr-ticket/internal/ticket"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	typedV1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type ListItem struct {
	secret *coreV1.Secret
	ticket *ticket.MaprTicket
}

type Lister struct {
	client    typedV1.SecretInterface
	namespace string

	filterOnlyExpired   bool
	filterOnlyUnexpired bool
	filterByMaprCluster *string
	filterByMaprUser    *string
	filterByUID         *uint32
	filterByGID         *uint32
}

type ListerOption func(*Lister)

func WithFilterByMaprCluster(cluster string) ListerOption {
	return func(l *Lister) {
		l.filterByMaprCluster = &cluster
	}
}

func WithFilterByMaprUser(user string) ListerOption {
	return func(l *Lister) {
		l.filterByMaprUser = &user
	}
}

func WithFilterByUID(uid uint32) ListerOption {
	return func(l *Lister) {
		l.filterByUID = &uid
	}
}

func WithFilterByGID(gid uint32) ListerOption {
	return func(l *Lister) {
		l.filterByGID = &gid
	}
}

func WithFilterOnlyExpired() ListerOption {
	return func(l *Lister) {
		l.filterOnlyExpired = true
	}
}

func WithFilterOnlyUnexpired() ListerOption {
	return func(l *Lister) {
		l.filterOnlyUnexpired = true
	}
}

// NewLister creates a new Lister
func NewLister(client kubernetes.Interface, namespace string, opts ...ListerOption) *Lister {
	const (
		defaultFilterOnlyExpired   = false
		defaultFilterOnlyUnexpired = false
	)

	l := &Lister{
		client:              client.CoreV1().Secrets(namespace),
		namespace:           namespace,
		filterOnlyExpired:   defaultFilterOnlyExpired,
		filterOnlyUnexpired: defaultFilterOnlyUnexpired,
	}

	for _, opt := range opts {
		opt(l)
	}

	return l
}

func (l *Lister) Run() ([]ListItem, error) {
	secrets, err := l.client.List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		return nil, err
	}

	// convert secrets to items, parse all tickets
	items := l.parseSecretsToItems(secrets.Items)

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

	// filter items to only tickets for the specified UID, if requested
	if l.filterByUID != nil {
		items = l.filterItemsByUID(items)
	}

	// filter items to only tickets for the specified GID, if requested
	if l.filterByGID != nil {
		items = l.filterItemsByGID(items)
	}

	return items, nil
}

// filterSecretsWithMaprTicketKey filters secrets to only those that contain a MapR ticket key
func (l *Lister) filterSecretsWithMaprTicketKey(secrets []coreV1.Secret) []coreV1.Secret {
	var filtered []coreV1.Secret

	for _, secret := range secrets {
		if ticket.SecretContainsMaprTicket(&secret) {
			filtered = append(filtered, secret)
		}
	}

	return filtered
}

// parseSecretsToItems parses secrets to items, ignoring secrets that don't contain a MapR ticket
func (l *Lister) parseSecretsToItems(secrets []coreV1.Secret) []ListItem {
	var items []ListItem

	for i := range l.filterSecretsWithMaprTicketKey(secrets) {
		ticket, err := ticket.NewTicketFromSecret(&secrets[i])
		if err != nil {
			continue
		}

		items = append(items, ListItem{
			secret: &secrets[i],
			ticket: ticket,
		})
	}

	return items
}

// filterItemsOnlyExpired filters items to only tickets that are expired already
func (l *Lister) filterItemsOnlyExpired(items []ListItem) []ListItem {
	var filtered []ListItem

	for _, item := range items {
		if item.ticket.IsExpired() {
			filtered = append(filtered, item)
		}
	}

	return filtered
}

// filterItemsOnlyUnexpired filters items to only tickets that are not expired yet
func (l *Lister) filterItemsOnlyUnexpired(items []ListItem) []ListItem {
	var filtered []ListItem

	for _, item := range items {
		if !item.ticket.IsExpired() {
			filtered = append(filtered, item)
		}
	}

	return filtered
}

// filterItemsByMaprCluster filters items to only tickets for the specified MapR cluster
func (l *Lister) filterItemsByMaprCluster(items []ListItem) []ListItem {
	var filtered []ListItem

	for _, item := range items {
		if item.ticket.Cluster == *l.filterByMaprCluster {
			filtered = append(filtered, item)
		}
	}

	return filtered
}

// filterItemsByMaprUser filters items to only tickets for the specified MapR user
func (l *Lister) filterItemsByMaprUser(items []ListItem) []ListItem {
	var filtered []ListItem

	for _, item := range items {
		if item.ticket.UserCreds.GetUserName() == *l.filterByMaprUser {
			filtered = append(filtered, item)
		}
	}

	return filtered
}

// filterItemsByUID filters items to only tickets for the specified UID
func (l *Lister) filterItemsByUID(items []ListItem) []ListItem {
	var filtered []ListItem

	for _, item := range items {
		if *item.ticket.UserCreds.Uid == *l.filterByUID {
			filtered = append(filtered, item)
		}
	}

	return filtered
}

// filterItemsByGID filters items to only tickets for the specified GID
func (l *Lister) filterItemsByGID(items []ListItem) []ListItem {
	var filtered []ListItem

	for _, item := range items {
		// check if GID is in the list of GIDs
		for _, gid := range item.ticket.UserCreds.Gids {
			if gid == *l.filterByGID {
				filtered = append(filtered, item)
				break
			}
		}
	}

	return filtered
}
