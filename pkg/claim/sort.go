package claim

import (
	"fmt"
	"sort"

	"github.com/nobbs/kubectl-mapr-ticket/pkg/types"
)

type SortOptions string

const (
	SortByNamespace       SortOptions = "namespace"
	SortByName            SortOptions = "name"
	SortBySecretNamespace SortOptions = "secretNamespace"
	SortBySecretName      SortOptions = "secretName"
	SortByVolumeName      SortOptions = "volumeName"
	SortByVolumePath      SortOptions = "volumePath"
	SortByVolumeHandle    SortOptions = "volumeHandle"
	SortByExpiryTime      SortOptions = "expiryTime"
	SortByAge             SortOptions = "age"
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
		SortByExpiryTime.String(),
		SortByAge.String(),
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
		case SortByNamespace.String():
		case SortByName.String():
		case SortBySecretNamespace.String():
		case SortBySecretName.String():
		case SortByVolumeName.String():
		case SortByVolumePath.String():
		case SortByVolumeHandle.String():
		case SortByExpiryTime.String():
		case SortByAge.String():
		default:
			return fmt.Errorf("invalid sort option: %s. Must be one of: namespace|name|secretNamespace|secretName|volumeName|volumePath|volumeHandle|expiryTime|age", sortOption)
		}
	}
	return nil
}

func sortByName(claims []types.VolumeClaim) {
	sort.Slice(claims, func(i, j int) bool {
		return claims[i].Claim.GetName() < claims[j].Claim.GetName()
	})
}

func sortByNamespace(claims []types.VolumeClaim) {
	sort.Slice(claims, func(i, j int) bool {
		return claims[i].Claim.GetNamespace() < claims[j].Claim.GetNamespace()
	})
}

func sortBySecretNamespace(claims []types.VolumeClaim) {
	sort.Slice(claims, func(i, j int) bool {
		return claims[i].Volume.GetSecretNamespace() < claims[j].Volume.GetSecretNamespace()
	})
}

func sortBySecretName(claims []types.VolumeClaim) {
	sort.Slice(claims, func(i, j int) bool {
		return claims[i].Volume.GetSecretName() < claims[j].Volume.GetSecretName()
	})
}

func sortByVolumeName(claims []types.VolumeClaim) {
	sort.Slice(claims, func(i, j int) bool {
		return claims[i].Volume.GetName() < claims[j].Volume.GetName()
	})
}

func sortByVolumePath(claims []types.VolumeClaim) {
	sort.Slice(claims, func(i, j int) bool {
		return claims[i].Volume.GetVolumePath() < claims[j].Volume.GetVolumePath()
	})
}

func sortByVolumeHandle(claims []types.VolumeClaim) {
	sort.Slice(claims, func(i, j int) bool {
		return claims[i].Volume.GetVolumeHandle() < claims[j].Volume.GetVolumeHandle()
	})
}

func sortByExpiryTime(claims []types.VolumeClaim) {
	sort.Slice(claims, func(i, j int) bool {
		return claims[i].Ticket.GetExpirationTime().Before(claims[j].Ticket.GetExpirationTime())
	})
}

func sortByAge(claims []types.VolumeClaim) {
	sort.Slice(claims, func(i, j int) bool {
		return claims[i].Claim.CreationTimestamp.Before(&claims[j].Claim.CreationTimestamp)
	})
}

// sort sorts the items by the specified sort options, in reverse order of the
// order in which they are specified. This makes for a more natural sort result
// when using multiple sort options.
func (l *Lister) sort() *Lister {
	// reverse the order of the sort options
	order := make([]SortOptions, len(l.sortBy))
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
		case SortByExpiryTime:
			sortByExpiryTime(l.volumeClaims)
		case SortByAge:
			sortByAge(l.volumeClaims)
		}
	}

	return l
}
