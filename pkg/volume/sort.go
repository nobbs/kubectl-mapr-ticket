package volume

import (
	"fmt"
	"slices"
	"sort"

	"github.com/nobbs/kubectl-mapr-ticket/pkg/types"
	"github.com/nobbs/kubectl-mapr-ticket/pkg/util"
)

type SortOption string

const (
	SortByName            SortOption = "name"
	SortBySecretNamespace SortOption = "secret.namespace"
	SortBySecretName      SortOption = "secret.name"
	SortByClaimNamespace  SortOption = "claim.namespace"
	SortByClaimName       SortOption = "claim.name"
	SortByVolumePath      SortOption = "volume.path"
	SortByVolumeHandle    SortOption = "volume.handle"
	SortByExpiration      SortOption = "expiration"
	SortByAge             SortOption = "age"
)

var (
	// SortOptionsList is the list of valid sort options
	SortOptionsList = []string{
		SortByName.String(),
		SortBySecretNamespace.String(),
		SortBySecretName.String(),
		SortByClaimNamespace.String(),
		SortByClaimName.String(),
		SortByVolumePath.String(),
		SortByVolumeHandle.String(),
		SortByExpiration.String(),
		SortByAge.String(),
	}

	// DefaultSortBy is the default sort order
	DefaultSortBy = []SortOption{
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

func sortByName(volumes []types.Volume) {
	sort.Slice(volumes, func(i, j int) bool {
		return volumes[i].Volume.GetName() < volumes[j].Volume.GetName()
	})
}

func sortBySecretNamespace(volumes []types.Volume) {
	sort.Slice(volumes, func(i, j int) bool {
		return volumes[i].Volume.GetSecretNamespace() < volumes[j].Volume.GetSecretNamespace()
	})
}

func sortBySecretName(volumes []types.Volume) {
	sort.Slice(volumes, func(i, j int) bool {
		return volumes[i].Volume.GetSecretName() < volumes[j].Volume.GetSecretName()
	})
}

func sortByClaimNamespace(volumes []types.Volume) {
	sort.Slice(volumes, func(i, j int) bool {
		return volumes[i].Volume.GetClaimNamespace() < volumes[j].Volume.GetClaimNamespace()
	})
}

func sortByClaimName(volumes []types.Volume) {
	sort.Slice(volumes, func(i, j int) bool {
		return volumes[i].Volume.GetClaimName() < volumes[j].Volume.GetClaimName()
	})
}

func sortByVolumePath(volumes []types.Volume) {
	sort.Slice(volumes, func(i, j int) bool {
		return volumes[i].Volume.GetVolumePath() < volumes[j].Volume.GetVolumePath()
	})
}

func sortByVolumeHandle(volumes []types.Volume) {
	sort.Slice(volumes, func(i, j int) bool {
		return volumes[i].Volume.GetVolumeHandle() < volumes[j].Volume.GetVolumeHandle()
	})
}

func sortByExpiration(volumes []types.Volume) {
	sort.Slice(volumes, func(i, j int) bool {
		if volumes[i].Ticket == nil {
			return true
		} else if volumes[j].Ticket == nil {
			return false
		}

		return volumes[i].Ticket.GetExpirationTime().Before(volumes[j].Ticket.GetExpirationTime())
	})
}

func sortByAge(volumes []types.Volume) {
	sort.Slice(volumes, func(i, j int) bool {
		return volumes[i].Volume.CreationTimestamp.Before(&volumes[j].Volume.CreationTimestamp)
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
		case SortByName:
			sortByName(l.volumes)
		case SortBySecretNamespace:
			sortBySecretNamespace(l.volumes)
		case SortBySecretName:
			sortBySecretName(l.volumes)
		case SortByClaimNamespace:
			sortByClaimNamespace(l.volumes)
		case SortByClaimName:
			sortByClaimName(l.volumes)
		case SortByVolumePath:
			sortByVolumePath(l.volumes)
		case SortByVolumeHandle:
			sortByVolumeHandle(l.volumes)
		case SortByExpiration:
			sortByExpiration(l.volumes)
		case SortByAge:
			sortByAge(l.volumes)
		}
	}

	return l
}
