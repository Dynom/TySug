package finder

// Option is the type accepted by finder to set specific options
type Option func(sug *Finder)

// WithAlgorithm allows you to set any algorithm
func WithAlgorithm(alg Algorithm) Option {
	return func(s *Finder) {
		s.Alg = alg
	}
}

// WithLengthTolerance defines a percentage of length above we no longer consider a length difference a typo, but
// instead we consider it as "completely wrong". A value of 0.2 specifies a tolerance of at most ~20% difference in
// size, with a minimum of 1 character. A value of 0 (the default) disables this feature.
func WithLengthTolerance(t float64) Option {
	return func(s *Finder) {
		s.LengthTolerance = t
	}
}

// WithPrefixBuckets splits the reference list into buckets by their first letter. At a trade-off that the first
// character must be correct, this will significantly improve performance as it has a much smaller list to consider
func WithPrefixBuckets(enable bool) Option {
	return func(s *Finder) {
		if enable {
			s.bucketChars = 1
		}
	}
}
