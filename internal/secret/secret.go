package secret

import (
	"context"
	"time"

	"github.com/nobbs/kubectl-mapr-ticket/internal/ticket"
	"github.com/nobbs/kubectl-mapr-ticket/internal/volume"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type ListItem struct {
	Secret *coreV1.Secret     `json:"originalSecret"`
	Ticket *ticket.MaprTicket `json:"parsedTicket"`
	InUse  uint32             `json:"inUse"`
}

type Lister struct {
	client    kubernetes.Interface
	namespace string

	filterOnlyExpired   bool
	filterOnlyUnexpired bool
	filterByMaprCluster *string
	filterByMaprUser    *string
	filterByUID         *uint32
	filterByGID         *uint32
	filterByInUse       bool
	filterExpiresBefore time.Duration
	showInUse           bool
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

func WithFilterByInUse() ListerOption {
	return func(l *Lister) {
		l.filterByInUse = true
	}
}

func WithFilterExpiresBefore(expiresBefore time.Duration) ListerOption {
	return func(l *Lister) {
		l.filterExpiresBefore = expiresBefore
	}
}

func WithShowInUse() ListerOption {
	return func(l *Lister) {
		l.showInUse = true
	}
}

// NewLister creates a new Lister
func NewLister(client kubernetes.Interface, namespace string, opts ...ListerOption) *Lister {
	const (
		defaultFilterOnlyExpired   = false
		defaultFilterOnlyUnexpired = false
	)

	l := &Lister{
		client:              client,
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
	secrets, err := l.client.CoreV1().Secrets(l.namespace).List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		return nil, err
	}

	// convert secrets to items, parse all tickets
	items := parseSecretsToItems(secrets.Items)

	// filter items to only expired tickets, if requested
	if l.filterOnlyExpired {
		items = filterItemsOnlyExpired(items)
	}

	// filter items to only unexpired tickets, if requested
	if l.filterOnlyUnexpired {
		items = filterItemsOnlyUnexpired(items)
	}

	// filter items to only tickets for the specified MapR cluster, if requested
	if l.filterByMaprCluster != nil && *l.filterByMaprCluster != "" {
		items = filterItemsByMaprCluster(items, *l.filterByMaprCluster)
	}

	// filter items to only tickets for the specified MapR user, if requested
	if l.filterByMaprUser != nil && *l.filterByMaprUser != "" {
		items = filterItemsByMaprUser(items, *l.filterByMaprUser)
	}

	// filter items to only tickets for the specified UID, if requested
	if l.filterByUID != nil {
		items = filterItemsByUID(items, *l.filterByUID)
	}

	// filter items to only tickets for the specified GID, if requested
	if l.filterByGID != nil {
		items = filterItemsByGID(items, *l.filterByGID)
	}

	// filter items to only tickets that expire before the specified duration from now, if requested
	if l.filterExpiresBefore > 0 {
		items = filterExpiresBefore(items, l.filterExpiresBefore)
	}

	// enrich items with an InUse condition, if requested
	if l.showInUse || l.filterByInUse {
		items, err = l.enrichItemsWithInUseCondition(items)
		if err != nil {
			return nil, err
		}
	}

	// filter items to only tickets that are in use by a persistent volume, if requested
	if l.filterByInUse {
		items = filterItemsToOnlyInUse(items)
	}

	return items, nil
}

// filterSecretsWithMaprTicketKey filters secrets to only those that contain a MapR ticket key
func filterSecretsWithMaprTicketKey(secrets []coreV1.Secret) []coreV1.Secret {
	var filtered []coreV1.Secret

	for i := range secrets {
		secret := secrets[i]

		if ticket.SecretContainsMaprTicket(&secret) {
			filtered = append(filtered, secret)
		}
	}

	return filtered
}

// parseSecretsToItems parses secrets to items, ignoring secrets that don't contain a MapR ticket
func parseSecretsToItems(secrets []coreV1.Secret) []ListItem {
	var items []ListItem

	filtered := filterSecretsWithMaprTicketKey(secrets)

	for i := range filtered {
		secret := filtered[i]

		ticket, err := ticket.NewMaprTicketFromSecret(&secret)
		if err != nil {
			continue
		}

		items = append(items, ListItem{
			Secret: &secret,
			Ticket: ticket,
		})
	}

	return items
}

// filterItemsOnlyExpired filters items to only tickets that are expired already
func filterItemsOnlyExpired(items []ListItem) []ListItem {
	var filtered []ListItem

	for _, item := range items {
		if item.Ticket.IsExpired() {
			filtered = append(filtered, item)
		}
	}

	return filtered
}

// filterItemsOnlyUnexpired filters items to only tickets that are not expired yet
func filterItemsOnlyUnexpired(items []ListItem) []ListItem {
	var filtered []ListItem

	for _, item := range items {
		if !item.Ticket.IsExpired() {
			filtered = append(filtered, item)
		}
	}

	return filtered
}

// filterItemsByMaprCluster filters items to only tickets for the specified MapR cluster
func filterItemsByMaprCluster(items []ListItem, cluster string) []ListItem {
	var filtered []ListItem

	for _, item := range items {
		if item.Ticket.Cluster == cluster {
			filtered = append(filtered, item)
		}
	}

	return filtered
}

// filterItemsByMaprUser filters items to only tickets for the specified MapR user
func filterItemsByMaprUser(items []ListItem, user string) []ListItem {
	var filtered []ListItem

	for _, item := range items {
		if item.Ticket.UserCreds.GetUserName() == user {
			filtered = append(filtered, item)
		}
	}

	return filtered
}

// filterItemsByUID filters items to only tickets for the specified UID
func filterItemsByUID(items []ListItem, uid uint32) []ListItem {
	var filtered []ListItem

	for _, item := range items {
		if *item.Ticket.UserCreds.Uid == uid {
			filtered = append(filtered, item)
		}
	}

	return filtered
}

// filterItemsByGID filters items to only tickets for the specified GID
func filterItemsByGID(items []ListItem, gid uint32) []ListItem {
	var filtered []ListItem

	for _, item := range items {
		// check if GID is in the list of GIDs
		for _, gotGid := range item.Ticket.UserCreds.Gids {
			if gotGid == gid {
				filtered = append(filtered, item)
				break
			}
		}
	}

	return filtered
}

// enrichItemsWithInUseCondition enriches items with an InUse condition based on whether a
// persistent volume is using the ticket or not
func (l *Lister) enrichItemsWithInUseCondition(items []ListItem) ([]ListItem, error) {
	pvs, err := l.client.CoreV1().PersistentVolumes().List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		return nil, err
	}

	// Filter the volumes to only MapR CSI-based ones
	maprVolumes := volume.FilterVolumesToMaprCSI(pvs.Items)

	// check for each ticket if it is in use by a persistent volume
	for i := range items {
		for _, pv := range maprVolumes {
			if volume.UsesTicket(&pv, items[i].Secret.Name, items[i].Secret.Namespace) {
				items[i].InUse++
			}
		}
	}

	return items, nil
}

// filterItemsToOnlyInUse filters items to only tickets that are in use by a persistent volume
func filterItemsToOnlyInUse(items []ListItem) []ListItem {
	var filtered []ListItem

	for _, item := range items {
		if item.InUse > 0 {
			filtered = append(filtered, item)
		}
	}

	return filtered
}

// filterExpiresBefore filters items to only tickets that expire before the
// specified duration from now
func filterExpiresBefore(items []ListItem, expiresBefore time.Duration) []ListItem {
	var filtered []ListItem

	for _, item := range items {
		if item.Ticket.ExpiresBefore(expiresBefore) {
			filtered = append(filtered, item)
		}
	}

	return filtered
}
