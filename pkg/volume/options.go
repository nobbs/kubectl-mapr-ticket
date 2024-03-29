// Copyright (c) 2024 Alexej Disterhoft
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
//
// SPDX-License-Identifier: MIT

package volume

// ListerOption is a function that can be used to configure the volume lister.
type ListerOption func(*Lister)

// WithSortBy sets the sort order used by the Lister for output
func WithSortBy(sortBy []SortOption) ListerOption {
	return func(l *Lister) {
		l.sortBy = sortBy
	}
}

// WithSecretLister sets the secret lister used by the Lister to collect secrets and tickets
// referenced by the volumes
func WithSecretLister(secretLister secretLister) ListerOption {
	return func(l *Lister) {
		l.secretLister = secretLister
	}
}
