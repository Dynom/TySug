package finder

import (
	"context"
	"errors"
	"math"
)

type AlgWrapper func(a, b string) float64

type Scorer struct {
	referenceMap map[string]struct{}
	reference    []string
	Alg          AlgWrapper
}

var (
	ErrNoAlgorithmDefined = errors.New("no algorithm defined")
)

// New creates a new instance of Scorer. The order of the list is significant
func New(list []string, options ...Option) (*Scorer, error) {
	i := &Scorer{
		referenceMap: make(map[string]struct{}, len(list)),
		reference:    list,
	}

	for _, r := range list {
		i.referenceMap[r] = struct{}{}
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
func (t Scorer) Find(input string) (string, float64) {
	return t.FindCtx(context.Background(), input)
}

func (t Scorer) FindCtx(ctx context.Context, input string) (string, float64) {

	// Exact matches
	if _, exists := t.referenceMap[input]; exists {
		return input, 1
	}

	var hs = math.Inf(-1)
	var best string
	for _, ref := range t.reference {
		select {
		case <-ctx.Done():
			return input, 0
		default:
		}

		if d := t.Alg(input, ref); d > hs {
			hs = d
			best = ref
		}
	}

	return best, hs
}
