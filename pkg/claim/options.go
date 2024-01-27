// Copyright (c) 2024 Alexej Disterhoft
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
//
// SPDX-License-Identifier: MIT

package claim

// ListerOption is a function that can be used to configure the volume claim lister.
type ListerOption func(*Lister)

// WithSecretLister configures the volume claim lister to use the given secret lister.
func WithSecretLister(secretLister secretLister) ListerOption {
	return func(l *Lister) {
		l.secretLister = secretLister
	}
}

// WithSortBy configures the volume claim lister to sort the volume claims by the given sort options.
func WithSortBy(sortBy []SortOption) ListerOption {
	return func(l *Lister) {
		l.sortBy = sortBy
	}
}
