// Copyright (c) 2024 Alexej Disterhoft
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
//
// SPX-License-Identifier: MIT

package claim

type ListerOption func(*Lister)

func WithSecretLister(secretLister secretLister) ListerOption {
	return func(l *Lister) {
		l.secretLister = secretLister
	}
}

func WithSortBy(sortBy []SortOption) ListerOption {
	return func(l *Lister) {
		l.sortBy = sortBy
	}
}
