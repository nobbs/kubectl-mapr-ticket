package claim

import (
	"fmt"
	"slices"
	"sort"

	"github.com/nobbs/kubectl-mapr-ticket/pkg/types"
	"github.com/nobbs/kubectl-mapr-ticket/pkg/util"
)

type SortOption string

const (
	SortByNamespace       SortOption = "namespace"
	SortByName            SortOption = "name"
	SortBySecretNamespace SortOption = "secret.namespace"
	SortBySecretName      SortOption = "secret.name"
	SortByVolumeName      SortOption = "volume.name"
	SortByVolumePath      SortOption = "volume.path"
	SortByVolumeHandle    SortOption = "volume.handle"
	SortByExpiration      SortOption = "expiration"
	SortByAge             SortOption = "age"
)

var (
	// SortOptionsList is the list of valid sort options
	SortOptionsList = []string{
		SortByNamespace.String(),
		SortByName.String(),
		SortBySecretNamespace.String(),
		SortBySecretName.String(),
		SortByVolumeName.String(),
		SortByVolumePath.String(),
		SortByVolumeHandle.String(),
		SortByExpiration.String(),
		SortByAge.String(),
	}

	// DefaultSortBy is the default sort order
	DefaultSortBy = []SortOption{
		SortByNamespace,
		SortByName,
	}
)

func (s SortOption) String() string {
	return string(s)
}

// ValidateSortOptions validates the specified sort options
func ValidateSortOptions(sortOptions []string) error {
	for _, sortOption := range sortOptions {
		if !slices.Contains(SortOptionsList, sortOption) {
			return fmt.Errorf("invalid sort option: %s. Must be one of: (%s)", sortOption, util.StringSliceToCommaSeparatedString(SortOptionsList))
		}
	}

	return nil
}

func sortByName(claims []types.MaprVolumeClaim) {
	sort.Slice(claims, func(i, j int) bool {
		return claims[i].Claim.GetName() < claims[j].Claim.GetName()
	})
}

func sortByNamespace(claims []types.MaprVolumeClaim) {
	sort.Slice(claims, func(i, j int) bool {
		return claims[i].Claim.GetNamespace() < claims[j].Claim.GetNamespace()
	})
}

func sortBySecretNamespace(claims []types.MaprVolumeClaim) {
	sort.Slice(claims, func(i, j int) bool {
		return claims[i].Volume.GetSecretNamespace() < claims[j].Volume.GetSecretNamespace()
	})
}

func sortBySecretName(claims []types.MaprVolumeClaim) {
	sort.Slice(claims, func(i, j int) bool {
		return claims[i].Volume.GetSecretName() < claims[j].Volume.GetSecretName()
	})
}

func sortByVolumeName(claims []types.MaprVolumeClaim) {
	sort.Slice(claims, func(i, j int) bool {
		return claims[i].Volume.GetName() < claims[j].Volume.GetName()
	})
}

func sortByVolumePath(claims []types.MaprVolumeClaim) {
	sort.Slice(claims, func(i, j int) bool {
		return claims[i].Volume.GetVolumePath() < claims[j].Volume.GetVolumePath()
	})
}

func sortByVolumeHandle(claims []types.MaprVolumeClaim) {
	sort.Slice(claims, func(i, j int) bool {
		return claims[i].Volume.GetVolumeHandle() < claims[j].Volume.GetVolumeHandle()
	})
}

func sortByExpiration(claims []types.MaprVolumeClaim) {
	sort.Slice(claims, func(i, j int) bool {
		return claims[i].Ticket.GetExpirationTime().Before(claims[j].Ticket.GetExpirationTime())
	})
}

func sortByAge(claims []types.MaprVolumeClaim) {
	sort.Slice(claims, func(i, j int) bool {
		return claims[i].Claim.CreationTimestamp.Before(&claims[j].Claim.CreationTimestamp)
	})
}

// sort sorts the items by the specified sort options, in reverse order of the
// order in which they are specified. This makes for a more natural sort result
// when using multiple sort options.
func (l *Lister) sort() *Lister {
	// reverse the order of the sort options
	order := make([]SortOption, len(l.sortBy))
	for i, j := 0, len(l.sortBy)-1; i < len(l.sortBy); i, j = i+1, j-1 {
		order[i] = l.sortBy[j]
	}

	// sort the items by the specified sort options
	for _, sortOption := range order {
		switch sortOption {
		case SortByNamespace:
			sortByNamespace(l.volumeClaims)
		case SortByName:
			sortByName(l.volumeClaims)
		case SortBySecretNamespace:
			sortBySecretNamespace(l.volumeClaims)
		case SortBySecretName:
			sortBySecretName(l.volumeClaims)
		case SortByVolumeName:
			sortByVolumeName(l.volumeClaims)
		case SortByVolumePath:
			sortByVolumePath(l.volumeClaims)
		case SortByVolumeHandle:
			sortByVolumeHandle(l.volumeClaims)
		case SortByExpiration:
			sortByExpiration(l.volumeClaims)
		case SortByAge:
			sortByAge(l.volumeClaims)
		}
	}

	return l
}
