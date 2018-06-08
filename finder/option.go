package finder

type Option func(sug *Scorer)

func OptSetAlgorithm(alg AlgWrapper) Option {
	return func(s *Scorer) {
		s.Alg = alg
	}
}
