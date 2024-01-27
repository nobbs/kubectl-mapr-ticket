// Copyright (c) 2024 Alexej Disterhoft
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
//
// SPDX-License-Identifier: MIT

package volume

import (
	"sort"

	"github.com/nobbs/kubectl-mapr-ticket/pkg/types"
)

// SortOption is the type of a sort option, basically a wrapper around a string to provide
// type safety.
type SortOption string

// All valid sort options are defined here
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

// String returns the string representation of the sort option.
func (s SortOption) String() string {
	return string(s)
}

func sortByName(volumes []types.MaprVolume) {
	sort.Slice(volumes, func(i, j int) bool {
		return volumes[i].Volume.GetName() < volumes[j].Volume.GetName()
	})
}

func sortBySecretNamespace(volumes []types.MaprVolume) {
	sort.Slice(volumes, func(i, j int) bool {
		return volumes[i].Volume.GetSecretNamespace() < volumes[j].Volume.GetSecretNamespace()
	})
}

func sortBySecretName(volumes []types.MaprVolume) {
	sort.Slice(volumes, func(i, j int) bool {
		return volumes[i].Volume.GetSecretName() < volumes[j].Volume.GetSecretName()
	})
}

func sortByClaimNamespace(volumes []types.MaprVolume) {
	sort.Slice(volumes, func(i, j int) bool {
		return volumes[i].Volume.GetClaimNamespace() < volumes[j].Volume.GetClaimNamespace()
	})
}

func sortByClaimName(volumes []types.MaprVolume) {
	sort.Slice(volumes, func(i, j int) bool {
		return volumes[i].Volume.GetClaimName() < volumes[j].Volume.GetClaimName()
	})
}

func sortByVolumePath(volumes []types.MaprVolume) {
	sort.Slice(volumes, func(i, j int) bool {
		return volumes[i].Volume.GetVolumePath() < volumes[j].Volume.GetVolumePath()
	})
}

func sortByVolumeHandle(volumes []types.MaprVolume) {
	sort.Slice(volumes, func(i, j int) bool {
		return volumes[i].Volume.GetVolumeHandle() < volumes[j].Volume.GetVolumeHandle()
	})
}

func sortByExpiration(volumes []types.MaprVolume) {
	sort.Slice(volumes, func(i, j int) bool {
		if volumes[i].Ticket == nil {
			return true
		} else if volumes[j].Ticket == nil {
			return false
		}

		return volumes[i].Ticket.GetExpirationTime().Before(volumes[j].Ticket.GetExpirationTime())
	})
}

func sortByAge(volumes []types.MaprVolume) {
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
