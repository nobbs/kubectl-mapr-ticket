// Copyright (c) 2024 Alexej Disterhoft
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
//
// SPX-License-Identifier: MIT

package secret

import (
	"context"
	"time"

	"github.com/nobbs/kubectl-mapr-ticket/pkg/ticket"
	"github.com/nobbs/kubectl-mapr-ticket/pkg/types"

	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type volumeLister interface {
	List() ([]types.MaprVolume, error)
}

type Lister struct {
	client       kubernetes.Interface
	volumeLister volumeLister

	namespace           string
	filterOnlyExpired   bool
	filterOnlyUnexpired bool
	filterByMaprCluster *string
	filterByMaprUser    *string
	filterByUID         *uint32
	filterByGID         *uint32
	filterByInUse       bool
	filterExpiresBefore time.Duration
	showInUse           bool
	sortBy              []SortOption

	tickets []types.MaprSecret
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
		sortBy:              DefaultSortBy,
		filterOnlyExpired:   defaultFilterOnlyExpired,
		filterOnlyUnexpired: defaultFilterOnlyUnexpired,
	}

	for _, opt := range opts {
		opt(l)
	}

	return l
}

func (l *Lister) List() ([]types.MaprSecret, error) {
	if err := l.getSecretsWithTickets(); err != nil {
		return nil, err
	}

	// run all filters and sorts
	l.filterTicketsOnlyExpired().
		filterTicketsOnlyUnexpired().
		filterTicketsByMaprCluster().
		filterTicketsByMaprUser().
		filterTicketsByUID().
		filterTicketsByGID().
		filterTicketsExpiresBefore().
		collectPVsUsingTickets().
		filterTicketsInUse().
		Sort()

	return l.tickets, nil
}

// getSecretsWithTickets retrieves the list of ticket secrets
func (l *Lister) getSecretsWithTickets() error {
	secrets, err := l.client.CoreV1().Secrets(l.namespace).List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		return err
	}

	// convert secrets to items, parse all tickets
	l.tickets = parseTicketsFromSecrets(secrets.Items)

	return nil
}

// rejectSecretsWithoutTicket filters secrets to only those that contain a MapR ticket key
func rejectSecretsWithoutTicket(secrets []coreV1.Secret) []coreV1.Secret {
	var filtered []coreV1.Secret

	for i := range secrets {
		secret := secrets[i]

		if ticket.SecretContainsMaprTicket(&secret) {
			filtered = append(filtered, secret)
		}
	}

	return filtered
}

// parseTicketsFromSecrets parses secrets to items, ignoring secrets that don't contain a MapR ticket
func parseTicketsFromSecrets(secrets []coreV1.Secret) []types.MaprSecret {
	items := make([]types.MaprSecret, 0, len(secrets))

	filtered := rejectSecretsWithoutTicket(secrets)

	for i := range filtered {
		s := filtered[i]

		ticket, err := ticket.NewMaprTicketFromSecret(&s)
		if err != nil {
			continue
		}

		items = append(items, types.MaprSecret{
			Secret: (*types.Secret)(&s),
			Ticket: ticket,
		})
	}

	return items
}

// filterTicketsOnlyExpired filters tickets to only those that are expired
func (l *Lister) filterTicketsOnlyExpired() *Lister {
	// if the filter is not enabled, we can skip this step
	if !l.filterOnlyExpired {
		return l
	}

	var filtered []types.MaprSecret

	for _, item := range l.tickets {
		if item.Ticket.IsExpired() {
			filtered = append(filtered, item)
		}
	}

	l.tickets = filtered

	return l
}

// filterTicketsOnlyUnexpired filters tickets to only those that are not expired
func (l *Lister) filterTicketsOnlyUnexpired() *Lister {
	// if the filter is not enabled, we can skip this step
	if !l.filterOnlyUnexpired {
		return l
	}

	var filtered []types.MaprSecret

	for _, item := range l.tickets {
		if !item.Ticket.IsExpired() {
			filtered = append(filtered, item)
		}
	}

	l.tickets = filtered

	return l
}

