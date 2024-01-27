// Copyright (c) 2024 Alexej Disterhoft
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
//
// SPX-License-Identifier: MIT

package secret

import (
	"sort"

	"github.com/nobbs/kubectl-mapr-ticket/pkg/types"
)

// SortOption is the type of a sort option, basically a wrapper around a string to provide
// type safety.
type SortOption string

// All valid sort options are defined here
const (
	SortByName        SortOption = "name"
	SortByNamespace   SortOption = "namespace"
	SortByMaprCluster SortOption = "mapr.cluster"
	SortByMaprUser    SortOption = "mapr.user"
	SortByAge         SortOption = "age"
	SortByExpiration  SortOption = "expiration"
	SortByNumPVCs     SortOption = "npvcs"
)

var (
	// SortOptionsList is the list of valid sort options
	SortOptionsList = []string{
		SortByName.String(),
		SortByNamespace.String(),
		SortByMaprCluster.String(),
		SortByMaprUser.String(),
		SortByAge.String(),
		SortByExpiration.String(),
		SortByNumPVCs.String(),
	}

	// DefaultSortBy is the default sort order
	DefaultSortBy = []SortOption{
		SortByNamespace,
		SortByName,
	}
)

// String returns the string representation of the sort option.
func (s SortOption) String() string {
	return string(s)
}

// sortByName sorts the items by secret name
func sortByName(items []types.MaprSecret) {
	sort.Slice(items, func(i, j int) bool {
		return items[i].GetSecretName() < items[j].GetSecretName()
	})
}

// sortByNamespace sorts the items by secret namespace
func sortByNamespace(items []types.MaprSecret) {
	sort.Slice(items, func(i, j int) bool {
		return items[i].GetSecretNamespace() < items[j].GetSecretNamespace()
	})
}

// sortByMaprCluster sorts the items by MapR cluster that the ticket is for
func sortByMaprCluster(items []types.MaprSecret) {
	sort.Slice(items, func(i, j int) bool {
		return items[i].GetCluster() < items[j].GetCluster()
	})
}

// sortByMaprUser sorts the items by MapR user that the ticket is for
func sortByMaprUser(items []types.MaprSecret) {
	sort.Slice(items, func(i, j int) bool {
		return items[i].GetUser() < items[j].GetUser()
	})
}

// sortByAge sorts the items by creation timestamp of the ticket
func sortByAge(items []types.MaprSecret) {
	sort.Slice(items, func(i, j int) bool {
		return items[i].GetCreationTime().Before(items[j].GetCreationTime())
	})
}

// sortByExpiration sorts the items by expiry time of the ticket
func sortByExpiration(items []types.MaprSecret) {
	sort.Slice(items, func(i, j int) bool {
		return items[i].GetExpirationTime().Before(items[j].GetExpirationTime())
	})
}

// sortByNumPVCs sorts the items by the number of persistent volumes that are
// using the secret
func sortByNumPVCs(items []types.MaprSecret) {
	sort.Slice(items, func(i, j int) bool {
		return items[i].NumPVC < items[j].NumPVC
	})
}

// Sort sorts the items by the specified sort options, in reverse order of the
// order in which they are specified. This makes for a more natural sort result
// when using multiple sort options.
func (l *Lister) Sort() *Lister {
	// reverse the order of the sort options
	order := make([]SortOption, len(l.sortBy))
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
		case SortByAge:
			sortByAge(l.tickets)
		case SortByExpiration:
			sortByExpiration(l.tickets)
		case SortByNumPVCs:
			sortByNumPVCs(l.tickets)
		}
	}

	return l
}
