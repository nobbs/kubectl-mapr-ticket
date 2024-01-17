package volume

import (
	"fmt"
	"sort"

	apiVolume "github.com/nobbs/kubectl-mapr-ticket/pkg/api/volume"
)

type SortOptions string

const (
	SortByName            SortOptions = "name"
	SortBySecretNamespace SortOptions = "secretNamespace"
	SortBySecretName      SortOptions = "secretName"
	SortByClaimNamespace  SortOptions = "claimNamespace"
	SortByClaimName       SortOptions = "claimName"
	SortByVolumePath      SortOptions = "volumePath"
	SortByVolumeHandle    SortOptions = "volumeHandle"
	SortByExpiryTime      SortOptions = "expiryTime"
	SortByAge             SortOptions = "age"
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
		SortByExpiryTime.String(),
		SortByAge.String(),
	}

	// DefaultSortBy is the default sort order
	DefaultSortBy = []SortOptions{
		SortBySecretNamespace,
		SortBySecretName,
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
		case SortBySecretNamespace.String():
		case SortBySecretName.String():
		case SortByClaimNamespace.String():
		case SortByClaimName.String():
		case SortByVolumePath.String():
		case SortByVolumeHandle.String():
		case SortByExpiryTime.String():
		case SortByAge.String():
		default:
			return fmt.Errorf("invalid sort option: %s. Must be one of: name|secretNamespace|secretName|claimNamespace|claimName|volumePath|volumeHandle|expiryTime|age", sortOption)
		}
	}
	return nil
}

func sortByName(volumes []apiVolume.Volume) {
	sort.Slice(volumes, func(i, j int) bool {
		return volumes[i].Name() < volumes[j].Name()
	})
}

func sortBySecretNamespace(volumes []apiVolume.Volume) {
	sort.Slice(volumes, func(i, j int) bool {
		return volumes[i].SecretNamespace() < volumes[j].SecretNamespace()
	})
}

func sortBySecretName(volumes []apiVolume.Volume) {
	sort.Slice(volumes, func(i, j int) bool {
		return volumes[i].SecretName() < volumes[j].SecretName()
	})
}

func sortByClaimNamespace(volumes []apiVolume.Volume) {
	sort.Slice(volumes, func(i, j int) bool {
		return volumes[i].ClaimNamespace() < volumes[j].ClaimNamespace()
	})
}

func sortByClaimName(volumes []apiVolume.Volume) {
	sort.Slice(volumes, func(i, j int) bool {
		return volumes[i].ClaimName() < volumes[j].ClaimName()
	})
}

func sortByVolumePath(volumes []apiVolume.Volume) {
	sort.Slice(volumes, func(i, j int) bool {
		return volumes[i].VolumePath() < volumes[j].VolumePath()
	})
}

func sortByVolumeHandle(volumes []apiVolume.Volume) {
	sort.Slice(volumes, func(i, j int) bool {
		return volumes[i].VolumeHandle() < volumes[j].VolumeHandle()
	})
}

func sortByExpiryTime(volumes []apiVolume.Volume) {
	sort.Slice(volumes, func(i, j int) bool {
		if volumes[i].Ticket == nil {
			return true
		} else if volumes[j].Ticket == nil {
			return false
		}

		return volumes[i].Ticket.Ticket.ExpirationTime().Before(volumes[j].Ticket.Ticket.ExpirationTime())
	})
}

func sortByAge(volumes []apiVolume.Volume) {
	sort.Slice(volumes, func(i, j int) bool {
		return volumes[i].Volume.CreationTimestamp.Time.Before(volumes[j].Volume.CreationTimestamp.Time)
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
		case SortByExpiryTime:
			sortByExpiryTime(l.volumes)
		case SortByAge:
			sortByAge(l.volumes)
		}
	}

	return l
}
