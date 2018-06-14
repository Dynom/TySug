package finder

// Option is the type accepted by finder to set specific options
type Option func(sug *Finder)

// WithAlgorithm allows you to set any algorithm
func WithAlgorithm(alg Algorithm) Option {
	return func(s *Finder) {
		s.Alg = alg
	}
}
