package finder

// Option is the type accepted by finder to set specific options
type Option func(sug *Scorer)

// OptSetAlgorithm allows you to set any algorithm
func OptSetAlgorithm(alg Algorithm) Option {
	return func(s *Finder) {
		s.Alg = alg
	}
}