// filterTicketsByMaprCluster filters tickets to only those that match the specified MapR cluster
func (l *Lister) filterTicketsByMaprCluster() *Lister {
	// if the filter is not enabled, we can skip this step
	if l.filterByMaprCluster == nil || *l.filterByMaprCluster == "" {
		return l
	}

	var filtered []types.MaprSecret

	for _, item := range l.tickets {
		if item.Ticket.Cluster == *l.filterByMaprCluster {
			filtered = append(filtered, item)
		}
	}

	l.tickets = filtered

	return l
}

// filterTicketsByMaprUser filters tickets to only those that match the specified MapR user
func (l *Lister) filterTicketsByMaprUser() *Lister {
	// if the filter is not enabled, we can skip this step
	if l.filterByMaprUser == nil || *l.filterByMaprUser == "" {
		return l
	}

	var filtered []types.MaprSecret

	for _, item := range l.tickets {
		if item.Ticket.UserCreds.GetUserName() == *l.filterByMaprUser {
			filtered = append(filtered, item)
		}
	}

	l.tickets = filtered

	return l
}

// filterTicketsByUID filters tickets to only those that match the specified UID
func (l *Lister) filterTicketsByUID() *Lister {
	// if the filter is not enabled, we can skip this step
	if l.filterByUID == nil {
		return l
	}

	var filtered []types.MaprSecret

	for _, item := range l.tickets {
		if *item.Ticket.UserCreds.Uid == *l.filterByUID {
			filtered = append(filtered, item)
		}
	}

	l.tickets = filtered

	return l
}

// filterTicketsByGID filters tickets to only those that match the specified GID
func (l *Lister) filterTicketsByGID() *Lister {
	// if the filter is not enabled, we can skip this step
	if l.filterByGID == nil {
		return l
	}

	var filtered []types.MaprSecret

	for _, item := range l.tickets {
		for _, gid := range item.Ticket.UserCreds.Gids {
			if gid == *l.filterByGID {
				filtered = append(filtered, item)
				break
			}
		}
	}

	l.tickets = filtered

	return l
}

// filterTicketsExpiresBefore filters tickets to only those that expire before the specified
// duration
func (l *Lister) filterTicketsExpiresBefore() *Lister {
	// if the filter is not enabled, we can skip this step
	if l.filterExpiresBefore <= 0 {
		return l
	}

	var filtered []types.MaprSecret

	for _, item := range l.tickets {
		if item.Ticket.ExpiresBefore(l.filterExpiresBefore) {
			filtered = append(filtered, item)
		}
	}

	l.tickets = filtered

	return l
}

// filterTicketsInUse filters tickets to only those that are in use by a persistent volume
func (l *Lister) filterTicketsInUse() *Lister {
	// if the filter is not enabled, we can skip this step
	if !l.filterByInUse {
		return l
	}

	var filtered []types.MaprSecret

	for _, item := range l.tickets {
		if item.NumPVC > 0 {
			filtered = append(filtered, item)
		}
	}

	l.tickets = filtered

	return l
}

// collectPVsUsingTickets enriches the ticket items with the number of PVCs using the ticket
func (l *Lister) collectPVsUsingTickets() *Lister {
	// if we don't have a volume lister, we need to skip this step
	if l.volumeLister == nil {
		return l
	}

	// if we don't need to show in use, or filter by in use, we can skip this step
	if !l.showInUse && !l.filterByInUse {
		return l
	}

	// get all persistent volumes
	pvs, err := l.volumeLister.List()
	if err != nil {
		return l
	}

	// check for each ticket if it is in use by a persistent volume
	for i := range l.tickets {
		for _, volume := range pvs {
			if volume.Volume.UsesSecret(l.tickets[i].Secret.Namespace, l.tickets[i].Secret.Name) {
				l.tickets[i].NumPVC++
			}
		}
	}

	return l
}
