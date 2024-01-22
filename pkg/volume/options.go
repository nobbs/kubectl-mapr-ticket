package volume

type ListerOption func(*Lister)

// WithSortBy sets the sort order used by the Lister for output
func WithSortBy(sortBy []SortOptions) ListerOption {
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
