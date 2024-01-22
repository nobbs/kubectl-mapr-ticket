package secret

import "time"

type ListerOption func(*Lister)

func WithSortBy(sortBy []SortOptions) ListerOption {
	return func(l *Lister) {
		l.sortBy = sortBy
	}
}

func WithFilterByMaprCluster(cluster string) ListerOption {
	return func(l *Lister) {
		l.filterByMaprCluster = &cluster
	}
}

func WithFilterByMaprUser(user string) ListerOption {
	return func(l *Lister) {
		l.filterByMaprUser = &user
	}
}

func WithFilterByUID(uid uint32) ListerOption {
	return func(l *Lister) {
		l.filterByUID = &uid
	}
}

func WithFilterByGID(gid uint32) ListerOption {
	return func(l *Lister) {
		l.filterByGID = &gid
	}
}

func WithFilterOnlyExpired() ListerOption {
	return func(l *Lister) {
		l.filterOnlyExpired = true
	}
}

func WithFilterOnlyUnexpired() ListerOption {
	return func(l *Lister) {
		l.filterOnlyUnexpired = true
	}
}

func WithFilterByInUse() ListerOption {
	return func(l *Lister) {
		l.filterByInUse = true
	}
}

func WithFilterExpiresBefore(expiresBefore time.Duration) ListerOption {
	return func(l *Lister) {
		l.filterExpiresBefore = expiresBefore
	}
}

func WithShowInUse() ListerOption {
	return func(l *Lister) {
		l.showInUse = true
	}
}

func WithVolumeLister(volumeLister volumeLister) ListerOption {
	return func(l *Lister) {
		l.volumeLister = volumeLister
	}
}
