package secret

import (
	"fmt"
	"sort"
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

// ValidateSortOptions validates the specified sort options
func ValidateSortOptions(sortOptions []string) error {
	for _, sortOption := range sortOptions {
		switch sortOption {
		case string(SortByName):
		case string(SortByNamespace):
		case string(SortByMaprCluster):
		case string(SortByMaprUser):
		case string(SortByCreationTimestamp):
		case string(SortByExpiryTime):
		case string(SortByNumPVC):
		default:
			return fmt.Errorf("invalid sort option: %s. Must be one of: name|namespace|maprCluster|maprUser|creationTimestamp|expiryTime", sortOption)
		}
	}

	return nil
}

// sortByName sorts the items by secret name
func sortByName(items []ListItem) {
	sort.Slice(items, func(i, j int) bool {
		return items[i].Secret.Name < items[j].Secret.Name
	})
}

// sortByNamespace sorts the items by secret namespace
func sortByNamespace(items []ListItem) {
	sort.Slice(items, func(i, j int) bool {
		return items[i].Secret.Namespace < items[j].Secret.Namespace
	})
}

// sortByMaprCluster sorts the items by MapR cluster that the ticket is for
func sortByMaprCluster(items []ListItem) {
	sort.Slice(items, func(i, j int) bool {
		return items[i].Ticket.Cluster < items[j].Ticket.Cluster
	})
}

// sortByMaprUser sorts the items by MapR user that the ticket is for
func sortByMaprUser(items []ListItem) {
	sort.Slice(items, func(i, j int) bool {
		return items[i].Ticket.UserCreds.GetUserName() < items[j].Ticket.UserCreds.GetUserName()
	})
}

// sortByCreationTimestamp sorts the items by creation timestamp of the ticket
func sortByCreationTimestamp(items []ListItem) {
	sort.Slice(items, func(i, j int) bool {
		return items[i].Ticket.CreationTime().Before(items[j].Ticket.CreationTime())
	})
}

// sortByExpiryTime sorts the items by expiry time of the ticket
func sortByExpiryTime(items []ListItem) {
	sort.Slice(items, func(i, j int) bool {
		return items[i].Ticket.ExpirationTime().Before(items[j].Ticket.ExpirationTime())
	})
}

// sortByNumPVC sorts the items by the number of persistent volumes that are
// using the secret
func sortByNumPVC(items []ListItem) {
	sort.Slice(items, func(i, j int) bool {
		return items[i].NumPVC < items[j].NumPVC
	})
}

// Sort sorts the items by the specified sort options, in reverse order of the
// order in which they are specified. This makes for a more natural sort result
// when using multiple sort options.
func Sort(items []ListItem, sortOptions []SortOptions) {
	// reverse the order of the sort options
	order := make([]SortOptions, len(sortOptions))
	for i, j := 0, len(sortOptions)-1; i < len(sortOptions); i, j = i+1, j-1 {
		order[i] = sortOptions[j]
	}

	// sort the items by each sort option
	for _, sortOption := range order {
		switch sortOption {
		case SortByName:
			sortByName(items)
		case SortByNamespace:
			sortByNamespace(items)
		case SortByMaprCluster:
			sortByMaprCluster(items)
		case SortByMaprUser:
			sortByMaprUser(items)
		case SortByCreationTimestamp:
			sortByCreationTimestamp(items)
		case SortByExpiryTime:
			sortByExpiryTime(items)
		case SortByNumPVC:
			sortByNumPVC(items)
		}
	}
}
