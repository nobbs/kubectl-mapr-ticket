// Copyright (c) 2024 Alexej Disterhoft
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
//
// SPDX-License-Identifier: MIT

package secret

import "time"

// ListerOption is a function that can be used to configure the secret lister.
type ListerOption func(*Lister)

// WithSortBy configures the secret lister to sort the tickets by the given sort options.
func WithSortBy(sortBy []SortOption) ListerOption {
	return func(l *Lister) {
		l.sortBy = sortBy
	}
}

// WithFilterByMaprCluster configures the secret lister to only list tickets that are backed by a
// MapR cluster with the given name.
func WithFilterByMaprCluster(cluster string) ListerOption {
	return func(l *Lister) {
		l.filterByMaprCluster = &cluster
	}
}

// WithFilterByMaprUser configures the secret lister to only list tickets that are for a MapR user
// with the given name.
func WithFilterByMaprUser(user string) ListerOption {
	return func(l *Lister) {
		l.filterByMaprUser = &user
	}
}

// WithFilterByUID configures the secret lister to only list tickets that are for a user with the
// given UID.
func WithFilterByUID(uid uint32) ListerOption {
	return func(l *Lister) {
		l.filterByUID = &uid
	}
}

// WithFilterByGID configures the secret lister to only list tickets that include a group with the
// given GID.
func WithFilterByGID(gid uint32) ListerOption {
	return func(l *Lister) {
		l.filterByGID = &gid
	}
}

// WithFilterOnlyExpired configures the secret lister to only list tickets that have expired.
func WithFilterOnlyExpired() ListerOption {
	return func(l *Lister) {
		l.filterOnlyExpired = true
	}
}

// WithFilterOnlyUnexpired configures the secret lister to only list tickets that have not expired.
func WithFilterOnlyUnexpired() ListerOption {
	return func(l *Lister) {
		l.filterOnlyUnexpired = true
	}
}

// WithFilterByInUse configures the secret lister to only list tickets that are in use by a
// persistent volume.
func WithFilterByInUse() ListerOption {
	return func(l *Lister) {
		l.filterByInUse = true
	}
}

// WithFilterExpiresAfter configures the secret lister to only list tickets that expire during the
// given duration or have already expired.
func WithFilterExpiresBefore(expiresBefore time.Duration) ListerOption {
	return func(l *Lister) {
		l.filterExpiresBefore = expiresBefore
	}
}

// WithShowInUse configures the secret lister to show by how many persistent volumes a ticket is in
// use.
func WithShowInUse() ListerOption {
	return func(l *Lister) {
		l.showInUse = true
	}
}

// WithShowVolumes configures the secret lister to make use of the given volume lister to collect
// information about the persistent volumes that are using tickets.
func WithVolumeLister(volumeLister volumeLister) ListerOption {
	return func(l *Lister) {
		l.volumeLister = volumeLister
	}
}
