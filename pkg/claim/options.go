package claim

type ListerOption func(*Lister)

func WithSecretLister(secretLister secretLister) ListerOption {
	return func(l *Lister) {
		l.secretLister = secretLister
	}
}

func WithSortBy(sortBy []SortOptions) ListerOption {
	return func(l *Lister) {
		l.sortBy = sortBy
	}
}
