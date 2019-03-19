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
