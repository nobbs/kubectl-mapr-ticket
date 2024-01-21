package secret

import (
	"fmt"
	"sort"

	apiSecret "github.com/nobbs/kubectl-mapr-ticket/pkg/api/secret"
)

type SortOptions string

const (
	SortByName              SortOptions = "name"
	SortByNamespace         SortOptions = "namespace"
	SortByMaprCluster       SortOptions = "maprCluster"
	SortByMaprUser          SortOptions = "maprUser"
	SortByCreationTimestamp SortOptions = "creationTimestamp"
	SortByExpiryTime        SortOptions = "expiryTime"
	SortByNumPVC            SortOptions = "numPVC"
)

var (
	// SortOptionsList is the list of valid sort options
	SortOptionsList = []string{
		SortByName.String(),
		SortByNamespace.String(),
		SortByMaprCluster.String(),
		SortByMaprUser.String(),
		SortByCreationTimestamp.String(),
		SortByExpiryTime.String(),
		SortByNumPVC.String(),
	}

	// DefaultSortBy is the default sort order
	DefaultSortBy = []SortOptions{
		SortByNamespace,
		SortByName,
	}
)

func (s SortOptions) String() string {
	return string(s)
}

// ValidateSortOptions validates the specified sort options
func ValidateSortOptions(sortOptions []string) error {
	for _, sortOption := range sortOptions {
		switch sortOption {
		case SortByName.String():
		case SortByNamespace.String():
		case SortByMaprCluster.String():
		case SortByMaprUser.String():
		case SortByCreationTimestamp.String():
		case SortByExpiryTime.String():
		case SortByNumPVC.String():
		default:
			return fmt.Errorf("invalid sort option: %s. Must be one of: name|namespace|maprCluster|maprUser|creationTimestamp|expiryTime", sortOption)
		}
	}

	return nil
}

// sortByName sorts the items by secret name
func sortByName(items []apiSecret.TicketSecret) {
	sort.Slice(items, func(i, j int) bool {
		return items[i].Secret.Name < items[j].Secret.Name
	})
}

// sortByNamespace sorts the items by secret namespace
func sortByNamespace(items []apiSecret.TicketSecret) {
	sort.Slice(items, func(i, j int) bool {
		return items[i].Secret.Namespace < items[j].Secret.Namespace
	})
}

// sortByMaprCluster sorts the items by MapR cluster that the ticket is for
func sortByMaprCluster(items []apiSecret.TicketSecret) {
	sort.Slice(items, func(i, j int) bool {
		return items[i].Ticket.Cluster < items[j].Ticket.Cluster
	})
}

// sortByMaprUser sorts the items by MapR user that the ticket is for
func sortByMaprUser(items []apiSecret.TicketSecret) {
	sort.Slice(items, func(i, j int) bool {
		return items[i].Ticket.UserCreds.GetUserName() < items[j].Ticket.UserCreds.GetUserName()
	})
}

// sortByCreationTimestamp sorts the items by creation timestamp of the ticket
func sortByCreationTimestamp(items []apiSecret.TicketSecret) {
	sort.Slice(items, func(i, j int) bool {
		return items[i].Ticket.CreationTime().Before(items[j].Ticket.CreationTime())
	})
}

// sortByExpiryTime sorts the items by expiry time of the ticket
func sortByExpiryTime(items []apiSecret.TicketSecret) {
	sort.Slice(items, func(i, j int) bool {
		return items[i].Ticket.ExpirationTime().Before(items[j].Ticket.ExpirationTime())
	})
}

// sortByNumPVC sorts the items by the number of persistent volumes that are
// using the secret
func sortByNumPVC(items []apiSecret.TicketSecret) {
	sort.Slice(items, func(i, j int) bool {
		return items[i].NumPVC < items[j].NumPVC
	})
}

// Sort sorts the items by the specified sort options, in reverse order of the
// order in which they are specified. This makes for a more natural sort result
// when using multiple sort options.
func (l *Lister) Sort() *Lister {
	// reverse the order of the sort options
	order := make([]SortOptions, len(l.sortBy))
	for i, j := 0, len(l.sortBy)-1; i < len(l.sortBy); i, j = i+1, j-1 {
		order[i] = l.sortBy[j]
	}

	// sort the items by each sort option
	for _, sortOption := range order {
		switch sortOption {
		case SortByName:
			sortByName(l.tickets)
		case SortByNamespace:
			sortByNamespace(l.tickets)
		case SortByMaprCluster:
			sortByMaprCluster(l.tickets)
		case SortByMaprUser:
			sortByMaprUser(l.tickets)
		case SortByCreationTimestamp:
			sortByCreationTimestamp(l.tickets)
		case SortByExpiryTime:
			sortByExpiryTime(l.tickets)
		case SortByNumPVC:
			sortByNumPVC(l.tickets)
		}
	}

	return l
}
