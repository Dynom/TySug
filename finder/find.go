package finder

import "errors"

type AlgWrapper func(a, b string) float64

type TySug struct {
	reference map[string]struct{}
	Alg       AlgWrapper
}

var (
	ErrNoAlgorithmDefined = errors.New("no algorithm defined")
)

// New creates a new instance of TySug
func New(list []string, options ...Option) (*TySug, error) {
	i := &TySug{
		reference: make(map[string]struct{}, len(list)),
	}

	for _, r := range list {
		i.reference[r] = struct{}{}
	}

	for _, o := range options {
		o(i)
	}

	if i.Alg == nil {
		return i, ErrNoAlgorithmDefined
	}

	return i, nil
}



// Find returns the best alternative and a score. A score of 1 means a perfect match
func (t TySug) Find(input string) (string, float64) {

	// Exact matches
	if _, exists := t.reference[input]; exists {
		return input, 1
	}

	var hs float64
	var best string
	for ref := range t.reference {

		if d := t.Alg(input, ref); d > hs {
			hs = d
			best = ref
		}
	}

	return best, hs
}

